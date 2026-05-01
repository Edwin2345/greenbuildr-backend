package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"

	resend "github.com/resend/resend-go/v2"
)

type ResendSender struct {
	client       *resend.Client
	from         string
	templatesDir string
}

func NewResendSender(apiKey, from, templatesDir string) *ResendSender {
	return &ResendSender{
		client:       resend.NewClient(apiKey),
		from:         from,
		templatesDir: templatesDir,
	}
}

func (s *ResendSender) renderTemplate(name string, data any) (string, error) {
	path := filepath.Join(s.templatesDir, name)
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", fmt.Errorf("parse template %s: %w", name, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template %s: %w", name, err)
	}
	return buf.String(), nil
}

func (s *ResendSender) SendVerificationEmail(ctx context.Context, to, verifyURL string) error {
	html, err := s.renderTemplate("verify.html", map[string]string{"VerifyURL": verifyURL})
	if err != nil {
		return err
	}
	_, err = s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: "Verify your Greenbuildr account",
		Html:    html,
	})
	return err
}

func (s *ResendSender) SendPasswordResetEmail(ctx context.Context, to, resetURL string) error {
	html, err := s.renderTemplate("password_reset.html", map[string]string{"ResetURL": resetURL})
	if err != nil {
		return err
	}
	_, err = s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: "Reset your Greenbuildr password",
		Html:    html,
	})
	return err
}
