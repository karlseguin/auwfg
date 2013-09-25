package validation

import (
  "testing"
  "github.com/viki-org/gspec"
)

func TestInvalidOnEmptyRequiredString(t *testing.T) {
  spec := gspec.New(t)
  spec.Expect(Required().Verify("")).ToEqual(false)
}

func TestValidOnNonEmptyRequiredString(t *testing.T) {
  spec := gspec.New(t)
  spec.Expect(Required().Verify("over 9000")).ToEqual(true)
}
