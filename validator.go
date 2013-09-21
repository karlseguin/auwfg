package auwfg

import (
  "auwfg/validation"
)

func Validator() *validation.Validator {
  return new(validation.Validator)
}

func Define(id string) *validation.Definition {
  return validation.Define(id)
}
