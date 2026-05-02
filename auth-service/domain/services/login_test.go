package services

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/greenbuildr/auth-service/domain/entities"
	domainerrors "github.com/greenbuildr/auth-service/domain/errors"
)

func newTestLoginService(u *mockUserRepo, h *mockHasher, j *mockJWTGenerator) *LoginService {
	return NewLoginService(u, h, j)
}

func verifiedUser() *entities.User {
	verified := time.Now()
	return &entities.User{
		ID:              "user-123",
		Email:           "test@example.com",
		PasswordHash:    "hashed_pw",
		EmailVerifiedAt: &verified,
	}
}

func TestLogin_UnknownUser(t *testing.T) {
	userRepo := &mockUserRepo{}
	hasher := &mockHasher{}
	jwtGen := &mockJWTGenerator{}

	userRepo.On("GetUserByEmail", mock.Anything, "unknown@example.com").Return(nil, domainerrors.ErrUserNotFound)

	svc := newTestLoginService(userRepo, hasher, jwtGen)
	_, err := svc.Login(context.Background(), LoginInput{Email: "unknown@example.com", Password: "password123"})

	assert.Equal(t, domainerrors.ErrInvalidCredentials, err)
	hasher.AssertNotCalled(t, "Compare")
	jwtGen.AssertNotCalled(t, "GenerateJWT")
}

func TestLogin_EmailNotVerified(t *testing.T) {
	userRepo := &mockUserRepo{}
	hasher := &mockHasher{}
	jwtGen := &mockJWTGenerator{}

	unverifiedUser := &entities.User{
		ID:              "user-123",
		Email:           "test@example.com",
		PasswordHash:    "hashed_pw",
		EmailVerifiedAt: nil,
	}
	userRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(unverifiedUser, nil)

	svc := newTestLoginService(userRepo, hasher, jwtGen)
	_, err := svc.Login(context.Background(), LoginInput{Email: "test@example.com", Password: "password123"})

	assert.Equal(t, domainerrors.ErrEmailNotVerified, err)
	hasher.AssertNotCalled(t, "Compare")
	jwtGen.AssertNotCalled(t, "GenerateJWT")
}

func TestLogin_InvalidPassword(t *testing.T) {
	userRepo := &mockUserRepo{}
	hasher := &mockHasher{}
	jwtGen := &mockJWTGenerator{}

	userRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(verifiedUser(), nil)
	hasher.On("Compare", "hashed_pw", "wrongpassword").Return(domainerrors.ErrInvalidCredentials)

	svc := newTestLoginService(userRepo, hasher, jwtGen)
	_, err := svc.Login(context.Background(), LoginInput{Email: "test@example.com", Password: "wrongpassword"})

	assert.Equal(t, domainerrors.ErrInvalidCredentials, err)
	jwtGen.AssertNotCalled(t, "GenerateJWT")
}

func TestLogin_Success(t *testing.T) {
	userRepo := &mockUserRepo{}
	hasher := &mockHasher{}
	jwtGen := &mockJWTGenerator{}

	const testSecret = "test-secret"
	testToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   "user-123",
		"email": "test@example.com",
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, _ := testToken.SignedString([]byte(testSecret))

	userRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(verifiedUser(), nil)
	hasher.On("Compare", "hashed_pw", "password123").Return(nil)
	jwtGen.On("GenerateJWT", "user-123", "test@example.com").Return(tokenStr, nil)

	svc := newTestLoginService(userRepo, hasher, jwtGen)
	result, err := svc.Login(context.Background(), LoginInput{Email: "test@example.com", Password: "password123"})

	assert.NoError(t, err)

	parsed, parseErr := jwt.Parse(result, func(t *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})
	assert.NoError(t, parseErr)
	claims := parsed.Claims.(jwt.MapClaims)
	assert.Equal(t, "user-123", claims["sub"])
	assert.Equal(t, "test@example.com", claims["email"])
}
