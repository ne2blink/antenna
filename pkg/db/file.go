package db

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/ne2blink/antenna/pkg/storage"
)

type file struct {
	db *bolt.DB
}

func (f *file) CreateApp(app storage.App) (string, error) {
	err := f.db.Update(func(tx *bolt.Tx) error {
		// Open Apps Bucket
		b := tx.Bucket([]byte("Apps"))
		// 生成自增序列
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		fmt.Println(id)
		app.ID = strconv.FormatUint(id, 10)
		fmt.Println(app.ID)
		_, err = app.SetSecret("")
		if err != nil {
			return err
		}
		// Marshal user data into bytes.
		buf, err := app.ToJSON()
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return b.Put([]byte(app.ID), buf)
	})
	if err != nil {
		return app.ID, err
	}

	return app.ID, nil
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

func (f *file) Close() error {
	return f.db.Close()
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

	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	createChatsErr := db.Update(createChatsBucket)
	if createChatsErr != nil {
		return nil, createChatsErr
	}
	createAppsErr := db.Update(createAppsBucket)
	if createAppsErr != nil {
		return nil, createAppsErr
	}

	file := file{db: db}
	return &file, nil
}

func init() {
	storage.Register("file", newFile)
}
