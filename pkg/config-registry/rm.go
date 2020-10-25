package config_registry

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/mumoshu/config-registry/internal/printer"
)

// DeleteOp indicates intention to delete contexts.
type DeleteOp struct {
	Configs []string // NAME or '.' to indicate current-context.
}

// deleteContexts deletes context entries one by one.
func (op DeleteOp) Run(_, stderr io.Writer) error {
	for _, confName := range op.Configs {
		confPath, err := confFilePath(confName)
		if err != nil {
			return errors.Wrapf(err, "determining config file path for %q", confName)
		}

		if err := os.Remove(confPath); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("config %q does not exist", confName)
			}
		}

		printer.Success(stderr, `Deleted config %s.`, printer.SuccessColor.Sprint(confName))
	}
	return nil
}
