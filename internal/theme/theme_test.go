package theme

import (
	_ "embed"
	"io/fs"
	"os"
	"reflect"
	"testing"

	"github.com/Equationzhao/g/internal/global"
	"github.com/agiledragon/gomonkey/v2"

	colortool "github.com/gookit/color"
)

func TestAll(t *testing.T) {
	ColorLevel = colortool.Level16
	ConvertThemeColor()
	pl := func(m map[string]Style) {
		for key := range m {
			t.Logf("%s %s %s %s", m[key].Color, m[key].Icon, key, global.Reset)
		}
	}
	pl(DefaultAll.InfoTheme)
	pl(DefaultAll.Permission)
	pl(DefaultAll.Size)
	pl(DefaultAll.User)
	pl(DefaultAll.Group)
	pl(DefaultAll.Symlink)
	pl(DefaultAll.Git)
	pl(DefaultAll.Name)
	pl(DefaultAll.Special)
	pl(DefaultAll.Ext)
}

func TestColor(t *testing.T) {
	println(global.Green + "\uF48A " + global.Underline + global.Bold + "hello" + global.Red + " hello" + global.Reset)
	println(global.Green + "\uF48A " + global.Underline + "hello" + global.Red + " hello" + global.Reset)
}

func Test_genStyleField(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "Test_genStyleField",
			want: []string{
				"color",
				"icon",
				"underline",
				"bold",
				"faint",
				"italics",
				"blink",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genStyleField(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("genStyleField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStyle_FromReadable(t *testing.T) {
	type fields struct {
		Color     string
		Icon      string
		Underline bool
		Bold      bool
		Faint     bool
		Italics   bool
		Blink     bool
	}
	tests := []struct {
		name      string
		fields    fields
		wantErr   bool
		wantColor string
	}{
		{
			name: "Basic",
			fields: fields{
				Color: "BrightBlue",
			},
			wantErr:   false,
			wantColor: global.BrightBlue,
		},
		{
			name: "8bit",
			fields: fields{
				Color: "[128]@256",
			},
			wantErr:   false,
			wantColor: color256(128),
		},
		{
			name: "8bit without []",
			fields: fields{
				Color: "128@256",
			},
			wantErr:   false,
			wantColor: color256(128),
		},
		{
			name: "8bit error",
			fields: fields{
				Color: "256@256",
			},
			wantErr:   true,
			wantColor: "256@256",
		},
		{
			name: "8bit error",
			fields: fields{
				Color: "a@256",
			},
			wantErr:   true,
			wantColor: "a@256",
		},
		{
			name: "rgb",
			fields: fields{
				Color: "[128,128,128]@rgb",
			},
			wantErr:   false,
			wantColor: rgb(128, 128, 128),
		},
		{
			name: "rgb without []",
			fields: fields{
				Color: "128,128,128@rgb",
			},
			wantErr:   false,
			wantColor: rgb(128, 128, 128),
		},
		{
			name: "rgb error",
			fields: fields{
				Color: "128,128,abc@rgb",
			},
			wantErr:   true,
			wantColor: "128,128,abc@rgb",
		},
		{
			name: "rgb error",
			fields: fields{
				Color: "128,128@rgb",
			},
			wantErr:   true,
			wantColor: "128,128@rgb",
		},
		{
			name: "rgb error",
			fields: fields{
				Color: "128,128,256@rgb",
			},
			wantErr:   true,
			wantColor: "128,128,256@rgb",
		},
		{
			name: "hex",
			fields: fields{
				Color: "[#ff0000]@hex",
			},
			wantErr:   false,
			wantColor: rgb(255, 0, 0),
		},
		{
			name: "hex without []",
			fields: fields{
				Color: "#ff0000@hex",
			},
			wantErr:   false,
			wantColor: rgb(255, 0, 0),
		},
		{
			name: "hex error",
			fields: fields{
				Color: "#ff000@hex",
			},
			wantErr:   true,
			wantColor: "#ff000@hex",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Style{
				Color: tt.fields.Color,
			}
			if err := s.FromReadable(); (err != nil) != tt.wantErr {
				t.Errorf("FromReadable() error = %v, wantErr %v", err, tt.wantErr)
			}
			if s.Color != tt.wantColor {
				t.Errorf("FromReadable() got = %v, want %v", s.Color, tt.wantColor)
			}
		})
	}
}

func Test_color2str(t *testing.T) {
	type args struct {
		color string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "unknown",
			args: args{
				color: "unknown",
			},
			want: "",
		},
		{
			name: "8bit",
			args: args{
				color: color256(128),
			},
			want: "[128]@256",
		},
		{
			name: "rgb",
			args: args{
				color: rgb(128, 128, 128),
			},
			want: "[128,128,128]@rgb",
		},
		{
			name: "underline",
			args: args{
				color: global.Underline + global.Green,
			},
			want: color2str(global.Underline) + " + " + color2str(global.Green),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := color2str(tt.args.color); got != tt.want {
				t.Errorf("color2str() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_str2color(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "black",
			args: args{
				str: "black",
			},
			want: global.Black,
		},
		{
			name: "red",
			args: args{
				str: "red",
			},
			want: global.Red,
		},
		{
			name: "green",
			args: args{
				str: "green",
			},
			want: global.Green,
		},
		{
			name: "yellow",
			args: args{
				str: "yellow",
			},
			want: global.Yellow,
		},
		{
			name: "blue",
			args: args{
				str: "blue",
			},
			want: global.Blue,
		},
		{
			name: "purple",
			args: args{
				str: "purple",
			},
			want: global.Purple,
		},
		{
			name: "cyan",
			args: args{
				str: "cyan",
			},
			want: global.Cyan,
		},
		{
			name: "white",
			args: args{
				str: "white",
			},
			want: global.White,
		},
		{
			name: "reset",
			args: args{
				str: "reset",
			},
			want: global.Reset,
		},
		{
			name: "bright-red",
			args: args{
				str: "bright-red",
			},
			want: global.BrightRed,
		},
		{
			name: "bright-black",
			args: args{
				str: "bright-black",
			},
			want: global.BrightBlack,
		},
		{
			name: "bright-red",
			args: args{
				str: "bright-red",
			},
			want: global.BrightRed,
		},
		{
			name: "bright-green",
			args: args{
				str: "bright-green",
			},
			want: global.BrightGreen,
		},
		{
			name: "bright-yellow",
			args: args{
				str: "bright-yellow",
			},
			want: global.BrightYellow,
		},
		{
			name: "bright-blue",
			args: args{
				str: "bright-blue",
			},
			want: global.BrightBlue,
		},
		{
			name: "bright-purple",
			args: args{
				str: "bright-purple",
			},
			want: global.BrightPurple,
		},
		{
			name: "bright-cyan",
			args: args{
				str: "bright-cyan",
			},
			want: global.BrightCyan,
		},
		{
			name: "bright-white",
			args: args{
				str: "bright-white",
			},
			want: global.BrightWhite,
		},
		{
			name: "empty",
			args: args{
				str: "",
			},
			want: "",
		},
		{
			name: "8bit",
			args: args{
				str: "[128]@256",
			},
			want:    color256(128),
			wantErr: false,
		},
		{
			name: "rgb",
			args: args{
				str: "[128,128,128]@rgb",
			},
			want:    rgb(128, 128, 128),
			wantErr: false,
		},
		{
			name: "hex",
			args: args{
				str: "[#ff0000]@hex",
			},
			want:    rgb(255, 0, 0),
			wantErr: false,
		},
		{
			name: "hex",
			args: args{
				str: "[#0xff0000]@hex",
			},
			want:    rgb(255, 0, 0),
			wantErr: false,
		},
		{
			name: "hex",
			args: args{
				str: "[0xff0000]@hex",
			},
			want:    rgb(255, 0, 0),
			wantErr: false,
		},
		{
			name: "8bit error",
			args: args{
				str: "[256]@256",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "8bit error",
			args: args{
				str: "[a]@256",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "rgb error",
			args: args{
				str: "[128,128,256]@rgb",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "rgb error",
			args: args{
				str: "[128,128,abc]@rgb",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "rgb error",
			args: args{
				str: "[128,128]@rgb",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "hex error",
			args: args{
				str: "[#1a3gff]@hex",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := str2color(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("str2color() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("str2color() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidHexColor(t *testing.T) {
	type args struct {
		color string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid-8",
			args: args{
				color: "#ff0000",
			},
			want: true,
		},
		{
			name: "valid-8",
			args: args{
				color: "#1a3fFf",
			},
			want: true,
		},
		{
			name: "valid-3",
			args: args{
				color: "#f00",
			},
			want: true,
		},
		{
			name: "valid-3",
			args: args{
				color: "#abc",
			},
			want: true,
		},
		{
			name: "invalid",
			args: args{
				color: "#ff00",
			},
		},
		{
			name: "invalid",
			args: args{
				color: "#ff00000",
			},
		},
		{
			name: "invalid",
			args: args{
				color: "#1a3gff",
			},
		},
		{
			name: "0x",
			args: args{
				color: "0x1a3fFf",
			},
			want: true,
		},
		{
			name: "#0x",
			args: args{
				color: "#0x1a3fFf",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidHexColor(tt.args.color); got != tt.want {
				t.Errorf("IsValidHexColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

//go:embed default.json
var defaultTheme []byte

func TestGetTheme(t *testing.T) {
	patch := gomonkey.ApplyFunc(os.ReadFile, func(name string) ([]byte, error) {
		if name == "not-exist.json" {
			return nil, fs.ErrNotExist
		}
		if name == "fatal.json" {
			fatal := append(defaultTheme, []byte("invalid")...)
			return fatal, nil
		}
		return defaultTheme, nil
	})
	defer patch.Reset()
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "not exist",
			args: args{
				path: "not-exist.json",
			},
			wantErr: true,
		},
		{
			name: "fatal",
			args: args{
				path: "fatal.json",
			},
			wantErr: true,
		},
		{
			name: "default",
			args: args{
				path: "default.json",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GetTheme(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTheme() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				t.Log(err)
			}
		})
	}
}
