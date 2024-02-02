package shell

import (
	"bytes"
	_ "embed"
	"fmt"

	constval "github.com/Equationzhao/g/internal/const"
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

func Init(shell string) (string, error) {
	switch shell {
	case "zsh":
		return string(bytes.ReplaceAll(ZSHContent, []byte("\r\n"), []byte("\n"))), nil
	case "bash":
		return string(bytes.ReplaceAll(BASHContent, []byte("\r\n"), []byte("\n"))), nil
	case "fish":
		return string(bytes.ReplaceAll(FISHContent, []byte("\r\n"), []byte("\n"))), nil
	case "powershell", "pwsh":
		return string(bytes.ReplaceAll(PSContent, []byte("\r\n"), []byte("\n"))), nil
	case "nushell", "nu":
		return string(bytes.ReplaceAll(NUContent, []byte("\r\n"), []byte("\n"))), nil
		// replace os newline with unix newline
	}
	return "", fmt.Errorf("unsupported shell: %s \n %s[zsh|bash|fish|powershell|nushell]", shell, constval.Success)
}
