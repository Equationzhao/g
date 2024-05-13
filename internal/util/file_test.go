package util

import (
	"os"
	"testing"
	"time"

	"github.com/Equationzhao/g/internal/item"
	"github.com/spf13/afero"
	"github.com/zeebo/assert"
)

func TestMockFileInfo(t *testing.T) {
	// Create a new instance of MockFileInfo
	modTime := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	m := NewMockFileInfo(100, true, "test", os.ModeDir, modTime)

	// Test the Name method
	if m.Name() != "test" {
		t.Errorf("expected 'test', got %s", m.Name())
	}

	// Test the Size method
	if m.Size() != 100 {
		t.Errorf("expected 100, got %d", m.Size())
	}

	// Test the Mode method
	if m.Mode() != os.ModeDir {
		t.Errorf("expected os.ModeDir, got %v", m.Mode())
	}

	// Test the ModTime method
	if !m.ModTime().Equal(modTime) {
		t.Errorf("expected %s time, got %s time", modTime, m.ModTime())
	}

	// Test the IsDir method
	if !m.IsDir() {
		t.Errorf("expected true, got false")
	}

	// Test the Sys method
	if m.Sys() != nil {
		t.Errorf("expected nil, got %v", m.Sys())
	}
}

func TestIsSymLink(t *testing.T) {
	// Create a new instance of MockFileInfo
	modTime := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	m := NewMockFileInfo(100, true, "test", os.ModeSymlink, modTime)

	// Test the IsSymLink method
	if !IsSymLink(m) {
		t.Errorf("expected true, got false")
	}

	m.mode = os.ModeDir
	if IsSymLink(m) {
		t.Errorf("expected false, got true")
	}
}

func TestIsSymLinkMode(t *testing.T) {
	// Test the IsSymLinkMode method
	if !IsSymLinkMode(os.ModeSymlink) {
		t.Errorf("expected true, got false")
	}

	if IsSymLinkMode(os.ModeDir) {
		t.Errorf("expected false, got true")
	}
}

func TestIsExecutable(t *testing.T) {
	m := NewMockFileInfo(100, true, "test", os.ModePerm, time.Now())

	// Test the IsExecutable method
	if !IsExecutable(m) {
		t.Errorf("expected true, got false")
	}

	m.mode = os.ModeDir
	if IsExecutable(m) {
		t.Errorf("expected false, got true")
	}
}

func TestIsExecutableMode(t *testing.T) {
	// Test the IsExecutableMode method
	if !IsExecutableMode(os.ModePerm) {
		t.Errorf("expected true, got false")
	}

	if IsExecutableMode(os.ModeDir) {
		t.Errorf("expected false, got true")
	}
}

func TestRecursivelySizeOf(t *testing.T) {
	// create a new fs using afero.NewMemMapFs
	afs := afero.NewMemMapFs()
	// create some files and dir
	_ = afs.Mkdir("dir", 0o755)
	_ = afero.WriteFile(afs, "dir/file1-100bytes", GenRandomData(100), 0o644)
	_ = afero.WriteFile(afs, "dir/file2-1000bytes", GenRandomData(1000), 0o644)
	_ = afero.WriteFile(afs, "dir/file3-10000bytes", GenRandomData(10000), 0o644)
	_ = afs.Mkdir("dir/subdir", 0o755)
	_ = afero.WriteFile(afs, "dir/subdir/file4-100000bytes", GenRandomData(100000), 0o644)
	// get the size of the dir
	err := afero.Walk(afs, "dir", func(path string, info os.FileInfo, err error) error {
		t.Logf("%s: %d", path, info.Size())
		return nil
	})
	assert.NoError(t, err)
	dir, _ := afs.Stat("dir")
	dirItem, err := item.NewFileInfoWithOption(item.WithFileInfo(dir), item.WithAbsPath("dir"))
	if err != nil {
		return
	}
	size := RecursivelySizeOfGenerator(afs)(dirItem, -1)
	if size != 111184 {
		t.Errorf("expected 111184, got %d", size)
	}
	size = RecursivelySizeOfGenerator(afs)(dirItem, 0)
	if size != 42 {
		t.Errorf("expected 42, got %d", size)
	}
	size = RecursivelySizeOfGenerator(afs)(dirItem, 1)
	if size != 11184 {
		t.Errorf("expected 11184, got %d", size)
	}
	size = RecursivelySizeOfGenerator(afs)(dirItem, 2)
	if size != 111184 {
		t.Errorf("expected 111184, got %d", size)
	}
}
