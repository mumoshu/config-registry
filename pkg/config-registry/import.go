package config_registry

import (
	"github.com/mumoshu/config-registry/internal/printer"
	"io"
)

// ImportOp indicates intention to import existing kubeconf
type ImportOp struct {
	Path, Name string
}

func (op ImportOp) Run(_, stderr io.Writer) error {
	if err := importConf(op.Path, op.Name); err != nil {
		return err
	}

	printer.Success(stderr, "Config %s created.",
		printer.SuccessColor.Sprint(op.Name))

	return nil
}
