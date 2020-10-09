package kubeconf

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mumoshu/kubeconf/internal/testutil"
)

func Test_readConfName_nonExistingFile(t *testing.T) {
	s, err := readConfName(filepath.FromSlash("/non/existing/file"))
	if err == nil {
		t.Fatal("missing expected error")
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}

	if s != "" {
		t.Fatalf("expected empty string; got=%q", s)
	}
}

func Test_readConfName(t *testing.T) {
	path, cleanup := testutil.TempFile(t, "foo")
	defer cleanup()

	s, err := readConfName(path)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "foo"; s != expected {
		t.Fatalf("expected=%q; got=%q", expected, s)
	}
}

func Test_writeConfName_err(t *testing.T) {
	path := filepath.Join(os.DevNull, "foo", "bar")
	err := writeConfName(path, "foo")
	if err == nil {
		t.Fatal("got empty error")
	}
}

func Test_writeConfName(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "state-file-test")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "foo", "bar")

	if err := writeConfName(path, "conf1"); err != nil {
		t.Fatal(err)
	}

	v, err := readConfName(path)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "conf1"; v != expected {
		t.Fatalf("read wrong value=%q; expected=%q", v, expected)
	}
}

func Test_kubeconfPrevConfFile(t *testing.T) {
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", filepath.FromSlash("/foo/bar"))
	defer os.Setenv("HOME", origHome)

	expected := filepath.Join(filepath.FromSlash("/foo/bar"), ".kube", "kubeconf", "prev")
	v, err := kubeconfPrevConfFile()
	if err != nil {
		t.Fatal(err)
	}
	if v != expected {
		t.Fatalf("expected=%q got=%q", expected, v)
	}
}

func Test_kubeconfPrevConfFile_error(t *testing.T) {
	origHome := os.Getenv("HOME")
	origUserprofile := os.Getenv("USERPROFILE")
	os.Unsetenv("HOME")
	os.Unsetenv("USERPROFILE")
	defer os.Setenv("HOME", origHome)
	defer os.Setenv("USERPROFILE", origUserprofile)

	_, err := kubeconfPrevConfFile()
	if err == nil {
		t.Fatal(err)
	}
}

func Test_kubeconfCurrentConfFile(t *testing.T) {
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", filepath.FromSlash("/foo/bar"))
	defer os.Setenv("HOME", origHome)

	expected := filepath.Join(filepath.FromSlash("/foo/bar"), ".kube", "kubeconf", "curr")
	v, err := kubeconfCurrentConfFile()
	if err != nil {
		t.Fatal(err)
	}
	if v != expected {
		t.Fatalf("expected=%q got=%q", expected, v)
	}
}

func Test_kubeconfCurrentConfFile_error(t *testing.T) {
	origHome := os.Getenv("HOME")
	origUserprofile := os.Getenv("USERPROFILE")
	os.Unsetenv("HOME")
	os.Unsetenv("USERPROFILE")
	defer os.Setenv("HOME", origHome)
	defer os.Setenv("USERPROFILE", origUserprofile)

	_, err := kubeconfCurrentConfFile()
	if err == nil {
		t.Fatal(err)
	}
}
