package validation

type MinLenRule struct {
  min int
}

func MinLen(min int) Rule {
  return &MinLenRule{min}
}

func (r *MinLenRule) Verify(value interface{}) bool {
  switch value := value.(type) {
    case string: return len(value) >= r.min
    default: return false
  }
}
