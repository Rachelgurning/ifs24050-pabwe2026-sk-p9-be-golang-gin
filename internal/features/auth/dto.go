package auth

import "app/internal/database"

type RegisterRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Name     *string `json:"name"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Bio      *string `json:"bio"`
	Password *string `json:"password"`
}

type LoginResult struct {
	Token       string         `json:"token"`
	AccessToken string         `json:"access_token"`
	TokenType   string         `json:"token_type"`
	User        *database.User `json:"user"`
}
