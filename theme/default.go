package theme

var (
	dir   = BrightBlue
	pic   = BrightPurple
	video = color256(200)
	audio = color256(199)
	db    = color256(195)
	lang  = color256(158)
	text  = color256(153)
	doc   = color256(150)
	tar   = Red
	pkg   = BrightYellow
	bash  = BrightGreen
	lock  = color256(227)
	vim   = color256(41)
	key   = color256(214)
	conf  = White
	iso   = color256(88)
)

var DefaultInfoTheme = Theme{
	"inode": Style{
		Color: Purple,
	},
	"time": Style{
		Color: Blue,
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
}

// DefaultTheme the key should be lowercase
var DefaultTheme = Theme{
	"deb": Style{
		Color: pkg,
		Icon:  "\uF306",
	},
	"apk": Style{
		Color: Green,
		Icon:  "\uF17B",
	},
	"app": Style{
		Icon: "\ueb44",
	},
	"msi": Style{
		Color: pkg,
		Icon:  "\uF17A",
	},
	"dpkg": Style{
		Color: pkg,
		Icon:  "\ue77d",
	},
	"ssh": Style{
		Icon: "\ueba9",
	},
	"boot": Style{
		Icon: "\uf287",
	},
	"cache": Style{
		Icon: "\uf49b",
	},
	"exe": Style{
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
		Icon: "\uea62",
	},
	"link": Style{
		Color: Purple,
		Icon:  "\ueb15",
	},
	"links": Style{
		Color: Purple,
		Icon:  "\ueb15",
	},
	"lnk": Style{
		Color: Purple,
		Icon:  "\ueb15",
	},
	"go": Style{
		Color: lang,
		Icon:  "\uE626",
	},
	"rs": Style{
		Color: lang,
		Icon:  "\uE7A8",
	},
	"cargo": {
		Color: lang,
		Icon:  "\uE7A8",
	},
	"rustup": {
		Color: lang,
		Icon:  "\uE7A8",
	},
	"cargo.lock": {
		Color: lock,
		Icon:  "\uE7A8",
	},
	"cargo.toml": {
		Color: lang,
		Icon:  "\uE7A8",
	},
	"c": Style{
		Color: lang,
		Icon:  "\uE61E",
	},
	"cpp": Style{
		Color: lang,
		Icon:  "\uE61D",
	},
	"c++": Style{
		Color: lang,
		Icon:  "\uE61D",
	},
	"cc": Style{
		Color: lang,
		Icon:  "\uE61D",
	},
	"cxx": Style{
		Color: lang,
		Icon:  "\uE61D",
	},
	"h": Style{
		Color: lang,
		Icon:  "\uE61F",
	},
	"hpp": Style{
		Color: lang,
		Icon:  "\uE61F",
	},
	"h++": Style{
		Color: lang,
		Icon:  "\uE61F",
	},
	"hh": Style{
		Color: lang,
		Icon:  "\uE61F",
	},
	"hxx": Style{
		Color: lang,
		Icon:  "\uE61F",
	},
	"cs": Style{
		Color: lang,
		Icon:  "\U000F031B",
	},
	"scala": Style{
		Color: lang,
		Icon:  "\uE737",
	},
	"swift": Style{
		Color: lang,
		Icon:  "\uE755",
	},
	"kt": Style{
		Color: lang,
		Icon:  "\uE634",
	},
	"m": Style{
		Color: lang,
		Icon:  "\uE61E",
	},
	"py": Style{
		Color: lang,
		Icon:  "\uE606",
	},
	"ipynb": Style{
		Color: lang,
		Icon:  "\uE606",
	},
	"pyc": Style{
		Color: lang,
		Icon:  "\uE606",
	},
	"whl": Style{
		Color: lang,
		Icon:  "\uE606",
	},
	"rb": Style{
		Color: lang,
		Icon:  "\uE21E",
	},
	"js": Style{
		Color: lang,
		Icon:  "\uE781",
	},
	"ts": Style{
		Color: lang,
		Icon:  "\uE628",
	},
	"http": Style{
		Color: lang,
		Icon:  "\ueb01",
	},
	"npm": Style{
		Color: lang,
		Icon:  "\ue71e",
	},
	"jsx": Style{
		Color: lang,
		Icon:  "\ue7ba",
	},
	"htm": Style{
		Color: lang,
		Icon:  "\uF13B",
	},
	"html": Style{
		Color: lang,
		Icon:  "\uF13B",
	},
	"css": Style{
		Color: lang,
		Icon:  "\uE749",
	},
	"java": Style{
		Color: lang,
		Icon:  "\uE738",
	},
	"jar": Style{
		Color: lang,
		Icon:  "\uE738",
	},
	"json": Style{
		Color: conf,
		Icon:  "\uE60B",
	},
	"json5": Style{
		Color: conf,
		Icon:  "\uE60B",
	},
	"hson": Style{
		Color: conf,
		Icon:  "\uE60B",
	},
	"xml": Style{
		Color: conf,
		Icon:  "\ue796",
	},
	"cson": Style{
		Color: conf,
		Icon:  "\uE601",
	},
	"ttf": Style{
		Color: Cyan,
		Icon:  "\uf031",
	},
	"otf": Style{
		Color: Cyan,
		Icon:  "\uf031",
	},
	"woff": Style{
		Color: Cyan,
		Icon:  "\uf031",
	},
	"woff2": Style{
		Color: Cyan,
		Icon:  "\uf031",
	},
	"eot": Style{
		Color: Cyan,
		Icon:  "\uf031",
	},
	"properties": Style{
		Color: conf,
		Icon:  "\uE60B",
	},
	"asm": Style{
		Color: lang,
		Icon:  "\uFB19",
	},
	"groovy": Style{
		Color: lang,
		Icon:  "\ue775",
	},
	"s": Style{
		Color: lang,
		Icon:  "\ue637",
	},
	"gv": Style{
		Color: Cyan,
		Icon:  "\uE225",
	},
	"hs": Style{
		Color: lang,
		Icon:  "\ue777",
	},
	"d": Style{
		Color: lang,
		Icon:  "\ue7af",
	},
	"dart": Style{
		Color: lang,
		Icon:  "\uE798",
	},
	"erl": Style{
		Color: lang,
		Icon:  "\ue7b1",
	},
	"hrl": Style{
		Color: lang,
		Icon:  "\uE7B1",
	},
	"coffee": Style{
		Color: lang,
		Icon:  "\uE751",
	},
	"iso": Style{
		Color: iso,
		Icon:  "\uE271",
	},
	"dmg": Style{
		Color: iso,
		Icon:  "\uE271",
	},
	"img": Style{
		Color: iso,
		Icon:  "\uE271",
	},
	"vhd": Style{
		Color: iso,
		Icon:  "\uE271",
	},
	"vhdx": Style{
		Color: iso,
		Icon:  "\uE271",
	},
	"lua": Style{
		Color: lang,
		Icon:  "\uE620",
	},
	"vue": Style{
		Color: lang,
		Icon:  "\uE6A0",
	},
	"php": Style{
		Color: lang,
		Icon:  "\uE73D",
	},
	"cmake": Style{
		Color: lang,
		Icon:  "\uE20F",
	},
	"less": Style{
		Color: lang,
		Icon:  "\ue758",
	},
	"r": Style{
		Color: lang,
		Icon:  "\uF25D",
	},
	"f": Style{
		Color: lang,
		Icon:  "\uf121",
	},
	"atom": Style{
		Color: lang,
		Icon:  "\ue764",
	},
	"docker": Style{
		Color: Yellow,
		Icon:  "\ue7b0",
	},
	"md": Style{
		Color: Underline + BrightYellow,
		Icon:  "\uF48A",
	},
	"txt": Style{
		Color: text,
		Icon:  "\uF15C",
	},
	"zig": Style{
		Color: lang,
		Icon:  "\uF0E7",
	},
	"rss": Style{
		Color: lang,
		Icon:  "\uF09E",
	},
	"ko": Style{
		Color: lang,
		Icon:  "\uebc6",
	},
	"zip": Style{
		Color: tar,
		Icon:  "",
	},
	"zst": Style{
		Color: tar,
		Icon:  "",
	},
	"sitx": Style{
		Color: tar,
		Icon:  "",
	},
	"tar": Style{
		Color: tar,
		Icon:  "",
	},
	"gz": Style{
		Color: tar,
		Icon:  "",
	},
	"rar": Style{
		Color: tar,
		Icon:  "",
	},
	"7z": Style{
		Color: tar,
		Icon:  "",
	},
	"xz": Style{
		Color: tar,
		Icon:  "",
	},
	"jpg": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"jpeg": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"png": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"bmp": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"tif": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"tiff": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"gif": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"svg": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"webp": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"pcx": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"avif": Style{
		Color: pic,
		Icon:  "\uF1C5",
	},
	"psd": Style{
		Color: pic,
		Icon:  "\ue67f",
	},
	"mp4": Style{
		Color: video,
		Icon:  "",
	},
	"mkv": Style{
		Color: video,
		Icon:  "",
	},
	"avi": Style{
		Color: video,
		Icon:  "",
	},
	"m4v": Style{
		Color: video,
		Icon:  "",
	},
	"flv": Style{
		Color: video,
		Icon:  "",
	},
	"mov": Style{
		Color: video,
		Icon:  "",
	},
	"mpeg": Style{
		Color: video,
		Icon:  "",
	},
	"ai": Style{
		Color: Cyan,
		Icon:  "\uE7B4",
	},
	"mp3": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"m4a": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"mid": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"midi": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"flac": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"ape": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"alac": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"aac": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"wav": Style{
		Color: audio,
		Icon:  "\uf1c7",
	},
	"ogg": Style{
		Color: audio,
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
		Color: doc,
		Icon:  "\uf1c2",
	},
	"docx": Style{
		Color: doc,
		Icon:  "\uf1c2",
	},
	"docm": Style{
		Color: doc,
		Icon:  "\uf1c2",
	},
	"xls": Style{
		Color: doc,
		Icon:  "\uf1c3",
	},
	"xlsx": Style{
		Color: doc,
		Icon:  "\uf1c3",
	},
	"xlsm": Style{
		Color: doc,
		Icon:  "\uf1c3",
	},
	"numbers": Style{
		Color: doc,
		Icon:  "\uf1c3",
	},
	"csv": Style{
		Color: doc,
		Icon:  "\ue64a",
	},
	"ppt": Style{
		Color: doc,
		Icon:  "\uf1c4",
	},
	"pptm": Style{
		Color: doc,
		Icon:  "\uf1c4",
	},
	"pptx": Style{
		Color: doc,
		Icon:  "\uf1c4",
	},
	"dot": Style{
		Color: doc,
		Icon:  "\uf1c2",
	},
	"dotx": Style{
		Color: doc,
		Icon:  "\uf1c2",
	},
	"diff": Style{
		Color: Cyan,
		Icon:  "\uF440",
	},
	"ini": Style{
		Color: conf,
		Icon:  "\uE615",
	},
	"conf": Style{
		Color: conf,
		Icon:  "\uE615",
	},
	"cfg": Style{
		Color: conf,
		Icon:  "\uE615",
	},
	"yml": Style{
		Color: conf,
		Icon:  "\uF481",
	},
	"yaml": Style{
		Color: conf,
		Icon:  "\uF481",
	},
	"tex": Style{
		Color: lang,
		Icon:  "\uF034",
	},
	"typ": Style{
		Color: lang,
		Icon:  "∫",
	},
	"toml": Style{
		Color: conf,
		Icon:  "\uE615",
	},
	"db": Style{
		Color: db,
		Icon:  "\uF1C0",
	},
	"accdb": Style{
		Color: db,
		Icon:  "\uF1C0",
	},
	"sqlite": Style{
		Color: db,
		Icon:  "\uE7C4",
	},
	"sqlite3": Style{
		Color: db,
		Icon:  "\uE7C4",
	},
	"sql": Style{
		Color: db,
		Icon:  "\uE706",
	},
	"db3": Style{
		Color: db,
		Icon:  "\uF1C0",
	},
	"sh": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"bat": Style{
		Color: bash,
		Icon:  "\ue629",
	},
	"ps1": Style{
		Color: bash,
		Icon:  "\ue683",
	},
	"csh": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"fish": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"ksh": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"zsh": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"zsh_history": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"zshrc": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"bash": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"bashrc": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"bash_history": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"bash_profile": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"nu": Style{
		Color: bash,
		Icon:  "\uF489",
	},
	"vim": Style{
		Color: vim,
		Icon:  "\uE62B",
	},
	"viminfo": Style{
		Color: vim,
		Icon:  "\uE62B",
	},
	"vimrc": Style{
		Color: vim,
		Icon:  "\uE62B",
	},
	"gvimrc": Style{
		Color: vim,
		Icon:  "\uE62B",
	},
	"nvim": Style{
		Color: vim,
		Icon:  "\uE62B",
	},
	"nvimrc": Style{
		Color: vim,
		Icon:  "\uE62B",
	},
	"log": Style{
		Color: White,
		Icon:  "\uF18D",
	},
	"lock": Style{
		Color: lock,
		Icon:  "\uF023",
	},
	"code-workspace": Style{
		Color: lang,
		Icon:  "\uE70C",
	},
	"key": Style{
		Color: key,
		Icon:  "\uF084",
	},
	"pub": Style{
		Color: key,
		Icon:  "\uF084",
	},
	"id_rsa": Style{
		Color: key,
		Icon:  "\uF084",
	},
	"gpg": Style{
		Color: key,
		Icon:  "\uF084",
	},
	"cer": Style{
		Color: key,
		Icon:  "\uF084",
	},
	"crt": Style{
		Color: key,
		Icon:  "\uF084",
	},
	"pgp": Style{
		Color: key,
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
	"so": Style{
		Color: Cyan,
		Icon:  "\uF121",
	},
	"dll": Style{
		Color: Cyan,
		Icon:  "\uF121",
	},
}

// todo
// permission: match file permission
// size   : match file size unit
// user	  : match user name
// group  : match group name
// symlink: match symlink
// git    : match git status
// name   : match file name
// special: match file type: symlink, dir, executable
// ext    : match file extension

var Special = map[string]Style{
	"dir": {
		Color: dir,
		Icon:  "\uF115",
	},
	"pipe": {
		Color: Cyan,
		Icon:  "\uF124",
	},
	"socket": {
		Color: Cyan,
		Icon:  "\uF1E6",
	},
}

var Name = map[string]Style{
	"home": {
		Icon: "\ue617",
	},
	"desktop": {
		Icon: "\uf108",
	},
	"downloads": {
		Icon: "\uf498",
	},
	"trash": {
		Icon: "\uf1f8",
	},
	"searches": {
		Icon: "\uf422",
	},
	"microsoft": {
		Icon: "\uF17A",
	},
	"google": {
		Icon: "\uf1a0",
	},
	"onedrive": {
		Icon: "\ue762",
	},
	"onedrivetemp": {
		Icon: "\ue762",
	},
	"favorites": {
		Icon: "\ue623",
	},
	"azure": {
		Icon: "\uebd8",
	},
	"contacts": {
		Icon: "\uf0c0",
	},
	"users": {
		Icon: "\uf0c0",
	},
	"lib": {
		Icon: "\uF121",
	},
	"library": {
		Icon: "\uF121",
	},
	"bin": {
		Icon: "\uE5FC",
	},
	"share": {
		Icon: "\uf064",
	},
	"license": {
		Color: key,
		Icon:  "\uF084",
	},
	"shell": {
		Icon: "\uF489",
	},
	"config": {
		Color: conf,
		Icon:  "\uE615",
	},
	"ds_store": {
		Color: Black,
		Icon:  "\uf179",
	},
	"pkgbuild": {
		Color: pkg,
		Icon:  "\uf303",
	},
	".srcinfo": {
		Color: pkg,
		Icon:  "\uf303",
	},
	"applications": {
		Icon: "\ueb44",
	},
	"android": {
		Color: Cyan,
		Icon:  "\uE70E",
	},
	".idea": {
		Icon: "\uE7B5",
	},
	".github": {
		Icon: "\uF408",
	},
	".vscode": {
		Icon: "\uE70C",
	},
	"include": {
		Icon: "\ue5fc",
	},
	"dotnet": {
		Icon: "\ue72e",
	},
	"src": {
		Icon: "\ue796",
	},
	"node_modules": {
		Icon: "\ue5fa",
	},
	"font": {
		Icon: "\uf031",
	},
	"fonts": {
		Icon: "\uf031",
	},
	".git": {
		Color: dir,
		Icon:  "\uF1D3",
	},
	".gitignore": {
		Color: Black,
		Icon:  "\ue65d",
	},
	"cmakeLists.txt": {
		Color: lang,
		Icon:  "\uE20F",
	},
	"makefile": {
		Color: lang,
		Icon:  "\uE20F",
	},
	"justfile": {
		Color: lang,
		Icon:  "\uE20F",
	},
	"history": {
		Icon: "\uF1DA",
	},
	"recovery": {
		Icon: "\uF1DA",
	},
	"apple": {
		Color: Black,
		Icon:  "\uF179",
	},
	"dockerfile": {
		Color: Yellow,
		Icon:  "\ue7b0",
	},
	"readme": {
		Color: Underline + BrightYellow,
		Icon:  "\uF48A",
	},
	"readme.md": {
		Color: Underline + BrightYellow,
		Icon:  "\uF48A",
	},
	"sdk": {
		Icon: "\uF121",
	},
	"pictures": {
		Icon: "\uF1C5",
	},
	"videos": {
		Icon: "",
	},
	"movies": {
		Icon: "",
	},
	"music": {
		Icon: "\uF025",
	},
	"audio": {
		Icon: "\uF025",
	},
}

var Permisson = map[string]Style{
	"D": {
		Color: Cyan,
	},
	"d": {
		Color: Blue,
	},
	"l": {
		Color: Purple,
	},
	"b": {
		Color: Yellow,
	},
	"c": {
		Color: Yellow,
	},
	"p": {
		Color: Yellow,
	},
	"s": {
		Color: Yellow,
	},
	"r": {
		Color: Yellow,
	},
	"w": {
		Color: Red,
	},
	"x": {
		Color: Green,
	},
	"-": {
		Color: White,
	},
	"octal": {
		Color: color256(208),
	},
}

var Size = map[string]Style{
	"bit": {
		Color: rgb(20, 255, 100),
	},
	"B": {
		Color: rgb(20, 230, 100),
	},
	"KB": {
		Color: rgb(20, 207, 100),
	},
	"MB": {
		Color: rgb(20, 188, 100),
	},
	"GB": {
		Color: rgb(20, 170, 100),
	},
	"TB": {
		Color: rgb(20, 153, 100),
	},
	"PB": {
		Color: rgb(20, 138, 100),
	},
	"EB": {
		Color: rgb(20, 125, 100),
	},
	"ZB": {
		Color: rgb(20, 112, 100),
	},
	"YB": {
		Color: rgb(20, 100, 100),
	},
	"BB": {
		Color: rgb(20, 90, 100),
	},
	"NB": {
		Color: rgb(20, 70, 100),
	},
}

var Git = map[string]Style{
	"git_unmodified": {
		Color: Black,
	},
	"git_modified": {
		Color: Yellow,
	},
	"git_renamed": {
		Color: Blue,
	},
	"git_copied": {
		Color: Purple,
	},
	"git_deleted": {
		Color: Red,
	},
	"git_added": {
		Color: Green,
	},
	"git_untracked": {
		Color: BrightBlack,
	},
	"git_ignored": {
		Color: BrightRed,
	},
	"git_type_changed": {
		Color: Yellow,
	},
	"git_updated_but_unmerged": {
		Color: BrightYellow,
	},
}

var Owner = map[string]Style{
	"owner": {
		Color: Yellow,
	},
}

var Group = map[string]Style{
	"group": {
		Color: Yellow,
	},
}

var Symlink = map[string]Style{
	"symlink_path": {
		Color: Green,
	},
	"symlink_broken_path": {
		Color: Underline + Red,
	},
	"symlink": {
		Color: Purple,
		Icon:  "\ueb15",
	},
	"symlink_arrow": {
		Color: BrightWhite,
		Icon:  " ~> ",
	},
}
