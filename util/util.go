package util

import (
	"path/filepath"
	"strings"

	"github.com/Equationzhao/g/theme"
)

func RemoveSep(s string) string {
	return strings.TrimRight(s, string(filepath.Separator))
}

var escapeReplacer = strings.NewReplacer(
	"\t", reverseColor(`\t`),
	"\r", reverseColor(`\r`),
	"\n", reverseColor(`\n`),
	"\"", reverseColor(`\"`),
	"\\", reverseColor(`\`),
)

func reverseColor(s string) string {
	return theme.Reverse + s + theme.ReverseDone
}

// Escape
// * Tab is escaped as `\t`.
// * Carriage return is escaped as `\r`.
// * Line feed is escaped as `\n`.
// * Single quote is escaped as `\'`.
// * Double quote is escaped as `\"`.
// * Backslash is escaped as `\\`.
func Escape(a string) string {
	return escapeReplacer.Replace(a)
}
