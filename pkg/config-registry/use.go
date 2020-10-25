package config_registry

import (
	"fmt"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	"github.com/mumoshu/config-registry/internal/kubeconfig"
	"github.com/mumoshu/config-registry/internal/printer"
)

// SwitchOp indicates intention to switch contexts.
type SwitchOp struct {
	Target string // '-' for back and forth, or NAME
}

func (op SwitchOp) Run(_, stderr io.Writer) error {
	var newCtx string
	var err error
	if op.Target == "-" {
		newCtx, err = swapConfig()
	} else {
		newCtx, err = switchConfig(op.Target)
	}
	if err != nil {
		return errors.Wrap(err, "failed to switch config")
	}
	err = printer.Success(stderr, "Switched to config %q.", newCtx)
	return errors.Wrap(err, "print error")
}

// switchConfig switches to specified config.
func switchConfig(name string) (string, error) {
	currConfFile, err := kubeconfCurrentConfFile()
	if err != nil {
		return "", errors.Wrap(err, "failed to determine state file")
	}

	curr, err := readConfName(currConfFile)
	if err != nil {
		return "", errors.Wrapf(err, "reading current config name from %s", currConfFile)
	}

	if curr != name {
		prevConfFile, err := kubeconfPrevConfFile()
		if err != nil {
			return "", errors.Wrap(err, "failed to determine previous config file path")
		}

		nextConfPath, err := confFilePath(name)
		if err != nil {
			return "", errors.Wrapf(err, "determining config file path for %q", name)
		}

		writeConfPath, err := kubeconfig.KubeconfigPath()
		if err != nil {
			return "", errors.Wrap(err, "determining config path")
		}

		if _, err := os.Stat(nextConfPath); os.IsNotExist(err) {
			return "", fmt.Errorf("config %s does not exist", name)
		}

		if err := hardLinkIfRegistered(nextConfPath, writeConfPath); err != nil {
			return "", errors.Wrapf(err, "applying %s to %s", nextConfPath, writeConfPath)
		}

		if err := writeConfName(prevConfFile, curr); err != nil {
			return "", errors.Wrap(err, "failed to save previous previous name")
		}

		if err := writeConfName(currConfFile, name); err != nil {
			return "", errors.Wrap(err, "failed to save current config name")
		}
	}

	return name, nil
}

// swapConfig switches to previous config.
func swapConfig() (string, error) {
	prevCtxFile, err := kubeconfPrevConfFile()
	if err != nil {
		return "", errors.Wrap(err, "failed to determine state file")
	}
	prev, err := readConfName(prevCtxFile)
	if err != nil {
		return "", errors.Wrap(err, "failed to read previous config file")
	}
	if prev == "" {
		return "", errors.New("no previous config found")
	}
	return switchConfig(prev)
}

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}

	defer func() {
		cerr := out.Close()

		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return
	}

	err = out.Sync()

	return
}

func HardLink(src, dst string) (err error) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return xerrors.Errorf("open %s: no such file or directory", src)
	}

	if err := os.Link(src, dst); err != nil {
		return xerrors.Errorf("hard linking %s to %s: %w", src, dst, err)
	}

	return
}

func hardLinkIfRegistered(src, dst string) (err error) {
	dstInfo, err := os.Stat(dst)
	if err != nil {
		return xerrors.Errorf("stat %s: %w", src)
	}

	registryPath, err := confRegistryPath()
	if err != nil {
		return xerrors.Errorf("determining config registry path: %w", err)
	}

	infos, err := ioutil.ReadDir(registryPath)
	if err != nil {
		return xerrors.Errorf("reading directory: %w", err)
	}

	var registered bool

	for _, info := range infos {
		if os.SameFile(dstInfo, info) {
			registered = true

			break
		}
	}

	if !registered {
		return fmt.Errorf("%s is not registered yet. This operation is blocked to prevent the overwriting the unregistered file. Run `import %s somename` first", dst, dst)
	}

	// It's safe to delete ~/.kube/config as it is already known to be registered to ~/.kube/kubeconf/registry

	if err := os.Remove(dst); err != nil {
		return xerrors.Errorf("removing %s: %w", dst, err)
	}

	if err := HardLink(src, dst); err != nil {
		return xerrors.Errorf("switching to %s: %w", dst, err)
	}

	return nil
}
