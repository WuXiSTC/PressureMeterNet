package option

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"reflect"
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
	} else {
		printValue(reflect.ValueOf(i), "YAMLTest", func(i ...interface{}) {
			fmt.Println(i...)
		})
	}
}

func TestOption(t *testing.T) {
	opt := Generate(func(i ...interface{}) {
		t.Log(i...)
	})
	YAMLTest(t, &opt)
}
