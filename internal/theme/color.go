package theme

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Equationzhao/g/internal/const"

	colortool "github.com/gookit/color"
)

var ColorLevel = colortool.DetectColorLevel()

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
	case constval.Black:
		r, g, b = 0, 0, 0
	case constval.Red:
		r, g, b = 205, 0, 0
	case constval.Green:
		r, g, b = 0, 205, 0
	case constval.Yellow:
		r, g, b = 205, 205, 0
	case constval.Blue:
		r, g, b = 0, 0, 238
	case constval.Purple:
		r, g, b = 205, 0, 205
	case constval.Cyan:
		r, g, b = 0, 205, 205
	case constval.White:
		r, g, b = 229, 229, 229
	case constval.BrightBlack:
		r, g, b = 127, 127, 127
	case constval.BrightRed:
		r, g, b = 255, 0, 0
	case constval.BrightGreen:
		r, g, b = 0, 255, 0
	case constval.BrightYellow:
		r, g, b = 255, 255, 0
	case constval.BrightBlue:
		r, g, b = 92, 92, 255
	case constval.BrightPurple:
		r, g, b = 255, 0, 255
	case constval.BrightCyan:
		r, g, b = 0, 255, 255
	case constval.BrightWhite:
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
	return c256ToRGB(v)
}

func c256ToRGB(v uint8) (string, error) {
	c256ToRgb := colortool.C256ToRgb(v)
	r, g, b := c256ToRgb[0], c256ToRgb[1], c256ToRgb[2]
	return rgb(r, g, b), nil
}

// Color256ToBasic convert 256 color to basic color
func Color256ToBasic(str256 string) (string, error) {
	var v uint8
	_, err := fmt.Fscanf(strings.NewReader(str256), Color256Format, &v)
	if err != nil {
		return "", err
	}
	return c256ToBasic(v)
}

func c256ToBasic(v uint8) (string, error) {
	c256ToRgb := colortool.C256ToRgb(v)
	r, g, b := c256ToRgb[0], c256ToRgb[1], c256ToRgb[2]
	return RGBToBasicInt(r, g, b), nil
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

// RGBMultiply will multiply each r,g,b value by radio
// if radio < 0, return the original string
// if radio = 0, return the color of black
// if 0 < radio < 1, the color will be darker
// if radio > 1, the color will be lighter, and if the result is greater than 255, it will be 255
func RGBMultiply(rgbStr string, radio float64) string {
	// parse rgb
	var r, g, b uint8 = 0, 0, 0
	_, _ = fmt.Fscanf(strings.NewReader(rgbStr), RGBFormat, &r, &g, &b)

	rf := math.Ceil(float64(r) * radio)
	gf := math.Ceil(float64(g) * radio)
	bf := math.Ceil(float64(b) * radio)
	if rf > 255 {
		rf = 255
	}
	if gf > 255 {
		gf = 255
	}
	if bf > 255 {
		bf = 255
	}
	return rgb(uint8(rf), uint8(gf), uint8(bf))
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

func ConvertColor(to colortool.Level, src string) (string, error) {
	if to == None {
		return "", nil
	}

	switch src {
	case "":
		return "", nil
	case
		// 1.basic
		constval.Reset,
		constval.Red,
		constval.Green,
		constval.Yellow,
		constval.Blue,
		constval.Purple,
		constval.Cyan,
		constval.White,
		constval.BrightBlack,
		constval.BrightRed,
		constval.BrightGreen,
		constval.BrightYellow,
		constval.BrightBlue,
		constval.BrightPurple,
		constval.BrightCyan,
		constval.BrightWhite:
		switch to {
		case Ascii:
			return src, nil
		case C256:
			return BasicTo256(src), nil
		case TrueColor:
			return BasicToRGB(src), nil
		default:
			return "", ErrUnknownColorType{Level: to}
		}
	case constval.Underline:
		return src, nil
	default:

		strReader := strings.NewReader(src)

		// 2.8bit/256 color
		var c uint8
		_, err := fmt.Fscanf(strReader, Color256Format, &c)
		if err == nil {
			switch to {
			case Ascii:
				return c256ToBasic(c)
			case C256:
				return src, nil
			case TrueColor:
				return c256ToRGB(c)
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

func ConvertColorIfGreaterThanExpect(to colortool.Level, src string) (string, error) {
	if to == None {
		return "", nil
	}

	switch src {
	case "":
		return "", nil
	case
		// 1.basic
		constval.Reset,
		constval.Black,
		constval.Red,
		constval.Green,
		constval.Yellow,
		constval.Blue,
		constval.Purple,
		constval.Cyan,
		constval.White,
		constval.BrightBlack,
		constval.BrightRed,
		constval.BrightGreen,
		constval.BrightYellow,
		constval.BrightBlue,
		constval.BrightPurple,
		constval.BrightCyan,
		constval.BrightWhite:
		return src, nil
	case constval.Underline:
		return src, nil
	default:

		strReader := strings.NewReader(src)

		// 2.8bit/256 color
		var c uint8
		_, err := fmt.Fscanf(strReader, Color256Format, &c)
		if err == nil {
			switch to {
			case Ascii:
				return c256ToBasic(c)
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
