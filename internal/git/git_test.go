package git

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestParseShort(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		wantRes FileGits
	}{
		{
			name: "case 1",
			args: "AM internal/git/git_test.go\n!! .DS_Store\n!! .idea/\n!! build/\n",
			wantRes: FileGits{
				{Name: "internal/git/git_test.go", X: Added, Y: Modified},
				{Name: ".DS_Store", X: Ignored, Y: Ignored},
				{Name: ".idea", X: Ignored, Y: Ignored},
				{Name: "build", X: Ignored, Y: Ignored},
			},
		},
		{
			name: "case 2",
			args: "D  my_folder/my_file.txt\n",
			wantRes: FileGits{
				{Name: "my_folder/my_file.txt", X: Deleted, Y: Unmodified},
			},
		},
		{
			name: "case 3",
			args: " D my_folder/my_file.txt\n",
			wantRes: FileGits{
				{Name: "my_folder/my_file.txt", X: Unmodified, Y: Deleted},
			},
		},
		{
			name: "case 4",
			args: "T  my_folder/my_file.txt\n",
			wantRes: FileGits{
				{Name: "my_folder/my_file.txt", X: TypeChanged, Y: Unmodified},
			},
		},
		{
			name: "case 5",
			args: "M  my_folder/file1.txt\nD  my_folder/file2.txt\nAU my_folder/file3.txt\n",
			wantRes: FileGits{
				{Name: "my_folder/file1.txt", X: Modified, Y: Unmodified},
				{Name: "my_folder/file2.txt", X: Deleted, Y: Unmodified},
				{Name: "my_folder/file3.txt", X: Added, Y: UpdatedButUnmerged},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes := ParseShort(tt.args)
			for i := range gotRes {
				gotRes[i].Name = normalizePath(gotRes[i].Name)
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("ParseShort() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func normalizePath(path string) string {
	// normalize path according to the OS
	switch os := runtime.GOOS; os {
	case "windows":
		return strings.ReplaceAll(path, "/", "\\")
	default:
		return strings.ReplaceAll(path, "\\", "/")
	}
}
