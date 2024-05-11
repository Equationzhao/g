package theme

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/Equationzhao/g/internal/const"
)

type Style struct {
	// Color of the text.
	Color string `json:"color,omitempty"`
	// unicode Icon
	Icon      string `json:"icon,omitempty"`
	Underline bool   `json:"underline,omitempty"`
	Bold      bool   `json:"bold,omitempty"`
	Faint     bool   `json:"faint,omitempty"`
	Italics   bool   `json:"italics,omitempty"`
	Blink     bool   `json:"blink,omitempty"`
}

var (
	genOnceStyleField sync.Once
	styleField        []string
)

func genStyleField() []string {
	genOnceStyleField.Do(func() {
		// reflect on Style
		styleType := Style{}
		styleTypeValue := reflect.TypeOf(styleType)
		n := styleTypeValue.NumField()
		styleField = make([]string, 0, n)
		for i := 0; i < n; i++ {
			if j := styleTypeValue.Field(i).Tag.Get("json"); j != "" {
				jName, _, _ := strings.Cut(j, ",")
				if jName != "" {
					styleField = append(styleField, jName)
				}
			}
		}
	})
	return styleField
}

func (s *Style) UnmarshalJSON(bytes []byte) error {
	type styleAlias Style // avoid recursion during JSON unmarshalling.
	var alias styleAlias
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	// check if there is any unknown field
	fields := genStyleField()
	for key := range raw {
		if !slices.Contains(fields, key) {
			return fmt.Errorf("unknown field: '%s'", key)
		}
	}

	if err := json.Unmarshal(bytes, &alias); err != nil {
		return err
	}

	*s = Style(alias)
	return nil
}

func (s *Style) ToReadable() Style {
	r := *s
	r.Color = color2str(r.Color)
	return r
}

func (s *Style) FromReadable() error {
	c, err := str2color(s.Color)
	if err != nil {
		return err
	}
	s.Color = c
	return nil
}

type Theme map[string]Style

func (t Theme) UnmarshalJSON(bytes []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	for key := range raw {
		s := Style{}
		if err := json.Unmarshal(raw[key], &s); err != nil {
			return fmt.Errorf("failed at key '%s': %s", key, err.Error())
		}
		t[key] = s
	}
	return nil
}

func color2str(color string) string {
	switch color {
	case constval.Red:
		return "red"
	case constval.Green:
		return "green"
	case constval.Yellow:
		return "yellow"
	case constval.Blue:
		return "blue"
	case constval.Purple:
		return "purple"
	case constval.Cyan:
		return "cyan"
	case constval.White:
		return "white"
	case constval.Black:
		return "black"
	case constval.BrightRed:
		return "BrightPed"
	case constval.BrightGreen:
		return "BrightPreen"
	case constval.BrightYellow:
		return "BrightYellow"
	case constval.BrightBlue:
		return "BrightBlue"
	case constval.BrightPurple:
		return "BrightPurple"
	case constval.BrightCyan:
		return "BrightCyan"
	case constval.BrightWhite:
		return "BrightWhite"
	case constval.BrightBlack:
		return "BrightBlack"
	case constval.Reset:
		return "reset"
	case constval.Underline:
		return "underline"
	default:
		// detect format:
		strReader := strings.NewReader(color)

		// 1.8bit/256 color
		var c uint8
		_, err := fmt.Fscanf(strReader, Color256Format, &c)
		if err == nil {
			return fmt.Sprintf("[%d]@256", c)
		}
		// 2.rgb
		var (
			r uint8 = 0
			g uint8 = 0
			b uint8 = 0
		)
		strReader = strings.NewReader(color)
		_, err = fmt.Fscanf(strReader, RGBFormat, &r, &g, &b)
		if err == nil {
			return fmt.Sprintf("[%d,%d,%d]@rgb", r, g, b)
		}

		color = strings.ReplaceAll(color, " ", "")
		if strings.HasPrefix(color, constval.Underline) {
			return color2str(constval.Underline) + " + " + color2str(color[len(constval.Underline):])
		}
		return ""
	}
}

// str2color convert string to color
// support: red, green, yellow, blue, purple, cyan, white, black, and their bright version
// Underline
// any color with underline, should be in the format of "Underline + [color]"
// [value]@256
// [values]@rgb
// [values]@hex (will be turned to rgb)
func str2color(str string) (string, error) {
	switch str {
	case "":
		return "", nil
	case "black", "Black":
		return constval.Black, nil
	case "red", "Red":
		return constval.Red, nil
	case "green", "Green":
		return constval.Green, nil
	case "yellow", "Yellow":
		return constval.Yellow, nil
	case "blue", "Blue":
		return constval.Blue, nil
	case "purple", "Purple":
		return constval.Purple, nil
	case "cyan", "Cyan":
		return constval.Cyan, nil
	case "white", "White":
		return constval.White, nil
	case "bright-red", "BrightRed":
		return constval.BrightRed, nil
	case "bright-green", "BrightGreen":
		return constval.BrightGreen, nil
	case "bright-yellow", "BrightYellow":
		return constval.BrightYellow, nil
	case "bright-blue", "BrightBlue":
		return constval.BrightBlue, nil
	case "bright-purple", "BrightPurple":
		return constval.BrightPurple, nil
	case "bright-cyan", "BrightCyan":
		return constval.BrightCyan, nil
	case "bright-white", "BrightWhite":
		return constval.BrightWhite, nil
	case "bright-black", "BrightBlack":
		return constval.BrightBlack, nil
	case "reset", "Reset":
		return constval.Reset, nil
	default:
		// remove spaces
		str = strings.ReplaceAll(str, " ", "")

		// 256 color
		if strings.HasSuffix(str, "@256") {
			code, err := strconv.Atoi(strings.Trim(str[:len(str)-4], "[]"))
			if err != nil {
				return "", err
			}
			colorStr, err := Color256(code)
			if err != nil {
				return "", err
			}
			return colorStr, nil
		}

		// rgb color
		if strings.HasSuffix(str, "@rgb") {
			code := strings.Trim(str[:len(str)-4], "[]")
			rgb := strings.Split(code, ",")
			if len(rgb) != 3 {
				return "", errors.New("too many or too few rgb values")
			}
			r, err1 := strconv.ParseUint(rgb[0], 10, 8)
			g, err2 := strconv.ParseUint(rgb[1], 10, 8)
			b, err3 := strconv.ParseUint(rgb[2], 10, 8)
			if err1 != nil || err2 != nil || err3 != nil {
				return "", errors.New("rgb values must be numbers")
			}
			colorStr, err := RGB(uint8(r), uint8(g), uint8(b))
			if err != nil {
				return "", err
			}
			return colorStr, nil
		}

		// hex
		if strings.HasSuffix(str, "@hex") {
			code := strings.Trim(str[:len(str)-4], "[]")
			rgb := HexToRgb(code)
			colorStr, err := RGB(rgb[0], rgb[1], rgb[2])
			if err != nil {
				return "", errors.New("rgb values must be numbers")
			}
			return colorStr, nil
		}

		return constval.Reset, nil
	}
}

/*
......
// if using 256 color, you can use color code like this:
[0-255]@256
// if using rgb color, you can use color code like this:
[0-255,0-255,0-255]@rgb
// if using hex color, you can use color code like this:
[hex]@hex
*/

type ErrBadColor struct {
	name string
	error
}

func (e ErrBadColor) Error() string {
	return fmt.Sprintf("bad color for %s: %s", e.name, e.error.Error())
}

type ErrOpenTheme struct {
	error
}

func (e ErrOpenTheme) Error() string {
	return fmt.Sprintf("load theme error: %s", e.error.Error())
}

func GetTheme(path string) error {
	themeJson, err := os.ReadFile(path)
	if err != nil {
		return ErrOpenTheme{err}
	}
	a, err, fatal := getTheme(themeJson)
	if fatal != nil {
		return ErrOpenTheme{fatal}
	}
	DefaultAll = a
	return err
}

func getTheme(themeJson []byte) (theme All, errSum error, fatal error) {
	err := json.Unmarshal(themeJson, &theme)
	if err != nil {
		return All{}, nil, err
	}
	// check
	errSum = nil
	check := func(m Theme) {
		for key := range m {
			style := m[key]
			err := style.FromReadable()
			if err != nil {
				errSum = errors.Join(errSum, ErrBadColor{name: key, error: err})
				continue
			}
			m[key] = style
		}
	}
	theme.Apply(check)
	return theme, errSum, nil
}

func ConvertThemeColor() {
	convert := func(m Theme) {
		for key := range m {
			if key == "reset" {
				continue
			}
			color, err := ConvertColorIfGreaterThanExpect(ColorLevel, m[key].Color)
			if err != nil {
				continue
			}
			m[key] = Style{
				Icon:      m[key].Icon,
				Color:     color,
				Underline: m[key].Underline,
				Bold:      m[key].Bold,
				Faint:     m[key].Faint,
				Italics:   m[key].Italics,
				Blink:     m[key].Blink,
			}
		}
	}
	DefaultAll.Apply(convert)
}
