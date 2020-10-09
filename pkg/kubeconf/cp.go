package kubeconf

import (
	"github.com/mumoshu/kubeconf/internal/printer"
	"github.com/pkg/errors"
	"io"
)

// CopyOp indicates intention to copy configs.
type CopyOp struct {
	New string
	Old string
}

func (op CopyOp) Run(_, stderr io.Writer) error {
	var cur string

	if op.Old == "." {
		currConfFile, err := kubeconfCurrentConfFile()
		if err != nil {
			return errors.Wrap(err, "failed to determine state file")
		}

		cur, err = readConfName(currConfFile)
		if err != nil {
			return errors.Wrapf(err, "reading current config name from %s", currConfFile)
		}

		if cur == "" {
			return errors.New("current config is not set")
		}

		op.Old = cur
	}

	newConfPath, err := confFilePath(op.New)
	if err != nil {
		return errors.Wrapf(err, "determining config file path for %q", op.New)
	}

	oldConfPath, err := confFilePath(op.Old)
	if err != nil {
		return errors.Wrapf(err, "determining config path for %s", op.Old)
	}

	if err := CopyFile(oldConfPath, newConfPath); err != nil {
		return errors.Wrapf(err, "copying %s to %s", oldConfPath, newConfPath)
	}

	printer.Success(stderr, "Copied config %s to %s.",
		printer.SuccessColor.Sprint(op.Old),
		printer.SuccessColor.Sprint(op.New))
	return nil
}
