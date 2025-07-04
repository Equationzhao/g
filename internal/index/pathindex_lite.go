//go:build !fuzzy

package index

import "errors"

// Lite version without fuzzy search dependencies
// This significantly reduces binary size by removing goleveldb and fuzzy dependencies

var ErrNotSupported = errors.New("fuzzy search feature not available in this build (built without 'fuzzy' tag)")

func getDB() (interface{}, error) {
	return nil, ErrNotSupported
}

func Close() error {
	return nil // No-op in lite build
}

func closeDB() error {
	return nil // No-op in lite build
}

type ErrUpdate struct {
	key string
}

func (e ErrUpdate) Error() string {
	return "index update not supported in lite build: " + e.key
}

func Update(key string) error {
	return nil // No-op in lite build - don't return error to avoid breaking basic functionality
}

func Delete(key string) error {
	return nil // No-op in lite build
}

func RebuildIndex() error {
	return nil // No-op in lite build
}

func All() ([]string, error) {
	return nil, ErrNotSupported
}

func DeleteThose(keys ...string) error {
	return nil // No-op in lite build
}

func FuzzySearch(key string) (string, error) {
	// In lite build, just return the original key without fuzzy matching
	return key, nil
}
