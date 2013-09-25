package validation

type LenRule struct {
  min int
  max int
}

func Len(min, max int) Rule {
  return &LenRule{min, max}
}

func (r *LenRule) Verify(value interface{}) bool {
  switch value := value.(type) {
    case string:
      l := len(value)
      return l >= r.min && l <= r.max
    default: return false
  }
}
