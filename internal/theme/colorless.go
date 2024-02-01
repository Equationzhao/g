package theme

func SetClassic() {
	DefaultAll.Apply(setClassic)
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
