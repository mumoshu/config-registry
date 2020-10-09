package kubeconf

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

// CurrentOp prints the current context
type CurrentOp struct{}

func (_op CurrentOp) Run(stdout, _ io.Writer) error {
	currConfFile, err := kubeconfCurrentConfFile()
	if err != nil {
		return errors.Wrap(err, "failed to determine state file")
	}

	curr, err := readConfName(currConfFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("current config is not set. please run `init`")
		}
		return errors.Wrapf(err, "reading current config name from %s", currConfFile)
	}

	if curr == "" {
		return errors.New("current config is not set")
	}

	_, err = fmt.Fprintln(stdout, curr)

	return errors.Wrap(err, "write error")
}
