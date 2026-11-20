package domain

import (
	"fmt"
	"time"
	"unicode/utf8"
)

const MaxAvatarSize = 2 << 20

type User struct {
	ID                         string `validate:"required,uuid"`
	Email                      string `validate:"email"`
	Password                   string
	Role                       string `validate:"oneof=admin user"`
	CreatedAt                  time.Time
	Username                   string `validate:"min=4,max=15,username_validation"`
	DisplayName                string `validate:"omitempty,min=1,max=30"`
	AvatarURL                  string
	Status                     string `validate:"omitempty,max=100"`
	EmailVerified              bool
	FollowersCount             int
	FollowingCount             int
	PasswordResetToken         string
	PasswordResetExpiresAt     time.Time
	IsFollowing                bool
	EmailConfirmationToken     string
	EmailConfirmationExpiresAt time.Time
}

func (u User) Validate() error {
	if err := validate.Struct(u); err != nil {
		return fmt.Errorf("validator struct: %w: %w", err, ErrInvalidRequest)
	}

	return nil
}

func (u *User) UpdateDisplayName(displayName string) error {
	if utf8.RuneCountInString(displayName) == 0 || utf8.RuneCountInString(displayName) > 30 {
		return fmt.Errorf("invalid display_name=%s length: %w", displayName, ErrInvalidRequest)
	}

	u.DisplayName = displayName

	return nil
}

func (u *User) UpdateUsername(username string) error {
	if err := ValidateUsername(username); err != nil {
		return err
	}
	u.Username = username

	return nil
}

func ValidateUsername(username string) error {
	if utf8.RuneCountInString(username) < 4 || utf8.RuneCountInString(username) > 15 {
		return fmt.Errorf("invalid username=%s length: %w", username, ErrInvalidRequest)
	}

	if !usernameRe.MatchString(username) {
		return fmt.Errorf("invalid username=%s: %w", username, ErrInvalidRequest)
	}

	return nil
}

func (u *User) UpdateStatus(status string) error {
	if utf8.RuneCountInString(status) == 0 || utf8.RuneCountInString(status) > 100 {
		return fmt.Errorf("invalid status=%s length: %w", status, ErrInvalidRequest)
	}
	u.Status = status

	return nil
}
