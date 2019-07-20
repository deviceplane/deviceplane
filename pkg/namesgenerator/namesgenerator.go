package namesgenerator

import (
	"strings"

	"github.com/docker/docker/pkg/namesgenerator"
)

func GetRandomName() string {
	return strings.ReplaceAll(namesgenerator.GetRandomName(0), "_", "-")
}
