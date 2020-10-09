package main

import (
	"fmt"
	"github.com/mumoshu/kubeconf/pkg/kubeconf"
	"os"

	"github.com/fatih/color"
	"github.com/mumoshu/kubeconf/internal/env"
	"github.com/mumoshu/kubeconf/internal/printer"
)

func main() {
	cmd := kubeconf.New()
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
