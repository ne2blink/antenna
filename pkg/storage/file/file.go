package file

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/ne2blink/antenna/pkg/storage"
	"github.com/ne2blink/antenna/pkg/storage/file/models"
)

type file struct {
	path string
}

func (f file) CreateApp(app storage.App) (string, error) {
	db, err := bolt.Open(f.path, 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var a models.App
	a.FromStoreApp(app)
	err = db.Update(func(tx *bolt.Tx) error {
		// Open Apps Bucket
		b := tx.Bucket([]byte("Apps"))
		// ID Auto Increment
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		// Change type to string
		a.ID = strconv.FormatUint(id, 10)
		// Marshal app data into bytes.
		buf, err := a.ToJSON()
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return b.Put([]byte(a.ID), buf)
	})
	return a.ID, err
}

func (f file) UpdateApp(app storage.App) error {
	db, err := bolt.Open(f.path, 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var a models.App
	a.FromStoreApp(app)
	err = db.Update(func(tx *bolt.Tx) error {
		// Open Apps Bucket
		b := tx.Bucket([]byte("Apps"))
		// Marshal app data into bytes.
		buf, err := a.ToJSON()
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return b.Put([]byte(a.ID), buf)
	})
	return err
}

func (f file) GetApp(ID string) (storage.App, error) {
	db, err := bolt.Open(f.path, 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var a models.App
	err = db.View(func(tx *bolt.Tx) error {
		// Open Apps Bucket
		b := tx.Bucket([]byte("Apps"))
		v := b.Get([]byte(ID))
		if v == nil {
			return errors.New(ID + ": not found")
		}
		return a.FromJSON(v)
	})
	return a.ToStoreApp(), err
}

func (f file) DeleteApp(ID string) error {
	db, err := bolt.Open(f.path, 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		// Open Apps Bucket
		b := tx.Bucket([]byte("Apps"))
		return b.Delete([]byte(ID))
	})
	return err
}

func (f file) ListApps() ([]storage.App, error) {
	return []storage.App{}, nil
}

func (f file) ListSubscribers(ID string) ([]string, error) {
	return []string{}, nil
}

func (f file) ListSubscribedApps(ChatID int64) ([]storage.App, error) {
	return []storage.App{}, nil
}

func (f file) Subscribe(ChatID int64, AppID string) error {
	return nil
}

func (f file) Unsubscribe(ChatID int64, AppID string) error {
	return nil
}

func (f file) UnsubscribeAll(ChatID int64) error {
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

	file := file{path: path}
	return &file, nil
}

func init() {
	storage.Register("file", newFile)
}
