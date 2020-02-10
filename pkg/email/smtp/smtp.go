package smtp

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/deviceplane/deviceplane/pkg/email"
	"github.com/pkg/errors"
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

	certPool, err := x509.SystemCertPool()
	if err != nil {
		return err
	}

	conn, err := tls.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", e.Server, e.Port),
		&tls.Config{
			ServerName: e.Server,
			RootCAs:    certPool,
		},
	)
	if err != nil {
		return errors.WithMessage(err, "creating TLS connection")
	}

	c, err := smtp.NewClient(conn, e.Server)
	if err != nil {
		errors.WithMessage(err, "creating a client")
	}

	if err = c.Auth(smtp.PlainAuth(
		"", e.Username, e.Password, e.Server,
	)); err != nil {
		errors.WithMessage(err, "setting auth")
	}

	if err = c.Mail(from.Address); err != nil {
		return errors.WithMessage(err, "setting from address")
	}

	if err = c.Rcpt(to.Address); err != nil {
		return errors.WithMessage(err, "setting to address")
	}

	w, err := c.Data()
	if err != nil {
		return errors.WithMessage(err, "DATA command")
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
		return errors.WithMessage(err, "writing data")
	}

	if err = w.Close(); err != nil {
		return errors.WithMessage(err, "closing data writer")
	}

	err = c.Quit()
	if err != nil {
		return errors.WithMessage(err, "closing connection")
	}
	return nil
}
