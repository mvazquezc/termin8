package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/mvazquezc/termin8/pkg/run"
	"github.com/mvazquezc/termin8/pkg/utils"
	"github.com/spf13/cobra"
)

// runOptions holds the flag values for a single invocation of the run command.
type runOptions struct {
	kubeconfigFile   string
	namespaces       []string
	skipAPIResources []string
	extendedOutput   string
	dryRun           bool
}

func NewRunCommand() *cobra.Command {
	opts := &runOptions{}
	cmd := &cobra.Command{
		Use:          "run",
		Short:        "Terminates stuck namespaced resources in the specified namespaces",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.validate(); err != nil {
				return err
			}
			runResults, err := run.Execute(opts.kubeconfigFile, opts.namespaces, opts.skipAPIResources, opts.dryRun)
			if err != nil {
				return err
			}
			if len(runResults.NonAvailableApiServices) > 0 {
				fmt.Println()
				fmt.Println("WARNING: There are some API Services in 'Not Available' state, some objects may not be deleted. You may want to fix them and run this tool again.")
				for _, nonAvailableApiService := range runResults.NonAvailableApiServices {
					fmt.Printf("  - %s\n", nonAvailableApiService)
				}
			}
			switch opts.extendedOutput {
			case "yaml":
				if len(runResults.Results) > 0 {
					fmt.Println()
					if err := utils.WriteYamlOutput(runResults.Results); err != nil {
						return err
					}
				}
			case "json":
				if len(runResults.Results) > 0 {
					fmt.Println()
					if err := utils.WriteJsonOutput(runResults.Results); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
	opts.addFlags(cmd)
	return cmd
}

func (o *runOptions) addFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringVarP(&o.kubeconfigFile, "kubeconfig", "k", "", "Path to the kubeconfig file to be used. If not set, will default to KUBECONFIG env var")
	flags.StringSliceVarP(&o.namespaces, "namespaces", "n", nil, "List of namespaces where stuck objects will be terminated (comma separated) e.g: ns1,ns2")
	flags.StringSliceVarP(&o.skipAPIResources, "skip-api-resources", "s", nil, "List of namespaced api resources to skip (comma separated) e.g: myresource.group.example.com,myresource2.group2.example.com")
	flags.StringVarP(&o.extendedOutput, "extended-output", "o", "", "Extended output in an specific format. Usage: '-o [  yaml | json ]'")
	flags.BoolVarP(&o.dryRun, "dry-run", "d", false, "Will not terminate stuck resources, will output what would have been terminated")
	cmd.MarkFlagRequired("namespaces")
}

// validate checks that the flag values provided by the user are valid.
func (o *runOptions) validate() error {
	if o.kubeconfigFile != "" {
		if _, err := os.Stat(o.kubeconfigFile); err != nil {
			return fmt.Errorf("kubeconfig file %s does not exist", o.kubeconfigFile)
		}
	} else {
		if _, err := os.Stat(os.Getenv("KUBECONFIG")); err != nil {
			return fmt.Errorf("kubeconfig file %s does not exist", os.Getenv("KUBECONFIG"))
		}
	}
	for _, namespace := range o.namespaces {
		if namespace == "" {
			return errors.New("namespaces list contains an empty namespace")
		}
	}
	for _, apiResource := range o.skipAPIResources {
		if apiResource == "" {
			return errors.New("skip-api-resources list contains an empty entry")
		}
	}
	if o.extendedOutput != "" && o.extendedOutput != "yaml" && o.extendedOutput != "json" {
		return fmt.Errorf("unsupported extended output format %s", o.extendedOutput)
	}
	// Default dry-run output to yaml, but don't override an explicit -o choice.
	if o.dryRun && o.extendedOutput == "" {
		o.extendedOutput = "yaml"
	}

	return nil
}
