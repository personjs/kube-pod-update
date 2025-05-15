package compare

import "strings"

func IsOutdated(currentImage, remoteDigest string) bool {
	return strings.Contains(currentImage, "@sha256:") && !strings.Contains(currentImage, remoteDigest)
}
