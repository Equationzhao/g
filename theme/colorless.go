package theme

func RemoveAllColor() {
	DefaultAll.Apply(resetColor)
	DefaultAll.InfoTheme["reset"] = Style{
		Color: Reset,
	}
}

func resetColor(m Theme) {
	for k := range m {
		m[k] = Style{
			Icon:      m[k].Icon,
			Color:     "",
			Underline: false,
			Bold:      false,
			Faint:     false,
			Italics:   false,
		}
	}
}
