package option

import (
	"fmt"
	"testing"
)

func TestRegExp(t *testing.T) {
	PMFieldValue := "////123//12&(*3/////1^&23//23&(*e//wf/w*(e//wt//33_))+P:+_////"
	PMFieldValue = re.ReplaceAllString(PMFieldValue, "/") //去除多余斜杠
	PMFieldValue = rex.ReplaceAllString(PMFieldValue, "") //去除非法字符
	fmt.Println(PMFieldValue)
}
