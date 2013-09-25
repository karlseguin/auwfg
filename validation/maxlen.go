package validation

type MaxLenRule struct {
  max int
}

func MaxLen(max int) Rule {
  return &MaxLenRule{max}
}

func (r *MaxLenRule) Verify(value interface{}) bool {
  switch value := value.(type) {
    case string: return len(value) <= r.max
    default: return false
  }
}
