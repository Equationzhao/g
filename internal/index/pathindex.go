package index

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Equationzhao/g/internal/config"
	gutil "github.com/Equationzhao/g/internal/util"
	"github.com/junegunn/fzf/src/algo"
	"github.com/junegunn/fzf/src/util"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	db        *leveldb.DB
	initOnce  gutil.Once
	closeOnce gutil.Once
	indexPath string
)

func SetReadOnly() {
	db, err := getDB()
	if err != nil {
		return
	}
	_ = db.SetReadOnly()
	return
}

func getDB() (*leveldb.DB, error) {
	err := initOnce.Do(func() error {
		var err error
		indexPath, err = config.GetUserConfigDir()
		if err != nil {
			return err
		}
		indexPath = filepath.Join(indexPath, "index")
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

func All() ([]string, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0)
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		keys = append(keys, string(iter.Key()))
	}
	return keys, nil
}

func DeleteThose(keys ...string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	t, err := db.OpenTransaction()
	defer t.Discard()
	if err != nil {
		return nil
	}
	var errSum error
	for _, key := range keys {
		err := t.Delete([]byte(key), nil)
		if err != nil {
			errSum = errors.Join(errSum, err)
		}
	}
	errSum = errors.Join(errSum, t.Commit())
	return errSum
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
