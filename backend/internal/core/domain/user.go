package domain

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

var AllowedRoles = []string{"admin", "user"}

type User struct {
	ID        string `validate:"required,uuid"`
	Email     string `validate:"email"`
	Username  string `validate:"omitempty,min=4,max=15,username_validation"`
	Password  string
	Role      string `validate:"oneof=admin user"`
	CreatedAt time.Time
}

func (u *User) Validate() error {
	val := validator.New()

	err := val.RegisterValidation("username_validation", validateUsername)
	if err != nil {
		return fmt.Errorf("validator err: %w", err)
	}

	if err := val.Struct(u); err != nil {
		return fmt.Errorf("validator struct: %w", err)
	}

	return nil
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return re.MatchString(username)
}
