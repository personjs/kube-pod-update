package main

import (
	"fmt"
	"os"

	"github.com/personjs/kube-pod-update/internal/compare"
	"github.com/personjs/kube-pod-update/internal/k8s"
	"github.com/personjs/kube-pod-update/internal/notifier"
	"github.com/personjs/kube-pod-update/internal/registry"
	"github.com/spf13/cobra"
)

var update bool

var rootCmd = &cobra.Command{
	Use:   "kube-checker",
	Short: "Detect and optionally update outdated images in running pods",
	Run: func(cmd *cobra.Command, args []string) {
		client := k8s.NewClient()
		pods := k8s.ListPods(client)

		for _, pod := range pods {
			for _, container := range pod.Spec.Containers {
				remoteDigest, err := registry.GetRemoteDigest(container.Image)
				if err != nil {
					notifier.LogRegistryError(pod, container.Image, err)
					continue
				}
				if compare.IsOutdated(container.Image, remoteDigest) {
					if update {
						err := k8s.UpdatePodImage(client, pod, container.Name, remoteDigest)
						if err != nil {
							notifier.LogUpdateError(pod, container.Image, remoteDigest, err)
						} else {
							notifier.LogUpdated(pod, container.Image, remoteDigest)
						}
					} else {
						notifier.LogOutdatedImage(pod, container.Image, remoteDigest)
					}
				}
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&update, "update", "u", false, "Update deployments with outdated images")
}
