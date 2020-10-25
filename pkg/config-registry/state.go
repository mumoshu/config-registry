package config_registry

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/mumoshu/config-registry/internal/cmdutil"
)

func kubeconfPrevConfFile() (string, error) {
	return stateFilePath("prev")
}

func kubeconfCurrentConfFile() (string, error) {
	return stateFilePath("curr")
}

func confFilePath(name string) (string, error) {
	reg, err := confRegistryPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(reg, name), nil
}

func stateFilePath(id string) (string, error) {
	home := cmdutil.HomeDir()
	if home == "" {
		return "", errors.New("HOME or USERPROFILE environment variable not set")
	}
	return filepath.Join(home, ".kube", "kubeconf", id), nil
}

func confRegistryPath() (string, error) {
	home := cmdutil.HomeDir()
	if home == "" {
		return "", errors.New("HOME or USERPROFILE environment variable not set")
	}
	return filepath.Join(home, ".kube", "kubeconf", "registry"), nil
}

func listStateFiles() ([]string, error) {
	reg, err := confRegistryPath()
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(reg)
	if err != nil {
		return nil, err
	}

	var list []string

	for _, f := range files {
		list = append(list, f.Name())
	}

	return list, nil
}

// readConfName returns the saved config name
// if the state file exists, otherwise returns "".
func readConfName(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	//if os.IsNotExist(err) {
	//	return "", nil
	//}
	return string(b), errors.Wrapf(err, "reading config name stored in %s", path)
}

// writeConfName saves the specified value to the state file.
// It creates missing parent directories.
func writeConfName(path, value string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Wrap(err, "failed to create parent directories")
	}
	return ioutil.WriteFile(path, []byte(value), 0644)
}
