package file

import (
	"errors"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/ne2blink/antenna/pkg/storage"
	"github.com/ne2blink/antenna/pkg/storage/file/models"
	"github.com/ne2blink/antenna/pkg/utils"
)

type store struct {
	db *bolt.DB
}

func (s store) CreateApp(app storage.App) (string, error) {
	var mApp models.App
	mApp.FromStoreApp(app)
	err := s.db.Update(func(tx *bolt.Tx) error {
		// Open Apps Bucket
		bucket := tx.Bucket([]byte("Apps"))
		// ID Auto Increment
		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}
		// Change type to string
		mApp.ID = strconv.FormatUint(id, 10)
		// Marshal app data into bytes.
		buf, err := mApp.ToJSON()
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return bucket.Put([]byte(mApp.ID), buf)
	})
	return mApp.ID, err
}

func (s store) UpdateApp(app storage.App) error {
	var mApp models.App
	mApp.FromStoreApp(app)
	err := s.db.Update(func(tx *bolt.Tx) error {
		// Open Apps Bucket
		bucket := tx.Bucket([]byte("Apps"))
		// Marshal app data into bytes.
		buf, err := mApp.ToJSON()
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return bucket.Put([]byte(mApp.ID), buf)
	})
	return err
}

func (s store) GetApp(id string) (storage.App, error) {
	var mApp models.App
	err := s.db.View(func(tx *bolt.Tx) error {
		// Open Apps Bucket
		bucket := tx.Bucket([]byte("Apps"))
		value := bucket.Get([]byte(id))
		if value == nil {
			return errors.New(id + ": not found")
		}
		return mApp.FromJSON(value)
	})
	return mApp.ToStoreApp(), err
}

func (s store) DeleteApp(id string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		// Open Apps Bucket
		bucket := tx.Bucket([]byte("Apps"))
		return bucket.Delete([]byte(id))
	})
	return err
}

func (s store) ListApps() ([]storage.App, error) {
	var apps []storage.App
	err := s.db.View(func(tx *bolt.Tx) error {
		// Open Apps Bucket
		bucket := tx.Bucket([]byte("Apps"))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var mApp models.App
			err := mApp.FromJSON(v)
			if err != nil {
				return err
			}
			apps = append(apps, mApp.ToStoreApp())
		}
		return nil
	})
	return apps, err
}

func (s store) ListSubscribers(id string) ([]int64, error) {
	var mApp models.App
	err := s.db.View(func(tx *bolt.Tx) error {
		// Open Apps Bucket
		bucket := tx.Bucket([]byte("Apps"))
		value := bucket.Get([]byte(id))
		if value == nil {
			return errors.New(id + ": not found")
		}
		return mApp.FromJSON(value)
	})
	return mApp.SubscribedChatIDs, err
}

func (s store) ListSubscribedApps(chatID int64) ([]storage.App, error) {
	err := s.checkChatID(chatID)
	if err != nil {
		s.createChat(chatID)
	}
	var apps []storage.App
	err = s.db.View(func(tx *bolt.Tx) error {
		var chat models.Chat
		// Open Apps Bucket
		bucket := tx.Bucket([]byte("Chats"))
		value := bucket.Get(utils.Int64ToBytes(chatID))
		if value == nil {
			return errors.New("not found")
		}
		err := chat.FromJSON(value)
		if err != nil {
			return err
		}
		for _, ID := range chat.SubscribedAppIDs {
			bucket := tx.Bucket([]byte("Apps"))
			value := bucket.Get([]byte(ID))
			var mApp models.App
			err := mApp.FromJSON(value)
			if err != nil {
				return err
			}
			apps = append(apps, mApp.ToStoreApp())
		}
		return nil
	})
	return apps, err
}

func (s store) Subscribe(chatID int64, appID string) error {
	var chat models.Chat
	err := s.checkChatID(chatID)
	if err != nil {
		s.createChat(chatID)
	}
	err = s.checkAppID(appID)
	if err != nil {
		return errors.New(appID + ": not found")
	}
	err = s.db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Chats"))
		value := bucket.Get(utils.Int64ToBytes(chatID))
		if value == nil {
			return errors.New("not found")
		}
		err := chat.FromJSON(value)
		if err != nil {
			return err
		}
		if utils.ChickInString(chat.SubscribedAppIDs, appID) {
			return errors.New("already subscribed this app")
		}
		chat.SubscribedAppIDs = append(chat.SubscribedAppIDs, appID)
		buf, err := chat.ToJSON()
		if err != nil {
			return err
		}
		err = bucket.Put(utils.Int64ToBytes(chatID), []byte(buf))
		if err != nil {
			return err
		}
		bucket = tx.Bucket([]byte("Apps"))
		value = bucket.Get([]byte(appID))
		var mApp models.App
		err = mApp.FromJSON(value)
		if err != nil {
			return err
		}
		mApp.SubscribedChatIDs = append(mApp.SubscribedChatIDs, chatID)
		buf, err = mApp.ToJSON()
		if err != nil {
			return err
		}
		return bucket.Put([]byte(appID), buf)
	})
	return err
}

func (s store) Unsubscribe(chatID int64, appID string) error {
	var chat models.Chat
	err := s.checkChatID(chatID)
	if err != nil {
		s.createChat(chatID)
	}
	err = s.checkAppID(appID)
	if err != nil {
		return errors.New(appID + ": not found")
	}
	err = s.db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Chats"))
		value := bucket.Get(utils.Int64ToBytes(chatID))
		if value == nil {
			return errors.New("not found")
		}
		err := chat.FromJSON(value)
		if err != nil {
			return err
		}
		chat.SubscribedAppIDs = utils.ReuseString(chat.SubscribedAppIDs, appID)
		buf, err := chat.ToJSON()
		if err != nil {
			return err
		}
		err = bucket.Put(utils.Int64ToBytes(chatID), []byte(buf))
		if err != nil {
			return err
		}
		bucket = tx.Bucket([]byte("Apps"))
		value = bucket.Get([]byte(appID))
		var mApp models.App
		err = mApp.FromJSON(value)
		if err != nil {
			return err
		}
		mApp.SubscribedChatIDs = utils.ReuseInt64(mApp.SubscribedChatIDs, chatID)
		buf, err = mApp.ToJSON()
		if err != nil {
			return err
		}
		return bucket.Put([]byte(appID), buf)
	})
	return err
}

func (s store) UnsubscribeAll(chatID int64) error {
	var chat models.Chat
	err := s.checkChatID(chatID)
	if err != nil {
		s.createChat(chatID)
	}
	err = s.db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Chats"))
		value := bucket.Get(utils.Int64ToBytes(chatID))
		if value == nil {
			return errors.New("not found")
		}
		err := chat.FromJSON(value)
		if err != nil {
			return err
		}
		tempChat := chat
		tempChat.SubscribedAppIDs = []string{}
		buf, err := tempChat.ToJSON()
		if err != nil {
			return err
		}
		err = bucket.Put(utils.Int64ToBytes(chatID), []byte(buf))
		if err != nil {
			return err
		}
		for _, AppID := range chat.SubscribedAppIDs {
			bucket := tx.Bucket([]byte("Apps"))
			value := bucket.Get([]byte(AppID))
			var mApp models.App
			err := mApp.FromJSON(value)
			if err != nil {
				return err
			}
			mApp.SubscribedChatIDs = utils.ReuseInt64(mApp.SubscribedChatIDs, chatID)
			buf, err := mApp.ToJSON()
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(AppID), buf)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (s store) Close() error {
	return s.db.Close()
}

func (s store) checkAppID(id string) error {
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Apps"))
		value := bucket.Get([]byte(id))
		if value == nil {
			return errors.New("not found")
		}
		return nil
	})
	return err
}

func (s store) checkChatID(chatID int64) error {
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Chats"))
		value := bucket.Get(utils.Int64ToBytes(chatID))
		if value == nil {
			return errors.New("not found")
		}
		return nil
	})
	return err
}

func (s store) createChat(chatID int64) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		chat := models.Chat{ID: chatID}
		buf, err := chat.ToJSON()
		if err != nil {
			return err
		}
		bucket := tx.Bucket([]byte("Chats"))
		return bucket.Put(utils.Int64ToBytes(chatID), []byte(buf))
	})
	return err
}
