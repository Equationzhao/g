package util

import (
	"os"
	"testing"
	"time"
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
