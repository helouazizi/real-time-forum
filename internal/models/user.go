package models

import "time"

type User struct {
	ID        int    `json:"id"`
	Age       int    `json:"age"`
	Nickname  string `json:"nickname"`
	Gender    string `json:"gender"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"` // hide in JSON
	// Bio               string    `json:"bio"`
	// AvatarURL         string    `json:"avatar_url"`
	// Role              string    `json:"role"`
	IsActive bool `json:"is_active"`
	// EmailVerified     bool      `json:"email_verified"`
	// VerificationToken string    `json:"verification_token,omitempty"`
	// SessionToken string `json:"session_token,omitempty"`
	// SessionExpiresAt  time.Time `json:"session_expires_at,omitempty"`
	// LastLoginAt       time.Time `json:"last_login_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserLogin struct {
	LoginId      string `json:"login_id"`
	Password     string `json:"password"`
	SessionToken string `json:"token"`
}

type UserInputErrors struct {
	HasError       bool
	Nickname       string
	Email          string
	Pass           string
	Age            string
	LastName       string
	FirstName      string
	Gender         string
	PostTilte      string
	PostContent    string
	Postcategories string
}

type SuccesMessage struct {
	Message string
}
