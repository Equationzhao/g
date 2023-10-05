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

// permission: match file permission
// size: match file size unit
// user: match username
// group: match group name
// symlink: match symlink
// git: match git status
// name: match file name
// special: match file type: symlink, dir, executable
// ext: match file extension

type All struct {
	InfoTheme  Theme `json:"info,omitempty"`
	Permission Theme `json:"permission,omitempty"`
	Size       Theme `json:"size,omitempty"`
	User       Theme `json:"user,omitempty"`
	Group      Theme `json:"group,omitempty"`
	Symlink    Theme `json:"symlink,omitempty"`
	Git        Theme `json:"git,omitempty"`
	Name       Theme `json:"name,omitempty"`
	Special    Theme `json:"special,omitempty"`
	Ext        Theme `json:"ext,omitempty"`
}

func (a *All) Apply(f func(theme Theme)) {
	f(a.InfoTheme)
	f(a.Permission)
	f(a.Size)
	f(a.User)
	f(a.Group)
	f(a.Symlink)
	f(a.Git)
	f(a.Name)
	f(a.Special)
	f(a.Ext)
}

var InfoTheme = Theme{
	"inode": Style{
		Color: Purple,
	},
	"time": Style{
		Color: Blue,
	},
	"reset": Style{
		Color: Reset,
	},
	"-": {
		Color: White,
	},
}

var Ext = Theme{
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
	"rpm": Style{
		Color: pkg,
		Icon:  "",
	},
	"exe": Style{
		Color: Green,
		Icon:  "\uF17A",
	},
	"rdp": Style{
		Color: White,
		Icon:  "\ueb39",
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
		Icon:  "",
	},
	"hpp": Style{
		Color: lang,
		Icon:  "",
	},
	"h++": Style{
		Color: lang,
		Icon:  "",
	},
	"hh": Style{
		Color: lang,
		Icon:  "",
	},
	"hxx": Style{
		Color: lang,
		Icon:  "",
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
		Icon:  "\ue637",
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
	"clj": Style{
		Color: lang,
		Icon:  "",
	},
	"bzl": Style{
		Color: lang,
		Icon:  "",
	},
	"avro": Style{
		Color: lang,
		Icon:  "\ue60b",
	},
	"svelte": Style{
		Color: lang,
		Icon:  "",
	},
	"mustache": Style{
		Color: lang,
		Icon:  "\ue60f",
	},
	"sass": Style{
		Color: lang,
		Icon:  "\ue603",
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
	"twig": Style{
		Color: lang,
		Icon:  "",
	},
	"fs": Style{
		Color: lang,
		Icon:  "",
	},
	"ex": Style{
		Color: lang,
		Icon:  "",
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
		Color: BrightYellow,
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
	"config": Style{
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
	"jl": Style{
		Color: lang,
		Icon:  "",
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
	"rdb": Style{
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
	"patch": Style{
		Color: White,
		Icon:  "",
	},
	"torrent": Style{
		Color: Green,
		Icon:  "\ueac2",
	},
	"service": Style{
		Color: White,
		Icon:  "",
	},
}

var Special = map[string]Style{
	"dir": {
		Color: dir,
		Icon:  "\uF115",
	},
	"dir-prompt": {
		Color: Yellow,
		Icon:  "► ",
	},
	"pipe": {
		Color: Cyan,
		Icon:  "\uF124",
	},
	"socket": {
		Color: Cyan,
		Icon:  "\uF1E6",
	},
	"exe": {
		Color: Green,
		Icon:  "\uF17A",
	},
	"link": {
		Color: Purple,
	},
	"file": {
		Color: White,
		Icon:  "\uF016",
	},
	"mounts": {
		Color: BrightBlack,
	},
}

var Name = map[string]Style{
	"home": {
		Icon:  "\ue617",
		Color: White,
	},
	"desktop": {
		Icon:  "\uf108",
		Color: White,
	},
	"documents": {
		Icon:  "\uf0f6",
		Color: White,
	},
	"links": {
		Icon:  "\uf0f6",
		Color: White,
	},
	"downloads": {
		Icon:  "\uf498",
		Color: White,
	},
	"trash": {
		Icon:  "\uf1f8",
		Color: BrightBlack,
	},
	"searches": {
		Icon:  "\uf422",
		Color: White,
	},
	"microsoft": {
		Icon:  "\uF17A",
		Color: White,
	},
	"google": {
		Icon:  "\uf1a0",
		Color: White,
	},
	"onedrive": {
		Icon:  "\ue762",
		Color: Blue,
	},
	"onedrivetemp": {
		Icon:  "\ue762",
		Color: Blue,
	},
	"favorites": {
		Icon:  "\ue623",
		Color: Red,
	},
	"saved games": {
		Icon:  "\uf11b",
		Color: Red,
	},
	".wakatime": {
		Icon:  "\ue641",
		Color: White,
	},
	".azure": {
		Icon:  "\uebd8",
		Color: White,
	},
	"contacts": {
		Icon:  "\uf0c0",
		Color: White,
	},
	"users": {
		Icon:  "\uf0c0",
		Color: White,
	},
	"lib": {
		Icon:  "\uF121",
		Color: White,
	},
	"library": {
		Icon:  "\uF121",
		Color: White,
	},
	"bin": {
		Icon:  "\uE5FC",
		Color: Green,
	},
	"share": {
		Icon:  "\uf064",
		Color: White,
	},
	"license": {
		Icon:  "",
		Color: key,
	},
	"shell": {
		Icon:  "\uF489",
		Color: bash,
	},
	"config": {
		Icon:  "\uE615",
		Color: conf,
	},
	".ds_store": {
		Color: BrightBlack,
		Icon:  "\uf179",
	},
	".gitconfig": {
		Color: conf,
		Icon:  "\uf1d3",
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
		Icon:  "\ueb44",
		Color: Green,
	},
	"android": {
		Icon:  "\uE70E",
		Color: White,
	},
	".idea": {
		Icon:  "\uE7B5",
		Color: White,
	},
	".github": {
		Icon:  "\uF408",
		Color: White,
	},
	".vscode": {
		Icon:  "\uE70C",
		Color: White,
	},
	"include": {
		Icon:  "\ue5fc",
		Color: White,
	},
	".dotnet": {
		Icon:  "\ue72e",
		Color: White,
	},
	"src": {
		Icon:  "\ue796",
		Color: White,
	},
	"node_modules": {
		Icon:  "\ue5fa",
		Color: White,
	},
	"font": {
		Icon:  "\uf031",
		Color: White,
	},
	"fonts": {
		Icon:  "\uf031",
		Color: White,
	},
	".git": {
		Color: White,
		Icon:  "\uF1D3",
	},
	".gitignore": {
		Color: BrightBlack,
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
		Icon:  "\uF1DA",
		Color: BrightBlack,
	},
	".history": {
		Icon:  "\uF1DA",
		Color: BrightBlack,
	},
	"recovery": {
		Icon:  "\uF1DA",
		Color: BrightBlack,
	},
	"apple": {
		Color: BrightBlack,
		Icon:  "\uF179",
	},
	"dockerfile": {
		Color: BrightYellow,
		Icon:  "\ue7b0",
	},
	"readme": {
		Color:     Yellow,
		Icon:      "\uF48A",
		Underline: true,
		Bold:      true,
	},
	"readme.md": {
		Color:     Yellow,
		Icon:      "\uF48A",
		Underline: true,
		Bold:      true,
	},
	"sdk": {
		Icon:  "\uF121",
		Color: White,
	},
	"pictures": {
		Icon:  "\uF1C5",
		Color: White,
	},
	"videos": {
		Icon:  "",
		Color: White,
	},
	"movies": {
		Icon:  "",
		Color: White,
	},
	"music": {
		Icon:  "\uF025",
		Color: White,
	},
	"audio": {
		Icon:  "\uF025",
		Color: White,
	},
	"cargo.lock": {
		Color: lock,
		Icon:  "\uE7A8",
	},
	"cargo.toml": {
		Color: lang,
		Icon:  "\uE7A8",
	},
	"known_hosts": {
		Color: White,
		Icon:  "\uEB39",
	},
	"repo": {
		Color: White,
		Icon:  "\uea62",
	},
	".ssh": {
		Icon:  "\ueba9",
		Color: White,
	},
	"boot": {
		Icon:  "\uf287",
		Color: White,
	},
	"cache": {
		Icon:  "\uf49b",
		Color: White,
	},
	"passwd": {
		Color: Cyan,
		Icon:  "\uF023",
	},
	"vagrantfile": {
		Color: lang,
		Icon:  "\uE20F",
	},
}

var Permission = map[string]Style{
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
	"-": {
		Color: White,
	},
	"block": {
		Color: rgb(20, 255, 100),
	},
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
		Color: BrightBlack,
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
	"root": {
		Color: Red,
	},
}

var Group = map[string]Style{
	"group": {
		Color: Yellow,
	},
	"root": {
		Color: Red,
	},
}

var Symlink = map[string]Style{
	"symlink_path": {
		Color: Green,
	},
	"symlink_broken_path": {
		Color:     Red,
		Underline: true,
	},
	"symlink": {
		Color: Purple,
		Icon:  "\ueb15",
	},
	"symlink_arrow": {
		Color: BrightWhite,
		Icon:  " => ",
	},
}

var (
	DefaultAll All
	_init      bool
)

func init() {
	if !_init {
		DefaultAll = All{
			InfoTheme:  InfoTheme,
			Permission: Permission,
			Size:       Size,
			User:       Owner,
			Group:      Group,
			Symlink:    Symlink,
			Git:        Git,
			Name:       Name,
			Special:    Special,
			Ext:        Ext,
		}
	}
}
