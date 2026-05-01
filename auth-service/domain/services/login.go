package services

import (
	"context"
	"errors"
	"log"

	"github.com/greenbuildr/auth-service/domain/entities"
	domainerrors "github.com/greenbuildr/auth-service/domain/errors"
	"github.com/greenbuildr/auth-service/ports"
)

type LoginService struct {
	userRepo     ports.UserRepository
	hasher       ports.PasswordHasher
	jwtGenerator ports.TokenGenerator
}

func NewLoginService(
	userRepo ports.UserRepository,
	hasher ports.PasswordHasher,
	jwtGenerator ports.TokenGenerator,
) *LoginService {
	return &LoginService{
		userRepo:     userRepo,
		hasher:       hasher,
		jwtGenerator: jwtGenerator,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

func (s *LoginService) Login(ctx context.Context, input LoginInput) (string, error) {
	log.Printf("login: starting login for email=%s", input.Email)

	//unknown user
	user, err := s.fetchUser(ctx, input.Email)
	if err != nil {
		return "", err
	}

	//email verificaiton not complete
	if err := s.checkEmailVerified(user); err != nil {
		return "", err
	}

	//invalid credentials
	if err := s.verifyPassword(user, input.Password); err != nil {
		return "", err
	}

	//return jwt with user email and id encrypted
	return s.issueToken(user)
}

func (s *LoginService) fetchUser(ctx context.Context, email string) (*entities.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domainerrors.ErrUserNotFound) {
			log.Printf("login: user not found email=%s", email)
			return nil, domainerrors.ErrInvalidCredentials
		}
		log.Printf("login: db error fetching user email=%s: %v", email, err)
		return nil, domainerrors.ErrInternalError
	}
	return user, nil
}

func (s *LoginService) verifyPassword(user *entities.User, plainPassword string) error {
	if err := s.hasher.Compare(user.PasswordHash, plainPassword); err != nil {
		log.Printf("login: invalid password for userID=%s", user.ID)
		return domainerrors.ErrInvalidCredentials
	}
	return nil
}

func (s *LoginService) checkEmailVerified(user *entities.User) error {
	if !user.IsEmailVerified() {
		log.Printf("login: email not verified for userID=%s", user.ID)
		return domainerrors.ErrEmailNotVerified
	}
	return nil
}

func (s *LoginService) issueToken(user *entities.User) (string, error) {
	token, err := s.jwtGenerator.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Printf("login: failed to generate JWT for userID=%s: %v", user.ID, err)
		return "", domainerrors.ErrInternalError
	}
	log.Printf("login: successful for userID=%s", user.ID)
	return token, nil
}
