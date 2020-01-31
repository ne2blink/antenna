package storage

import "io"

// AppStore provides app storage
type AppStore interface {
	CreateApp(App) (string, error)
	UpdateApp(App) error
	GetApp(id string) (App, error)
	DeleteApp(id string) error
	ListApps() ([]App, error)
	ListSubscribers(id string) ([]int64, error)
}

// SubscriberStore provides subscriber storage
type SubscriberStore interface {
	ListSubscribedApps(chatID int64) ([]App, error)
	Subscribe(chatID int64, appID string) error
	Unsubscribe(chatID int64, appID string) error
	UnsubscribeAll(chatID int64) error
}

// Store provides storage
type Store interface {
	AppStore
	SubscriberStore
	io.Closer
}
