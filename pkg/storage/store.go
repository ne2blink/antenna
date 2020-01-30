package storage

import "io"

// AppStore provides app storage
type AppStore interface {
	CreateApp(App) (string, error)
	UpdateApp(App) error
	GetApp(ID string) (App, error)
	DeleteApp(ID string) error
	ListApps() ([]App, error)
	ListSubscribers(ID string) ([]int64, error)
}

// SubscriberStore provides subscriber storage
type SubscriberStore interface {
	ListSubscribedApps(ChatID int64) ([]App, error)
	Subscribe(ChatID int64, AppID string) error
	Unsubscribe(ChatID int64, AppID string) error
	UnsubscribeAll(ChatID int64) error
}

// Store provides storage
type Store interface {
	AppStore
	SubscriberStore
	io.Closer
}
