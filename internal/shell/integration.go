package shell

import (
	"bytes"
	_ "embed"
)

//go:embed g.ps1
var PSContent []byte

//go:embed g.bash
var BASHContent []byte

//go:embed g.zsh
var ZSHContent []byte

//go:embed g.fish
var FISHContent []byte

//go:embed g.nu
var NUContent []byte

func Init() {
	// replace os newline with unix newline
	PSContent = bytes.ReplaceAll(PSContent, []byte("\r\n"), []byte("\n"))
	BASHContent = bytes.ReplaceAll(BASHContent, []byte("\r\n"), []byte("\n"))
	ZSHContent = bytes.ReplaceAll(ZSHContent, []byte("\r\n"), []byte("\n"))
	FISHContent = bytes.ReplaceAll(FISHContent, []byte("\r\n"), []byte("\n"))
	NUContent = bytes.ReplaceAll(NUContent, []byte("\r\n"), []byte("\n"))
}
