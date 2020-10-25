package config_registry

import (
	"github.com/mumoshu/config-registry/internal/kubeconfig"
	"github.com/mumoshu/config-registry/internal/printer"
	"github.com/pkg/errors"
	"io"
	"os"
)

// InitOp indicates intention to init kubeconf
type InitOp struct {
}

func (op InitOp) Run(_, stderr io.Writer) error {
	kubeconfigPath, err := kubeconfig.KubeconfigPath()
	if err != nil {
		return err
	}

	confName := "default"

	if err := importConf(kubeconfigPath, confName); err != nil {
		return err
	}

	printer.Success(stderr, "Config %s created.",
		printer.SuccessColor.Sprint(confName))

	currConfFile, err := kubeconfCurrentConfFile()
	if err != nil {
		return errors.Wrap(err, "failed to determine state file")
	}

	cur, err := readConfName(currConfFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {

			err = writeConfName(currConfFile, confName)

			return errors.Wrap(err, "failed to save current config name")
		}

		return errors.Wrapf(err, "reading current config name from %s", currConfFile)
	} else if cur == "" {
		err = writeConfName(currConfFile, confName)

		return errors.Wrap(err, "failed to save current config name")
	}

	return nil
}

func importConf(kubeconfigPath, confName string) error {
	defaultPath, err := confFilePath(confName)
	if err != nil {
		return err
	}

	regPath, err := confRegistryPath()
	if err != nil {
		return errors.Wrap(err, "determining config registry path")
	}

	if err := os.MkdirAll(regPath, 0755); err != nil {
		return errors.Wrapf(err, "initializing %s", regPath)
	}

	if err := HardLink(kubeconfigPath, defaultPath); err != nil {
		return errors.Wrapf(err, "copying %s to %s", kubeconfigPath, defaultPath)
	}

	return nil
}
