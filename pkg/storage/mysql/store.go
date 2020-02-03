package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/ne2blink/antenna/pkg/storage"
	"github.com/ne2blink/antenna/pkg/storage/mysql/models"
)

type store struct {
	db *gorm.DB
}

func (s store) CreateApp(app storage.App) (string, error) {
	mApp := models.App{Name: app.Name, Secret: app.Secret, Private: app.Private}
	err := s.db.Create(&mApp).Error
	if err != nil {
		return "", err
	}
	return mApp.ToStoreApp().ID, nil
}

func (s store) UpdateApp(app storage.App) error {
	mApp := models.App{}
	mApp.FromStoreApp(app)
	return s.db.Save(&mApp).Error
}

func (s store) GetApp(id string) (storage.App, error) {
	app := storage.App{ID: id}
	err := s.db.First(&app).Error
	if err != nil {
		return storage.App{}, err
	}
	return app, nil
}

func (s store) DeleteApp(id string) (err error) {
	mApp := models.App{}
	mApp.FromStoreApp(storage.App{ID: id})
	tx := s.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Where("app_id = ?", id).Delete(models.Sub{}).Error; err != nil {
		return err
	}
	if err := tx.Delete(&mApp).Error; err != nil {
		return err
	}
	return tx.Commit().Error
}

func (s store) ListApps() ([]storage.App, error) {
	apps := []storage.App{}
	err := s.db.Find(&apps).Error
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (s store) ListSubscribers(id string) ([]int64, error) {
	chats := []int64{}
	subs := []models.Sub{}
	err := s.db.Where("app_id = ?", id).Find(&subs).Error
	if err != nil {
		return nil, err
	}
	for _, appSub := range subs {
		chats = append(chats, appSub.ChatID)
	}
	return chats, nil
}

func (s store) ListSubscribedApps(chatID int64) ([]storage.App, error) {
	subs := []models.Sub{}
	apps := []storage.App{}
	err := s.db.Where("chat_id = ?", chatID).Find(&subs).Error
	if err != nil {
		return []storage.App{}, err
	}
	for _, chatSub := range subs {
		app := storage.App{ID: chatSub.AppID}
		err := s.db.First(&app).Error
		if err != nil {
			return []storage.App{}, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (s store) Subscribe(chatID int64, appID string) error {
	sub := models.Sub{ChatID: chatID, AppID: appID}
	return s.db.Create(&sub).Error

}

func (s store) Unsubscribe(chatID int64, appID string) error {
	return s.db.Where("chat_id = ?", chatID).Where("app_id = ?", appID).Delete(&models.Sub{}).Error
}

func (s store) UnsubscribeAll(chatID int64) error {
	return s.db.Where("chat_id = ?", chatID).Delete(&models.Sub{}).Error
}

func (s store) Close() error {
	return s.db.Close()
}
