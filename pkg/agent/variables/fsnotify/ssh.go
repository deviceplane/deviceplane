package fsnotify

import (
	"strings"

	"golang.org/x/crypto/ssh"
)

func parseAuthorizedKeysFile(in []byte) ([]ssh.PublicKey, error) {
	var authorizedKeys []ssh.PublicKey
	rest := in

	for {
		var authorizedKey ssh.PublicKey
		var err error
		authorizedKey, _, _, rest, err = ssh.ParseAuthorizedKey(rest)
		if err != nil {
			return nil, err
		}

		authorizedKeys = append(authorizedKeys, authorizedKey)

		if len(strings.TrimSpace(string(rest))) == 0 {
			break
		}
	}

	return authorizedKeys, nil
}
