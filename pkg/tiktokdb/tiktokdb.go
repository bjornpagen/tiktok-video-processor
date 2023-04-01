package tiktokdb

import (
	"errors"

	schema "github.com/bjornpagen/tiktok-video-processor/pkg/tiktokdb/schemas/tiktokdb"
	flatbuffers "github.com/google/flatbuffers/go"

	"os"

	"sync"

	"github.com/rs/zerolog"
	"wellquite.org/golmdb"
)

const (
	usersDb = "users"
)

type TikTokDB struct {
	client *golmdb.LMDBClient
	wg     sync.WaitGroup
}

func Open(path string) (*TikTokDB, error) {
	logger := zerolog.Nop()
	mode := os.FileMode(0644)
	numReaders := uint(8)
	numDBs := uint(1)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("database does not exist")
	}

	client, err := golmdb.NewLMDB(logger, path, mode, numReaders, numDBs, golmdb.EnvironmentFlag(0), 1)
	if err != nil {
		return nil, err
	}

	return &TikTokDB{client: client}, nil
}

func New(path string) (*TikTokDB, error) {
	logger := zerolog.Nop()
	mode := os.FileMode(0644)
	numReaders := uint(8)
	numDBs := uint(1)

	client, err := golmdb.NewLMDB(logger, path, mode, numReaders, numDBs, golmdb.EnvironmentFlag(0x40000), 1)
	if err != nil {
		return nil, err
	}
	defer client.TerminateSync()

	db := &TikTokDB{client: client}
	var wg sync.WaitGroup
	wg.Add(1)

	err = db.client.Update(func(txn *golmdb.ReadWriteTxn) error {
		_, err := txn.DBRef(usersDb, golmdb.DatabaseFlag(0x40000))
		defer wg.Done()
		return err
	})
	if err != nil {
		return nil, err
	}

	wg.Wait()

	return db, nil
}

func (db *TikTokDB) Close() {
	db.wg.Wait()
	db.client.TerminateSync()
}

func (db *TikTokDB) GetUser(userID string) (*schema.TikTokUser, error) {
	var user *schema.TikTokUser

	err := db.client.View(func(txn *golmdb.ReadOnlyTxn) error {
		db.wg.Add(1)
		defer db.wg.Done()

		dbRef, err := txn.DBRef(usersDb, golmdb.DatabaseFlag(0))
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
		db.wg.Add(1)
		defer db.wg.Done()

		dbRef, err := txn.DBRef(usersDb, golmdb.DatabaseFlag(0))
		if err != nil {
			return err
		}

		key := []byte(userID)
		value := user.Table().Bytes

		return txn.Put(dbRef, key, value, golmdb.PutFlag(0))
	})
}
