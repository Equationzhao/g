package cli

import (
	"sync"
	"testing"

	"github.com/Equationzhao/g/internal/filter"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/util"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/afero"
	"github.com/zeebo/assert"
)

func Test_dive(t *testing.T) {
	afs := afero.NewMemMapFs()
	// create files and dirs
	_ = afs.MkdirAll("a/b/c", 0755)
	_ = afero.WriteFile(afs, "a/b/c/d", util.GenRandomData(10), 0644)
	_ = afero.WriteFile(afs, "a/b/e", util.GenRandomData(10), 0644)
	_ = afero.WriteFile(afs, "a/f", util.GenRandomData(10), 0644)
	_ = afero.WriteFile(afs, "a/g", util.GenRandomData(10), 0644)
	_ = afero.WriteFile(afs, "a/h", util.GenRandomData(10), 0644)
	_ = afero.WriteFile(afs, "i", util.GenRandomData(10), 0644)

	var err error
	pool, err = ants.NewPool(ants.DefaultAntsPoolSize)
	assert.NoError(t, err)

	// test
	uslice := util.NewSlice[*item.FileInfo](10)
	errSlice := util.NewSlice[error](1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	filters := filter.NewItemFilter()
	dive(".", 1, -1, uslice, errSlice, &wg, filters, afs)
	wg.Wait()
	if uslice.Len() != 9 {
		t.Errorf("expect 7, got %d", uslice.Len())
		for _, info := range *uslice.GetRaw() {
			t.Logf("%s", info.FullPath)
		}
	}
	if errSlice.Len() != 0 {
		t.Errorf("expect 0, got %d", errSlice.Len())
	}
}
