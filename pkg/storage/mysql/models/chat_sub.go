package models

type ChatSub struct {
	ID     uint   `gorm:"primary_key"`
	ChatID int64  `gorm:"index"`
	AppID  string `gorm:"index"`
}
