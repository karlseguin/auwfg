package validation

var definitions = make(map[string]*Definition)

type Definition struct {
  id string
  field string
  message string
}

func Define(id string) *Definition {
  definition := &Definition{id: id,}
  definitions[id] = definition
  return definition

}

func (d *Definition) Field(field string) *Definition {
  d.field = field
  return d
}

func (d *Definition) Message(message string) *Definition {
  d.message = message
  return d
}
