package kubeconf

import (
	"github.com/fatih/color"
	"github.com/mumoshu/kubeconf/internal/cmdutil"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// UnsupportedOp indicates an unsupported flag.
type UnsupportedOp struct{ Err error }

func (op UnsupportedOp) Run(_, _ io.Writer) error {
	return op.Err
}

func getDefaultOp() Op {
	if cmdutil.IsInteractiveMode(os.Stdout) {
		return InteractiveSwitchOp{SelfCmd: os.Args[0]}
	}
	return ListOp{}
}

// New looks at flags (excl. executable name, i.e. argv[0])
// and decides which operation should be taken.
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use: selfName(),
		RunE: func(cmd *cobra.Command, args []string) error {
			//return getDefaultOp().Run(color.Output, color.Error)
			return getDefaultOp().Run(color.Output, color.Error)
		},
		SilenceUsage: true,
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "init",
			Short: "initialize kubeconf",
			Long:  "Initialize kubeconf by importing the current kubeconfig as the 'default' config",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return InitOp{}.Run(color.Output, color.Error)
			},
		},
		&cobra.Command{
			Use:   "import PATH NAME",
			Short: "import existing kubeconfig at PATH as NAME",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return ImportOp{Path: args[0], Name: args[1]}.Run(color.Output, color.Error)
			},
		},
		&cobra.Command{
			Use:   "rm",
			Short: "delete config NAME",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return DeleteOp{Configs: args}.Run(color.Output, color.Error)
			},
		},
		&cobra.Command{
			Use:   "current",
			Short: "show the current config name",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return CurrentOp{}.Run(color.Output, color.Error)
			},
		},
		&cobra.Command{
			Use:   "locate",
			Short: "print the path to config NAME",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return LocateOp{Name: args[0]}.Run(color.Output, color.Error)
			},
		},
		&cobra.Command{
			Use:   "ls",
			Short: "list the configs",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return ListOp{}.Run(color.Output, color.Error)
			},
		},
		&cobra.Command{
			Use:     "mv OLD NEW",
			Short:   "rename config OLD to NEW",
			Example: "`mv . NEW` renames current config to NEW",
			Args:    cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return RenameOp{New: args[1], Old: args[0]}.Run(color.Output, color.Error)
			},
		},
		&cobra.Command{
			Use:   "cp OLD NEW",
			Short: "copy config OLD to NEW",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return CopyOp{New: args[1], Old: args[0]}.Run(color.Output, color.Error)
			},
		},
		&cobra.Command{
			Use:     "use NAME",
			Short:   "switch to config NAME",
			Example: "`use -` switches to the previous config",
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return SwitchOp{Target: args[0]}.Run(color.Output, color.Error)
			},
		},
	)

	return cmd
}
