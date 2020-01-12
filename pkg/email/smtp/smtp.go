package smtp

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/deviceplane/deviceplane/pkg/email"
)

type Email struct {
	Server   string
	Port     int
	Username string
	Password string
}

func NewEmail(server string, port int, username, password string) *Email {
	return &Email{
		Server:   server,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (e *Email) Send(request email.Request) error {
	from := mail.Address{
		Name:    request.FromName,
		Address: request.FromAddress,
	}
	to := mail.Address{
		Name:    request.ToName,
		Address: request.ToAddress,
	}

	conn, err := tls.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", e.Server, e.Port),
		&tls.Config{
			ServerName: e.Server,
		})
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, e.Server)
	if err != nil {
		return err
	}

	if err = c.Auth(smtp.PlainAuth(
		"", e.Username, e.Password, e.Server,
	)); err != nil {
		return err
	}

	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	if _, err = w.Write([]byte(
		strings.Join([]string{
			strings.Join([]string{
				fmt.Sprintf("From: %s", from.String()),
				fmt.Sprintf("To: %s", to.String()),
				fmt.Sprintf("Subject: %s", request.Subject),
			}, "\r\n"),
			request.Body,
		}, "\r\n\r\n"),
	)); err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	return c.Quit()
}
