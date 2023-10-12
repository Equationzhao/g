package display

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/Equationzhao/g/display/tree"
	"github.com/Equationzhao/g/item"
	"github.com/Equationzhao/g/theme"
	"github.com/Equationzhao/g/util"
	"github.com/acarl005/stripansi"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mattn/go-runewidth"
	"github.com/olekukonko/ts"
	"github.com/valyala/bytebufferpool"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

var Output io.Writer = os.Stdout

func RawPrint(toPrint ...any) (n int, err error) {
	return fmt.Fprint(Output, toPrint...)
}

// print style control

type hook struct {
	BeforePrint   []func(Printer, ...*item.FileInfo)
	AfterPrint    []func(Printer, ...*item.FileInfo)
	disableBefore bool
	disableAfter  bool
}

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
		BeforePrint: make([]func(Printer, ...*item.FileInfo), 0, 5),
		AfterPrint:  make([]func(Printer, ...*item.FileInfo), 0, 5),
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

type Hook interface {
	AddBeforePrint(...func(Printer, ...*item.FileInfo))
	AddAfterPrint(...func(Printer, ...*item.FileInfo))
	DisablePreHook()
	EnablePreHook()
	DisablePostHook()
	EnablePostHook()
}

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
		_ = b.WriteByte('\n')
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
	f.printColumns(&s)

	if !f.disableAfter {
		fire(f.AfterPrint, f, i...)
	}
}

func (f *FitTerminal) printColumns(strs *[]string) {
	stringsArray := *strs

	termWidth := getTermWidth()

	maxColumnWidths := 0
	for i := 0; i < len(stringsArray); i++ {
		width := WidthLen(stringsArray[i])
		if width > maxColumnWidths {
			maxColumnWidths = width
		}
	}

	columnSpacing := 2

	maxTotalWidth := 0
	maxTotalWidth += maxColumnWidths
	numColumns := (termWidth + columnSpacing) / (maxTotalWidth + columnSpacing)
	if numColumns < 1 {
		numColumns = 1
	}

	numRows := (len(stringsArray) + numColumns - 1) / numColumns

	for i := 0; i < numRows; i++ {
		for j := 0; j < numColumns; j++ {
			index := j*numRows + i
			if index >= len(stringsArray) {
				break
			}
			s := stringsArray[index]
			width := WidthLen(s)
			_, _ = f.WriteString(s)
			padding := maxColumnWidths - width + columnSpacing
			_, _ = f.WriteString(strings.Repeat(" ", padding))
		}
		_ = f.WriteByte('\n')
	}
}

var IncludeHyperlink = false

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
	if IncludeHyperlink {
		name, other, ok := parseLink(str)
		if ok {
			str = other + name
		}
	}
	colorless := stripansi.Strip(str)
	// len() is insufficient here, as it counts emojis as 4 characters each
	length := runewidth.StringWidth(colorless)

	return length
}

func WidthNoHyperLinkLen(str string) int {
	colorless := stripansi.Strip(str)
	// len() is insufficient here, as it counts emojis as 4 characters each
	length := runewidth.StringWidth(colorless)

	return length
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
	c.printRowWithNoSpace(&s)
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
	a.printRow(&s)
	if !a.disableAfter {
		fire(a.AfterPrint, a, items...)
	}
}

func (a *Across) printRowWithNoSpace(strs *[]string) {
	defer a.Flush()
	width := getTermWidth()

	maxLength := 0
	for _, str := range *strs {
		maxLength += WidthLen(str)
		if maxLength <= width {
			_, _ = a.WriteString(str)
		} else {
			_, _ = a.WriteString("\n" + str)
			maxLength = WidthLen(str)
		}
	}
	_ = a.WriteByte('\n')
}

func (a *Across) printRow(strs *[]string) {
	stringsArray := *strs

	termWidth := getTermWidth()

	maxColumnWidths := 0
	for i := 0; i < len(stringsArray); i++ {
		width := WidthLen(stringsArray[i])
		if width > maxColumnWidths {
			maxColumnWidths = width
		}
	}

	columnSpacing := 2

	maxTotalWidth := 0
	maxTotalWidth += maxColumnWidths
	numColumns := (termWidth + columnSpacing) / (maxTotalWidth + columnSpacing)
	if numColumns < 1 {
		numColumns = 1
	}

	for i := 0; i < len(stringsArray); i += numColumns {
		for j := 0; j < numColumns; j++ {
			if i+j >= len(stringsArray) {
				break
			}
			s := stringsArray[i+j]
			width := WidthLen(s)
			_, _ = a.WriteString(s)

			padding := maxColumnWidths - width + columnSpacing
			_, _ = a.WriteString(strings.Repeat(" ", padding))
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
	Extra []any
}

func NewJsonPrinter() Printer {
	return &JsonPrinter{
		Writer: bufio.NewWriter(Output),
		hook:   newHook(),
		Extra:  make([]any, 0),
	}
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
			if name := v.Key(); name == "Name" {
				order = append(order, orderItem{name: name, content: v.Value().String(), no: v.Value().NO()})
			} else {
				// remove all leading spaces
				order = append(
					order, orderItem{name: name, content: strings.TrimSpace(v.Value().String()), no: v.Value().NO()},
				)
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
		Extra   []any                                    `json:"Extra,omitempty"`
		Content []*orderedmap.OrderedMap[string, string] `json:"Content,omitempty"`
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

type HTMLPrinter struct {
	*TablePrinter
}

func NewHTMLPrinter() Printer {
	h := &HTMLPrinter{}
	h.TablePrinter = NewTablePrinter(DefaultTB).(*TablePrinter)
	return h
}

func (p *HTMLPrinter) Print(s ...*item.FileInfo) {
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
	level := make(map[string][]*item.FileInfo)
	for _, v := range s {
		level[string(v.Cache["level"])] = append(level[string(v.Cache["level"])], v)
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

	root := level["0"][0]
	buildTree.Root.Meta = root
	nodeMap[root.FullPath] = buildTree.Root

	for i := 1; i < l; i++ {
		for _, v := range level[strconv.Itoa(i)] {
			node := nodeMap[string(v.Cache["parent"])]
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
		Empty = "    "
	} else if DefaultTreeStyle == TreeRectangle {
		LastChild = "└── "
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
		_, _ = t.WriteString(theme.Faint)
		for _, c := range node.Connectors {
			if c == "" {
				_, _ = t.WriteString(Empty)
			} else {
				_, _ = t.WriteString(c)
			}
		}
		_, _ = t.WriteString(theme.Reset)
		_, _ = t.WriteString(name)
		_ = t.WriteByte('\n')
	}
	buildTree.Root.ApplyThis(p)
	buildTree.Root.Apply2Child(p)
	if !t.hook.disableAfter {
		fire(t.AfterPrint, t, s...)
	}
}
