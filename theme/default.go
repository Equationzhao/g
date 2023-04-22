package theme

const (
	Black  = "\033[1;30m"
	Red    = "\033[1;31m"
	Green  = "\033[1;32m"
	Yellow = "\033[1;33m"
	Blue   = "\033[1;34m"
	Purple = "\033[1;35m"
	Cyan   = "\033[1;36m"
	White  = "\033[1;37m"
	Reset  = "\033[0m"
)

var DefaultInfoTheme = Theme{
	"d": Style{
		Color: Blue,
	},
	"l": Style{
		Color: Purple,
	},
	"b": Style{
		Color: Yellow,
	},
	"c": Style{
		Color: Yellow,
	},
	"p": Style{
		Color: Yellow,
	},
	"s": Style{
		Color: Yellow,
	},
	"r": Style{
		Color: Yellow,
	},
	"w": Style{
		Color: Red,
	},
	"x": Style{
		Color: Green,
	},
	"-": Style{
		Color: White,
	},
	"time": Style{
		Color: Blue,
	},
	"size": Style{
		Color: Green,
	},
	"owner": Style{
		Color: Yellow,
	},
	"group": Style{
		Color: Yellow,
	},
	"reset": Style{
		Color: Reset,
	},
	"root": Style{
		Color: Red,
	},
}

var DefaultTheme = Theme{
	"dir": {
		Color: Blue,
		Icon:  "📁",
	},
	"exe": Style{
		Color: Green,
		Icon:  "🚀",
	},
	"file": Style{
		Color: White,
		Icon:  "📄",
	},
	"symlink": Style{
		Color: Purple,
		Icon:  "🔗",
	},
	"link": Style{
		Color: Purple,
		Icon:  "🔗",
	},
	"go": Style{
		Color: Cyan,
		Icon:  "🐹",
	},
	"rs": Style{
		Color: Cyan,
		Icon:  "🦀",
	},
	"c": Style{
		Color: Cyan,
		Icon:  "\uE61E",
	},
	"cpp": Style{
		Color: Cyan,
		Icon:  "\uE61D",
	},
	"c++": Style{
		Color: Cyan,
		Icon:  "\uE61D",
	},
	"cc": Style{
		Color: Cyan,
		Icon:  "\uE61D",
	},
	"cxx": Style{
		Color: Cyan,
		Icon:  "\uE61D",
	},
	"h": Style{
		Color: Cyan,
		Icon:  "\uE61F",
	},
	"hpp": Style{
		Color: Cyan,
		Icon:  "\uE61F",
	},
	"h++": Style{
		Color: Cyan,
		Icon:  "\uE61F",
	},
	"hh": Style{
		Color: Cyan,
		Icon:  "\uE61F",
	},
	"hxx": Style{
		Color: Cyan,
		Icon:  "\uE61F",
	},
	"cs": Style{
		Color: Cyan,
		Icon:  "\uF81A",
	},
	"scala": Style{
		Color: Cyan,
		Icon:  "\uE737",
	},
	"swift": Style{
		Color: Cyan,
		Icon:  "\uE755",
	},
	"kt": Style{
		Color: Cyan,
		Icon:  "\uE634",
	},
	"m": Style{
		Color: Cyan,
		Icon:  "ﬧ",
	},
	"deb": Style{
		Color: Cyan,
		Icon:  "\uF306",
	},
	"dpkg": Style{
		Color: Cyan,
		Icon:  "\uF17C",
	},
	"android": Style{
		Color: Cyan,
		Icon:  "\uE70E",
	},
	"py": Style{
		Color: Cyan,
		Icon:  "🐍",
	},
	"rb": Style{
		Color: Cyan,
		Icon:  "🐇",
	},
	"js": Style{
		Color: Cyan,
		Icon:  "\uE781",
	},
	"ts": Style{
		Color: Cyan,
		Icon:  "ﯤ",
	},
	"jsx": Style{
		Color: Cyan,
		Icon:  "\ue7ba",
	},
	"html": Style{
		Color: Cyan,
		Icon:  "🌐",
	},
	"css": Style{
		Color: Cyan,
		Icon:  "🌐",
	},
	"java": Style{
		Color: Cyan,
		Icon:  "\uE738",
	},
	"json": Style{
		Color: Cyan,
		Icon:  "\uE60B",
	},
	"cson": Style{
		Color: Cyan,
		Icon:  "\uE601",
	},
	"font": Style{
		Color: Cyan,
		Icon:  "\uf031",
	},
	"properties": Style{
		Color: Cyan,
		Icon:  "\uE60B",
	},
	"gitignore": Style{
		Color: Cyan,
		Icon:  "\uF1D3",
	},
	"git": Style{
		Color: Cyan,
		Icon:  "\uF1D3",
	},
	"asm": Style{
		Color: Cyan,
		Icon:  "\uFB19",
	},
	"groovy": Style{
		Color: Cyan,
		Icon:  "\ue775",
	},
	"s": Style{
		Color: Cyan,
		Icon:  "\uFB19",
	},
	"gv": Style{
		Color: Cyan,
		Icon:  "\uE225",
	},
	"hs": Style{
		Color: Cyan,
		Icon:  "\ue777",
	},
	"d": Style{
		Color: Cyan,
		Icon:  "\ue7af",
	},
	"dart": Style{
		Color: Cyan,
		Icon:  "\uE798",
	},
	"coffee": Style{
		Color: Cyan,
		Icon:  "\uE751",
	},
	"iso": Style{
		Color: Cyan,
		Icon:  "\uF7C9",
	},
	"lua": Style{
		Color: Cyan,
		Icon:  "\uE620",
	},
	"vue": Style{
		Color: Cyan,
		Icon:  "\ufd42",
	},
	"makefile": Style{
		Color: Cyan,
		Icon:  "\uE20F",
	},
	"Makefile": Style{
		Color: Cyan,
		Icon:  "\uE20F",
	},
	"less": Style{
		Color: Cyan,
		Icon:  "\ue758",
	},
	"r": Style{
		Color: Cyan,
		Icon:  "ﳒ",
	},
	"f": Style{
		Color: Cyan,
		Icon:  "\uF794",
	},
	"audio": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"php": Style{
		Color: Cyan,
		Icon:  "🐘",
	},
	"apple": Style{
		Color: Cyan,
		Icon:  "\uF179",
	},
	"dockerfile": Style{
		Color: Cyan,
		Icon:  "\ue7b0",
	},
	"Dockerfile": Style{
		Color: Cyan,
		Icon:  "\ue7b0",
	},
	"md": Style{
		Color: Cyan,
		Icon:  "\uF48A",
	},
	"txt": Style{
		Color: Cyan,
		Icon:  "\uF15C",
	},
	"sh": Style{
		Color: Cyan,
		Icon:  "🐚",
	},
	"bat": Style{
		Color: Cyan,
		Icon:  "\uF17A",
	},
	"ps1": Style{
		Color: Cyan,
		Icon:  "🐚",
	},
	"zig": Style{
		Color: Cyan,
		Icon:  "\uF0E7",
	},
	"rss": Style{
		Color: Cyan,
		Icon:  "\uF09E",
	},
	"ko": Style{
		Color: Cyan,
		Icon:  "\uebc6",
	},
	"zip": Style{
		Color: Cyan,
		Icon:  "🗜",
	},
	"tar": Style{
		Color: Cyan,
		Icon:  "🗜",
	},
	"gz": Style{
		Color: Cyan,
		Icon:  "🗜",
	},
	"rar": Style{
		Color: Cyan,
		Icon:  "🗜",
	},
	"7z": Style{
		Color: Cyan,
		Icon:  "🗜",
	},
	"jpg": Style{
		Color: Cyan,
		Icon:  "🖼",
	},
	"jpeg": Style{
		Color: Cyan,
		Icon:  "🖼",
	},
	"png": Style{
		Color: Cyan,
		Icon:  "🖼",
	},
	"gif": Style{
		Color: Cyan,
		Icon:  "🖼",
	},
	"mp4": Style{
		Color: Cyan,
		Icon:  "🎥",
	},
	"mkv": Style{
		Color: Cyan,
		Icon:  "🎥",
	},
	"avi": Style{
		Color: Cyan,
		Icon:  "🎥",
	},
	"flv": Style{
		Color: Cyan,
		Icon:  "🎥",
	},
	"mov": Style{
		Color: Cyan,
		Icon:  "🎥",
	},
	"ai": Style{
		Color: Cyan,
		Icon:  "🎨",
	},
	"mp3": Style{
		Color: Cyan,
		Icon:  "🎵",
	},
	"m4a": Style{
		Color: Cyan,
		Icon:  "🎵",
	},
	"flac": Style{
		Color: Cyan,
		Icon:  "🎵",
	},
	"ape": Style{
		Color: Cyan,
		Icon:  "🎵",
	},
	"alac": Style{
		Color: Cyan,
		Icon:  "🎵",
	},
	"aac": Style{
		Color: Cyan,
		Icon:  "🎵",
	},
	"wav": Style{
		Color: Cyan,
		Icon:  "🎵",
	},
	"ogg": Style{
		Color: Cyan,
		Icon:  "🎵",
	},
	"pdf": Style{
		Color: Cyan,
		Icon:  "\uF1C1",
	},
	"doc": Style{
		Color: Cyan,
		Icon:  "\uf1c2",
	},
	"diff": Style{
		Color: Cyan,
		Icon:  "\uF440",
	},
	"ppt": Style{
		Color: Cyan,
		Icon:  "\uf1c4",
	},
	"ini": Style{
		Color: Cyan,
		Icon:  "⚙️",
	},
	"conf": Style{
		Color: Cyan,
		Icon:  "⚙️",
	},
	"cfg": Style{
		Color: Cyan,
		Icon:  "⚙️",
	},
	"config": Style{
		Color: Cyan,
		Icon:  "⚙️",
	},
	"yml": Style{
		Color: Cyan,
		Icon:  "⚙️",
	},
	"yaml": Style{
		Color: Cyan,
		Icon:  "⚙️",
	},
	"tex": Style{
		Color: Cyan,
		Icon:  "∫",
	},
	"toml": Style{
		Color: Cyan,
		Icon:  "⚙️",
	},
	"db": Style{
		Color: Cyan,
		Icon:  "🗄",
	},
	"sqlite": Style{
		Color: Cyan,
		Icon:  "\uE7C4",
	},
	"sqlite3": Style{
		Color: Cyan,
		Icon:  "\uE7C4",
	},
	"sql": Style{
		Color: Cyan,
		Icon:  "\uE706",
	},
	"db3": Style{
		Color: Cyan,
		Icon:  "🗄",
	},
	"zsh": Style{
		Color: Cyan,
		Icon:  "🐚",
	},
	"bash": Style{
		Color: Cyan,
		Icon:  "🐚",
	},
	"vim": Style{
		Color: Cyan,
		Icon:  "\uE62B",
	},
	"vimrc": Style{
		Color: Cyan,
		Icon:  "\uE62B",
	},
	"gvimrc": Style{
		Color: Cyan,
		Icon:  "\uE62B",
	},
	"nvim": Style{
		Color: Cyan,
		Icon:  "\uE62B",
	},
	"nvimrc": Style{
		Color: Cyan,
		Icon:  "\uE62B",
	},
	"log": Style{
		Color: White,
		Icon:  "\uF18D",
	},
	"lock": Style{
		Color: Cyan,
		Icon:  "\uF023",
	},
	"github": Style{
		Color: Cyan,
		Icon:  "\uF09B",
	},
	"vscode": Style{
		Color: Cyan,
		Icon:  "\uE70C",
	},
	"code-workspace": Style{
		Color: Cyan,
		Icon:  "\uE70C",
	},
	"key": Style{
		Color: Cyan,
		Icon:  "\uF084",
	},
	"passwd": Style{
		Color: Cyan,
		Icon:  "\uF023",
	},
	"maintainers": Style{
		Color: Cyan,
		Icon:  "\uF0C0",
	},
}
