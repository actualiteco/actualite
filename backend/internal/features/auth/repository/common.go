package auth_models

import (
	"database/sql"
	"time"
)

type UserModel struct {
	ID        string
	Username  string
	Email     string
	Password  string
	Role      string
	AvatarURL sql.NullString
	Bio       sql.NullString
	CreatedAt time.Time
}
