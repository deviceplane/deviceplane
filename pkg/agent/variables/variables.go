package variables

import (
	"golang.org/x/crypto/ssh"
)

const (
	DisableSSH            = "disable-ssh"
	AuthorizedSSHKeys     = "authorized-ssh-keys"
	WhitelistedImages     = "whitelisted-images"
	DisableCustomCommands = "disable-custom-commands"
)

type Interface interface {
	GetDisableSSH() bool
	GetAuthorizedSSHKeys() []ssh.PublicKey
	GetWhitelistedImages() []string
	GetDisableCustomCommands() bool
}
