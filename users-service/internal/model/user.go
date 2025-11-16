package model

import "time"

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID           int64
	Email        string
	Name         string
	PasswordHash string
	Role         string // "user" или "admin"
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateUserRequest struct {
	Email    string
	Name     string
	Password string
}
