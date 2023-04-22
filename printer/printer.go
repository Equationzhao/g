package printer

import (
	"fmt"
	"os"
	"strings"
)

// print style control

type Printer interface {
	Print(s ...string)
}

type Inline struct{}

func (i Inline) Print(s ...string) {
	fmt.Print(strings.Join(s, " "), "\n")
}

type Block struct{}

func (b Block) Print(s ...string) {
	fmt.Println(strings.Join(s, " "))
}

type Byline struct{}

func (b Byline) Print(s ...string) {
	fmt.Println(strings.Join(s, "\n"))
}

type FitTerminal struct{}

func (f FitTerminal) Print(s ...string) {
	w := int(os.Stdout.Fd())
	max := 0
	for _, str := range s {
		if len(str) > max {
			max = len(str)
		}
	}
	if max*2 > w {
		fmt.Println(strings.Join(s, "\n"))
	} else {
		n := max / w
		for i, si := range s {
			fmt.Print(si)
			if i%n == 0 {
				fmt.Println()
			} else {
				fmt.Print(strings.Repeat("", n))
			}
		}
	}
}
