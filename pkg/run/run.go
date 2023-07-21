package run

import (
	"flag"
	"fmt"
	"time"

	"github.com/chelnak/ysmrr"
	"github.com/mvazquezc/termin8/pkg/utils"
	"github.com/mvazquezc/termin8/pkg/version"
	"k8s.io/client-go/discovery"
	"k8s.io/klog/v2"
)

// Can terminate specific objects in a given namespace
// Must check that objects to be deleted have a terminationTimestamp

func RunCommandRun(kubeconfigFile string, namespaces []string, skipAPIResources []string, dryRun bool) ([]utils.RunResult, error) {
	// Disable klog
	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "false")
	flag.Parse()
	var runResults []utils.RunResult
	var terminatedResources []string
	terminatedResourcesCount := 0
	sm := ysmrr.NewSpinnerManager()
	termin8Spinner := sm.AddSpinner(version.GetBinaryName() + " started")
	sm.Start()
	// Get a kubeclient
	client, _, clientSet, err := utils.NewKubeClients(kubeconfigFile)
	if err != nil {
		termin8Spinner.UpdateMessage("Couldn't connect to the cluster")
		termin8Spinner.Error()
		sm.Stop()
		return runResults, err
	}
	for _, namespace := range namespaces {
		stuckResources, err := utils.GetNamespacedStuckResources(namespace, skipAPIResources, clientSet, client, termin8Spinner)
		if err != nil {
			termin8Spinner.UpdateMessagef("Couldn't get stuck resources in namespace %s", namespace)
			termin8Spinner.Error()
			sm.Stop()
			if discovery.IsGroupDiscoveryFailedError(err) {
				fmt.Printf("\tWARNING: There are some API Services in 'Not Available' state, so you need to clean them up before continue terminating the namespaces")
				fmt.Printf("\n\tTo do that you just need to execute 'kubectl get apiservice | grep ServiceNotFound', then you can backthem up and delete them:\n")
				fmt.Printf("\n\tkubectl get apiservice <API SERVICE> -o yaml > ./apiservice_<name>_bck.yaml")
				fmt.Printf("\n\tkubectl delete apiservice <API SERVICE>\n\n")
			}
			return runResults, err
		}

		if len(stuckResources) > 0 {
			termin8Spinner.UpdateMessagef("Terminating stuck resources in namespace %s", namespace)
			for _, resource := range stuckResources {
				termin8Spinner.UpdateMessagef("Terminating resource: %s/%s", resource.ResourceType, resource.ResourceName)
				if !dryRun {
					err = utils.RemoveFinalizer(resource, client)
					if err != nil {
						termin8Spinner.UpdateMessagef("Couldn't terminate resource: %s/%s", resource.ResourceType, resource.ResourceName)
						termin8Spinner.Error()
						sm.Stop()
						return runResults, err
					}
				}
				terminatedResources = append(terminatedResources, resource.ResourceType+"/"+resource.ResourceName)
			}
			runResults = append(runResults, utils.RunResult{
				Namespace:           namespace,
				TerminatedResources: terminatedResources,
			})
			termin8Spinner.UpdateMessagef("%d stuck resources in namespace %s have been terminated", len(stuckResources), namespace)
			time.Sleep(2 * time.Second)
		}
		terminatedResourcesCount += len(stuckResources)
		// Clean terminatedresources for each iteration
		terminatedResources = nil
	}
	if !dryRun {
		termin8Spinner.UpdateMessagef("%s completed. %d stuck resources terminated.", version.GetBinaryName(), terminatedResourcesCount)
	} else {
		termin8Spinner.UpdateMessagef("%s completed. %d stuck resources would have been terminated.", version.GetBinaryName(), terminatedResourcesCount)
	}

	termin8Spinner.Complete()
	sm.Stop()
	return runResults, nil
}
