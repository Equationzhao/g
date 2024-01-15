package theme

import "github.com/Equationzhao/g/internal/const"

func SetClassic() {
	DefaultAll.Apply(setClassic)
	DefaultAll.InfoTheme["reset"] = Style{
		Color: constval.Reset,
	}
}

func setClassic(m Theme) {
	for k := range m {
		m[k] = Style{
			Icon:      m[k].Icon,
			Color:     "",
			Underline: false,
			Bold:      false,
			Faint:     false,
			Italics:   false,
			Blink:     false,
		}
	}
}
