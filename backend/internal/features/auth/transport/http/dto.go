package auth_transport

import "time"

type LoginDTO struct {
	Email    string `json:"email" example:"test.user@email.com" validate:"required"`
	Password string `json:"password" example:"hashpassword" validate:"required,min=8,max=64,ascii"`
}

type TokensDTOResponse struct {
	AccessToken  string `json:"access_token" example:"new_access_token"`
	RefreshToken string `json:"refresh_token" example:"new_refresh_token"`
	TokenType    string `json:"token_type" example:"Bearer"`
	ExpiresIn    int64  `json:"expires_in" example:"600"`
}

type RegisterDTO struct {
	Email    string `json:"email" validate:"required,email" example:"m4rkek@gmail.com" extensions:"x-order=0"`
	Username string `json:"username" example:"m4rkek" extensions:"x-order=1"`
	Password string `json:"password" validate:"required,min=8,max=64,ascii" example:"supersecretpass" extensions:"x-order=2"`
	Role     string `json:"role" validate:"required" example:"admin" extensions:"x-order=3"`
}

type RegisterDTOResponse struct {
	ID        string    `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6" extensions:"x-order=0"`
	Username  string    `json:"username" example:"m4rkek" extensions:"x-order=1"`
	Email     string    `json:"email" example:"m4rkek@gmail.com" extensions:"x-order=2"`
	Role      string    `json:"role" example:"user" extensions:"x-order=3"`
	CreatedAt time.Time `json:"createdAt" example:"2026-03-25T12:00:41.267Z" extensions:"x-order=4"`
}

type ResetPasswordDTORequest struct {
	Email string `json:"email" example:"test.user@email.com" validate:"required,email"`
}

type ConfirmPasswordResetDTORequest struct {
	Token       string `json:"token" example:"reset-token" validate:"required"`
	NewPassword string `json:"new_password" example:"new_password" validate:"required,min=8,max=64,ascii"`
}

func toTokensDTOResponse(accessToken, refreshToken string, expiresIn int64) TokensDTOResponse {
	return TokensDTOResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}
}
