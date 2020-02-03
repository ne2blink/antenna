package models

type AppSub struct {
	ID     uint   `gorm:"primary_key"`
	AppID  string `gorm:"index"`
	ChatID int64  `gorm:"index"`
}
