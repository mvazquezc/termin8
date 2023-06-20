package main

import (
	"github.com/mvazquezc/termin8/cmd/cli"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	command := newRootCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}

}

// newRootCommand implements the root command of example-ci
func newRootCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "termin8",
		Short: "Terminates stuck resources in Kubernetes namespaces",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}

	c.AddCommand(cli.NewRunCommand())
	c.AddCommand(cli.NewVersionCommand())

	return c
}
