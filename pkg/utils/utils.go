package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chelnak/ysmrr"
	"gopkg.in/yaml.v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
)

func Contains(s []string, e string, matchCase bool) bool {
	for _, a := range s {
		if matchCase {
			if a == e {
				return true
			}
		} else {
			if strings.ToLower(a) == strings.ToLower(e) {
				return true
			}
		}

	}
	return false
}

func GetNamespacedStuckResources(namespace string, skipAPIResources []string, clientSet *kubernetes.Clientset, client *dynamic.DynamicClient, spinner *ysmrr.Spinner) ([]StuckResource, error) {
	var stuckResources []StuckResource
	currentProgress := 0
	// If namespace doesn't exist we skip it
	_, err := clientSet.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return stuckResources, nil
		} else {
			return stuckResources, err
		}

	}
	apiResources, err := clientSet.Discovery().ServerPreferredNamespacedResources()
	if err != nil {
		return stuckResources, err
	}
	for i, apiResourceList := range apiResources {
		currentProgress = (i * 100) / len(apiResources)
		spinner.UpdateMessagef("Getting stuck resources in namespace %s. Progress: %d%%", namespace, currentProgress)
		// apiResourceList has all the resources from a group, so we need to iterate through them
		groupVersion := strings.Split(apiResourceList.GroupVersion, "/")
		group := ""
		version := ""
		if len(groupVersion) <= 1 {
			version = groupVersion[0]
		} else {
			group = groupVersion[0]
			version = groupVersion[1]
		}

		for _, apiResource := range apiResourceList.APIResources {
			if Contains(skipAPIResources, apiResource.Name+"."+group, false) {
				continue
			}
			gvr := schema.GroupVersionResource{
				Group:    group,
				Version:  version,
				Resource: apiResource.Name,
			}
			// Get the resources for a given gvr in the namespace
			// We skip resources that don't support list, get or delete
			if Contains(apiResource.Verbs, "list", true) && Contains(apiResource.Verbs, "get", true) && Contains(apiResource.Verbs, "delete", true) {
				resourceList, err := client.Resource(gvr).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					return nil, err
				}
				for _, resource := range resourceList.Items {
					// Get resource, check if it has a deletionTimestamp, if it has, check if it has finalizers, if it has add it to the stuckresource list
					resourceData, err := client.Resource(gvr).Namespace(namespace).Get(context.TODO(), resource.GetName(), metav1.GetOptions{})
					if err != nil {
						return nil, err
					}
					if resourceData.GetDeletionTimestamp() != nil && resourceData.GetFinalizers() != nil {
						stuckResource := StuckResource{
							ResourceName:      resource.GetName(),
							ResourceType:      apiResource.Name,
							ResourceNamespace: resource.GetNamespace(),
							ResourceGroup:     group,
							ResourceVersion:   version,
						}
						stuckResources = append(stuckResources, stuckResource)
					}
				}
			}
		}
	}

	return stuckResources, err
}

func RemoveFinalizer(sr StuckResource, client *dynamic.DynamicClient) error {
	//Retrieve the resource
	gvr := schema.GroupVersionResource{
		Group:    sr.ResourceGroup,
		Version:  sr.ResourceVersion,
		Resource: sr.ResourceType,
	}
	resource, err := client.Resource(gvr).Namespace(sr.ResourceNamespace).Get(context.TODO(), sr.ResourceName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	resource.SetFinalizers(nil)
	_, err = client.Resource(gvr).Namespace(sr.ResourceNamespace).Update(context.TODO(), resource, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return err
}

func NewKubeClients(kubeconfig string) (*dynamic.DynamicClient, *discovery.DiscoveryClient, *kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if kubeconfig == "" {
		envKubeconfig := os.Getenv("KUBECONFIG")
		if len(envKubeconfig) > 0 {
			config, err = clientcmd.BuildConfigFromFlags("", envKubeconfig)
			if err != nil {
				return nil, nil, nil, err
			}
		} else {
			return nil, nil, nil, errors.New("No kubeconfig file was provided and KUBECONFIG env var is unset")
		}

	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}
	dClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}
	return client, dClient, clientSet, err
}

func WriteJsonOutput(rr []RunResult) {
	o, _ := json.MarshalIndent(rr, "", "    ")
	fmt.Println(string(o))
}

func WriteYamlOutput(rr []RunResult) {
	o, _ := yaml.Marshal(rr)
	fmt.Println(string(o))
}
