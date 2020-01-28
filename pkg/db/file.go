package db

import "github.com/ne2blink/antenna/pkg/storage"

type file struct {
	path string
}

func (f *file) CreateApp(name string) (storage.App, error) {
	return storage.App{}, nil
}

func (f *file) UpdateApp(storage.App) error {
	return nil
}

func (f *file) GetApp(ID string) (storage.App, error) {
	return storage.App{}, nil
}

func (f *file) DeleteApp(ID string) error {
	return nil
}

func (f *file) ListApps() ([]storage.App, error) {
	return []storage.App{}, nil
}

func (f *file) ListSubscribers(ID string) ([]string, error) {
	return []string{}, nil
}

func (f *file) ListSubscribedApps(ChatID string) ([]storage.App, error) {
	return []storage.App{}, nil
}

func (f *file) Subscribe(ChatID, AppID string) error {
	return nil
}

func (f *file) Unsubscribe(ChatID, AppID string) error {
	return nil
}

func (f *file) UnsubscribeAll(ChatID string) error {
	return nil
}

func newFile(options map[string]interface{}) (storage.Store, error) {
	// panic("todo")
	path := "./file.db"
	for k, v := range options {
		switch k {
		case "path":
			if s, ok := v.(string); ok {
				path = s
			}
		}
	}
	file := file{path: path}
	return &file, nil
}

func init() {
	storage.Register("file", newFile)
}
