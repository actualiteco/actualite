package domain

var (
	ErrInvalidRequest          = NewCustomError("INVALID_REQUEST", "invalid request")
	ErrInternal                = NewCustomError("INTERNAL_ERROR", "internal server error")
	ErrUnauthorized            = NewCustomError("UNAUTHORIZED", "invalid authorization")
	ErrEmailAlreadyUsed        = NewCustomError("EMAIL_ALREADY_USED", "email is already in use")
	ErrForbidden               = NewCustomError("FORBIDDEN", "forbidden")
	ErrNotFound                = NewCustomError("NOT_FOUND", "not found")
	ErrConflict                = NewCustomError("CONFLICT", "conflict")
	ErrUsernameAlreadyUsed     = NewCustomError("USERNAME_ALREADY_USED", "username is already in use")
	ErrTokenExpired            = NewCustomError("TOKEN_EXPIRED", "access-token expired")
	ErrSessionRevokedOrExpired = NewCustomError("SESSION_REVOKED_OR_EXPIRED", "session was revoked or expired")
	ErrRequestIsTooLarge       = NewCustomError("REQUEST IS TOO LARGE", "request is too large")
	ErrUnsupportedImageType    = NewCustomError("UNSUPPORTED IMAGE TYPE", "image type is invalid")
)

type CustomError struct {
	Code    string `json:"code"`
	Message string `json:"message" example:"invalid request"`
}

type InternalError struct {
	Code    string `json:"code" enums:"INTERNAL_ERROR"`
	Message string `json:"message"`
}

func NewCustomError(code string, msg string) CustomError {
	return CustomError{
		Code:    code,
		Message: msg,
	}
}

func (e CustomError) Error() string {
	return e.Message
}
