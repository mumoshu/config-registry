package config_registry

import (
	"fmt"
	"io"

	"facette.io/natsort"
	"github.com/pkg/errors"

	"github.com/mumoshu/config-registry/internal/printer"
)

// ListOp describes listing contexts.
type ListOp struct{}

func (_ ListOp) Run(stdout, stderr io.Writer) error {
	files, err := listStateFiles()
	if err != nil {
		return err
	}

	natsort.Sort(files)

	currConfFile, err := kubeconfCurrentConfFile()
	if err != nil {
		return errors.Wrap(err, "failed to determine state file")
	}

	cur, err := readConfName(currConfFile)
	if err != nil {
		return errors.Wrapf(err, "reading current config name from %s", currConfFile)
	}

	for _, c := range files {
		s := c
		if c == cur {
			s = printer.ActiveItemColor.Sprint(c)
		}
		fmt.Fprintf(stdout, "%s\n", s)
	}
	return nil
}
