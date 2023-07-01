package theme

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
	"inode": Style{
		Color: Purple,
	},
	"time": Style{
		Color: Blue,
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
	"link": Style{
		Color: Purple,
	},
	"bit": Style{
		Color: rgb(20, 255, 100),
	},
	"B": Style{
		Color: rgb(20, 230, 100),
	},
	"KB": Style{
		Color: rgb(20, 207, 100),
	},
	"MB": Style{
		Color: rgb(20, 188, 100),
	},
	"GB": Style{
		Color: rgb(20, 170, 100),
	},
	"TB": Style{
		Color: rgb(20, 153, 100),
	},
	"PB": Style{
		Color: rgb(20, 138, 100),
	},
	"EB": Style{
		Color: rgb(20, 125, 100),
	},
	"ZB": Style{
		Color: rgb(20, 112, 100),
	},
	"YB": Style{
		Color: rgb(20, 100, 100),
	},
	"BB": Style{
		Color: rgb(20, 90, 100),
	},
	"NB": Style{
		Color: rgb(20, 70, 100),
	},
	"git_unmodified": Style{
		Color: BrightBlack,
	},
	"git_modified": Style{
		Color: Yellow,
	},
	"git_renamed": Style{
		Color: Blue,
	},
	"git_copied": Style{
		Color: Purple,
	},
	"git_deleted": Style{
		Color: Red,
	},
	"git_added": Style{
		Color: Green,
	},
	"git_untracked": Style{
		Color: BrightBlack,
	},
	"git_ignored": Style{
		Color: BrightRed,
	},
	"git_type_changed": Style{
		Color: Yellow,
	},
	"git_updated_but_unmerged": Style{
		Color: BrightYellow,
	},
	"symlink_path": Style{
		Color: Green,
	},
	"symlink_broken_path": Style{
		Color: Red,
	},
}

var DefaultTheme = Theme{
	"dir": Style{
		Color: BrightBlue,
		Icon:  "\uF115",
	},
	"home": Style{
		Color: BrightBlue,
		Icon:  "\ue617",
	},
	"Desktop": Style{
		Color: BrightBlue,
		Icon:  "\uf108",
	},
	"desktop": Style{
		Color: BrightBlue,
		Icon:  "\uf108",
	},
	"downloads": Style{
		Color: BrightBlue,
		Icon:  "\uf498",
	},
	"Downloads": Style{
		Color: BrightBlue,
		Icon:  "\uf498",
	},
	"Trash": Style{
		Color: BrightBlue,
		Icon:  "\uf1f8",
	},
	"Searches": Style{
		Color: BrightBlue,
		Icon:  "\uf422",
	},
	"DS_Store": Style{
		Color: BrightBlue,
		Icon:  "\uf179",
	},
	"ssh": Style{
		Color: BrightBlue,
		Icon:  "\ueba9",
	},
	"boot": Style{
		Color: BrightBlue,
		Icon:  "\uf287",
	},
	"cache": Style{
		Color: BrightBlue,
		Icon:  "\uf49b",
	},
	"exe": Style{
		Color: Green,
		Icon:  "\uF17A",
	},
	"EXE": Style{
		Color: Green,
		Icon:  "\uF17A",
	},
	"file": Style{
		Color: White,
		Icon:  "\uF016",
	},
	"rdp": Style{
		Color: White,
		Icon:  "\ueb39",
	},
	"known_hosts": Style{
		Color: White,
		Icon:  "\uEB39",
	},
	"repo": Style{
		Color: White,
		Icon:  "\uea62",
	},
	"link": Style{
		Color: Purple,
		Icon:  "\ueb15",
	},
	"lnk": Style{
		Color: Purple,
		Icon:  "\ueb15",
	},
	"go": Style{
		Color: Cyan,
		Icon:  "\uE626",
	},
	"rs": Style{
		Color: Cyan,
		Icon:  "\uE7A8",
	},
	"cargo": {
		Color: Cyan,
		Icon:  "\uE7A8",
	},
	"rustup": {
		Color: Cyan,
		Icon:  "\uE7A8",
	},
	"cargo.lock": {
		Color: Cyan,
		Icon:  "\uE7A8",
	},
	"cargo.toml": {
		Color: Cyan,
		Icon:  "\uE7A8",
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
	"include": Style{
		Color: Cyan,
		Icon:  "\ue5fc",
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
		Icon:  "\U000F031B",
	},
	"dotnet": Style{
		Color: Cyan,
		Icon:  "\ue72e",
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
		Icon:  "\uE61E",
	},
	"deb": Style{
		Color: Cyan,
		Icon:  "\uF306",
	},
	"apk": Style{
		Color: Green,
		Icon:  "\uF17B",
	},
	"PKGBUILD": Style{
		Color: Green,
		Icon:  "\uf303",
	},
	"app": Style{
		Color: Green,
		Icon:  "\ueb44",
	},
	"Applications": Style{
		Color: Green,
		Icon:  "\ueb44",
	},
	"msi": Style{
		Color: Green,
		Icon:  "\ue70f",
	},
	"dpkg": Style{
		Color: Cyan,
		Icon:  "\ue77d",
	},
	"android": Style{
		Color: Cyan,
		Icon:  "\uE70E",
	},
	"src": Style{
		Color: Cyan,
		Icon:  "\ue796",
	},
	"py": Style{
		Color: Cyan,
		Icon:  "\uE606",
	},
	"ipynb": Style{
		Color: Cyan,
		Icon:  "\uE606",
	},
	"pyc": Style{
		Color: Cyan,
		Icon:  "\uE606",
	},
	"whl": Style{
		Color: Cyan,
		Icon:  "\uE606",
	},
	"rb": Style{
		Color: Cyan,
		Icon:  "\uE21E",
	},
	"js": Style{
		Color: Cyan,
		Icon:  "\uE781",
	},
	"ts": Style{
		Color: Cyan,
		Icon:  "\uE628",
	},
	"http": Style{
		Color: Cyan,
		Icon:  "\ueb01",
	},
	"node_modules": Style{
		Color: Cyan,
		Icon:  "\ue5fa",
	},
	"npm": Style{
		Color: Cyan,
		Icon:  "\ue71e",
	},
	"jsx": Style{
		Color: Cyan,
		Icon:  "\ue7ba",
	},
	"htm": Style{
		Color: Cyan,
		Icon:  "\uF13B",
	},
	"html": Style{
		Color: Cyan,
		Icon:  "\uF13B",
	},
	"css": Style{
		Color: Cyan,
		Icon:  "\uE749",
	},
	"java": Style{
		Color: Cyan,
		Icon:  "\uE738",
	},
	"jar": Style{
		Color: Cyan,
		Icon:  "\uE738",
	},
	"json": Style{
		Color: Cyan,
		Icon:  "\uE60B",
	},
	"xml": Style{
		Color: Cyan,
		Icon:  "\ue796",
	},
	"cson": Style{
		Color: Cyan,
		Icon:  "\uE601",
	},
	"font": Style{
		Color: Cyan,
		Icon:  "\uf031",
	},
	"ttf": Style{
		Color: Cyan,
		Icon:  "\uf031",
	},
	"properties": Style{
		Color: Cyan,
		Icon:  "\uE60B",
	},
	"git": Style{
		Color: Cyan,
		Icon:  "\uF1D3",
	},
	"gitignore": {
		Color: Black,
		Icon:  "\uf1d3",
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
	"erl": Style{
		Color: Cyan,
		Icon:  "\ue7b1",
	},
	"hrl": Style{
		Color: Cyan,
		Icon:  "\uE7B1",
	},
	"coffee": Style{
		Color: Cyan,
		Icon:  "\uE751",
	},
	"iso": Style{
		Color: Cyan,
		Icon:  "\uF7C9",
	},
	"dmg": Style{
		Color: Cyan,
		Icon:  "\uF7C9",
	},
	"img": Style{
		Color: Cyan,
		Icon:  "\uF7C9",
	},
	"vhd": Style{
		Color: Cyan,
		Icon:  "\uE61C",
	},
	"vhdx": Style{
		Color: Cyan,
		Icon:  "\uE61C",
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
	"justfile": Style{
		Color: Cyan,
		Icon:  "\uE20F",
	},
	"cmake": Style{
		Color: Cyan,
		Icon:  "\uE20F",
	},
	"less": Style{
		Color: Cyan,
		Icon:  "\ue758",
	},
	"r": Style{
		Color: Cyan,
		Icon:  "\uF25D",
	},
	"f": Style{
		Color: Cyan,
		Icon:  "\uF794",
	},
	"history": Style{
		Color: Cyan,
		Icon:  "\uF1DA",
	},
	"audio": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"php": Style{
		Color: Cyan,
		Icon:  "\uE73D",
	},
	"apple": Style{
		Color: Cyan,
		Icon:  "\uF179",
	},
	"atom": Style{
		Color: Cyan,
		Icon:  "\ue764",
	},
	"docker": Style{
		Color: Cyan,
		Icon:  "\ue7b0",
	},
	"dockerfile": Style{
		Color: Cyan,
		Icon:  "\ue7b0",
	},
	"Dockerfile": Style{
		Color: Cyan,
		Icon:  "\ue7b0",
	},
	"Docker": Style{
		Color: Cyan,
		Icon:  "\ue7b0",
	},
	"md": Style{
		Color: BrightYellow,
		Icon:  "\uF48A",
	},
	"README": Style{
		Color: BrightYellow,
		Icon:  "\uF48A",
	},
	"txt": Style{
		Color: Cyan,
		Icon:  "\uF15C",
	},
	"sdk": Style{
		Color: Cyan,
		Icon:  "\uF121",
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
		Icon:  "",
	},
	"sitx": Style{
		Color: Cyan,
		Icon:  "",
	},
	"tar": Style{
		Color: Cyan,
		Icon:  "",
	},
	"gz": Style{
		Color: Cyan,
		Icon:  "",
	},
	"rar": Style{
		Color: Cyan,
		Icon:  "",
	},
	"7z": Style{
		Color: Cyan,
		Icon:  "",
	},
	"Pictures": Style{
		Color: Purple,
		Icon:  "\uF1C5",
	},
	"jpg": Style{
		Color: Purple,
		Icon:  "\uF1C5",
	},
	"jpeg": Style{
		Color: Purple,
		Icon:  "\uF1C5",
	},
	"png": Style{
		Color: Purple,
		Icon:  "\uF1C5",
	},
	"bmp": Style{
		Color: Purple,
		Icon:  "\uF1C5",
	},
	"tif": Style{
		Color: Purple,
		Icon:  "\uF1C5",
	},
	"tiff": Style{
		Color: Purple,
		Icon:  "\uF1C5",
	},
	"gif": Style{
		Color: Purple,
		Icon:  "\uF1C5",
	},
	"Videos": {
		Color: Cyan,
		Icon:  "",
	},
	"Movies": {
		Color: Cyan,
		Icon:  "",
	},
	"mp4": Style{
		Color: Cyan,
		Icon:  "",
	},
	"mkv": Style{
		Color: Cyan,
		Icon:  "",
	},
	"avi": Style{
		Color: Cyan,
		Icon:  "",
	},
	"flv": Style{
		Color: Cyan,
		Icon:  "",
	},
	"mov": Style{
		Color: Cyan,
		Icon:  "\uF03D",
	},
	"ai": Style{
		Color: Cyan,
		Icon:  "\uE7B4",
	},
	"music": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"Music": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"mp3": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"m4a": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"mid": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"midi": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"mpeg": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"flac": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"ape": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"alac": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"aac": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"wav": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"ogg": Style{
		Color: Cyan,
		Icon:  "\uf1c7",
	},
	"pdf": Style{
		Color: Cyan,
		Icon:  "\uF1C1",
	},
	"epub": Style{
		Color: Cyan,
		Icon:  "\uE28A",
	},
	"mobi": Style{
		Color: Cyan,
		Icon:  "\uF12D",
	},
	"azw": Style{
		Color: Cyan,
		Icon:  "\uF12D",
	},
	"azw3": Style{
		Color: Cyan,
		Icon:  "\uF12D",
	},
	"doc": Style{
		Color: Cyan,
		Icon:  "\uf1c2",
	},
	"docx": Style{
		Color: Cyan,
		Icon:  "\uf1c2",
	},
	"docm": Style{
		Color: Cyan,
		Icon:  "\uf1c2",
	},
	"xls": Style{
		Color: Cyan,
		Icon:  "\uf1c3",
	},
	"xlsx": Style{
		Color: Cyan,
		Icon:  "\uf1c3",
	},
	"xlsm": Style{
		Color: Cyan,
		Icon:  "\uf1c3",
	},
	"numbers": Style{
		Color: Cyan,
		Icon:  "\uf1c3",
	},
	"csv": Style{
		Color: Cyan,
		Icon:  "\ue64a",
	},
	"ppt": Style{
		Color: Cyan,
		Icon:  "\uf1c4",
	},
	"pptm": Style{
		Color: Cyan,
		Icon:  "\uf1c4",
	},
	"pptx": Style{
		Color: Cyan,
		Icon:  "\uf1c4",
	},
	"dot": Style{
		Color: Cyan,
		Icon:  "\uf1c2",
	},
	"dotx": Style{
		Color: Cyan,
		Icon:  "\uf1c2",
	},
	"diff": Style{
		Color: Cyan,
		Icon:  "\uF440",
	},
	"ini": Style{
		Color: Cyan,
		Icon:  "\uE615",
	},
	"conf": Style{
		Color: Cyan,
		Icon:  "\uE615",
	},
	"cfg": Style{
		Color: Cyan,
		Icon:  "\uE615",
	},
	"config": Style{
		Color: Cyan,
		Icon:  "\uE615",
	},
	"yml": Style{
		Color: Cyan,
		Icon:  "\uE615",
	},
	"yaml": Style{
		Color: Cyan,
		Icon:  "\uF481",
	},
	"tex": Style{
		Color: Cyan,
		Icon:  "\uF034",
	},
	"typ": Style{
		Color: Cyan,
		Icon:  "∫",
	},
	"toml": Style{
		Color: Cyan,
		Icon:  "\uF615",
	},
	"db": Style{
		Color: Cyan,
		Icon:  "\uF1C0",
	},
	"accdb": Style{
		Color: Cyan,
		Icon:  "\uF1C0",
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
		Icon:  "\uF1C0",
	},
	"sh": Style{
		Color: BrightGreen,
		Icon:  "\uF489",
	},
	"bat": Style{
		Color: BrightGreen,
		Icon:  "\ue629",
	},
	"ps1": Style{
		Color: BrightGreen,
		Icon:  "\uF489",
	},
	"csh": Style{
		Color: BrightGreen,
		Icon:  "\uF489",
	},
	"fish": Style{
		Color: BrightGreen,
		Icon:  "\uF489",
	},
	"ksh": Style{
		Color: BrightGreen,
		Icon:  "\uF489",
	},
	"zsh": Style{
		Color: BrightGreen,
		Icon:  "\uF489",
	},
	"zsh_history": Style{
		Color: Cyan,
		Icon:  "\uF489",
	},
	"zshrc": Style{
		Color: Cyan,
		Icon:  "\uF489",
	},
	"bash": Style{
		Color: BrightGreen,
		Icon:  "\uF489",
	},
	"bashrc": Style{
		Color: Cyan,
		Icon:  "\uF489",
	},
	"bash_history": Style{
		Color: Cyan,
		Icon:  "\uF489",
	},
	"bash_profile": Style{
		Color: Cyan,
		Icon:  "\uF489",
	},
	"vim": Style{
		Color: Cyan,
		Icon:  "\uE62B",
	},
	"viminfo": Style{
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
		Icon:  "\uf470",
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
	"pub": Style{
		Color: Cyan,
		Icon:  "\uF084",
	},
	"id_rsa": Style{
		Color: Cyan,
		Icon:  "\uF084",
	},
	"gpg": Style{
		Color: Cyan,
		Icon:  "\uF084",
	},
	"cer": Style{
		Color: Cyan,
		Icon:  "\uF084",
	},
	"crt": Style{
		Color: Cyan,
		Icon:  "\uF084",
	},
	"pgp": Style{
		Color: Cyan,
		Icon:  "\uF084",
	},
	"license": Style{
		Color: Cyan,
		Icon:  "\uF084",
	},
	"LICENSE": Style{
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
	"lib": Style{
		Color: Cyan,
		Icon:  "\uF121",
	},
	"Library": Style{
		Color: Cyan,
		Icon:  "\uF121",
	},
	"bin": Style{
		Color: Cyan,
		Icon:  "\uE5FC",
	},
	"share": Style{
		Color: Cyan,
		Icon:  "\uf064",
	},
	"idea": Style{
		Color: Cyan,
		Icon:  "\uE7B5",
	},
	"so": Style{
		Color: Cyan,
		Icon:  "\uF121",
	},
	"dll": Style{
		Color: Cyan,
		Icon:  "\uF121",
	},
	"pipe": Style{
		Color: Cyan,
		Icon:  "\uF124",
	},
	"socket": Style{
		Color: Cyan,
		Icon:  "\uF1E6",
	},
	"symlink": Style{
		Color: Purple,
		Icon:  "\ueb15",
	},
	"symlink_arrow": Style{
		Color: BrightWhite,
		Icon:  " ~> ",
	},
	"Users": Style{
		Color: Cyan,
		Icon:  "\uf0c0",
	},
}
