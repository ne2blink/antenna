package storage

// AppStore provides app storage
type AppStore interface {
	CreateApp(name string) (App, error)
	UpdateApp(App) error
	GetApp(ID string) (App, error)
	DeleteApp(ID string) error
	ListApps() ([]App, error)
	ListSubscribers(ID string) ([]string, error)
}

// SubscriberStore provides subscriber storage
type SubscriberStore interface {
	ListSubscribedApps(ChatID string) ([]App, error)
	Subscribe(ChatID, AppID string) error
	Unsubscribe(ChatID, AppID string) error
	UnsubscribeAll(ChatID string) error
}

// Store provides storage
type Store interface {
	AppStore
	SubscriberStore
}
