package validation

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestInvalidWhenStringIsTooShort(t *testing.T) {
	spec := gspec.New(t)
	rule := MinLen(4)
	for _, str := range []string{"1", "12", "123"} {
		spec.Expect(rule.Verify(str)).ToEqual(false)
	}
}

func TestValidWhenStringIsLongerThanMinLen(t *testing.T) {
	spec := gspec.New(t)
	rule := MinLen(4)
	for _, str := range []string{"1234", "12345", "123456"} {
		spec.Expect(rule.Verify(str)).ToEqual(true)
	}
}
