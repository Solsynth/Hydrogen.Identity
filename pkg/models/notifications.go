package models

import (
	"gorm.io/datatypes"
	"time"
)

type Notification struct {
	BaseModel

	Subject     string                                `json:"subject"`
	Content     string                                `json:"content"`
	Links       datatypes.JSONSlice[NotificationLink] `json:"links"`
	IsImportant bool                                  `json:"is_important"`
	ReadAt      *time.Time                            `json:"read_at"`
	SenderID    *uint                                 `json:"sender_id"`
	RecipientID uint                                  `json:"recipient_id"`
}

// NotificationLink Used to embed into notify and render actions
type NotificationLink struct {
	Label string `json:"label"`
	Url   string `json:"url"`
}

const (
	NotifySubscriberFirebase = "firebase"
)

type NotificationSubscriber struct {
	BaseModel

	UserAgent string `json:"user_agent"`
	Provider  string `json:"provider"`
	DeviceID  string `json:"device_id" gorm:"uniqueIndex"`
	AccountID uint   `json:"account_id"`
}
