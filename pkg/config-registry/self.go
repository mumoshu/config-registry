package config_registry

import (
	"os"
	"path/filepath"
	"strings"
)

// selfName guesses how the user invoked the program.
func selfName() string {
	me := filepath.Base(os.Args[0])
	pluginPrefix := "kubectl-"
	if strings.HasPrefix(me, pluginPrefix) {
		return "kubectl " + strings.TrimPrefix(me, pluginPrefix)
	}
	return "config-registry"
}
