package notifier

import (
	"context"
	"fmt"
	"gopkg.in/mail.v2"
)

type Client interface {
	NotifyTo(context.Context, NotifyToOptions) error
}

type Options struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewClient(options Options) Client {
	return client{
		sender: options.Username,
		dialer: mail.NewDialer(
			options.Host,
			options.Port,
			options.Username,
			options.Password,
		),
	}
}

// Dialer is an abstraction for mail.Dialer making it mockeable
type Dialer interface {
	DialAndSend(m ...*mail.Message) error
}

// Client represents a mail sender, all the messages are sent from the same address.
type client struct {
	sender string
	dialer Dialer
}

type NotifyToOptions struct {
	To      string
	Subject string
	Body    string // Body must be HTML formatted
}

func (c client) NotifyTo(_ context.Context, options NotifyToOptions) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", c.sender)
	msg.SetHeader("To", options.To)
	msg.SetHeader("Subject", options.Subject)
	msg.SetBody("text/html", options.Body)

	if err := c.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("unexpected error sending mail due to: %w", err)
	}

	return nil
}
