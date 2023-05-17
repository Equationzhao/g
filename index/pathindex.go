package index

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/xrash/smetrics"
)

var (
	db        *leveldb.DB
	initOnce  sync.Once
	closeOnce sync.Once
	errInit   error
	errClose  error
	indexPath string
)

func getDB() (*leveldb.DB, error) {
	initOnce.Do(func() {
		errInit = nil
		var err error
		indexPath, err = os.UserConfigDir()
		if err != nil {
			errInit = err
			initOnce = sync.Once{}
			return
		}
		indexPath = filepath.Join(indexPath, "g", "index")
		err = os.MkdirAll(indexPath, os.ModePerm)
		if err != nil {
			errInit = err
			initOnce = sync.Once{}
			return
		}
		db, err = leveldb.OpenFile(indexPath, nil)
		if err != nil {
			errInit = err
			initOnce = sync.Once{}
			return
		}
	})
	return db, errInit
}

func close() error {
	closeOnce.Do(func() {
		errClose = nil
		if db != nil {
			err := db.Close()
			if err != nil {
				errClose = err
				closeOnce = sync.Once{}
				return
			}
		}
	})
	return errClose
}

type ErrUpdate struct {
	key string
}

func (e ErrUpdate) Error() string {
	return fmt.Sprint("failed to update `", string(e.key), "`")
}

func FuzzySearch(key string) (string, error) {
	db, err := getDB()
	if err != nil {
		return "", err
	}
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	result := key
	times := 0
	highest := 0.0
	for iter.Next() {
		base := filepath.Base(string(iter.Key()))
		score := smetrics.JaroWinkler(key, base, 0.7, 4)
		times, err = strconv.Atoi(string(iter.Value()))
		if err != nil {
			continue
		}
		score *= math.Sqrt(float64(times))
		if smetrics.Soundex(key) == smetrics.Soundex(base) {
			score *= 1.5
		}
		if score > highest {
			highest = score
			result = string(iter.Key())
		}
	}

	return result, nil
}

func Update(key string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	times := 0
	data, err := db.Get([]byte(key), nil)
	if err != nil {
		err := db.Put([]byte(key), []byte("1"), nil)
		if err != nil {
			return err
		}
	} else {
		times, err = strconv.Atoi(string(data))
		if err != nil {
			return err
		}
	}
	return db.Put([]byte(key), []byte(strconv.Itoa(times+1)), nil)
}

func Delete(key string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	return db.Delete([]byte(key), nil)
}

func RebuildIndex() error {
	err := db.Close()
	if err != nil {
		return err
	}
	err = os.RemoveAll(indexPath)
	if err != nil {
		return err
	}
	return nil
}

func All() ([]string, []string, error) {
	db, err := getDB()
	if err != nil {
		return nil, nil, err
	}
	keys, values := make([]string, 0), make([]string, 0)
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		keys = append(keys, string(iter.Key()))
		values = append(values, string(iter.Value()))
	}
	return keys, values, nil
}
