package smtpmailhogclient

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

func NewMust() *SMTPClient {
	cfg, err := newConfig()
	if err != nil {
		panic(err)
	}

	client, err := new(cfg)
	if err != nil {
		panic(err)
	}

	return client
}

func new(cfg config) (*SMTPClient, error) {
	client, err := mail.NewClient(cfg.Host,
		mail.WithPort(cfg.Port),
		mail.WithTLSPolicy(mail.NoTLS),
		mail.WithTimeout(cfg.Timeout),
	)
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

func (c *SMTPClient) SendLinkToConfirmEmail(ctx context.Context, to, link string) error {
	return c.send(ctx, to, confirmEmailSubject, link)
}

func (c *SMTPClient) GetSendEmailTimeout() time.Duration {
	return c.cfg.Timeout
}

func (c *SMTPClient) send(ctx context.Context, to string, subject string, body string) error {
	msg := mail.NewMsg()

	if err := msg.From(c.cfg.From); err != nil {
		return fmt.Errorf("set from: %w", err)
	}

	if err := msg.To(to); err != nil {
		return fmt.Errorf("set to: %w", err)
	}

	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextPlain, body)

	if err := c.client.DialAndSendWithContext(ctx, msg); err != nil {
		return fmt.Errorf("send mail to %q: %w", to, err)
	}
	return nil
}
