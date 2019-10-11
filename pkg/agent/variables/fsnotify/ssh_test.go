package fsnotify

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

var (
	rawKeys = []string{
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC5TLZVo6MkNVzXEWQGOB4hhaKVSz18LmAWllDXfadorxobV47chg4YE/8K+rUtbG/gaXVaZ0ENnu1wbojofETzdnFfKIiWBLkQFDWVBG+xJSEQWr/udi0kiSZ6AS0vCu7iVwgCjbKilOeRrKniQneGVMeti0YXuJpBuNzOMNxQOIDI6a45l6+EaVmsFcesPRg2cvvuizc9ejsm5JJ3XfNhmP7ovvk9vzwyHUbbywvki8m1I5/lGL3n3NoTvfKrI7uX24a5MQR3/cNRMfev9Nlf7Ss/uyLfW0afFBF+XoDfqGbVmyyDV812+kpSIcj833/5C84G3vrxeAZMgMZLZNW1 josh@x",
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDhrCl7MabBByB75H1zLXlY9idvtPcx7MU3bNCVL6MnoOfp9Tqm/vroEJ6dXUJmkKvsWpJBQ4pY9+EX3xBZxPPosEaaCyFU9LkeX/IMcJIHLEbDlo/x490wjwaJgt63Ul+2qJfL7SCB8e3qaFDdJKTY25ScoLviWFvbYJG5mTf67kv3annVn6NXJy3/pYBIIJWkQIlBVMGEOMCX6oELaCm9y1MTyc6n6YWiQzLydkkmfOdupWW/lrw5/M/kfumV5P6msfH+zpZMY/WuwZNi+RgnQ45i/LvULIQcfe0C8OFVOcXQ0dKRnoYCf/6W1Y8+Rxzg1JJ2TOzZ5+ouSJhNssGd josh@x",
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCi16gbYQITUde13lrZMwTcNZWvt7jlG+LKz0HRLOl3EoWqfEt7i0N1OYTXe76zKj6IGpIw68tR6wW2PSxNGnwdhrCKP5kHJsAsWrqT4zeW7nJkXddddadtH2Vw9nTQBUs2hpkgsGbcWtg/cQpqBKEIAiRkN3g3W9FTnLhsbRkdJ0yEoWF5zDWLHmGFopngHVmempCLpqJzWaGYoCsiQv9fU2ThF9ivGuzkqvPl6z1Rz3/r5i6s0gl38XGDoMjam48R5DqXsnHphk1dCvkj9B0/kDOfK7J0wby3xMexHseCDv21xXkwi4hcvZ8Dowud/3aqGzW/LktlfqYYF8Cm01Jx josh@x",
	}
)

func TestParseAuthorizedKeysFile(t *testing.T) {
	var keys []ssh.PublicKey
	for _, rawKey := range rawKeys {
		key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(rawKey))
		if err != nil {
			t.Fatal(err)
		}
		keys = append(keys, key)
	}

	for i := 1; i < len(rawKeys); i++ {
		for j := 0; j < 10; j++ {
			// Add some trailing newlines to make sure we parse correctly still
			trailingCharacters := strings.Repeat("\n", j)

			parsedKeys, err := parseAuthorizedKeysFile([]byte(strings.Join(rawKeys[:i], "\n") + trailingCharacters))

			require.NoError(t, err)
			require.Len(t, parsedKeys, i)
			for k, parsedKey := range parsedKeys {
				require.Equal(t, keys[k], parsedKey)
			}
		}
	}
}
