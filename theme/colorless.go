package theme

var ColorlessInfo = Theme{}
var Colorless = Theme{}

func init() {
	SyncColorlessWithTheme()
}

func SyncColorlessWithTheme() {
	for k := range DefaultInfoTheme {
		ColorlessInfo[k] = Style{
			Icon:  DefaultInfoTheme[k].Icon,
			Color: "",
		}
	}
	for k := range DefaultTheme {
		Colorless[k] = Style{
			Icon:  DefaultTheme[k].Icon,
			Color: "",
		}
	}
}
