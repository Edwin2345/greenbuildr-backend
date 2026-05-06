package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/greenbuildr/auth-service/domain/entities"
	domainerrors "github.com/greenbuildr/auth-service/domain/errors"
)

func newTestRegisterService(u *mockUserRepo, tr *mockTokenRepo, h *mockHasher, e *mockEmailSender) *RegisterService {
	return NewRegisterService(u, tr, h, e, "http://localhost:3000")
}

func TestRegister_PasswordTooShort(t *testing.T) {
	userRepo := &mockUserRepo{}
	svc := newTestRegisterService(userRepo, &mockTokenRepo{}, &mockHasher{}, &mockEmailSender{})

	err := svc.Register(context.Background(), RegisterInput{
		Email:    "test@example.com",
		Password: "short",
	})

	assert.Equal(t, domainerrors.ErrPasswordTooShort, err)
	userRepo.AssertNotCalled(t, "CreateUser")
}

func TestRegister_DuplicateEmail(t *testing.T) {
	userRepo := &mockUserRepo{}
	hasher := &mockHasher{}
	tokenRepo := &mockTokenRepo{}
	emailSender := &mockEmailSender{}

	hasher.On("Hash", mock.Anything).Return("hashed_pw", nil)
	userRepo.On("CreateUser", mock.Anything, mock.Anything).Return(domainerrors.ErrEmailAlreadyExists)

	svc := newTestRegisterService(userRepo, tokenRepo, hasher, emailSender)

	err := svc.Register(context.Background(), RegisterInput{
		Email:    "existing@example.com",
		Password: "password123",
	})

	assert.Equal(t, domainerrors.ErrEmailAlreadyExists, err)
	tokenRepo.AssertNotCalled(t, "CreateToken")
	emailSender.AssertNotCalled(t, "SendVerificationEmail")
}

func TestVerifyEmail_ExpiredToken(t *testing.T) {
	tokenRepo := &mockTokenRepo{}
	userRepo := &mockUserRepo{}

	expiredToken := &entities.Token{
		ID:        "token-id",
		UserID:    "user-id",
		TokenHash: "somehash",
		TokenType: entities.EmailVerification,
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	tokenRepo.On("GetTokenByHash", mock.Anything, mock.Anything).Return(expiredToken, nil)

	svc := newTestRegisterService(userRepo, tokenRepo, &mockHasher{}, &mockEmailSender{})

	err := svc.VerifyEmail(context.Background(), "somerawtoken")

	assert.Equal(t, domainerrors.ErrInvalidToken, err)
	userRepo.AssertNotCalled(t, "VerifyUserEmail")
	tokenRepo.AssertNotCalled(t, "DeleteTokensByUserAndType")
}

func TestVerifyEmail_InvalidToken(t *testing.T) {
	tokenRepo := &mockTokenRepo{}
	userRepo := &mockUserRepo{}

	tokenRepo.On("GetTokenByHash", mock.Anything, mock.Anything).Return(nil, domainerrors.ErrInvalidToken)

	svc := newTestRegisterService(userRepo, tokenRepo, &mockHasher{}, &mockEmailSender{})
	err := svc.VerifyEmail(context.Background(), "somerawtoken")

	assert.Equal(t, domainerrors.ErrInvalidToken, err)
	userRepo.AssertNotCalled(t, "VerifyUserEmail")
	tokenRepo.AssertNotCalled(t, "DeleteTokensByUserAndType")
}

func TestRegister_Success(t *testing.T) {
	userRepo := &mockUserRepo{}
	tokenRepo := &mockTokenRepo{}
	hasher := &mockHasher{}
	emailSender := &mockEmailSender{}

	hasher.On("Hash", mock.Anything).Return("hashed_pw", nil)
	userRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)
	tokenRepo.On("CreateToken", mock.Anything, mock.Anything).Return(nil)
	emailSender.On("SendVerificationEmail", mock.Anything, "test@example.com", mock.Anything).Return(nil)

	svc := newTestRegisterService(userRepo, tokenRepo, hasher, emailSender)

	err := svc.Register(context.Background(), RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	userRepo.AssertCalled(t, "CreateUser", mock.Anything, mock.Anything)
	tokenRepo.AssertCalled(t, "CreateToken", mock.Anything, mock.Anything)
	emailSender.AssertCalled(t, "SendVerificationEmail", mock.Anything, "test@example.com", mock.Anything)
}

func TestResendVerificationEmail_AccountNotFound(t *testing.T) {
	userRepo := &mockUserRepo{}
	tokenRepo := &mockTokenRepo{}
	emailSender := &mockEmailSender{}

	userRepo.On("GetUserByEmail", mock.Anything, "ghost@example.com").Return(nil, domainerrors.ErrUserNotFound)

	svc := newTestRegisterService(userRepo, tokenRepo, &mockHasher{}, emailSender)
	err := svc.ResendVerificationEmail(context.Background(), "ghost@example.com")

	assert.NoError(t, err)
	emailSender.AssertNotCalled(t, "SendVerificationEmail")
}

func TestResendVerificationEmail_AlreadyVerified(t *testing.T) {
	userRepo := &mockUserRepo{}

	verifiedAt := time.Now().Add(-24 * time.Hour)
	verifiedUser := &entities.User{ID: "user-123", Email: "user@example.com", EmailVerifiedAt: &verifiedAt}
	userRepo.On("GetUserByEmail", mock.Anything, "user@example.com").Return(verifiedUser, nil)

	emailSender := &mockEmailSender{}
	svc := newTestRegisterService(userRepo, &mockTokenRepo{}, &mockHasher{}, emailSender)
	err := svc.ResendVerificationEmail(context.Background(), "user@example.com")

	assert.NoError(t, err)
	emailSender.AssertNotCalled(t, "SendVerificationEmail")
}

func TestResendVerificationEmail_Success(t *testing.T) {
	userRepo := &mockUserRepo{}
	tokenRepo := &mockTokenRepo{}
	emailSender := &mockEmailSender{}

	unverifiedUser := &entities.User{ID: "user-123", Email: "user@example.com", EmailVerifiedAt: nil}
	userRepo.On("GetUserByEmail", mock.Anything, "user@example.com").Return(unverifiedUser, nil)
	tokenRepo.On("DeleteTokensByUserAndType", mock.Anything, "user-123", entities.EmailVerification).Return(nil)
	tokenRepo.On("CreateToken", mock.Anything, mock.Anything).Return(nil)
	emailSender.On("SendVerificationEmail", mock.Anything, "user@example.com", mock.Anything).Return(nil)

	svc := newTestRegisterService(userRepo, tokenRepo, &mockHasher{}, emailSender)
	err := svc.ResendVerificationEmail(context.Background(), "user@example.com")

	assert.NoError(t, err)
	tokenRepo.AssertCalled(t, "DeleteTokensByUserAndType", mock.Anything, "user-123", entities.EmailVerification)
	emailSender.AssertCalled(t, "SendVerificationEmail", mock.Anything, "user@example.com", mock.Anything)
}

func TestVerifyEmail_Success(t *testing.T) {
	tokenRepo := &mockTokenRepo{}
	userRepo := &mockUserRepo{}

	validToken := &entities.Token{
		ID:        "token-id",
		UserID:    "user-123",
		TokenHash: "somehash",
		TokenType: entities.EmailVerification,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	tokenRepo.On("GetTokenByHash", mock.Anything, "rawtoken").Return(validToken, nil)
	userRepo.On("VerifyUserEmail", mock.Anything, "user-123").Return(nil)
	tokenRepo.On("DeleteTokensByUserAndType", mock.Anything, "user-123", entities.EmailVerification).Return(nil)

	svc := newTestRegisterService(userRepo, tokenRepo, &mockHasher{}, &mockEmailSender{})

	err := svc.VerifyEmail(context.Background(), "rawtoken")

	assert.NoError(t, err)
	userRepo.AssertCalled(t, "VerifyUserEmail", mock.Anything, "user-123")
	tokenRepo.AssertCalled(t, "DeleteTokensByUserAndType", mock.Anything, "user-123", entities.EmailVerification)
}
