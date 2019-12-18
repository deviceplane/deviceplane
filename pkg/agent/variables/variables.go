package variables

import (
	"golang.org/x/crypto/ssh"
)

const (
	DisableSSH            = "disable-ssh"
	AuthorizedSSHKeys     = "authorized-ssh-keys"
	RegistryAuth          = "registry-auth"
	WhitelistedImages     = "whitelisted-images"
	DisableCustomCommands = "disable-custom-commands"
)

type Interface interface {
	GetDisableSSH() bool
	GetAuthorizedSSHKeys() []ssh.PublicKey
	GetRegistryAuth() string
	GetWhitelistedImages() []string
	GetDisableCustomCommands() bool
}
