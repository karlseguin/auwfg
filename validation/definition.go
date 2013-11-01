package validation

var definitions = make(map[string]*Definition)

type Definition struct {
  id string
  field string
  message string
  rule Rule
}
func Define(id, field, message string, rule Rule) *Definition {
  definition := &Definition{id, field, message, rule}
  definitions[id] = definition
  return definition
}
