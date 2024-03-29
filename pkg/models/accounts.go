package models

import (
	"path/filepath"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type Account struct {
	BaseModel

	Name              string                   `json:"name" gorm:"uniqueIndex"`
	Nick              string                   `json:"nick"`
	Description       string                   `json:"description"`
	Avatar            string                   `json:"avatar"`
	Banner            string                   `json:"banner"`
	Profile           AccountProfile           `json:"profile"`
	Sessions          []AuthSession            `json:"sessions"`
	Challenges        []AuthChallenge          `json:"challenges"`
	Factors           []AuthFactor             `json:"factors"`
	Contacts          []AccountContact         `json:"contacts"`
	Events            []ActionEvent            `json:"events"`
	MagicTokens       []MagicToken             `json:"-" gorm:"foreignKey:AssignTo"`
	ThirdClients      []ThirdClient            `json:"clients"`
	Notifications     []Notification           `json:"notifications" gorm:"foreignKey:RecipientID"`
	NotifySubscribers []NotificationSubscriber `json:"notify_subscribers"`
	ConfirmedAt       *time.Time               `json:"confirmed_at"`
	PowerLevel        int                      `json:"power_level"`
}

func (v Account) GetPrimaryEmail() AccountContact {
	val, _ := lo.Find(v.Contacts, func(item AccountContact) bool {
		return item.Type == EmailAccountContact && item.IsPrimary
	})
	return val
}

func (v Account) GetAvatarPath() string {
	basepath := viper.GetString("content")
	return filepath.Join(basepath, v.Avatar)
}

func (v Account) GetBannerPath() string {
	basepath := viper.GetString("content")
	return filepath.Join(basepath, v.Banner)
}

type AccountContactType = int8

const (
	EmailAccountContact = AccountContactType(iota)
)

type AccountContact struct {
	BaseModel

	Type       int8       `json:"type"`
	Content    string     `json:"content" gorm:"uniqueIndex"`
	IsPublic   bool       `json:"is_public"`
	IsPrimary  bool       `json:"is_primary"`
	VerifiedAt *time.Time `json:"verified_at"`
	AccountID  uint       `json:"account_id"`
}
