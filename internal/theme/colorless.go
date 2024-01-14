package theme

func SetClassic() {
	DefaultAll.Apply(setClassic)
	DefaultAll.InfoTheme["reset"] = Style{
		Color: Reset,
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
