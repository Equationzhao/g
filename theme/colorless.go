package theme

func RemoveAllColor() {
	resetColor(DefaultAll.InfoTheme)
	resetColor(DefaultAll.Permission)
	resetColor(DefaultAll.Size)
	resetColor(DefaultAll.User)
	resetColor(DefaultAll.Group)
	resetColor(DefaultAll.Symlink)
	resetColor(DefaultAll.Git)
	resetColor(DefaultAll.Name)
	resetColor(DefaultAll.Special)
	resetColor(DefaultAll.Ext)
	DefaultAll.InfoTheme["reset"] = Style{
		Color: Reset,
	}
}

func resetColor(m Theme) {
	for k := range m {
		m[k] = Style{
			Icon:      InfoTheme[k].Icon,
			Color:     "",
			Underline: false,
			Bold:      false,
		}
	}
}
