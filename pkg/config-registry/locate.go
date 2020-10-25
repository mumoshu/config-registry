package config_registry

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
)

// LocateOp indicates intention to copy configs.
type LocateOp struct {
	Name string
}

func (op LocateOp) Run(stdout, _ io.Writer) error {
	confPath, err := confFilePath(op.Name)
	if err != nil {
		return errors.Wrapf(err, "determining config file path for %q", op.Name)
	}

	_, err = readConfName(confPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("config %s does not exist", op.Name)
		}
		return errors.Wrapf(err, "reading config file %s", confPath)
	}

	stdout.Write([]byte(confPath))

	return nil
}
