package db

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	mysqldrv "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/greenbuildr/auth-service/domain/entities"
	domainerrors "github.com/greenbuildr/auth-service/domain/errors"
)

type MySQLRepository struct {
	db *sqlx.DB
}

func NewMySQLRepository(db *sqlx.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

// dbUser and dbToken are adapter-internal only — never leave this package.
// sql.NullTime handles the nullable email_verified_at column.
type dbUser struct {
	ID              string       `db:"id"`
	Email           string       `db:"email"`
	PasswordHash    string       `db:"password_hash"`
	EmailVerifiedAt sql.NullTime `db:"email_verified_at"`
	CreatedAt       time.Time    `db:"created_at"`
	UpdatedAt       time.Time    `db:"updated_at"`
}

type dbToken struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	TokenType string    `db:"token_type"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

// functions to map db rows back to domain entities
func (r dbUser) toDomain() *entities.User {
	u := &entities.User{
		ID:           r.ID,
		Email:        r.Email,
		PasswordHash: r.PasswordHash,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
	if r.EmailVerifiedAt.Valid {
		u.EmailVerifiedAt = &r.EmailVerifiedAt.Time
	}
	return u
}
func (r dbToken) toDomain() *entities.Token {
	return &entities.Token{
		ID:        r.ID,
		UserID:    r.UserID,
		TokenHash: r.TokenHash,
		TokenType: entities.TokenType(r.TokenType),
		ExpiresAt: r.ExpiresAt,
		CreatedAt: r.CreatedAt,
	}
}

func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysqldrv.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

// --- UserRepository ---

func (r *MySQLRepository) CreateUser(ctx context.Context, user *entities.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, email, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		user.ID, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		if isDuplicateEntry(err) {
			return domainerrors.ErrEmailAlreadyExists
		}
		return domainerrors.ErrInternalError
	}
	return nil
}

func (r *MySQLRepository) GetUserById(ctx context.Context, id string) (*entities.User, error) {
	var row dbUser
	err := r.db.GetContext(ctx, &row, `SELECT * FROM users WHERE id = ?`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainerrors.ErrUserNotFound
		}
		return nil, domainerrors.ErrInternalError
	}
	return row.toDomain(), nil
}

func (r *MySQLRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var row dbUser
	err := r.db.GetContext(ctx, &row, `SELECT * FROM users WHERE email = ?`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainerrors.ErrUserNotFound
		}
		return nil, domainerrors.ErrInternalError
	}
	return row.toDomain(), nil
}

func (r *MySQLRepository) UpdateUserPassword(ctx context.Context, id string, newPassword string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash = ? WHERE id = ?`,
		newPassword, id,
	)
	if err != nil {
		return domainerrors.ErrInternalError
	}
	return nil
}

func (r *MySQLRepository) VerifyUserEmail(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET email_verified_at = NOW() WHERE id = ?`,
		id,
	)
	if err != nil {
		return domainerrors.ErrInternalError
	}
	return nil
}

// --- TokenRepository ---

func (r *MySQLRepository) CreateToken(ctx context.Context, token *entities.Token) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO auth_tokens (id, user_id, token_hash, token_type, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		token.ID, token.UserID, token.TokenHash, string(token.TokenType), token.ExpiresAt, token.CreatedAt,
	)
	if err != nil {
		return domainerrors.ErrInternalError
	}
	return nil
}

func (r *MySQLRepository) GetTokenByHash(ctx context.Context, rawToken string) (*entities.Token, error) {
	var row dbToken
	err := r.db.GetContext(ctx, &row,
		`SELECT * FROM auth_tokens WHERE token_hash = ?`,
		hashToken(rawToken),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainerrors.ErrInvalidToken
		}
		return nil, domainerrors.ErrInternalError
	}
	return row.toDomain(), nil
}

func (r *MySQLRepository) DeleteToken(ctx context.Context, rawToken string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM auth_tokens WHERE token_hash = ?`,
		hashToken(rawToken),
	)
	if err != nil {
		return domainerrors.ErrInternalError
	}
	return nil
}

func (r *MySQLRepository) DeleteTokensByUserAndType(ctx context.Context, userID string, tokenType entities.TokenType) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM auth_tokens WHERE user_id = ? AND token_type = ?`,
		userID, string(tokenType),
	)
	if err != nil {
		return domainerrors.ErrInternalError
	}
	return nil
}
