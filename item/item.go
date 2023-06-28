package item

type Item interface {
	String() string
	NO() int // NO return the No. of item
	SetPrefix(string)
	SetSuffix(string)
	GetPrefix() string
	GetSuffix() string
}
