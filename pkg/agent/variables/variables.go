package variables

import (
	"golang.org/x/crypto/ssh"
)

const (
	DisableSSH        = "disable-ssh"
	AuthorizedSSHKeys = "authorized-ssh-keys"
)

type Interface interface {
	GetDisableSSH() bool
	GetAuthorizedSSHKeys() []ssh.PublicKey
}
