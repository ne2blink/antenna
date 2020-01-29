package db

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/ne2blink/antenna/pkg/storage"
)

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

func createChatsBucket(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists([]byte("Chats"))
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}
	return nil
}

func createAppsBucket(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists([]byte("Apps"))
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}
	return nil
}

func newFile(options map[string]interface{}) (storage.Store, error) {
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

	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createChatsErr := db.Update(createChatsBucket)
	if createChatsErr != nil {
		return nil, createChatsErr
	}

	createAppsErr := db.Update(createAppsBucket)
	if createAppsErr != nil {
		return nil, createAppsErr
	}

	return &file, nil
}

func init() {
	storage.Register("file", newFile)
}
