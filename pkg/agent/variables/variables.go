package variables

import (
	"golang.org/x/crypto/ssh"
)

const (
	DisableSSH            = "disable-ssh"
	AuthorizedSSHKeys     = "authorized-ssh-keys"
	HostSignerKey         = "host-signer-key"
	RegistryAuth          = "registry-auth"
	WhitelistedImages     = "whitelisted-images"
	DisableCustomCommands = "disable-custom-commands"
)

type Interface interface {
	GetDisableSSH() bool
	GetAuthorizedSSHKeys() []ssh.PublicKey
	GetHostSignerKey() string
	GetRegistryAuth() string
	GetWhitelistedImages() []string
	GetDisableCustomCommands() bool
}
