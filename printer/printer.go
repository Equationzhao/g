package printer

import (
	"fmt"
	"strings"
)

// print style control

type Printer interface {
	Print(s ...string)
}

type Byline struct{}

func (b Byline) Print(s ...string) {
	fmt.Println(strings.Join(s, "\n"))
}

type FitTerminal struct{}

func (f FitTerminal) Print(s ...string) {
	PrintColumns(&s, 4)
}
