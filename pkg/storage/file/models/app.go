package models

import (
	"encoding/json"

	"github.com/ne2blink/antenna/pkg/storage"
)

// App represnets an application
type App struct {
	storage.App
	SubscribedChatIDs []int64 `json:"subscribed_chat_ids,omitempty"`
}

// FromJSON is decoding json to App
func (a *App) FromJSON(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, a)
}

// ToJSON is encoding App to json
func (a App) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

// FromStoreApp is decoding storage.App to models.App
func (a *App) FromStoreApp(app storage.App) {
	a.ID = app.ID
	a.Name = app.Name
	a.Secret = app.Secret
	a.Private = app.Private
}

// ToStoreApp is encoding models.App to storage.App
func (a App) ToStoreApp() storage.App {
	app := storage.App{
		ID:      a.ID,
		Name:    a.Name,
		Secret:  a.Secret,
		Private: a.Private,
	}
	return app
}
