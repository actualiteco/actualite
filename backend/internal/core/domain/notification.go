package domain

import (
	"time"
	"uuid"
)

type NotificationType string

const (
	NotificationTypeFollow  NotificationType = "follow"
	NotificationTypeLike    NotificationType = "like"
	NotificationTypeComment NotificationType = "comment"
)

type Notification struct {
	ID          uuid.UUID
	RecipientID uuid.UUID
	ActorID     uuid.UUID
	Type        NotificationType
	PostID      *uuid.UUID
	CommentID   *uuid.UUID
	IsRead      bool
	CreatedAt   time.Time
}

func NewFollowNotification(recipientID, actorID uuid.UUID) Notification {
	return Notification{
		ID:          uuid.New(),
		RecipientID: recipientID,
		ActorID:     actorID,
		Type:        NotificationTypeFollow,
	}
}
