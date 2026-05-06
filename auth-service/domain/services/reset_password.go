package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/greenbuildr/auth-service/domain/entities"
	domainerrors "github.com/greenbuildr/auth-service/domain/errors"
	"github.com/greenbuildr/auth-service/ports"
)

type ResetPasswordService struct {
	userRepo    ports.UserRepository
	tokenRepo   ports.TokenRepository
	hasher      ports.PasswordHasher
	emailSender ports.EmailSender
	frontendURL string
}

func NewResetPasswordService(
	userRepo ports.UserRepository,
	tokenRepo ports.TokenRepository,
	hasher ports.PasswordHasher,
	emailSender ports.EmailSender,
	frontendURL string,
) *ResetPasswordService {
	return &ResetPasswordService{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		hasher:      hasher,
		emailSender: emailSender,
		frontendURL: frontendURL,
	}
}

// Initial Password Reset Request
// ---------------------------------------------------------------------------------------
func (s *ResetPasswordService) RequestPasswordReset(ctx context.Context, email string) error {
	log.Printf("reset password: request password reset for email=%s", email)

	foundUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		//don't show to client that account doesnt exist -> fail silently
		if err == domainerrors.ErrUserNotFound {
			return nil
		}
		log.Printf("reset password: failed to find account using email=%s: %v", email, err)
		return err
	}

	resetPasswordToken, err := s.createResetPasswordToken(ctx, foundUser.ID)
	if err != nil {
		return err
	}

	return s.sendResetPasswordEmail(ctx, foundUser.Email, resetPasswordToken)
}

func (s *ResetPasswordService) createResetPasswordToken(ctx context.Context, userID string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Printf("reset password request: failed to generate random token: %v", err)
		return "", domainerrors.ErrInternalError
	}
	rawToken := hex.EncodeToString(b)

	token := &entities.Token{
		ID:        uuid.NewString(),
		UserID:    userID,
		TokenHash: hashForStorage(rawToken),
		TokenType: entities.PasswordReset,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.tokenRepo.CreateToken(ctx, token); err != nil {
		log.Printf("reset password request: failed to store token for userID=%s: %v", userID, err)
		return "", domainerrors.ErrInternalError
	}

	log.Printf("reset password request: token created userID=%s", userID)
	return rawToken, nil
}

func (s *ResetPasswordService) sendResetPasswordEmail(ctx context.Context, email, rawToken string) error {
	resetPasswordURL := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, rawToken)
	log.Printf("reset password request: sending email to=%s url=%s", email, resetPasswordURL)

	if err := s.emailSender.SendPasswordResetEmail(ctx, email, resetPasswordURL); err != nil {
		log.Printf("reset password request: failed to send email to=%s: %v", email, err)
		return domainerrors.ErrInternalError
	}

	log.Printf("reset password request: email sent to=%s", email)
	return nil
}

// Reset Password
// ---------------------------------------------------------------------------------------
func (s *ResetPasswordService) ResetPassword(ctx context.Context, resetPasswordToken string, newPassword string) error {
	//validate password requirement
	if len(newPassword) < 8 {
		return domainerrors.ErrPasswordTooShort
	}

	log.Printf("reset password: verifying token")
	token, err := s.tokenRepo.GetTokenByHash(ctx, resetPasswordToken)
	if err != nil {
		log.Printf("reset password: token lookup failed: %v", err)
		return err
	}
	if token.IsExpired() {
		log.Printf("reset password: token expired for userID=%s", token.UserID)
		return domainerrors.ErrInvalidToken
	}

	passwordHash, err := s.hasher.Hash(newPassword)
	if err != nil {
		log.Printf("reset password: failed to hash password: %v", err)
		return domainerrors.ErrInternalError
	}

	if err = s.userRepo.UpdateUserPassword(ctx, token.UserID, passwordHash); err != nil {
		log.Printf("reset password: failed update password for userID %s: %v", token.UserID, err)
		return err
	}

	return s.tokenRepo.DeleteTokensByUserAndType(ctx, token.UserID, entities.PasswordReset)
}
