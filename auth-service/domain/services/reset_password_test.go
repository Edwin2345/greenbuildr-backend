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

func newTestResetPasswordService(u *mockUserRepo, tr *mockTokenRepo, h *mockHasher, e *mockEmailSender) *ResetPasswordService {
	return NewResetPasswordService(u, tr, h, e, "http://localhost:3000")
}

func TestRequestPasswordReset_AccountNotFound(t *testing.T) {
	userRepo := &mockUserRepo{}
	tokenRepo := &mockTokenRepo{}
	emailSender := &mockEmailSender{}

	userRepo.On("GetUserByEmail", mock.Anything, "ghost@example.com").Return(nil, domainerrors.ErrUserNotFound)

	svc := newTestResetPasswordService(userRepo, tokenRepo, &mockHasher{}, emailSender)
	err := svc.RequestPasswordReset(context.Background(), "ghost@example.com")

	//should not return any error to user
	assert.NoError(t, err)
	tokenRepo.AssertNotCalled(t, "CreateToken")
	emailSender.AssertNotCalled(t, "SendPasswordResetEmail")
}

func TestRequestPasswordReset_Success(t *testing.T) {
	userRepo := &mockUserRepo{}
	tokenRepo := &mockTokenRepo{}
	emailSender := &mockEmailSender{}

	foundUser := &entities.User{ID: "user-123", Email: "user@example.com"}
	userRepo.On("GetUserByEmail", mock.Anything, "user@example.com").Return(foundUser, nil)
	tokenRepo.On("CreateToken", mock.Anything, mock.Anything).Return(nil)
	emailSender.On("SendPasswordResetEmail", mock.Anything, "user@example.com", mock.Anything).Return(nil)

	svc := newTestResetPasswordService(userRepo, tokenRepo, &mockHasher{}, emailSender)
	err := svc.RequestPasswordReset(context.Background(), "user@example.com")

	assert.NoError(t, err)
	tokenRepo.AssertCalled(t, "CreateToken", mock.Anything, mock.Anything)
	emailSender.AssertCalled(t, "SendPasswordResetEmail", mock.Anything, "user@example.com", mock.Anything)
}

func TestResetPassword_PasswordTooShort(t *testing.T) {
	svc := newTestResetPasswordService(&mockUserRepo{}, &mockTokenRepo{}, &mockHasher{}, &mockEmailSender{})

	err := svc.ResetPassword(context.Background(), "sometoken", "short")

	assert.Equal(t, domainerrors.ErrPasswordTooShort, err)
}

func TestResetPassword_ExpiredToken(t *testing.T) {
	tokenRepo := &mockTokenRepo{}

	expiredToken := &entities.Token{
		ID:        "token-id",
		UserID:    "user-123",
		TokenType: entities.PasswordReset,
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	tokenRepo.On("GetTokenByHash", mock.Anything, mock.Anything).Return(expiredToken, nil)

	svc := newTestResetPasswordService(&mockUserRepo{}, tokenRepo, &mockHasher{}, &mockEmailSender{})
	err := svc.ResetPassword(context.Background(), "sometoken", "newpassword123")

	assert.Equal(t, domainerrors.ErrInvalidToken, err)
}

func TestResetPassword_Success(t *testing.T) {
	tokenRepo := &mockTokenRepo{}
	userRepo := &mockUserRepo{}
	hasher := &mockHasher{}

	validToken := &entities.Token{
		ID:        "token-id",
		UserID:    "user-123",
		TokenType: entities.PasswordReset,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	tokenRepo.On("GetTokenByHash", mock.Anything, "resettoken").Return(validToken, nil)
	hasher.On("Hash", "newpassword123").Return("hashed_pw", nil)
	userRepo.On("UpdateUserPassword", mock.Anything, "user-123", "hashed_pw").Return(nil)
	tokenRepo.On("DeleteTokensByUserAndType", mock.Anything, "user-123", entities.PasswordReset).Return(nil)

	svc := newTestResetPasswordService(userRepo, tokenRepo, hasher, &mockEmailSender{})
	err := svc.ResetPassword(context.Background(), "resettoken", "newpassword123")

	assert.NoError(t, err)
	userRepo.AssertCalled(t, "UpdateUserPassword", mock.Anything, "user-123", "hashed_pw")
	tokenRepo.AssertCalled(t, "DeleteTokensByUserAndType", mock.Anything, "user-123", entities.PasswordReset)
}
