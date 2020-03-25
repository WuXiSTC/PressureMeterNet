package option

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"testing"
)

func YAMLTest(t *testing.T, i *Option) {
	is, err := yaml.Marshal(i)
	if err != nil {
		t.Log(err)
	} else {
		fmt.Println(string(is))
	}
	if err = yaml.Unmarshal(is, i); err != nil {
		t.Log(err)
	}
}

func TestOption(t *testing.T) {
	opt, _ := Generate(func(i ...interface{}) {
		t.Log(i...)
	})
	YAMLTest(t, &opt)
}
