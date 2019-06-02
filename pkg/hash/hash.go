package hash

import (
	"crypto/sha256"
	"fmt"
)

func Hash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}

func ShortHash(s string) string {
	hash := Hash(s)
	return hash[:6]
}
