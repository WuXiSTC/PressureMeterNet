package option

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"testing"
)

func YAMLTest(t *testing.T, i interface{}) {
	is, err := yaml.Marshal(i)
	if err != nil {
		t.Log(err)
	} else {
		fmt.Println(string(is))
	}
	if err = yaml.Unmarshal(is, i); err != nil {
		t.Log(err)
	} else {
		printStructOption(i, "YAMLTest", func(i ...interface{}) {
			fmt.Println(i...)
		})
	}
}

func TestOption(t *testing.T) {
	opt := GenerateOption(func(i ...interface{}) {
		t.Log(i...)
	})
	YAMLTest(t, opt)
}
