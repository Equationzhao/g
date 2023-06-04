package index

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/junegunn/fzf/src/algo"
	"github.com/junegunn/fzf/src/util"
	"github.com/syndtr/goleveldb/leveldb"
)

type once struct {
	m    sync.Mutex
	done uint32
}

func (o *once) Do(fn func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		return nil
	}
	return o.doSlow(fn)
}

func (o *once) doSlow(fn func() error) error {
	o.m.Lock()
	defer o.m.Unlock()
	var err error
	if o.done == 0 {
		err = fn()
		if err == nil {
			atomic.StoreUint32(&o.done, 1)
		}
	}
	return err
}

var (
	db        *leveldb.DB
	initOnce  once
	closeOnce once
	indexPath string
)

func getDB() (*leveldb.DB, error) {
	err := initOnce.Do(func() error {
		var err error
		indexPath, err = os.UserConfigDir()
		if err != nil {
			return err
		}
		indexPath = filepath.Join(indexPath, "g", "index")
		err = os.MkdirAll(indexPath, os.ModePerm)
		if err != nil {
			return err
		}
		db, err = leveldb.OpenFile(indexPath, nil)
		if err != nil {
			return err
		}
		return nil
	})
	return db, err
}

func Close() error {
	return closeDB()
}

func closeDB() error {
	err := closeOnce.Do(func() error {
		if db != nil {
			err := db.Close()
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

type ErrUpdate struct {
	key string
}

func (e ErrUpdate) Error() string {
	return fmt.Sprint("failed to update `", e.key, "`")
}

func Update(key string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	_, err = db.Get([]byte(key), nil)
	if err != nil {
		err := db.Put([]byte(key), []byte("1"), nil)
		if err != nil {
			return err
		}
	} else {
		return nil
	}
	return nil
}

func Delete(key string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	return db.Delete([]byte(key), nil)
}

func RebuildIndex() error {
	db, err := getDB()
	if err != nil {
		goto remove
	}
	err = db.Close()
	if err != nil {
		return err
	}
remove:
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

func FuzzySearch(key string) (string, error) {
	db, err := getDB()
	if err != nil {
		return "", err
	}
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	result := key
	// times := 0
	highest := 0
	for iter.Next() {
		input := util.ToChars([]byte(strings.ToLower(string(iter.Key()))))
		pattern := algo.NormalizeRunes([]rune(strings.ToLower(key)))
		res, _ := algo.FuzzyMatchV2(false, true, true, &input, pattern, true, nil)
		score := res.Score
		// base := filepath.Base(string(iter.Key()))
		// score := smetrics.JaroWinkler(key, base, 0.7, 4)
		// times, err = strconv.Atoi(string(iter.Value()))
		// if err != nil {
		// 	continue
		// }
		// score *= math.Sqrt(float64(times))
		// if smetrics.Soundex(key) == smetrics.Soundex(base) {
		// 	score *= 1.5
		// }
		if score > highest {
			highest = score
			result = string(iter.Key())
		}
	}

	return result, nil
}
