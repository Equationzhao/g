package theme

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Equationzhao/g/internal/global"

	colortool "github.com/gookit/color"
)

var ColorLevel = colortool.TermColorLevel()

const (
	BasicFormat    = "\033[0;%dm"
	Color256Format = "\033[38;5;%dm"
	RGBFormat      = "\033[38;2;%d;%d;%dm"
)

func Color256(color int) (string, error) {
	if color < 0 || color > 255 {
		return "", fmt.Errorf("color must between 0 and 255")
	}
	return color256(color), nil
}

func color256(color int) string {
	return fmt.Sprintf(Color256Format, color)
}

func rgb(r, g, b uint8) string {
	return fmt.Sprintf(RGBFormat, r, g, b)
}

func RGB(r, g, b uint8) (string, error) {
	return rgb(r, g, b), nil
}

// BasicTo256 convert basic color to 256 color
func BasicTo256(b string) string {
	// parse basic color
	var color int
	_, _ = fmt.Fscanf(strings.NewReader(b), BasicFormat, &color)
	return fmt.Sprintf(Color256Format, color-30)
}

func BasicToRGBInt(basic string) (r, g, b uint8) {
	switch basic {
	case global.Black:
		r, g, b = 0, 0, 0
	case global.Red:
		r, g, b = 205, 0, 0
	case global.Green:
		r, g, b = 0, 205, 0
	case global.Yellow:
		r, g, b = 205, 205, 0
	case global.Blue:
		r, g, b = 0, 0, 238
	case global.Purple:
		r, g, b = 205, 0, 205
	case global.Cyan:
		r, g, b = 0, 205, 205
	case global.White:
		r, g, b = 229, 229, 229
	case global.BrightBlack:
		r, g, b = 127, 127, 127
	case global.BrightRed:
		r, g, b = 255, 0, 0
	case global.BrightGreen:
		r, g, b = 0, 255, 0
	case global.BrightYellow:
		r, g, b = 255, 255, 0
	case global.BrightBlue:
		r, g, b = 92, 92, 255
	case global.BrightPurple:
		r, g, b = 255, 0, 255
	case global.BrightCyan:
		r, g, b = 0, 255, 255
	case global.BrightWhite:
		r, g, b = 255, 255, 255
	default:
		r, g, b = 0, 0, 0
	}
	return
}

// BasicToRGB convert basic color to rgb color
func BasicToRGB(basic string) string {
	r, g, b := BasicToRGBInt(basic)
	res := rgb(r, g, b)
	return res
}

// RGBTo256Int convert rgb color to 256 color
func RGBTo256Int(r, g, b uint8) int {
	return 16 + 36*int(math.Round(float64(r)/51)) + 6*int(math.Round(float64(g)/51)) + int(math.Round(float64(b)/51))
}

// RGBTo256 convert rgb color to 256 color
func RGBTo256(r, g, b uint8) string {
	return fmt.Sprintf(Color256Format, RGBTo256Int(r, g, b))
}

func RGBToBasicInt(r, g, b uint8) string {
	return fmt.Sprintf(BasicFormat, colortool.Rgb2ansi(r, g, b, false))
}

// RGBToBasic convert rgb color to basic color
func RGBToBasic(r, g, b uint8) string {
	return RGBToBasicInt(r, g, b)
}

// Color256ToRGB convert 256 color to RGB color
func Color256ToRGB(str256 string) (string, error) {
	// parse 256
	var v uint8
	_, err := fmt.Fscanf(strings.NewReader(str256), Color256Format, &v)
	if err != nil {
		return "", err
	}
	return c256ToRGB(v), nil
}

func c256ToRGB(v uint8) string {
	c256ToRgb := colortool.C256ToRgb(v)
	r, g, b := c256ToRgb[0], c256ToRgb[1], c256ToRgb[2]
	return rgb(r, g, b)
}

// Color256ToBasic convert 256 color to basic color
func Color256ToBasic(str256 string) (string, error) {
	var v uint8
	_, err := fmt.Fscanf(strings.NewReader(str256), Color256Format, &v)
	if err != nil {
		return "", err
	}
	return c256ToBasic(v), nil
}

func c256ToBasic(v uint8) string {
	c256ToRgb := colortool.C256ToRgb(v)
	r, g, b := c256ToRgb[0], c256ToRgb[1], c256ToRgb[2]
	return RGBToBasicInt(r, g, b)
}

// HexToRgb convert hex color to rgb color
// from gookit/color
// github.com/gookit/color
func HexToRgb(hex string) (rgb []uint8) {
	rgb = make([]uint8, 3)
	hex = strings.TrimSpace(hex)
	if hex == "" {
		return
	}

	// like from css. eg "#ccc" "#ad99c0"
	if hex[0] == '#' {
		hex = hex[1:]
	}

	hex = strings.ToLower(hex)
	switch len(hex) {
	case 3: // "ccc"
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	case 8: // "0xad99c0"
		hex = strings.TrimPrefix(hex, "0x")
	}

	// recheck
	if len(hex) != 6 {
		return
	}

	// convert string to int64
	if i64, err := strconv.ParseInt(hex, 16, 32); err == nil {
		color := int(i64)
		// parse int
		rgb = make([]uint8, 3)
		rgb[0] = uint8(color >> 16)
		rgb[1] = uint8((color & 0x00FF00) >> 8)
		rgb[2] = uint8(color & 0x0000FF)
	}
	return
}

// Convert HSL (Hue, Saturation, Lightness) to RGB (Red, Green, Blue)
func HslToRgb(h, s, l float64) (uint8, uint8, uint8) {
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2
	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	r = (r + m) * 255
	g = (g + m) * 255
	b = (b + m) * 255
	return uint8(math.Round(r)), uint8(math.Round(g)), uint8(math.Round(b))
}

const (
	None      = colortool.LevelNo
	Ascii     = colortool.Level16
	C256      = colortool.Level256
	TrueColor = colortool.LevelRgb
)

type ErrUnknownColorType struct {
	colortool.Level
}

func (e ErrUnknownColorType) Error() string {
	return fmt.Sprintf("unknown color type:%d", e.Level)
}

func ConvertColorIfGreaterThanExpect(to colortool.Level, src string) (string, error) {
	if to == None {
		return "", nil
	}

	switch src {
	case "":
		return "", nil
	case
		// 1.basic
		global.Reset,
		global.Black,
		global.Red,
		global.Green,
		global.Yellow,
		global.Blue,
		global.Purple,
		global.Cyan,
		global.White,
		global.BrightBlack,
		global.BrightRed,
		global.BrightGreen,
		global.BrightYellow,
		global.BrightBlue,
		global.BrightPurple,
		global.BrightCyan,
		global.BrightWhite,
		global.Underline:
		return src, nil
	default:

		strReader := strings.NewReader(src)

		// 2.8bit/256 color
		var c uint8
		_, err := fmt.Fscanf(strReader, Color256Format, &c)
		if err == nil {
			switch to {
			case Ascii:
				return c256ToBasic(c), nil
			case C256:
				return src, nil
			case TrueColor:
				return src, nil
			default:
				return "", ErrUnknownColorType{Level: to}
			}
		}

		// 3.rgb
		var (
			r uint8 = 0
			g uint8 = 0
			b uint8 = 0
		)
		strReader = strings.NewReader(src)

		_, err = fmt.Fscanf(strReader, RGBFormat, &r, &g, &b)
		if err == nil {
			switch to {
			case Ascii:
				return RGBToBasic(r, g, b), nil
			case C256:
				return RGBTo256(r, g, b), nil
			case TrueColor:
				return src, nil
			default:
				return "", ErrUnknownColorType{Level: to}
			}
		}
		return "", errors.New("unknown color type")
	}
}
