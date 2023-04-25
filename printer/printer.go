package printer

import (
	"bufio"
	"github.com/acarl005/stripansi"
	"github.com/olekukonko/ts"
	"io"
	"math"
	"os"
	"strings"
	"unicode/utf8"
)

var Output io.Writer = os.Stdout

// print style control

type Printer interface {
	Print(s ...string)
}

type Byline struct {
	*bufio.Writer
}

func NewByline() Printer {
	return &Byline{bufio.NewWriter(Output)}
}

func (b *Byline) Print(s ...string) {
	for _, v := range s {
		_, _ = b.WriteString(v)
		_ = b.WriteByte('\n')
	}
	_ = b.Flush()
}

// Modified from github.com/acarl005/textcol

type FitTerminal struct {
	*bufio.Writer
}

func NewFitTerminal() Printer {
	return &FitTerminal{Writer: bufio.NewWriter(Output)}
}

func (f *FitTerminal) Print(s ...string) {
	f.printColumns(&s, 4)
}

func (f *FitTerminal) printColumns(strs *[]string, margin int) {
	defer f.Flush()

	maxLength := 0
	marginStr := strings.Repeat(" ", margin)
	// also keep track of each individual length to easily calculate padding
	var lengths []int
	for _, str := range *strs {
		colorless := stripansi.Strip(str)
		// len() is insufficient here, as it counts emojis as 4 characters each
		length := utf8.RuneCountInString(colorless)
		maxLength = max(maxLength, length)
		lengths = append(lengths, length)
	}

	// see how wide the terminal is
	width := getTermWidth()
	// calculate the dimensions of the columns
	numCols, numRows := calculateTableSize(width, margin, maxLength, len(*strs))

	// if we're forced into a single column, fall back to simple printing (one per line)
	if numCols == 1 {
		for _, str := range *strs {
			_, _ = f.WriteString(str)
		}
		return
	}

	// `i` will be a left-to-right index. this will need to get converted to a top-to-bottom index
	for i := 0; i < numCols*numRows; i++ {
		// treat output like a "table" with (x, y) coordinates as an intermediate representation
		// first calculate (x, y) from i
		x, y := rowIndexToTableCoords(i, numCols)
		// then convey (x, y) to `j`, the top-to-bottom index
		j := tableCoordsToColIndex(x, y, numRows)

		// try to access the array, but the table might have more cells than array elements, so only try to access if within bounds
		strLen := 0
		str := ""
		if j < len(lengths) {
			strLen = lengths[j]
			str = (*strs)[j]
		}

		// calculate the amount of padding required
		numSpacesRequired := maxLength - strLen
		spaceStr := strings.Repeat(" ", numSpacesRequired)

		// print the item itself
		_, _ = f.WriteString(str)
		// if we're at the last column, print a line break
		if x+1 == numCols {
			_ = f.WriteByte('\n')
		} else {
			_, _ = f.WriteString(spaceStr)
			_, _ = f.WriteString(marginStr)
		}
	}
}

// getTermWidth returns the width of the terminal in characters
// this is a modified version
func getTermWidth() int {
	size, _ := ts.GetSize()
	return size.Col()
}

/*
	original version

	func getTermWidth() int {
		cmd := exec.Command("stty", "size")
		cmd.Stdin = os.Stdin
		out, err1 := cmd.Output()
		check(err1)
		numsStr := strings.Trim(string(out), "\n ")
		width, err2 := strconv.Atoi(strings.Split(numsStr, " ")[1])
		check(err2)
		return width
	}
*/

func calculateTableSize(width, margin, maxLength, numCells int) (int, int) {
	numCols := (width + margin) / (maxLength + margin)
	if numCols == 0 {
		numCols = 1
	}
	numRows := int(math.Ceil(float64(numCells) / float64(numCols)))
	return numCols, numRows
}

func rowIndexToTableCoords(i, numCols int) (int, int) {
	x := i % numCols
	y := i / numCols
	return x, y
}

func tableCoordsToColIndex(x, y, numRows int) int {
	return y + numRows*x
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
