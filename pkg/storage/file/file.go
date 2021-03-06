package file

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/ne2blink/antenna/pkg/storage"
)

var (
	bucket = []string{"Apps", "Chats"} //boltdb buckets
)

func newFileStore(options map[string]interface{}) (storage.Store, error) {
	// init options
	path, _ := options["path"].(string)
	if path == "" {
		path = "file.db"
	}
	// open database
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	//init file struct
	file := &store{db: db}
	// create boltdb buckets
	for _, v := range bucket {
		err := file.db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(v))
			return err
		})
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

func init() {
	storage.Register("file", newFileStore)
}
