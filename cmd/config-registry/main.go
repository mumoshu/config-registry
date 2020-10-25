package main

import (
	"fmt"
	config_registry "github.com/mumoshu/config-registry/pkg/config-registry"
	"os"

	"github.com/fatih/color"
	"github.com/mumoshu/config-registry/internal/env"
	"github.com/mumoshu/config-registry/internal/printer"
)

func main() {
	cmd := config_registry.New()
	cmd.SetArgs(os.Args[1:])

	// Instead of letting cobra print the error, we do our own by using printer.Error below.
	cmd.SilenceErrors = true

	if err := cmd.Execute(); err != nil {
		printer.Error(color.Error, err.Error())

		if _, ok := os.LookupEnv(env.EnvDebug); ok {
			// print stack trace in verbose mode
			fmt.Fprintf(color.Error, "[DEBUG] error: %+v\n", err)
		}
		defer os.Exit(1)
	}
}
