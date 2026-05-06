package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/greenbuildr/auth-service/domain/entities"
	domainerrors "github.com/greenbuildr/auth-service/domain/errors"
	"github.com/greenbuildr/auth-service/ports"
)

type RegisterService struct {
	userRepo    ports.UserRepository
	tokenRepo   ports.TokenRepository
	hasher      ports.PasswordHasher
	emailSender ports.EmailSender
	frontendURL string
}

func NewRegisterService(
	userRepo ports.UserRepository,
	tokenRepo ports.TokenRepository,
	hasher ports.PasswordHasher,
	emailSender ports.EmailSender,
	frontendURL string,
) *RegisterService {
	return &RegisterService{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		hasher:      hasher,
		emailSender: emailSender,
		frontendURL: frontendURL,
	}
}

type RegisterInput struct {
	Email    string
	Password string
}

// Initial Registration
// ---------------------------------------------------------------------------------------
func (s *RegisterService) Register(ctx context.Context, input RegisterInput) error {
	log.Printf("register: starting registration for email=%s", input.Email)

	user, err := s.createUser(ctx, input)
	if err != nil {
		return err
	}

	rawToken, err := s.createVerificationToken(ctx, user.ID)
	if err != nil {
		return err
	}

	return s.sendVerificationEmail(ctx, user.Email, rawToken)
}

func (s *RegisterService) createUser(ctx context.Context, input RegisterInput) (*entities.User, error) {
	if len(input.Password) < 8 {
		log.Printf("register: rejected as password too short for email: %v", input.Email)
		return nil, domainerrors.ErrPasswordTooShort
	}

	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		log.Printf("register: failed to hash password: %v", err)
		return nil, domainerrors.ErrInternalError
	}

	user := &entities.User{
		ID:           uuid.NewString(),
		Email:        input.Email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		log.Printf("register: failed to create user email=%s: %v", input.Email, err)
		return nil, err
	}

	log.Printf("register: user created id=%s email=%s", user.ID, user.Email)
	return user, nil
}

func hashForStorage(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

func (s *RegisterService) createVerificationToken(ctx context.Context, userID string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Printf("register: failed to generate random token: %v", err)
		return "", domainerrors.ErrInternalError
	}
	rawToken := hex.EncodeToString(b)

	token := &entities.Token{
		ID:        uuid.NewString(),
		UserID:    userID,
		TokenHash: hashForStorage(rawToken),
		TokenType: entities.EmailVerification,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.tokenRepo.CreateToken(ctx, token); err != nil {
		log.Printf("register: failed to store verification token userID=%s: %v", userID, err)
		return "", domainerrors.ErrInternalError
	}

	log.Printf("register: verification token created userID=%s", userID)
	return rawToken, nil
}

func (s *RegisterService) sendVerificationEmail(ctx context.Context, email, rawToken string) error {
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", s.frontendURL, rawToken)
	log.Printf("register: sending verification email to=%s url=%s", email, verifyURL)

	if err := s.emailSender.SendVerificationEmail(ctx, email, verifyURL); err != nil {
		log.Printf("register: failed to send verification email to=%s: %v", email, err)
		return domainerrors.ErrInternalError
	}

	log.Printf("register: verification email sent to=%s", email)
	return nil
}

func (s *RegisterService) ResendVerificationEmail(ctx context.Context, email string) error {
	log.Printf("register: resending verification email to=%s", email)

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domainerrors.ErrUserNotFound {
			return nil
		}
		log.Printf("resendVerification: failed to find account email=%s: %v", email, err)
		return domainerrors.ErrInternalError
	}

	if user.IsEmailVerified() {
		return nil
	}

	//Delete old emnail vrificaiton tokens and create new one
	if err := s.tokenRepo.DeleteTokensByUserAndType(ctx, user.ID, entities.EmailVerification); err != nil {
		log.Printf("resendVerification: failed to delete old tokens userID=%s: %v", user.ID, err)
		return domainerrors.ErrInternalError
	}
	rawToken, err := s.createVerificationToken(ctx, user.ID)
	if err != nil {
		return err
	}

	return s.sendVerificationEmail(ctx, user.Email, rawToken)
}

// Email Verification
// ---------------------------------------------------------------------------------------
func (s *RegisterService) VerifyEmail(ctx context.Context, rawToken string) error {
	log.Printf("verifyEmail: verifying token")

	token, err := s.tokenRepo.GetTokenByHash(ctx, rawToken)
	if err != nil {
		log.Printf("verifyEmail: token lookup failed: %v", err)
		return err
	}

	if token.IsExpired() {
		log.Printf("verifyEmail: token expired userID=%s", token.UserID)
		return domainerrors.ErrInvalidToken
	}

	if err := s.userRepo.VerifyUserEmail(ctx, token.UserID); err != nil {
		log.Printf("verifyEmail: failed to verify email userID=%s: %v", token.UserID, err)
		return err
	}

	log.Printf("verifyEmail: email verified userID=%s", token.UserID)

	if err := s.tokenRepo.DeleteTokensByUserAndType(ctx, token.UserID, entities.EmailVerification); err != nil {
		log.Printf("verifyEmail: failed to delete tokens userID=%s: %v", token.UserID, err)
		return err
	}

	return nil
}
