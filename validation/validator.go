package validation

import(
  "fmt"
)

type Validator struct {
  Errors map[string][]*Definition
}

func New() *Validator {
  return new(Validator)
}

func (v *Validator) Required(value, definition string) *Validator {
  if len(value) == 0 { v.AddError(definition) }
  return v
}

func (v *Validator) Len(value string, min, max int, definition string) *Validator {
  l := len(value)
  if l < min || l > max { v.AddError(definition) }
  return v
}

func (v *Validator) MinLen(value string, min int, definition string) *Validator {
  if len(value) < min { v.AddError(definition) }
  return v
}

func (v *Validator) MaxLen(value string, max int, definition string) *Validator {
  if len(value) > max { v.AddError(definition) }
  return v
}

func (v *Validator) Same(a, b, definition string) *Validator {
  if a != b { v.AddError(definition) }
  return v
}

func (v *Validator) IsValid() bool {
  return v.Errors == nil
}

func (v *Validator) Response() (*InvalidResponse, bool) {
  if v.Errors == nil { return nil, true }
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
