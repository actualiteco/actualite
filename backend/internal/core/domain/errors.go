package domain

var (
	ErrInvalidRequest   = NewCustomError("INVALID_REQUEST", "invalid request")
	ErrInternal         = NewCustomError("INTERNAL_ERROR", "internal server error")
	ErrUnauthorized     = NewCustomError("UNAUTHORIZED", "invalid authorization")
	ErrEmailAlreadyUsed = NewCustomError("EMAIL_ALREADY_USED", "email is already in use")
	ErrForbidden        = NewCustomError("FORBIDDEN", "...")
	ErrNotFound         = NewCustomError("NOT_FOUND", "not found")
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
