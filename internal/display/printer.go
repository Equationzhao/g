// Package display
// This package controls the display format, like grid/across/byline/...
package display

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/Equationzhao/g/internal/config"
	"github.com/Equationzhao/g/internal/display/tree"
	"github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/olekukonko/ts"
	"github.com/valyala/bytebufferpool"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

var Output io.Writer = os.Stdout

// RawPrint used for print statistics
// like:
//
//	underwent 898 microseconds
//	statistic: 15 file(s), 5 dir(s), 0 link(s)
func RawPrint(toPrint ...any) (n int, err error) {
	return fmt.Fprint(Output, toPrint...)
}

// all printers contain a Hook.
// hook implements the Hook interface
type hook struct {
	BeforePrint   []func(Printer, ...*item.FileInfo)
	AfterPrint    []func(Printer, ...*item.FileInfo)
	disableBefore bool
	disableAfter  bool
}

// fire fires the hook 🔥
func fire(h []func(Printer, ...*item.FileInfo), p Printer, i ...*item.FileInfo) {
	for _, fn := range h {
		if fn == nil {
			continue
		}
		fn(p, i...)
	}
}

func newHook() *hook {
	return &hook{
		BeforePrint: make([]func(Printer, ...*item.FileInfo), 0, global.DefaultHookLen),
		AfterPrint:  make([]func(Printer, ...*item.FileInfo), 0, global.DefaultHookLen),
	}
}

func (h *hook) AddBeforePrint(f ...func(Printer, ...*item.FileInfo)) {
	h.BeforePrint = append(h.BeforePrint, f...)
}

func (h *hook) AddAfterPrint(f ...func(Printer, ...*item.FileInfo)) {
	h.AfterPrint = append(h.AfterPrint, f...)
}

func (h *hook) DisablePreHook() {
	h.disableBefore = true
}

func (h *hook) EnablePreHook() {
	h.disableBefore = false
}

func (h *hook) DisablePostHook() {
	h.disableAfter = true
}

func (h *hook) EnablePostHook() {
	h.disableAfter = false
}

// Hook is used to add pre-content and post-content
// like adding a header or footer:
//
//	Permissions Size Owner Group Time Modified Name
type Hook interface {
	AddBeforePrint(...func(Printer, ...*item.FileInfo))
	AddAfterPrint(...func(Printer, ...*item.FileInfo))
	DisablePreHook()
	EnablePreHook()
	DisablePostHook()
	EnablePostHook()
}

// Printer is the interface of all printers.
// all printers should implement this interface
type Printer interface {
	Print(s ...*item.FileInfo)
	Hook
	io.Writer
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

func (b *Byline) Print(i ...*item.FileInfo) {
	if !b.disableBefore {
		fire(b.BeforePrint, b, i...)
	}
	defer b.Flush()
	for _, v := range i {
		_, _ = b.WriteString(v.OrderedContent(" "))
		_ = b.WriteByte('\n') // byline means a new line :)
		if len(v.AfterLines) != 0 {
			items := v.ValuesByOrdered()
			res := bytebufferpool.Get()
			for j, it := range items[:len(items)-1] {
				_, _ = res.WriteString(it.String())
				if j != len(items)-1 {
					_, _ = res.WriteString(" ")
				}
			}
			space := util.WidthNoHyperLinkLen(res.String())
			bytebufferpool.Put(res)
			_, _ = b.WriteString(strings.Join(append(v.Extended(space), ""), "\n"))
		}
	}
	if !b.disableAfter {
		fire(b.AfterPrint, b, i...)
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

func (f *FitTerminal) Print(i ...*item.FileInfo) {
	if !f.disableBefore {
		fire(f.BeforePrint, f, i...)
	}
	defer f.Flush()
	s := make([]string, 0, len(i))
	for _, v := range i {
		s = append(s, v.OrderedContent(" "))
	}
	f.printColumns(s, global.Space)

	if !f.disableAfter {
		fire(f.AfterPrint, f, i...)
	}
}

func (f *FitTerminal) printColumns(stringsArray []string, space int) {
	termWidth := getTermWidth() - 5
	n := len(stringsArray)

	// Calculate the maximum width of each string
	maxWidth := 0
	widths := make([]int, n)
	for i, s := range stringsArray {
		width := util.WidthLen(s)
		widths[i] = width
		if width > maxWidth {
			maxWidth = width
		}
	}

	// Determine the number of columns
	numColumns := 1
	for col := 1; col <= n; col++ {
		totalWidth := -space
		for i := 0; i < col; i++ {
			colMaxWidth := 0
			for j := i; j < n; j += col {
				if widths[j] > colMaxWidth {
					colMaxWidth = widths[j]
				}
			}
			totalWidth += colMaxWidth + space
		}
		if totalWidth <= termWidth {
			numColumns = col
		} else {
			break
		}
	}

	numRows := (n + numColumns - 1) / numColumns

	for i := 0; i < numRows; i++ {
		for j := 0; j < numColumns; j++ {
			index := j*numRows + i
			if index >= n {
				break
			}
			s := stringsArray[index]
			width := widths[index]
			_, _ = f.WriteString(s)
			if j < numColumns-1 { // Avoid adding spaces after the last column
				maxColWidth := 0
				for k := j * numRows; k < (j+1)*numRows && k < n; k++ {
					if widths[k] > maxColWidth {
						maxColWidth = widths[k]
					}
				}
				padding := maxColWidth - width + space
				_, _ = f.WriteString(strings.Repeat(" ", padding))
			}
		}
		_ = f.WriteByte('\n')
	}
}

var (
	getTermWidthOnce util.Once
	size             ts.Size
	CustomTermSize   uint
)

// getTermWidth returns the width of the terminal in characters
// this is a modified version
func getTermWidth() int {
	if CustomTermSize != 0 {
		return int(CustomTermSize)
	}

	if err := getTermWidthOnce.Do(
		func() error {
			var err error
			size, err = ts.GetSize()
			if err != nil {
				return err
			}
			return nil
		},
	); err != nil {
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

func (c *CommaPrint) Print(items ...*item.FileInfo) {
	if !c.disableBefore {
		fire(c.BeforePrint, c, items...)
	}
	defer c.Flush()
	s := make([]string, 0, len(items))
	for i, v := range items {
		if i != len(items)-1 {
			s = append(s, v.OrderedContent(" ")+",")
		} else {
			s = append(s, v.OrderedContent(" "))
		}
	}
	c.printRowWithNoSpace(s)
	if !c.disableAfter {
		fire(c.AfterPrint, c, items...)
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

func (a *Across) Print(items ...*item.FileInfo) {
	if !a.disableBefore {
		fire(a.BeforePrint, a, items...)
	}
	defer a.Flush()
	s := make([]string, 0, len(items))
	for _, v := range items {
		s = append(s, v.OrderedContent(" "))
	}
	a.printRow(s, global.Space)
	if !a.disableAfter {
		fire(a.AfterPrint, a, items...)
	}
}

func (a *Across) printRowWithNoSpace(strs []string) {
	defer a.Flush()
	width := getTermWidth() - 5

	maxLength := 0
	for _, str := range strs {
		maxLength += util.WidthLen(str)
		if maxLength <= width {
			_, _ = a.WriteString(str)
		} else {
			_, _ = a.WriteString("\n" + str)
			maxLength = util.WidthLen(str)
		}
	}
	_ = a.WriteByte('\n')
}

func (a *Across) printRow(stringsArray []string, space int) {
	termWidth := getTermWidth() - 2
	n := len(stringsArray)

	// Calculate the maximum width of each string
	maxWidth := 0
	widths := make([]int, n)
	for i, s := range stringsArray {
		width := util.WidthLen(s)
		widths[i] = width
		if width > maxWidth {
			maxWidth = width
		}
	}

	// Determine the number of columns
	numColumns := 1
	for col := 1; col <= n; col++ {
		totalWidth := -space
		for i := 0; i < col; i++ {
			colMaxWidth := 0
			for j := i; j < n; j += col {
				if widths[j] > colMaxWidth {
					colMaxWidth = widths[j]
				}
			}
			totalWidth += colMaxWidth + space
		}
		if totalWidth <= termWidth {
			numColumns = col
		} else {
			break
		}
	}

	for i := 0; i < len(stringsArray); i += numColumns {
		for j := 0; j < numColumns; j++ {
			if i+j >= len(stringsArray) {
				break
			}
			s := stringsArray[i+j]
			width := widths[i+j]
			_, _ = a.WriteString(s)

			if j < numColumns-1 { // Avoid adding spaces after the last column
				maxColWidth := 0
				for k := j; k < len(stringsArray); k += numColumns {
					if widths[k] > maxColWidth {
						maxColWidth = widths[k]
					}
				}
				padding := maxColWidth - width + space
				_, _ = a.WriteString(strings.Repeat(" ", padding))
			}
		}
		_ = a.WriteByte('\n')
	}
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

func (z *Zero) Print(items ...*item.FileInfo) {
	if !z.disableBefore {
		fire(z.BeforePrint, z, items...)
	}
	defer z.Flush()
	for _, v := range items {
		_, _ = z.WriteString(v.OrderedContent(" "))
	}
	if !z.disableAfter {
		fire(z.AfterPrint, z, items...)
	}
}

type JsonPrinter struct {
	*bufio.Writer
	*hook
	Extra []any // RawPrint can't output in json, so we need an extra field
}

func NewJsonPrinter() Printer {
	return &JsonPrinter{
		Writer: bufio.NewWriter(Output),
		hook:   newHook(),
		Extra:  make([]any, 0),
	}
}

var makeJsonFieldNameReplacer = strings.NewReplacer(" ", "_", "-", "_")

func makeJsonFieldName(s string) string {
	return strings.ToLower(makeJsonFieldNameReplacer.Replace(s))
}

func (j *JsonPrinter) Print(items ...*item.FileInfo) {
	if !j.disableBefore {
		fire(j.BeforePrint, j, items...)
	}
	defer j.Flush()

	list := make([]*orderedmap.OrderedMap[string, string], 0, len(items))
	for _, v := range items {
		all := v.Meta.Pairs()

		type orderItem struct {
			name    string
			content string
			no      int
		}

		order := make([]orderItem, 0, len(all))

		// sort by v.Content.No
		for _, v := range all {
			if name := v.Key(); name != "#" {
				order = append(order, orderItem{name: makeJsonFieldName(name), content: v.Value().String(), no: v.Value().NO()})
			}
		}

		slices.SortFunc(
			order, func(a, b orderItem) int {
				return a.no - b.no
			},
		)

		s := orderedmap.New[string, string](
			orderedmap.WithCapacity[string, string](len(order)),
		)

		list = append(list, s)

		for _, v := range order {
			s.Set(v.name, v.content)
		}
	}

	wrap := &struct {
		Extra   []any                                    `json:"extra,omitempty"`
		Content []*orderedmap.OrderedMap[string, string] `json:"entries,omitempty"`
	}{
		Extra:   j.Extra,
		Content: list,
	}

	pretty, err := json.MarshalIndent(wrap, "", "	")
	if err != nil {
		_, _ = j.WriteString(err.Error() + "\n")
		return
	}
	_, _ = j.Write(pretty)
	_, _ = j.WriteString("\n")
	if !j.disableAfter {
		fire(j.AfterPrint, j, items...)
	}
}

type PrettyPrinter interface {
	SetTitle(title string)
	AddHeader(headers string)
	AddFooter(footer string)
	Printer
}

type TablePrinter struct {
	*bufio.Writer
	*hook
	header, footer table.Row
	w              table.Writer
}

func (t *TablePrinter) SetTitle(title string) {
	t.w.SetTitle("path: %s", title)
}

func (t *TablePrinter) AddHeader(headers string) {
	t.header = append(t.header, headers)
}

func (t *TablePrinter) AddFooter(footer string) {
	t.footer = append(t.footer, footer)
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

func (t *TablePrinter) PrintBase(fn func() string, s ...*item.FileInfo) {
	if !t.disableBefore {
		fire(t.BeforePrint, t, s...)
	}
	defer t.Flush()
	t.w.ResetRows()
	t.setTB(s...)
	if len(t.header) != 0 {
		t.w.AppendHeader(t.header)
	}
	if len(t.footer) != 0 {
		t.w.AppendFooter(t.footer)
	}
	fn()
	t.w.ResetHeaders()
	t.w.ResetFooters()

	// empty header and footer
	t.header = t.header[:0]
	t.footer = t.footer[:0]

	if !t.disableAfter {
		fire(t.AfterPrint, t, s...)
	}
}

func (t *TablePrinter) Print(s ...*item.FileInfo) {
	t.PrintBase(t.w.Render, s...)
}

func (t *TablePrinter) setTB(s ...*item.FileInfo) {
	for _, v := range s {
		all := v.ValuesByOrdered()
		row := make(table.Row, 0, len(all))
		for _, v := range all {
			row = append(row, strings.TrimLeft(v.String(), " "))
		}
		t.w.AppendRow(row)
	}
}

const (
	_ = iota
	TreeUnicode
	TreeASCII
	TreeRectangle
)

var (
	UNICODEStyle     = table.StyleRounded
	ASCIIStyle       = table.StyleDefault
	DefaultTBStyle   = ASCIIStyle
	DefaultTreeStyle = TreeUnicode
)

func DefaultTB(w table.Writer) {
	w.SetAllowedRowLength(getTermWidth())
	w.SetStyle(DefaultTBStyle)
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

func (m *MDPrinter) Print(s ...*item.FileInfo) {
	m.PrintBase(m.w.RenderMarkdown, s...)
}

type CSVPrinter struct {
	*TablePrinter
}

func NewCSVPrinter() Printer {
	c := &CSVPrinter{}
	c.TablePrinter = NewTablePrinter(DefaultTB).(*TablePrinter)
	return c
}

func (c *CSVPrinter) Print(s ...*item.FileInfo) {
	c.PrintBase(c.w.RenderCSV, s...)
}

type TSVPrinter struct {
	*TablePrinter
}

func NewTSVPrinter() Printer {
	t := &TSVPrinter{}
	t.TablePrinter = NewTablePrinter(DefaultTB).(*TablePrinter)
	return t
}

func (t *TSVPrinter) Print(s ...*item.FileInfo) {
	t.PrintBase(t.w.RenderTSV, s...)
}

type TreePrinter struct {
	*bufio.Writer
	*hook
	NO bool
}

func NewTreePrinter() *TreePrinter {
	return &TreePrinter{
		Writer: bufio.NewWriter(Output),
		hook:   newHook(),
	}
}

func (t *TreePrinter) Print(s ...*item.FileInfo) {
	if !t.hook.disableBefore {
		fire(t.BeforePrint, t, s...)
	}
	defer t.Flush()

	// split by full path
	// the item sharing the same dir will be grouped together
	// and the order is the same as the input
	total := len(s)

	buildTree := tree.NewTree(tree.WithCap(total / 2))
	level := make(map[int][]*item.FileInfo)
	for _, v := range s {
		level[(v.Cache["level"].(int))] = append(level[(v.Cache["level"].(int))], v)
	}

	prefixAndName := func(info *item.FileInfo) (prefix string, name string) {
		v := info.ValuesByOrdered()
		pb := bytebufferpool.Get()
		defer bytebufferpool.Put(pb)
		name = ""
		vv := v[:len(v)-1]
		for _, s := range vv {
			_, _ = pb.WriteString(s.String())
			_ = pb.WriteByte(' ')
		}
		prefix = pb.String()
		name = v[len(v)-1].String()
		return
	}

	// root
	l := len(level)
	nodeMap := make(map[string]*tree.Node, l)

	root := level[0][0]
	buildTree.Root.Meta = root
	nodeMap[root.FullPath] = buildTree.Root

	for i := 1; i < l; i++ {
		for _, v := range level[i] {
			node := nodeMap[v.Cache["parent"].(string)]
			c := &tree.Node{
				Parent:     node,
				Child:      make([]*tree.Node, 0, 10),
				Level:      i,
				Meta:       v,
				Connectors: make([]string, i),
			}
			nodeMap[v.FullPath] = c
			node.AddChild(c)
		}
	}

	Child := "├── "
	LastChild := "╰── "
	Mid := "│   "
	Empty := "    "

	if DefaultTreeStyle == TreeASCII {
		Child = "|---- "
		LastChild = "|---- "
		Mid = "|     "
	} else if DefaultTreeStyle == TreeRectangle {
		LastChild = "└── "
	} else if config.Default.CustomTreeStyle.IsEnabled() {
		Child = config.Default.CustomTreeStyle.Child
		LastChild = config.Default.CustomTreeStyle.LastChild
		Mid = config.Default.CustomTreeStyle.Mid
		Empty = config.Default.CustomTreeStyle.Empty
	}

	// print
	// the number of the prefixes is the level of the node
	// the length of prefix is 4

	applyConnectors := func(nodes []*tree.Node) {
		l := len(nodes)
		for i, n := range nodes {
			if i != l-1 {
				n.Connectors[n.Level-1] = Child
				n.Apply2Child(
					func(node *tree.Node) {
						node.Connectors[n.Level-1] = Mid
					},
				)
			} else {
				n.Connectors[n.Level-1] = LastChild
			}
		}
	}

	buildTree.Root.Apply2ChildSlice(applyConnectors)

	counter := 0
	totalLen := len(strconv.Itoa(total))
	// print
	p := func(node *tree.Node) {
		if t.NO {
			no := &ItemContent{
				No:      -1,
				Content: StringContent(strconv.Itoa(counter)),
			}
			no.SetSuffix(strings.Repeat(" ", totalLen-len(strconv.Itoa(counter))))
			counter++
			node.Meta.Set("#", no)
		}
		prefix, name := prefixAndName(node.Meta)
		_, _ = t.WriteString(prefix)
		_, _ = t.WriteString(global.Faint)
		for _, c := range node.Connectors {
			if c == "" {
				_, _ = t.WriteString(Empty)
			} else {
				_, _ = t.WriteString(c)
			}
		}
		_, _ = t.WriteString(global.Reset)
		_, _ = t.WriteString(name)
		_ = t.WriteByte('\n')
	}
	buildTree.Root.ApplyThis(p)
	buildTree.Root.Apply2Child(p)
	if !t.hook.disableAfter {
		fire(t.AfterPrint, t, s...)
	}
}
