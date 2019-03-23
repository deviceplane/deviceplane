package sendgrid

import (
	"github.com/deviceplane/deviceplane/pkg/email"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Email struct {
	client *sendgrid.Client
}

func NewEmail(client *sendgrid.Client) *Email {
	return &Email{
		client: client,
	}
}

func (e *Email) Send(request email.Request) error {
	from := mail.NewEmail(request.FromName, request.FromAddress)
	to := mail.NewEmail(request.ToName, request.ToAddress)
	message := mail.NewSingleEmail(from, request.Subject, to, request.PlainTextContent, request.HTMLContent)
	_, err := e.client.Send(message)
	return err
}
