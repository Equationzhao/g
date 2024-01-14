package theme

import "testing"

func TestColor2String(t *testing.T) {
	c := Underline + Red
	s := color2str(c)
	println(s)

	c, _ = str2color(s)
	println(c + "test")
}
