package validation

var definitions = make(map[string]*Definition)

type Definition struct {
  id string
  field string
  message string
}

func Define(id, field, message string) *Definition {
  definition := &Definition{
    id: id,
    field: field,
    message: message,
  }
  definitions[id] = definition
  return definition
}