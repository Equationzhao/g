package util

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/mattn/go-runewidth"

	"github.com/Equationzhao/g/internal/global"
)

func RemoveSep(s string) string {
	return strings.TrimRight(s, string(filepath.Separator))
}

var escapeReplacer = strings.NewReplacer(
	"\t", reverseColor(`\t`),
	"\r", reverseColor(`\r`),
	"\n", reverseColor(`\n`),
	"\"", reverseColor(`\"`),
	"\\", reverseColor(`\\`),
)

func reverseColor(s string) string {
	return global.Reverse + s + global.ReverseDone
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

func MakeLink(abs string, name string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", abs, name)
}

// SplitNumberAndUnit splits a string like
// "10bit" to 10 and "bit"
//
//	"12.3ml" to 12.4 and "ml"
//
// "-1,234,213kg" to -1234213 and "kg"
func SplitNumberAndUnit(input string) (float64, string) {
	var number float64
	var unit string

	// Find the index of the first non-digit character
	i := 0
	for i < len(input) && (input[i] >= '0' && input[i] <= '9' || input[i] == '.' || input[i] == '-' || input[i] == ',') {
		i++
	}

	// Parse the number part
	numberPart := input[:i]
	number, _ = strconv.ParseFloat(strings.ReplaceAll(numberPart, ",", ""), 64)

	// Extract the unit part
	unit = input[i:]

	return number, unit
}

func parseLink(link string) (name, other string, ok bool) {
	re := regexp.MustCompile(`\033]8;;(.*?)\033\\(.*?)\033]8;;\033\\`)
	matches := re.FindStringSubmatch(link)

	if len(matches) == 3 {
		// if matches, get the other content and the link
		other = strings.Replace(link, matches[0], "", 1)
		return matches[2], other, true
	}
	return "", "", false
}

func WidthLen(str string) int {
	if global.IncludeHyperlink {
		name, other, ok := parseLink(str)
		if ok {
			str = other + name
		}
	}
	colorless := stripansi.Strip(str)
	// len() is not proper here, as it counts emojis as 4 characters each
	length := runewidth.StringWidth(colorless)

	return length
}

func WidthNoHyperLinkLen(str string) int {
	colorless := stripansi.Strip(str)
	// len() is insufficient here, as it counts emojis as 4 characters each
	length := runewidth.StringWidth(colorless)
	return length
}
