package tiktokdb

import (
	"errors"

	tiktokdmschema "github.com/bjornpagen/tiktok-video-processor/autogen/tiktokdb"
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

func Create(path string) error {
	logger := zerolog.Nop()
	mode := os.FileMode(0644)
	numReaders := uint(8)
	numDBs := uint(1)

	if _, err := os.Stat(path); err == nil {
		return errors.New("directory already exists")
	}

	err := os.Mkdir(path, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	client, err := golmdb.NewLMDB(logger, path, mode, numReaders, numDBs, golmdb.EnvironmentFlag(0x40000), 1)
	if err != nil {
		return err
	}

	db := &TikTokDB{client: client}
	defer db.Close()

	err = db.client.Update(func(txn *golmdb.ReadWriteTxn) error {
		db.wg.Add(1)
		defer db.wg.Done()

		_, err := txn.DBRef(usersDb, golmdb.DatabaseFlag(0x40000))
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *TikTokDB) Close() {
	db.wg.Wait()
	db.client.TerminateSync()
}

func (db *TikTokDB) GetUser(userID string) (*tiktokdmschema.TikTokUser, error) {
	var user *tiktokdmschema.TikTokUser

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

		user = new(tiktokdmschema.TikTokUser)
		user.Init(value, flatbuffers.GetUOffsetT(value))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *TikTokDB) SetUser(userID string, user *tiktokdmschema.TikTokUser) error {
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
