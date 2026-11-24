package smtpsandboxclient

import (
	"context"
	"fmt"
	"time"

	"github.com/wneessen/go-mail"
)

const (
	resetPasswordTokenSubject = "Actualite. Reset Password Link"
	confirmEmailSubject       = "Actualite. Email Confirmation Link"
)

type SMTPClient struct {
	cfg    config
	client *mail.Client
}

func New(cfg config) (*SMTPClient, error) {
	client, err := mail.NewClient(cfg.Host, mail.WithPort(cfg.Port), mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(cfg.Username), mail.WithPassword(cfg.Password), mail.WithTimeout(cfg.Timeout))
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	return &SMTPClient{
		cfg:    cfg,
		client: client,
	}, nil
}

func (c *SMTPClient) SendLinkToResetPassword(ctx context.Context, to, link string) error {
	return c.send(ctx, to, resetPasswordTokenSubject, link)
}

func (c *SMTPClient) GetSendEmailTimeout() time.Duration {
	return c.cfg.Timeout
}

func (c *SMTPClient) SendLinkToConfirmEmail(ctx context.Context, to, link string) error {
	return c.send(ctx, to, confirmEmailSubject, link)
}

func (c *SMTPClient) send(ctx context.Context, to string, subject string, body string) error {
	msg := mail.NewMsg()

	if err := msg.From(c.cfg.From); err != nil {
		return err
	}

	if err := msg.To(to); err != nil {
		return err
	}

	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextPlain, body)

	return c.client.DialAndSendWithContext(ctx, msg)
}
