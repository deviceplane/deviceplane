package image

import "strings"

func ToCanonical(image string) string {
	parts := strings.SplitN(image, "/", 3)

	switch len(parts) {
	case 1:
		parts = []string{
			"docker.io",
			"library",
			parts[0],
		}
	case 2:
		parts = []string{
			"docker.io",
			parts[0],
			parts[1],
		}
	}

	return strings.Join(parts, "/")

}
