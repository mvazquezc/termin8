package cli

import (
	"fmt"
	"github.com/mvazquezc/termin8/pkg/version"
	"github.com/spf13/cobra"
)

var (
	short bool
)

func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !short {
				fmt.Printf("Cli version: %s\n", version.PrintVersion())
				fmt.Printf("Build time: %s\n", version.GetBuildTime())
				fmt.Printf("Git commit: %s\n", version.GetGitCommit())
				fmt.Printf("Go version: %s\n", version.GetGoVersion())
				fmt.Printf("Go compiler: %s\n", version.GetGoCompiler())
				fmt.Printf("Go Platform: %s\n", version.GetGoPlatform())
			} else {
				fmt.Printf("%s\n", version.PrintVersion())
			}
			return nil
		},
	}
	addVersionFlags(cmd)
	return cmd
}

func addVersionFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.BoolVar(&short, "short", false, "show only the version number")
}
