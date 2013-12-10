package validation

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestInvalidWhenStringLengthIsOutsideLen(t *testing.T) {
	spec := gspec.New(t)
	rule := Len(4, 6)
	for _, str := range []string{"123", "12", "1234567", "12345678"} {
		spec.Expect(rule.Verify(str)).ToEqual(false)
	}
}

func TestValidWhenStringLengthIsWithinLen(t *testing.T) {
	spec := gspec.New(t)
	rule := Len(4, 6)
	for _, str := range []string{"1234", "12345", "123456"} {
		spec.Expect(rule.Verify(str)).ToEqual(true)
	}
}
