package models

// Sub is chat subscribed app for mysql table
type Sub struct {
	ChatID int64  `gorm:"primary_key"`
	AppID  string `gorm:"primary_key"`
}
