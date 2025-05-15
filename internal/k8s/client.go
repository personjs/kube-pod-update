package k8s

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClient() *kubernetes.Clientset {
	var config *rest.Config
	var err error

	// Use in-cluster config if available
	if _, ok := os.LookupEnv("KUBERNETES_SERVICE_HOST"); ok {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Error loading in-cluster config: %v", err)
		}
	} else {
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("Error loading kubeconfig: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}
	return clientset
}

// ListPods retrieves all pods across all namespaces
func ListPods(client *kubernetes.Clientset) []v1.Pod {
	podList, err := client.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing pods: %v", err)
	}
	return podList.Items
}

func UpdatePodImage(client *kubernetes.Clientset, pod v1.Pod, containerName, newDigest string) error {
	deploymentsClient := client.AppsV1().Deployments(pod.Namespace)

	// Assumes the pod belongs to a deployment with same name prefix
	// You could improve this with an actual ownerReference check
	deployName := getDeploymentNameFromPod(pod)
	if deployName == "" {
		return fmt.Errorf("unable to determine deployment for pod %s", pod.Name)
	}

	deploy, err := deploymentsClient.Get(context.Background(), deployName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get deployment failed: %w", err)
	}

	// Update container image in the deployment spec
	updated := false
	for i, c := range deploy.Spec.Template.Spec.Containers {
		if c.Name == containerName {
			refParts := strings.Split(c.Image, "@")
			if len(refParts) < 1 {
				return fmt.Errorf("cannot parse current image: %s", c.Image)
			}
			deploy.Spec.Template.Spec.Containers[i].Image = fmt.Sprintf("%s@%s", refParts[0], newDigest)
			updated = true
		}
	}

	if !updated {
		return fmt.Errorf("container %s not found in deployment", containerName)
	}

	_, err = deploymentsClient.Update(context.Background(), deploy, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

func getDeploymentNameFromPod(pod v1.Pod) string {
	for _, owner := range pod.OwnerReferences {
		if owner.Kind == "ReplicaSet" {
			// Remove hash suffix: e.g. nginx-6dfb498665 -> nginx
			return strings.Join(strings.Split(owner.Name, "-")[:len(strings.Split(owner.Name, "-"))-1], "-")
		}
	}
	return ""
}
