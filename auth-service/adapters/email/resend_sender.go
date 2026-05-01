package email

import (
	"context"
	"fmt"

	resend "github.com/resend/resend-go/v2"
)

type ResendSender struct {
	client *resend.Client
	from   string
}

func NewResendSender(apiKey, from string) *ResendSender {
	return &ResendSender{
		client: resend.NewClient(apiKey),
		from:   from,
	}
}

func (s *ResendSender) SendVerificationEmail(ctx context.Context, to, verifyURL string) error {
	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: "Verify your Greenbuildr account",
		Html:    fmt.Sprintf(`Thanks for signing up for GreenBuildr!<br><br><a href="%s">Click here to verify your email</a><br><br>This link expires in 1 hour.`, verifyURL),
	})
	return err
}

func (s *ResendSender) SendPasswordResetEmail(ctx context.Context, to, resetURL string) error {
	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: "Reset your Greenbuildr password",
		Html:    fmt.Sprintf(`<p>We received a request to reset your password. <a href="%s">Click here to set a new password</a>. This link expires in 15 minutes.</p><p>If you didn't request this, ignore this email.</p>`, resetURL),
	})
	return err
}
