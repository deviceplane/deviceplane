package email

type Interface interface {
	Send(request Request) error
}

type Request struct {
	FromName    string
	FromAddress string
	ToName      string
	ToAddress   string
	Subject     string
	Body        string
}
