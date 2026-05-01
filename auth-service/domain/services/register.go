package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
	//insert inital user into db (email_validated_at is null)
	user, err := s.createUser(ctx, input)
	if err != nil {
		return err
	}

	//generate verfication email
	rawToken, err := s.createVerificationToken(ctx, user.ID)
	if err != nil {
		return err
	}
	return s.sendVerificationEmail(ctx, user.Email, rawToken)
}

func (s *RegisterService) createUser(ctx context.Context, input RegisterInput) (*entities.User, error) {
	if len(input.Password) < 8 {
		return nil, domainerrors.ErrPasswordTooShort
	}

	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
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
		return nil, err
	}

	return user, nil
}

func hashForStorage(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

func (s *RegisterService) createVerificationToken(ctx context.Context, userID string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
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
		return "", domainerrors.ErrInternalError
	}

	return rawToken, nil
}

func (s *RegisterService) sendVerificationEmail(ctx context.Context, email, rawToken string) error {
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", s.frontendURL, rawToken)

	if err := s.emailSender.SendVerificationEmail(ctx, email, verifyURL); err != nil {
		return domainerrors.ErrInternalError
	}
	return nil
}

// Email Verification
// ---------------------------------------------------------------------------------------
func (s *RegisterService) VerifyEmail(ctx context.Context, rawToken string) error {
	//check that veriicaiton token exists and is valid
	token, err := s.tokenRepo.GetTokenByHash(ctx, rawToken)
	if err != nil {
		return err
	}
	if token.IsExpired() {
		return domainerrors.ErrInvalidToken
	}

	//verify user given id in token
	if err := s.userRepo.VerifyUserEmail(ctx, token.UserID); err != nil {
		return err
	}

	//delete consumed token
	return s.tokenRepo.DeleteToken(ctx, rawToken)
}
