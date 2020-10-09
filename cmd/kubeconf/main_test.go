package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/mumoshu/kubeconf/pkg/kubeconf"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

const (
	VarRunMainForTesting          = "RUN_MAIN_FOR_TESTING"
	VarRunMainForTestingArgPrefix = "RUN_MAIN_FOR_TESTING_ARG_"
	TraceServerURLEnv             = "TRACE_SERVER_URL"
)

type cond struct {
	subject        string
	input          string
	args           []string
	wantStdout     string
	wantStderr     string
	wantExitCode   int
	wantStdoutFile string
	commands       []cond
	env            []string
}

func TestKubeconfMain(t *testing.T) {
	if os.Getenv(VarRunMainForTesting) == "1" {
		args := envToArgs(os.Environ())

		os.Args = args

		tracePost(args)

		// We DO call hcledit's main() here. So this looks like a normal `hcledit` process.
		main()

		// If `main()` did not call os.Exit(0) explicitly, we assume there was no error hence it's safe to call os.Exit(0)
		// on behalf of go runtime.
		os.Exit(0)

		// As main() or this block calls os.Exit, we never reach this line.
		// But the test called this block of code catches and verifies the exit code.
		return
	}

	runCheck(t, "-h should not fail", cond{
		args:           []string{"kubeconf", "-h"},
		wantExitCode:   0,
		wantStdoutFile: "testdata/help",
	})

	runCheck(t, "help should not fail", cond{
		args:           []string{"kubeconf", "help"},
		wantExitCode:   0,
		wantStdoutFile: "testdata/help",
	})

	t.Run("acceptance", func(t *testing.T) {
		dir := t.TempDir()

		homeEnv := "HOME=" + dir

		check(t, cond{
			args:         []string{"kubeconf", "current"},
			env:          []string{homeEnv},
			wantStderr:   "error: current config is not set. please run `init`\n",
			wantExitCode: 1,
		})

		check(t, cond{
			args:         []string{"kubeconf", "ls"},
			env:          []string{homeEnv},
			wantStderr:   "error: open " + dir + "/.kube/kubeconf/registry: no such file or directory\n",
			wantExitCode: 1,
		})

		check(t, cond{
			args:         []string{"kubeconf", "init"},
			env:          []string{homeEnv},
			wantExitCode: 1,
			wantStderr:   "error: copying " + dir + "/.kube/config to " + dir + "/.kube/kubeconf/registry/default: open " + dir + "/.kube/config: no such file or directory\n",
		})

		if err := kubeconf.CopyFile("testdata/config1", filepath.Join(dir, ".kube", "config")); err != nil {
			t.Fatal(err)
		}

		traces := withTracing(t, func(url string) {
			check(t, cond{
				args:         []string{"kubeconf", "init"},
				env:          []string{homeEnv, TraceServerURLEnv + "=" + url},
				wantExitCode: 0,
				wantStderr:   "✔ Config default created.\n",
			})
		})

		if numTraces := len(traces); numTraces != 1 {
			t.Fatalf("unexpected number of traces: %v", numTraces)
		}

		checkFileContentEquality(t, "testdata/config1", filepath.Join(dir, ".kube/config"))

		check(t, cond{
			args:       []string{"kubeconf", "ls"},
			env:        []string{homeEnv},
			wantStdout: "default\n",
		})

		check(t, cond{
			args:         []string{"kubeconf", "import", "testdata/invalid", "config2"},
			env:          []string{homeEnv},
			wantStderr:   "error: copying testdata/invalid to " + dir + "/.kube/kubeconf/registry/config2: open testdata/invalid: no such file or directory\n",
			wantExitCode: 1,
		})

		check(t, cond{
			args:         []string{"kubeconf", "locate", "config2"},
			env:          []string{homeEnv},
			wantStderr:   "error: config config2 does not exist\n",
			wantExitCode: 1,
		})

		check(t, cond{
			args:         []string{"kubeconf", "import", "testdata/config2", "config2"},
			env:          []string{homeEnv},
			wantStderr:   "✔ Config config2 created.\n",
			wantExitCode: 0,
		})

		check(t, cond{
			args:         []string{"kubeconf", "locate", "config2"},
			env:          []string{homeEnv},
			wantStdout:   dir + "/.kube/kubeconf/registry/config2",
			wantExitCode: 0,
		})

		check(t, cond{
			args:       []string{"kubeconf", "ls"},
			env:        []string{homeEnv},
			wantStdout: "config2\ndefault\n",
		})

		check(t, cond{
			args:       []string{"kubeconf", "current"},
			env:        []string{homeEnv},
			wantStdout: "default\n",
		})

		check(t, cond{
			args:         []string{"kubeconf", "use", "invalid"},
			env:          []string{homeEnv},
			wantStderr:   "error: failed to switch config: config invalid does not exist\n",
			wantExitCode: 1,
		})

		checkFileContentEquality(t, "testdata/config1", filepath.Join(dir, ".kube/config"))

		check(t, cond{
			args:         []string{"kubeconf", "use", "config2"},
			env:          []string{homeEnv},
			wantStderr:   "✔ Switched to config \"config2\".\n",
			wantExitCode: 0,
		})

		checkFileContentEquality(t, "testdata/config2", filepath.Join(dir, ".kube/config"))
	})
}

func checkFileContentEquality(t *testing.T, wantPath string, gotPath string) {
	t.Helper()

	want, err := ioutil.ReadFile(wantPath)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ioutil.ReadFile(gotPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) == 0 {
		t.Fatal("unexpected empty .kube/config written")
	}

	if d := cmp.Diff(string(want), string(got)); d != "" {
		t.Fatalf("unexpected config diff:\n%s", d)
	}
}

type trace struct {
	Args []string `json:"args"`
}

func tracePost(args []string) {
	url := os.Getenv(TraceServerURLEnv)
	if url == "" {
		return
	}

	trace := trace{
		Args: args,
	}

	js, err := json.Marshal(trace)
	if err != nil {
		panic(err)
	}

	_, err = http.Post(url, "application/json", strings.NewReader(string(js)))
	if err != nil {
		panic(err)
	}
}

func withTracing(t *testing.T, f func(string)) []trace {
	t.Helper()

	var traces []trace

	mux := &http.ServeMux{}
	path := "/traces"
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			// Decode the JSON in the body and overwrite 'tom' with it
			d := json.NewDecoder(r.Body)
			tr := trace{}
			err := d.Decode(&tr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			traces = append(traces, tr)

			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "unsupported method")
		}
	})

	l, err := net.Listen("tcp", ":18080")
	if err != nil {
		t.Fatal(err)
	}

	srv := &http.Server{Handler: mux}

	go func() {
		srv.Serve(l)
	}()
	defer srv.Shutdown(context.Background())

	f(fmt.Sprintf("http://127.0.0.1:18080%s", path))

	return traces
}

func runCheck(t *testing.T, subject string, tc cond) {
	t.Helper()

	t.Run(subject, func(t *testing.T) {
		t.Helper()

		check(t, tc)
	})
}

func runChecks(t *testing.T, subject string, tcs ...cond) {
	t.Helper()

	t.Run(subject, func(t *testing.T) {
		t.Helper()

		checkAll(t, tcs)
	})
}

func checkAll(t *testing.T, commands []cond) {
	t.Helper()

	for _, cmd := range commands {
		check(t, cmd)
	}
}

func check(t *testing.T, tc cond) {
	t.Helper()

	// Do a second run of this specific test(TestKubeconfMain) with RUN_MAIN_FOR_TESTING=1 set,
	// So that the second run is able to run main() and this first run can verify the exit status returned by that.
	//
	// This technique originates from https://talks.golang.org/2014/testing.slide#23.
	cmd := exec.Command(os.Args[0], "-test.run=TestKubeconfMain")
	cmd.Env = append(
		cmd.Env,
		os.Environ()...,
	)
	cmd.Env = append(
		cmd.Env,
		tc.env...,
	)
	cmd.Env = append(
		cmd.Env,
		VarRunMainForTesting+"=1",
	)
	cmd.Env = append(
		cmd.Env,
		argsToEnv(tc.args)...,
	)

	stdin := strings.NewReader(tc.input)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()

	wantStdout := tc.wantStdout

	if tc.wantStdoutFile != "" {
		bs, err := ioutil.ReadFile(tc.wantStdoutFile)
		if err != nil {
			t.Fatalf("reading %s: %v", tc.wantStdoutFile, err)
		}

		wantStdout = string(bs)
	}

	got := stdout.String()
	if d := cmp.Diff(wantStdout, got); d != "" {
		t.Errorf("Unexpected stdout: want (-), got (+):\n%s", d)
	}

	wantStderr := tc.wantStderr
	if d := cmp.Diff(wantStderr, stderr.String()); d != "" {
		t.Errorf("Unexpected stderr: want (-), got (+):\n%s", d)
	}

	wantExitCode := tc.wantExitCode
	gotExitCode := 0
	exiterr, ok := err.(*exec.ExitError)

	if ok {
		gotExitCode = exiterr.ExitCode()
	}

	if gotExitCode != wantExitCode {
		t.Errorf("Unexpected exit code: want %d, got %d", wantExitCode, gotExitCode)
	}
}

func argsToEnv(args []string) []string {
	var env []string

	for i, arg := range args {
		env = append(env, fmt.Sprintf("%s%d=%s", VarRunMainForTestingArgPrefix, i, arg))
	}

	return env
}

func envToArgs(env []string) []string {
	var envvars []string

	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, VarRunMainForTestingArgPrefix) {
			envvars = append(envvars, kv)
		}
	}

	sort.Strings(envvars)

	var args []string

	for _, kv := range envvars {
		args = append(args, strings.Split(kv, "=")[1])
	}

	return args
}
