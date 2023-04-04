package tiktokdb

import (
	"bytes"
	"encoding/gob"
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

const (
	numDBs = uint(2)
)

type TikTokDB struct {
	Lmdb *lmdb.LMDBClient
	wg   sync.WaitGroup
	path string
}

func New(path string) *TikTokDB {
	db := &TikTokDB{
		Lmdb: nil,
		wg:   sync.WaitGroup{},
		path: path,
	}
	return db
}

func (db *TikTokDB) Open() error {
	logger := zerolog.Nop()
	mode := os.FileMode(0644)
	numReaders := uint(8)

	// check if directory exists, if not create it
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		err = os.MkdirAll(db.path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	client, err := lmdb.NewLMDB(logger, db.path, mode, numReaders, numDBs, lmdb.EnvironmentFlag(0), 1)
	if err != nil {
		return err
	}
	db.Lmdb = client

	return nil
}

func (db *TikTokDB) Close() {
	db.wg.Wait()
	db.Lmdb.TerminateSync()
}

func (db *TikTokDB) GetUser(userID string) (*scraperapi.User, error) {
	var user *scraperapi.User

	err := db.Lmdb.View(func(txn *lmdb.ReadOnlyTxn) error {
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
	return db.Lmdb.Update(func(txn *lmdb.ReadWriteTxn) error {
		db.wg.Add(1)
		defer db.wg.Done()

		dbRef, err := txn.DBRef(usersDb, lmdb.DatabaseFlag(0x40000))
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

func (db *TikTokDB) GetAwemeList(userID string) ([]scraperapi.Aweme, error) {
	var awemeList []scraperapi.Aweme

	err := db.Lmdb.View(func(txn *lmdb.ReadOnlyTxn) error {
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
		return nil, err
	}

	return awemeList, nil
}

func (db *TikTokDB) SetAwemeList(userID string, awemeList []scraperapi.Aweme) error {
	return db.Lmdb.Update(func(txn *lmdb.ReadWriteTxn) error {
		db.wg.Add(1)
		defer db.wg.Done()

		dbRef, err := txn.DBRef(awemeDb, lmdb.DatabaseFlag(0x40000))
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
