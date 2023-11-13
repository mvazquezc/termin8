package run

import (
	"flag"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"time"

	"github.com/chelnak/ysmrr"
	"github.com/mvazquezc/termin8/pkg/utils"
	"github.com/mvazquezc/termin8/pkg/version"
	"k8s.io/klog/v2"
)

// Can terminate specific objects in a given namespace
// Must check that objects to be deleted have a terminationTimestamp

func RunCommandRun(kubeconfigFile string, namespaces []string, skipAPIResources []string, dryRun bool) (utils.RunResults, error) {
	// Disable klog
	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "false")
	flag.Parse()
	var runResults utils.RunResults
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
		stuckResources, nonAvailableApiServices, err := utils.GetNamespacedStuckResources(namespace, skipAPIResources, clientSet, client, termin8Spinner)
		if err != nil {
			// Resource may have been deleted while we were iterating, if that's the case we don't err
			if !apierrors.IsNotFound(err) {
				termin8Spinner.UpdateMessagef("Couldn't get stuck resources in namespace %s", namespace)
				termin8Spinner.Error()
				sm.Stop()
				return runResults, err
			}
		}

		if len(stuckResources) > 0 {
			termin8Spinner.UpdateMessagef("Terminating stuck resources in namespace %s", namespace)
			for _, resource := range stuckResources {
				termin8Spinner.UpdateMessagef("Terminating resource: %s/%s", resource.ResourceType, resource.ResourceName)
				if !dryRun {
					err = utils.RemoveFinalizer(resource, client)
					if err != nil {
						// Resource may have been deleted while we were iterating, if that's the case we don't err
						if !apierrors.IsNotFound(err) {
							termin8Spinner.UpdateMessagef("Couldn't terminate resource: %s/%s", resource.ResourceType, resource.ResourceName)
							termin8Spinner.Error()
							sm.Stop()
							return runResults, err
						}
					}
				}
				terminatedResources = append(terminatedResources, resource.ResourceType+"/"+resource.ResourceName)
			}
			runResults.Results = append(runResults.Results, utils.RunResult{
				Namespace:           namespace,
				TerminatedResources: terminatedResources,
			})
			termin8Spinner.UpdateMessagef("%d stuck resources in namespace %s have been terminated", len(stuckResources), namespace)
			time.Sleep(2 * time.Second)
		}
		terminatedResourcesCount += len(stuckResources)
		// Clean terminatedresources for each iteration
		terminatedResources = nil
		runResults.NonAvailableApiServices = nonAvailableApiServices
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
