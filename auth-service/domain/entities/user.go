package entities

import "time"

type User struct {
	ID               string
	Email            string
	PasswordHash     string
	EmailVerifiedAt  *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (u *User) IsEmailVerified() bool {
	return u.EmailVerifiedAt != nil
}
