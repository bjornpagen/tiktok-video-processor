package tiktokdb

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"sync"

	"github.com/bjornpagen/tiktok-video-processor/pkg/tiktok/scraperapi"
	"github.com/rs/zerolog"
	lmdb "wellquite.org/golmdb"
)

const (
	usersDb = "users"
	awemeDb = "awemes"
)

type TikTokDB struct {
	client *lmdb.LMDBClient
	wg     sync.WaitGroup
	path   string
}

func New(path string) (*TikTokDB, error) {
	logger := zerolog.Nop()
	mode := os.FileMode(0644)
	numReaders := uint(8)
	numDBs := uint(1)

	if _, err := os.Stat(path); err == nil {
		return nil, errors.New("directory already exists")
	}

	err := os.Mkdir(path, 0755)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	client, err := lmdb.NewLMDB(logger, path, mode, numReaders, numDBs, lmdb.EnvironmentFlag(0x40000), 1)
	if err != nil {
		return nil, err
	}

	db := &TikTokDB{
		client: client,
		wg:     sync.WaitGroup{},
		path:   path,
	}
	defer db.Close()

	err = db.client.Update(func(txn *lmdb.ReadWriteTxn) error {
		db.wg.Add(1)
		defer db.wg.Done()

		_, err := txn.DBRef(usersDb, lmdb.DatabaseFlag(0x40000))
		if err != nil {
			return err
		}

		_, err = txn.DBRef(awemeDb, lmdb.DatabaseFlag(0x40000))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *TikTokDB) Open() error {
	logger := zerolog.Nop()
	mode := os.FileMode(0644)
	numReaders := uint(8)
	numDBs := uint(1)

	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		return errors.New("database does not exist")
	}

	client, err := lmdb.NewLMDB(logger, db.path, mode, numReaders, numDBs, lmdb.EnvironmentFlag(0), 1)
	if err != nil {
		return err
	}
	db.client = client

	return nil
}

func (db *TikTokDB) Close() {
	db.wg.Wait()
	db.client.TerminateSync()
}

func (db *TikTokDB) GetUser(userID string) (*scraperapi.User, error) {
	var user *scraperapi.User

	err := db.client.View(func(txn *lmdb.ReadOnlyTxn) error {
		db.wg.Add(1)
		defer db.wg.Done()

		dbRef, err := txn.DBRef(usersDb, lmdb.DatabaseFlag(0))
		if err != nil {
			return err
		}

		value, err := txn.Get(dbRef, []byte(userID))
		if err != nil {
			return err
		}

		decoder := gob.NewDecoder(bytes.NewReader(value))
		user = &scraperapi.User{}
		err = decoder.Decode(user)
		return err
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *TikTokDB) SetUser(userID string, user *scraperapi.User) error {
	return db.client.Update(func(txn *lmdb.ReadWriteTxn) error {
		db.wg.Add(1)
		defer db.wg.Done()

		dbRef, err := txn.DBRef(usersDb, lmdb.DatabaseFlag(0))
		if err != nil {
			return err
		}

		key := []byte(userID)

		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		err = encoder.Encode(user)
		if err != nil {
			return err
		}

		return txn.Put(dbRef, key, buf.Bytes(), lmdb.PutFlag(0))
	})
}

func (db *TikTokDB) GetAwemeList(userID string) []scraperapi.Aweme {
	var awemeList []scraperapi.Aweme

	err := db.client.View(func(txn *lmdb.ReadOnlyTxn) error {
		db.wg.Add(1)
		defer db.wg.Done()

		dbRef, err := txn.DBRef(awemeDb, lmdb.DatabaseFlag(0))
		if err != nil {
			return err
		}

		value, err := txn.Get(dbRef, []byte(userID))
		if err != nil {
			return err
		}

		decoder := gob.NewDecoder(bytes.NewReader(value))
		err = decoder.Decode(&awemeList)
		return err
	})

	if err != nil {
		return nil
	}

	return awemeList
}

func (db *TikTokDB) SetAwemeList(userID string, awemeList []scraperapi.Aweme) error {
	return db.client.Update(func(txn *lmdb.ReadWriteTxn) error {
		db.wg.Add(1)
		defer db.wg.Done()

		dbRef, err := txn.DBRef(awemeDb, lmdb.DatabaseFlag(0))
		if err != nil {
			return err
		}

		key := []byte(userID)

		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		err = encoder.Encode(awemeList)
		if err != nil {
			return err
		}

		return txn.Put(dbRef, key, buf.Bytes(), lmdb.PutFlag(0))
	})
}
