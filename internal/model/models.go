package model

import (
	"time"
)

type User struct {
	ID           int        `json:"id"`
	Login        string     `json:"login"`
	PasswordHash string     `json:"-"`
	FullName     string     `json:"full_name"`
	Phone        string     `json:"phone"`
	Email        string     `json:"email"`
	IsAdmin      bool       `json:"is_admin"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type Application struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	CourseName    string    `json:"course_name"`
	StartDate     time.Time `json:"start_date"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	Review        *string   `json:"review,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	User          *User     `json:"user,omitempty"`
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ApplicationRequest struct {
	CourseName    string `json:"course_name"`
	StartDate     string `json:"start_date"`
	PaymentMethod string `json:"payment_method"`
}

type ReviewRequest struct {
	Review string `json:"review"`
}

type StatusUpdateRequest struct {
	Status string `json:"status"`
}
