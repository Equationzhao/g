package util

import (
	"os"
	"testing"
)

func TestSplitNumberAndUnit(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 string
	}{
		{
			name: "123bit",
			args: args{
				s: "123bit",
			},
			want:  123,
			want1: "bit",
		},
		{
			name: "123",
			args: args{
				s: "123",
			},
			want:  123,
			want1: "",
		},
		{
			name: "1,234.321bit",
			args: args{
				s: "1,234.321bit",
			},
			want:  1234.321,
			want1: "bit",
		},
		{
			name: "-1,234.321bit",
			args: args{
				s: "-1,234.321bit",
			},
			want:  -1234.321,
			want1: "bit",
		},
		{
			name: "bit",
			args: args{
				s: "bit",
			},
			want:  0,
			want1: "bit",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SplitNumberAndUnit(tt.args.s)
			if got != tt.want {
				t.Errorf("SplitNumberAndUnit() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SplitNumberAndUnit() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMakeLink(t *testing.T) {
	link := MakeLink("abs", "name")
	if link != "\033]8;;abs\033\\name\033]8;;\033\\" {
		t.Errorf("MakeLink failed")
	}
}

func TestRemoveSep(t *testing.T) {
	sep := RemoveSep("a/b/c")
	if sep != "a/b/c" {
		t.Errorf("RemoveSep failed")
	}
	sep = RemoveSep("a/b/c/")
	if sep != "a/b/c" {
		t.Errorf("RemoveSep failed")
	}
}

func TestEscape(t *testing.T) {
	type args struct {
		a string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "tab",
			args: args{
				a: "\t",
			},
			want: "\x1b[7m\\t\x1b[27m",
		},
		{
			name: "carriage return",
			args: args{
				a: "\r",
			},
			want: "\x1b[7m\\r\x1b[27m",
		},
		{
			name: "line feed",
			args: args{
				a: "\n",
			},
			want: "\x1b[7m\\n\x1b[27m",
		},
		{
			name: "double quote",
			args: args{
				a: "\"",
			},
			want: "\x1b[7m\\\"\x1b[27m",
		},
		{
			name: "backslash",
			args: args{
				a: "\\",
			},
			want: "\x1b[7m\\\\\x1b[27m",
		},
		{
			name: "single quote",
			args: args{
				a: "'",
			},
			want: "'",
		},
		{
			name: "normal",
			args: args{
				a: "normal",
			},
			want: "normal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Escape(tt.args.a); got != tt.want {
				t.Errorf("Escape() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLocale(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		want     string
		saveEnvs map[string]string
	}{
		{
			name: "No env vars set",
			envVars: map[string]string{
				"LC_ALL":      "",
				"LC_MESSAGES": "",
				"LANG":        "",
			},
			want: "en_US",
		},
		{
			name: "Only LANG set",
			envVars: map[string]string{
				"LC_ALL":      "",
				"LC_MESSAGES": "",
				"LANG":        "fr_FR.UTF-8",
			},
			want: "fr_FR",
		},
		{
			name: "LC_MESSAGES takes precedence over LANG",
			envVars: map[string]string{
				"LC_ALL":      "",
				"LC_MESSAGES": "de_DE.UTF-8",
				"LANG":        "fr_FR.UTF-8",
			},
			want: "de_DE",
		},
		{
			name: "LC_ALL takes precedence over all",
			envVars: map[string]string{
				"LC_ALL":      "ja_JP.UTF-8",
				"LC_MESSAGES": "de_DE.UTF-8",
				"LANG":        "fr_FR.UTF-8",
			},
			want: "ja_JP",
		},
		{
			name: "No UTF-8 suffix",
			envVars: map[string]string{
				"LC_ALL": "ko_KR",
			},
			want: "ko_KR",
		},
	}

	// Save original env vars to restore later
	saveEnvs := make(map[string]string)
	for _, envVar := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		saveEnvs[envVar] = os.Getenv(envVar)
	}

	// Restore env vars after test completes
	defer func() {
		for k, v := range saveEnvs {
			os.Setenv(k, v)
		}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables for this test
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			if got := GetLocale(); got != tt.want {
				t.Errorf("GetLocale() = %v, want %v", got, tt.want)
			}
		})
	}
}
