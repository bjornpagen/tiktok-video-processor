package tiktokdb

import (
	schema "github.com/bjornpagen/tiktok-video-processor/pkg/tiktokdb/schemas/tiktokdb"
	flatbuffers "github.com/google/flatbuffers/go"

	"wellquite.org/golmdb"
)

const (
	dbName = "tiktokdb"
)

type TikTokDB struct {
	client *golmdb.LMDBClient
}

func (db *TikTokDB) GetUser(userID string) (*schema.TikTokUser, error) {
	var user *schema.TikTokUser

	err := db.client.View(func(txn *golmdb.ReadOnlyTxn) error {
		dbRef, err := txn.DBRef(dbName, golmdb.DatabaseFlag(0))
		if err != nil {
			return err
		}

		value, err := txn.Get(dbRef, []byte(userID))
		if err != nil {
			return err
		}

		user = new(schema.TikTokUser)
		user.Init(value, flatbuffers.GetUOffsetT(value))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *TikTokDB) SetUser(userID string, user *schema.TikTokUser) error {
	return db.client.Update(func(txn *golmdb.ReadWriteTxn) error {
		dbRef, err := txn.DBRef(dbName, golmdb.DatabaseFlag(0))
		if err != nil {
			return err
		}

		key := []byte(userID)
		value := user.Table().Bytes

		return txn.Put(dbRef, key, value, golmdb.PutFlag(0))
	})
}
