package util

import "testing"

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
