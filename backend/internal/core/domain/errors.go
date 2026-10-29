package domain

import (
	"errors"
)

var (
	INVALID_REQUEST    = NewCustomError("INVALID_REQUEST", "invalid request")
	INTERNAL_ERROR     = NewCustomError("INTERNAL_ERROR", "internal server error")
	UNAUTHORIZED       = NewCustomError("UNAUTHORIZED", "invalid authorization")
	EMAIL_ALREADY_USED = NewCustomError("EMAIL_ALREADY_USED", "email is already in use")
	FORBIDDEN          = NewCustomError("FORBIDDEN", "...")
	NOT_FOUND          = errors.New("not found")
)

type CustomError struct {
	Code    string `json:"code" enums:"INVALID_REQUEST, INTERNAL_ERROR, SCHEDULE_EXISTS, SLOT_ALREADY_BOOKED, FORBIDEN, NOT_FOUND"`
	Message string `json:"message" example:"invalid request"`
}

type InternalError struct {
	Code    string `json:"code" enums:"INTERNAL_ERROR"`
	Message string `json:"message" example:"internal server error"`
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
