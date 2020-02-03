package mysql

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // register mysql
	"github.com/ne2blink/antenna/pkg/storage"
	"github.com/ne2blink/antenna/pkg/storage/mysql/models"
)

func newStore(options map[string]interface{}) (storage.Store, error) {
	conn, _ := options["conn"].(string)
	if conn == "" {
		return nil, errors.New("mysql: missing connection string")
	}
	db, err := gorm.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.App{}, &models.AppSub{}, &models.ChatSub{})

	return &store{db: db}, nil
}

func init() {
	storage.Register("mysql", newStore)
}
