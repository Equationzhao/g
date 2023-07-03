package shell

import _ "embed"

//go:embed g.ps1
var PSContent []byte

//go:embed g.bash
var BASHContent []byte

//go:embed g.zsh
var ZSHContent []byte

//go:embed g.fish
var FISHContent []byte
