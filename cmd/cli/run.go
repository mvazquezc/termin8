package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/mvazquezc/termin8/pkg/run"
	"github.com/mvazquezc/termin8/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	kubeconfigFile   string
	namespaces       []string
	skipAPIResources []string
	extendedOutput   string
	dryRun           bool
)

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "run",
		Short:        "Terminates stuck namespaced resources in the specified namespaces",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate command Args
			err := validateRunCommandArgs()
			if err != nil {
				return err
			}
			runResults, err := run.RunCommandRun(kubeconfigFile, namespaces, skipAPIResources, dryRun)
			if err != nil {
				return err
			}
			switch {
			case extendedOutput == "yaml":
				if len(runResults) > 0 {
					fmt.Println()
					utils.WriteYamlOutput(runResults)
				}
			case extendedOutput == "json":
				if len(runResults) > 0 {
					fmt.Println()
					utils.WriteJsonOutput(runResults)
				}
			}
			return err
		},
	}
	addRunCommandFlags(cmd)
	return cmd
}

func addRunCommandFlags(cmd *cobra.Command) {

	flags := cmd.Flags()
	flags.StringVarP(&kubeconfigFile, "kubeconfig", "k", "", "Path to the kubeconfig file to be used. If not set, will default to KUBECONFIG env var")
	flags.StringSliceVarP(&namespaces, "namespaces", "n", []string{""}, "List of namespaces where stuck objects will be terminated (comma separated) e.g: ns1,ns2")
	flags.StringSliceVarP(&skipAPIResources, "skip-api-resources", "s", nil, "List of namespaced api resources to skip (comma separated) e.g: myresource.group.example.com,myresource2.group2.example.com")
	flags.StringVarP(&extendedOutput, "extended-output", "o", "", "Extended output in an specific format. Usage: '-o [  yaml | json ]'")
	flags.BoolVarP(&dryRun, "dry-run", "d", false, "Will not terminate stuck resources, will output what would have been terminated")
	cmd.MarkFlagRequired("namespaces")
}

// validateCommandArgs validates that arguments passed by the user are valid
func validateRunCommandArgs() error {

	if kubeconfigFile != "" {
		if _, err := os.Stat(kubeconfigFile); err != nil {
			return errors.New("Kubeconfig file " + kubeconfigFile + " does not exist.")
		}
	} else {
		if _, err := os.Stat(os.Getenv("KUBECONFIG")); err != nil {
			return errors.New("Kubeconfig file " + os.Getenv("KUBECONFIG") + " does not exist.")
		}
	}
	for _, namespace := range namespaces {
		if namespace == "" {
			return errors.New("Namespaces list contains spaces")
		}
	}
	if len(skipAPIResources) > 0 {
		for _, apiResource := range skipAPIResources {
			if apiResource == "" {
				return errors.New("skip-api-resources list contains spaces")
			}
		}
	}
	if extendedOutput != "" && extendedOutput != "yaml" && extendedOutput != "json" {
		return errors.New("Unsupported extended output format " + extendedOutput)
	}
	if dryRun {
		extendedOutput = "yaml"
	}

	return nil
}
