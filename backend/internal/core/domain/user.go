package domain

import (
	"fmt"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

var AllowedRoles = []string{"admin", "user"}

type User struct {
	ID                     string `validate:"required,uuid"`
	Email                  string `validate:"email"`
	Password               string `validate:"min=8,max=64,ascii"`
	Role                   string `validate:"oneof=admin user"`
	CreatedAt              time.Time
	Username               string `validate:"min=4,max=15,username_validation"`
	DisplayName            string `validate:"omitempty,min=1,max=30"`
	AvatarURL              string
	Status                 string `validate:"omitempty,max=100"`
	EmailVerified          bool
	FollowersCount         int
	FollowingCount         int
	PasswordResetToken     string
	PasswordResetExpiresAt time.Time
}

func (u *User) Validate() error {
	validator := validator.New()

	err := validator.RegisterValidation("username_validation", validateUsername)
	if err != nil {
		return fmt.Errorf("validator err: %w", err)
	}

	if err := validator.Struct(u); err != nil {
		return fmt.Errorf("validator struct: %w: %w", err, ErrInvalidRequest)
	}

	return nil
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return re.MatchString(username)
}

func (u *User) UpdateDisplayName(displayName string) error {
	if utf8.RuneCountInString(displayName) == 0 || utf8.RuneCountInString(displayName) > 30 {
		return fmt.Errorf("invalid display_name=%s length: %w", displayName, ErrInvalidRequest)
	}

	u.DisplayName = displayName

	return nil
}

func (u *User) UpdateUsername(username string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !re.MatchString(username) {
		return fmt.Errorf("invalid username=%s: %w", username, ErrInvalidRequest)
	}
	u.Username = username

	return nil
}

func (u *User) UpdateStatus(status string) error {
	if utf8.RuneCountInString(status) == 0 || utf8.RuneCountInString(status) > 100 {
		return fmt.Errorf("invalid status=%s length: %w", status, ErrInvalidRequest)
	}
	u.Status = status

	return nil
}
