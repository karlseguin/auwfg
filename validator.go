package auwfg

import (
	"github.com/karlseguin/auwfg/validation"
)

func Validator() *validation.Validator {
	return new(validation.Validator)
}
