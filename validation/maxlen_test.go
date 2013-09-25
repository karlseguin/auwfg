package validation

import (
  "testing"
  "github.com/viki-org/gspec"
)

func TestInvalidWhenStringIsTooLong(t *testing.T) {
  spec := gspec.New(t)
  rule := MaxLen(4)
  for _, str := range []string{"12345", "123456"} {
    spec.Expect(rule.Verify(str)).ToEqual(false)
  }
}

func TestValidWhenStringIsShorterThanMaxLen(t *testing.T) {
  spec := gspec.New(t)
  rule := MaxLen(4)
  for _, str := range []string{"1", "12", "123", "1234"} {
    spec.Expect(rule.Verify(str)).ToEqual(true)
  }
}
