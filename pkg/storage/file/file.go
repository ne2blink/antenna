package file

import (
	"errors"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/ne2blink/antenna/pkg/storage"
	"github.com/ne2blink/antenna/pkg/storage/file/models"
	"github.com/ne2blink/antenna/pkg/utils"
)

var (
	bucket []string = []string{"Apps", "Chats"} //boltdb buckets
)

type dbFunc func(db *bolt.DB) error

type file struct {
	path string
}

func (f file) CreateApp(app storage.App) (string, error) {
	var a models.App
	a.FromStoreApp(app)
	err := f.db(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
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
	})
	return a.ID, err
}

func (f file) UpdateApp(app storage.App) error {
	var a models.App
	a.FromStoreApp(app)
	err := f.db(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
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
	})
	return err
}

func (f file) GetApp(ID string) (storage.App, error) {
	var a models.App
	err := f.db(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			// Open Apps Bucket
			b := tx.Bucket([]byte("Apps"))
			v := b.Get([]byte(ID))
			if v == nil {
				return errors.New(ID + ": not found")
			}
			return a.FromJSON(v)
		})
	})
	return a.ToStoreApp(), err
}

func (f file) DeleteApp(ID string) error {
	err := f.db(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			// Open Apps Bucket
			b := tx.Bucket([]byte("Apps"))
			return b.Delete([]byte(ID))
		})
	})
	return err
}

func (f file) ListApps() ([]storage.App, error) {
	var apps []storage.App
	err := f.db(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			// Open Apps Bucket
			b := tx.Bucket([]byte("Apps"))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				var a models.App
				err := a.FromJSON(v)
				if err != nil {
					return err
				}
				apps = append(apps, a.ToStoreApp())
			}
			return nil
		})
	})
	return apps, err
}

func (f file) ListSubscribers(ID string) ([]int64, error) {
	var a models.App
	err := f.db(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			// Open Apps Bucket
			b := tx.Bucket([]byte("Apps"))
			v := b.Get([]byte(ID))
			if v == nil {
				return errors.New(ID + ": not found")
			}
			return a.FromJSON(v)
		})
	})
	return a.SubscribedChatIDs, err
}

func (f file) checkAppID(ID string) error {
	err := f.db(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Apps"))
			v := b.Get([]byte(ID))
			if v == nil {
				return errors.New("not found")
			}
			return nil
		})
	})
	return err
}

func (f file) checkChatID(ChatID int64) error {
	err := f.db(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Chats"))
			v := b.Get(utils.I64tob(ChatID))
			if v == nil {
				return errors.New("not found")
			}
			return nil
		})
	})
	return err
}

func (f file) createChat(ChatID int64) error {
	err := f.db(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			c := models.Chat{ID: ChatID}
			buf, err := c.ToJSON()
			if err != nil {
				return err
			}
			b := tx.Bucket([]byte("Chats"))
			return b.Put(utils.I64tob(ChatID), []byte(buf))
		})
	})
	return err
}

func (f file) ListSubscribedApps(ChatID int64) ([]storage.App, error) {
	err := f.checkChatID(ChatID)
	if err != nil {
		f.createChat(ChatID)
	}
	var apps []storage.App
	err = f.db(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			var c models.Chat
			// Open Apps Bucket
			b := tx.Bucket([]byte("Chats"))
			v := b.Get(utils.I64tob(ChatID))
			if v == nil {
				return errors.New("not found")
			}
			err := c.FromJSON(v)
			if err != nil {
				return err
			}
			for _, ID := range c.SubscribedAppIDs {
				b := tx.Bucket([]byte("Apps"))
				v := b.Get([]byte(ID))
				var a models.App
				err := a.FromJSON(v)
				if err != nil {
					return err
				}
				apps = append(apps, a.ToStoreApp())
			}
			return nil
		})
	})
	return apps, err
}

func (f file) Subscribe(ChatID int64, AppID string) error {
	var c models.Chat
	err := f.checkChatID(ChatID)
	if err != nil {
		f.createChat(ChatID)
	}
	err = f.checkAppID(AppID)
	if err != nil {
		return errors.New(AppID + ": not found")
	}
	err = f.db(func(db *bolt.DB) error {
		return db.Batch(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Chats"))
			v := b.Get(utils.I64tob(ChatID))
			if v == nil {
				return errors.New("not found")
			}
			err := c.FromJSON(v)
			if err != nil {
				return err
			}
			if utils.ChickInString(c.SubscribedAppIDs, AppID) {
				return errors.New("already subscribed this app")
			}
			c.SubscribedAppIDs = append(c.SubscribedAppIDs, AppID)
			buf, err := c.ToJSON()
			if err != nil {
				return err
			}
			err = b.Put(utils.I64tob(ChatID), []byte(buf))
			if err != nil {
				return err
			}
			b = tx.Bucket([]byte("Apps"))
			v = b.Get([]byte(AppID))
			var a models.App
			err = a.FromJSON(v)
			if err != nil {
				return err
			}
			a.SubscribedChatIDs = append(a.SubscribedChatIDs, ChatID)
			buf, err = a.ToJSON()
			if err != nil {
				return err
			}
			return b.Put([]byte(AppID), buf)
		})
	})
	return err
}

func (f file) Unsubscribe(ChatID int64, AppID string) error {
	var c models.Chat
	err := f.checkChatID(ChatID)
	if err != nil {
		f.createChat(ChatID)
	}
	err = f.checkAppID(AppID)
	if err != nil {
		return errors.New(AppID + ": not found")
	}
	err = f.db(func(db *bolt.DB) error {
		return db.Batch(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Chats"))
			v := b.Get(utils.I64tob(ChatID))
			if v == nil {
				return errors.New("not found")
			}
			err := c.FromJSON(v)
			if err != nil {
				return err
			}
			c.SubscribedAppIDs = utils.ReuseString(c.SubscribedAppIDs, AppID)
			buf, err := c.ToJSON()
			if err != nil {
				return err
			}
			err = b.Put(utils.I64tob(ChatID), []byte(buf))
			if err != nil {
				return err
			}
			b = tx.Bucket([]byte("Apps"))
			v = b.Get([]byte(AppID))
			var a models.App
			err = a.FromJSON(v)
			if err != nil {
				return err
			}
			a.SubscribedChatIDs = utils.ReuseInt64(a.SubscribedChatIDs, ChatID)
			buf, err = a.ToJSON()
			if err != nil {
				return err
			}
			return b.Put([]byte(AppID), buf)
		})
	})
	return err
}

func (f file) UnsubscribeAll(ChatID int64) error {
	var c models.Chat
	err := f.checkChatID(ChatID)
	if err != nil {
		f.createChat(ChatID)
	}
	err = f.db(func(db *bolt.DB) error {
		return db.Batch(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Chats"))
			v := b.Get(utils.I64tob(ChatID))
			if v == nil {
				return errors.New("not found")
			}
			err := c.FromJSON(v)
			if err != nil {
				return err
			}
			t := c
			t.SubscribedAppIDs = []string{}
			buf, err := t.ToJSON()
			if err != nil {
				return err
			}
			err = b.Put(utils.I64tob(ChatID), []byte(buf))
			if err != nil {
				return err
			}
			for _, AppID := range c.SubscribedAppIDs {
				b := tx.Bucket([]byte("Apps"))
				v := b.Get([]byte(AppID))
				var a models.App
				err := a.FromJSON(v)
				if err != nil {
					return err
				}
				a.SubscribedChatIDs = utils.ReuseInt64(a.SubscribedChatIDs, ChatID)
				buf, err := a.ToJSON()
				if err != nil {
					return err
				}
				err = b.Put([]byte(AppID), buf)
				if err != nil {
					return err
				}
			}
			return nil
		})
	})
	return err
}

func (f file) Close() error {
	return nil
}

func (f file) db(dbFunc dbFunc) error {
	db, err := bolt.Open(f.path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()
	return dbFunc(db)
}

func newFile(options map[string]interface{}) (storage.Store, error) {
	// init options
	path := "./file.db"
	for k, v := range options {
		switch k {
		case "path":
			if s, ok := v.(string); ok {
				path = s
			}
		}
	}
	//init file struct
	file := &file{path: path}

	// create boltdb buckets
	for _, v := range bucket {
		err := file.db(func(db *bolt.DB) error {
			return db.Update(func(tx *bolt.Tx) error {
				_, err := tx.CreateBucketIfNotExists([]byte(v))
				if err != nil {
					return err
				}
				return nil
			})
		})
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

func init() {
	storage.Register("file", newFile)
}
