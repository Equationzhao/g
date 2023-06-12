package display

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/Equationzhao/g/util"
	"github.com/acarl005/stripansi"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mattn/go-runewidth"
	"github.com/olekukonko/ts"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

const dot = '\uF111'

var Output io.Writer = os.Stdout

func RawPrint(toPrint ...any) (n int, err error) {
	return fmt.Fprint(Output, toPrint...)
}

// print style control

type hook struct {
	BeforePrint   []func(...Item)
	AfterPrint    []func(...Item)
	disableBefore bool
	disableAfter  bool
}

func fire(h []func(...Item), i ...Item) {
	for _, fn := range h {
		if fn == nil {
			continue
		}
		fn(i...)
	}
}

func newHook() *hook {
	return &hook{
		BeforePrint: make([]func(...Item), 0, 5),
		AfterPrint:  make([]func(...Item), 0, 5),
	}
}

func (h *hook) AddBeforePrint(f ...func(...Item)) {
	h.BeforePrint = append(h.BeforePrint, f...)
}

func (h *hook) AddAfterPrint(f ...func(...Item)) {
	h.AfterPrint = append(h.AfterPrint, f...)
}

func (h *hook) DisableHookBefore() {
	h.disableBefore = true
}

func (h *hook) EnableHookBefore() {
	h.disableBefore = false
}

func (h *hook) DisableHookAfter() {
	h.disableAfter = true
}

func (h *hook) EnableHookAfter() {
	h.disableAfter = false
}

type Hook interface {
	AddBeforePrint(...func(...Item))
	AddAfterPrint(...func(...Item))
	DisableHookBefore()
	EnableHookBefore()
	DisableHookAfter()
	EnableHookAfter()
}

type Printer interface {
	Print(s ...Item)
	Hook
}

type Byline struct {
	*bufio.Writer
	*hook
}

func NewByline() Printer {
	return &Byline{
		Writer: bufio.NewWriter(Output),
		hook:   newHook(),
	}
}

func (b *Byline) Print(i ...Item) {
	if !b.disableBefore {
		fire(b.BeforePrint, i...)
	}
	defer b.Flush()
	for _, v := range i {
		_, _ = b.WriteString(v.OrderedContent())
		_ = b.WriteByte('\n')
	}
	if !b.disableAfter {
		fire(b.AfterPrint, i...)
	}
}

// Modified from github.com/acarl005/textcol

type FitTerminal struct {
	*bufio.Writer
	*hook
}

func NewFitTerminal() Printer {
	return &FitTerminal{
		Writer: bufio.NewWriter(Output),
		hook:   newHook(),
	}
}

func (f *FitTerminal) Print(i ...Item) {
	if !f.disableBefore {
		fire(f.BeforePrint, i...)
	}

	s := make([]string, 0, len(i))
	for _, v := range i {
		s = append(s, v.OrderedContent())
	}
	f.printColumns(&s)

	if !f.disableAfter {
		fire(f.AfterPrint, i...)
	}
}

func (f *FitTerminal) printColumns(strs *[]string) {
	defer f.Flush()

	maxLength, lengths, numCols, numRows := calculateRowCol(strs, 6)

	// if we're forced into a single column, fall back to simple printing (one per line)
	if numCols == 1 {
		for _, str := range *strs {
			_, _ = f.WriteString(str)
			_ = f.WriteByte('\n')
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
		spaceStr := strings.Repeat(" ", numSpacesRequired+1)

		// print the item itself
		_, _ = f.WriteString(str)
		// if we're at the last column, print a line break
		if x+1 == numCols {
			_ = f.WriteByte('\n')
		} else {
			_, _ = f.WriteString(spaceStr)
		}
	}
}

// maxLength, maxLength, numCols, numRows
func calculateRowCol(strs *[]string, margin int) (maxLength int, lengths []int, numCols int, numRows int) {
	// also keep track of each individual length to easily calculate padding
	lengths = make([]int, 0, len(*strs))
	for _, str := range *strs {
		length := WidthLen(str)
		maxLength = max(maxLength, length)
		lengths = append(lengths, length)
	}

	// see how wide the terminal is
	width := getTermWidth()
	// calculate the dimensions of the columns
	numCols, numRows = calculateTableSize(width, margin, maxLength, len(*strs))
	return
}

func WidthLen(str string) int {
	colorless := stripansi.Strip(str)
	// len() is insufficient here, as it counts emojis as 4 characters each
	length := runewidth.StringWidth(colorless)
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		if strings.ContainsRune(colorless, dot) {
			length--
		}
	}
	return length
}

var (
	getTermWidthOnce util.Once
	size             ts.Size
)

// getTermWidth returns the width of the terminal in characters
// this is a modified version
func getTermWidth() int {
	if err := getTermWidthOnce.Do(func() error {
		var err error
		size, err = ts.GetSize()
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0
	}
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
	numCols := (width + margin) / (maxLength + 1)
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

type CommaPrint struct {
	*Across
	*hook
}

func NewCommaPrint() Printer {
	a := NewAcross()
	return &CommaPrint{
		Across: a.(*Across),
		hook:   newHook(),
	}
}

func (c *CommaPrint) Print(items ...Item) {
	if !c.disableBefore {
		fire(c.BeforePrint, items...)
	}
	s := make([]string, 0, len(items))
	for i, v := range items {
		if i != len(items)-1 {
			s = append(s, v.OrderedContent()+",")
		} else {
			s = append(s, v.OrderedContent())
		}
	}
	c.printRowWithNoSpace(&s)
	if !c.disableAfter {
		fire(c.AfterPrint, items...)
	}
}

type Across struct {
	*bufio.Writer
	*hook
}

func NewAcross() Printer {
	return &Across{
		Writer: bufio.NewWriter(Output),
		hook:   newHook(),
	}
}

func (a *Across) Print(items ...Item) {
	if !a.disableBefore {
		fire(a.BeforePrint, items...)
	}
	s := make([]string, 0, len(items))
	for _, v := range items {
		s = append(s, v.OrderedContent())
	}
	a.printRow(&s)
	if !a.disableAfter {
		fire(a.AfterPrint, items...)
	}
}

func (a *Across) printRowWithNoSpace(strs *[]string) {
	defer a.Flush()
	width := getTermWidth()

	maxLength := 0
	for _, str := range *strs {
		colorless := stripansi.Strip(str)
		maxLength += runewidth.StringWidth(stripansi.Strip(str))
		if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
			if strings.ContainsRune(colorless, dot) {
				maxLength--
			}
		}

		if maxLength <= width {
			_, _ = a.WriteString(str)
		} else {
			_, _ = a.WriteString("\n" + str)
			maxLength = runewidth.StringWidth(colorless)
			if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
				if strings.ContainsRune(colorless, dot) {
					maxLength--
				}
			}
		}
	}
	_ = a.WriteByte('\n')
}

func (a *Across) printRow(strs *[]string) {
	defer a.Flush()
	width := getTermWidth()

	const m = 1
	strLen := make([]int, len(*strs))

	maxLength := 0
	for i, str := range *strs {
		colorless := stripansi.Strip(str)
		strLen[i] = runewidth.StringWidth(colorless)
		if runtime.GOOS == "windows" {
			if strings.ContainsRune(colorless, dot) {
				strLen[i]--
			}
		}
		maxLength = max(maxLength, strLen[i])
	}

	cols := (width + m) / (maxLength + m)
	if cols == 0 {
		cols = 1
	}

	colWidth := (width+m)/cols - m

	for i := 0; i < len(*strs); i += cols {
		for j := 0; j < cols && i+j < len(*strs); j++ {
			index := i + j
			str := (*strs)[index]
			padding := colWidth - strLen[index]
			if padding < 0 {
				padding = 0
			}
			if j < cols-1 {
				_, _ = fmt.Fprintf(a, "%s%s", str, a.stringOf(' ', padding+m))
			} else {
				_, _ = fmt.Fprintf(a, "%s%s", str, a.stringOf(' ', padding))
			}
		}
		_ = a.WriteByte('\n')
	}
}

func (a *Across) stringOf(ch rune, count int) string {
	s := make([]rune, count)
	for i := 0; i < count; i++ {
		s[i] = ch
	}
	return string(s)
}

type Zero struct {
	*bufio.Writer
	*hook
}

func NewZero() Printer {
	return &Zero{
		Writer: bufio.NewWriter(Output),
		hook:   newHook(),
	}
}

func (z *Zero) Print(items ...Item) {
	if !z.disableBefore {
		fire(z.BeforePrint, items...)
	}
	defer z.Flush()
	for _, v := range items {
		v.Delimiter = ""
		_, _ = z.WriteString(v.OrderedContent())
	}
	if !z.disableAfter {
		fire(z.AfterPrint, items...)
	}
}

type JsonPrinter struct {
	*bufio.Writer
	*hook
}

func (j *JsonPrinter) pretty(data []byte) (string, error) {
	b := &bytes.Buffer{}
	err := json.Indent(b, data, "", "	")
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func NewJsonPrinter() Printer {
	return &JsonPrinter{
		Writer: bufio.NewWriter(Output),
		hook:   newHook(),
	}
}

func (j *JsonPrinter) Print(items ...Item) {
	if !j.disableBefore {
		fire(j.BeforePrint, items...)
	}
	defer j.Flush()

	list := make([]*orderedmap.OrderedMap[string, string], 0, len(items))
	for _, v := range items {
		all := v.GetAll()

		type orderItem struct {
			name    string
			content string
			no      int
		}

		order := make([]orderItem, 0, len(all))

		// sort by v.Content.No
		for name, v := range all {
			c := stripansi.Strip(v.Content.String())
			if name == "name" {
				order = append(order, orderItem{name: name, content: c, no: v.No})
			} else if name == "underwent" || name == "statistic" {
				order = append(order, orderItem{name: name, content: strings.TrimLeft(c, "\n "), no: v.No})
			} else if name == "total" {
				order = append(order, orderItem{name: name, content: strings.TrimPrefix(c, "  total "), no: v.No})
			} else {
				// remove all leading spaces
				order = append(order, orderItem{name: name, content: strings.TrimLeft(c, " "), no: v.No})
			}
		}

		sort.Slice(order, func(i, j int) bool {
			return order[i].no < order[j].no
		})

		s := orderedmap.New[string, string](
			orderedmap.WithCapacity[string, string](len(order)),
		)

		list = append(list, s)

		for _, v := range order {
			s.Set(v.name, v.content)
		}
	}
	pretty, err := json.MarshalIndent(list, "", "	")
	if err != nil {
		_, _ = j.WriteString(err.Error() + "\n")
		return
	}
	_, _ = j.Write(pretty)
	_, _ = j.WriteString("\n")
	if !j.disableAfter {
		fire(j.AfterPrint, items...)
	}
}

type PrettyPrinter interface {
	SetTitle(title string)
	AddHeader(headers string)
	Printer
}

type TablePrinter struct {
	*bufio.Writer
	*hook
	header table.Row
	w      table.Writer
}

func (t *TablePrinter) SetTitle(title string) {
	t.w.SetTitle("path: %s", title)
}

func (t *TablePrinter) AddHeader(headers string) {
	t.header = append(t.header, headers)
}

func NewTablePrinter(opts ...func(writer table.Writer)) Printer {
	t := &TablePrinter{
		Writer: bufio.NewWriter(Output),
		hook:   newHook(),
	}
	w := table.NewWriter()
	for _, opt := range opts {
		opt(w)
	}
	w.SetOutputMirror(t.Writer)
	t.w = w
	return t
}

func (t *TablePrinter) PrintBase(fn func() string, s ...Item) {
	if !t.disableBefore {
		fire(t.BeforePrint, s...)
	}
	defer t.Writer.Flush()
	t.w.ResetRows()
	t.setTB(s...)
	t.w.AppendHeader(t.header)
	fn()
	t.w.ResetHeaders()

	// empty header
	t.header = t.header[:0]

	if !t.disableAfter {
		fire(t.AfterPrint, s...)
	}
}

func (t *TablePrinter) Print(s ...Item) {
	t.PrintBase(t.w.Render, s...)
}

func (t *TablePrinter) setTB(s ...Item) {
	for _, v := range s {
		all := v.GetAllOrdered()
		row := make(table.Row, 0, len(all))
		for _, v := range all {
			row = append(row, strings.TrimLeft(v.Content.String(), " "))
		}
		t.w.AppendRow(row)
	}
}

func DefaultTB(w table.Writer) {
	w.SetAllowedRowLength(getTermWidth())
	if runtime.GOOS != "windows" {
		w.SetStyle(table.StyleRounded)
	}
	w.Style().Options.SeparateColumns = true
	w.Style().Options.SeparateFooter = true
	w.SetPageSize(1000)
}

type MDPrinter struct {
	*TablePrinter
}

func NewMDPrinter() Printer {
	m := &MDPrinter{}
	m.TablePrinter = NewTablePrinter(DefaultTB).(*TablePrinter)
	return m
}

func (m *MDPrinter) Print(s ...Item) {
	m.PrintBase(m.w.RenderMarkdown, s...)
}

type HTMLPrinter struct {
	*TablePrinter
}

func NewHTMLPrinter() Printer {
	h := &HTMLPrinter{}
	h.TablePrinter = NewTablePrinter(DefaultTB).(*TablePrinter)
	return h
}

func (p *HTMLPrinter) Print(s ...Item) {
	p.PrintBase(p.w.RenderHTML, s...)
}

type CSVPrinter struct {
	*TablePrinter
}

func NewCSVPrinter() Printer {
	c := &CSVPrinter{}
	c.TablePrinter = NewTablePrinter(DefaultTB).(*TablePrinter)
	return c
}

func (c *CSVPrinter) Print(s ...Item) {
	c.PrintBase(c.w.RenderCSV, s...)
}
