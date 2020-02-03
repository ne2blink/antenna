package models

import (
	"strconv"

	"github.com/ne2blink/antenna/pkg/storage"
)

// App is an application for mysql table
type App struct {
	ID      uint `gorm:"primary_key"`
	Name    string
	Secret  string
	Private bool
}

// FromStoreApp is decoding storage.App to models.App
func (a *App) FromStoreApp(app storage.App) error {
	id, err := strconv.ParseUint(app.ID, 10, 32)
	if err != nil {
		return err
	}
	a.ID = uint(id)
	a.Name = app.Name
	a.Secret = app.Secret
	a.Private = app.Private
	return nil
}

// ToStoreApp is encoding models.App to storage.App
func (a App) ToStoreApp() storage.App {
	app := storage.App{
		ID:      strconv.FormatUint(uint64(a.ID), 10),
		Name:    a.Name,
		Secret:  a.Secret,
		Private: a.Private,
	}
	return app
}
