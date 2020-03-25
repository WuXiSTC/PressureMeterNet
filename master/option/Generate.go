package option

import (
	"github.com/yindaheng98/cipg"
)

func Generate(logger func(i ...interface{})) (opt Option, exit bool) {
	opt = DefaultOption()
	exit = cipg.GenerateWithYAML(&opt, logger)
	return
}
