package notifier

import (
	"fmt"

	"github.com/fatih/color"
	v1 "k8s.io/api/core/v1"
)

var (
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	// cyan   = color.New(color.FgCyan).SprintFunc()
)

func LogOutdatedImage(pod v1.Pod, currentImage, latestDigest string) {
	fmt.Printf("%s %s/%s: %s → %s\n",
		yellow("[OUTDATED]"),
		pod.Namespace, pod.Name,
		currentImage,
		latestDigest,
	)
}

func LogUpdated(pod v1.Pod, currentImage, newDigest string) {
	fmt.Printf("%s %s/%s: %s → %s\n",
		green("[UPDATED]"),
		pod.Namespace, pod.Name,
		currentImage,
		newDigest,
	)
}

func LogUpdateError(pod v1.Pod, currentImage, newDigest string, err error) {
	fmt.Printf("%s %s/%s: failed to update %s → %s: %v\n",
		red("[ERROR]"),
		pod.Namespace, pod.Name,
		currentImage,
		newDigest,
		err,
	)
}

func LogRegistryError(pod v1.Pod, image string, err error) {
	fmt.Printf("%s %s/%s: failed to check image %s: %v\n",
		red("[ERROR]"),
		pod.Namespace, pod.Name,
		image,
		err,
	)
}
