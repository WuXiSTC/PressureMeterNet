package hsync

import (
	"fmt"
	"testing"
)

func TestRegExp(t *testing.T) {
	s := "  \n\t123e\t\n\n\n\n  1e23e2d\n\tefwev\t\n#\twvdcsdvf\n\n\n"
	fmt.Println(s)
	s = reg.ReplaceAllString(s, "")
	fmt.Println(s)
	s = regx.ReplaceAllString(s, " ")
	fmt.Println(s)
}
