package user

import "context"

type MailSender interface {
	Send(ctx context.Context, email, message string) error
}
