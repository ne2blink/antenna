package models

// Sub is chat subscribed app for mysql table
type Sub struct {
	ID     uint   `gorm:"primary_key"`
	ChatID int64  `gorm:"index"`
	AppID  string `gorm:"index"`
}
