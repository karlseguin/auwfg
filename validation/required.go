package validation

type RequiredRule struct{}

var req = new(RequiredRule)

func Required() Rule { return req }

func (r *RequiredRule) Verify(value interface{}) bool {
	switch value := value.(type) {
	case string:
		return len(value) > 0
	default:
		return false
	}
}
