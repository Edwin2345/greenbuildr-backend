package entities

import "time"

type TokenType string

const (
	EmailVerification TokenType = "email_verification"
	PasswordReset     TokenType = "password_reset"
)

type Token struct {
	ID        string
	UserID    string
	TokenHash string
	TokenType TokenType
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
