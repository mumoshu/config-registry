package kubeconf

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/mumoshu/kubeconf/internal/kubeconfig"
	"github.com/mumoshu/kubeconf/internal/printer"
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

		if err := CopyFile(nextConfPath, writeConfPath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return "", fmt.Errorf("config %s does not exist", name)
			}

			return "", errors.Wrapf(err, "copying %s to %s", nextConfPath, writeConfPath)
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
