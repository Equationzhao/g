package display

type Content interface {
	String() string
}

type StringContent string

func (s StringContent) String() string {
	return string(s)
}

type ItemContent struct {
	No             int
	Content        Content
	Prefix, Suffix string
}

func (i *ItemContent) SetPrefix(s string) {
	i.Prefix = s
}

func (i *ItemContent) SetSuffix(s string) {
	i.Suffix = s
}

func (i *ItemContent) AddPrefix(add string) {
	i.Prefix = add + i.Prefix
}

func (i *ItemContent) AddSuffix(add string) {
	i.Suffix = add + i.Suffix
}

func (i *ItemContent) String() string {
	return i.Prefix + i.Content.String() + i.Suffix
}

func (i *ItemContent) NO() int {
	return i.No
}
