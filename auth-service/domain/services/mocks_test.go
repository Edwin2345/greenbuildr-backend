package services

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/greenbuildr/auth-service/domain/entities"
)

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) CreateUser(ctx context.Context, user *entities.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *mockUserRepo) GetUserById(ctx context.Context, id string) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}
func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}
func (m *mockUserRepo) UpdateUserPassword(ctx context.Context, id, hash string) error {
	return m.Called(ctx, id, hash).Error(0)
}
func (m *mockUserRepo) VerifyUserEmail(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type mockTokenRepo struct{ mock.Mock }

func (m *mockTokenRepo) CreateToken(ctx context.Context, token *entities.Token) error {
	return m.Called(ctx, token).Error(0)
}
func (m *mockTokenRepo) GetTokenByHash(ctx context.Context, rawToken string) (*entities.Token, error) {
	args := m.Called(ctx, rawToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Token), args.Error(1)
}
func (m *mockTokenRepo) DeleteToken(ctx context.Context, rawToken string) error {
	return m.Called(ctx, rawToken).Error(0)
}
func (m *mockTokenRepo) DeleteTokensByUserAndType(ctx context.Context, userID string, tokenType entities.TokenType) error {
	return m.Called(ctx, userID, tokenType).Error(0)
}

type mockHasher struct{ mock.Mock }

func (m *mockHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}
func (m *mockHasher) Compare(hashed, plain string) error {
	return m.Called(hashed, plain).Error(0)
}

type mockJWTGenerator struct{ mock.Mock }

func (m *mockJWTGenerator) GenerateJWT(userID, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Error(1)
}

type mockEmailSender struct{ mock.Mock }

func (m *mockEmailSender) SendVerificationEmail(ctx context.Context, to, url string) error {
	return m.Called(ctx, to, url).Error(0)
}
func (m *mockEmailSender) SendPasswordResetEmail(ctx context.Context, to, url string) error {
	return m.Called(ctx, to, url).Error(0)
}
