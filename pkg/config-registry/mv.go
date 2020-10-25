package config_registry

import (
	"github.com/pkg/errors"
	"io"
	"os"

	"github.com/mumoshu/config-registry/internal/printer"
)

// RenameOp indicates intention to rename contexts.
type RenameOp struct {
	New string // NAME of New context
	Old string // NAME of Old context (or '.' for current-context)
}

// rename changes the old (NAME or '.' for current-context)
// to the "new" value. If the old refers to the current-context,
// current-context preference is also updated.
func (op RenameOp) Run(_, stderr io.Writer) error {
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

	if err := os.Rename(oldConfPath, newConfPath); err != nil {
		return errors.Wrapf(err, "renaming %s to %s", oldConfPath, newConfPath)
	}

	if cur != "" {
		_, err = switchConfig(op.New)
		return errors.Wrapf(err, "switching to %s", op.New)
	}

	printer.Success(stderr, "Renamed config %s to %s.",
		printer.SuccessColor.Sprint(op.Old),
		printer.SuccessColor.Sprint(op.New))
	return nil
}
