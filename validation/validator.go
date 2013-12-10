package validation

import (
	"fmt"
)

type Rule interface {
	Verify(value interface{}) bool
}

type Validator struct {
	Errors map[string][]*Definition
}

func New() *Validator {
	return new(Validator)
}

func (v *Validator) Validate(value interface{}, definitionId string) bool {
	definition, exists := definitions[definitionId]
	if exists == false {
		panic(fmt.Sprintf("unknown definition %v", definitionId))
	}
	valid := definition.rule.Verify(value)
	if valid == false {
		v.AddError(definitionId)
	}
	return valid
}

func (v *Validator) IsValid() bool {
	return v.Errors == nil
}

func (v *Validator) Response() (*InvalidResponse, bool) {
	if v.Errors == nil {
		return nil, true
	}
	return NewResponse(v.Errors), false
}

func (v *Validator) AddError(definitionId string) {
	definition, exists := definitions[definitionId]
	if exists == false {
		panic(fmt.Sprintf("unknown definition %v", definitionId))
	}
	if v.Errors == nil {
		v.Errors = make(map[string][]*Definition)
	}
	bucket, exists := v.Errors[definition.field]
	if exists == false {
		bucket = make([]*Definition, 0, 2)
		v.Errors[definition.field] = bucket
	}
	v.Errors[definition.field] = append(bucket, definition)
}
