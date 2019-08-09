package variables

const (
	DisableSSH = "disable-ssh"
)

type Interface interface {
	GetDisableSSH() bool
}
