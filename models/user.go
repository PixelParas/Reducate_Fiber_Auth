package models

import (
	"fmt"
	"regexp"
	"time"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"type:varchar(20);default:'user';not null" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}

type AuthInput struct {
	Email    string `json: "email"`
	Password string `json: "password"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (a *AuthInput) Validate() error {
	if a.Email == "" {
		return fmt.Errorf("Email is a required field")
	}
	if len(a.Password) < 6 {
		return fmt.Errorf("Password should be at least 6 letter")
	}

	return nil
}
