package ports

import (
	"context"

	"github.com/greenbuildr/auth-service/domain/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserById(ctx context.Context, id string) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error) //for login route
	UpdateUserPassword(ctx context.Context, id string, newPassword string) error
	VerifyUserEmail(ctx context.Context, id string) error
}

type TokenRepository interface {
	CreateToken(ctx context.Context, token *entities.Token) error
	GetTokenByHash(ctx context.Context, tokenHash string) (*entities.Token, error)
	DeleteToken(ctx context.Context, tokenHash string) error                                          //to consume token after use
	DeleteTokensByUserAndType(ctx context.Context, userID string, tokenType entities.TokenType) error //invalidate related tokens before issuing new one
}
