package validation

import (
  "testing"
  "github.com/viki-org/gspec"
)

func init() {
  Define("generic").Field("x").Message("whyuno")
}

func TestPanicsWhenValidationFailsForAnUnknownDefinition(t *testing.T) {
  spec := gspec.New(t)
  defer func() {
    if r := recover(); r != nil {
      spec.Expect(r.(string)).ToEqual("unknown definition invalid")
    } else {
      t.Error("expected panic never happened")
    }
  }()
  New().Required("", "invalid")
}

// Required
func TestInvalidOnEmptyRequiredString(t *testing.T) {
  validator := New().Required("", "generic")
  assertGenericInvalid(t, validator)
}

func TestValidOnNonEmptyRequiredString(t *testing.T) {
  validator := New().Required("over 9000", "generic")
  assertGenericValid(t, validator)
}

// Len
func TestInvalidWhenStringLengthIsOutsideLen(t *testing.T) {
  for _, str := range []string{"123", "12", "1234567", "12345678"} {
    validator := New().Len(str, 4, 6, "generic")
    assertGenericInvalid(t, validator)
  }
}

func TestValidWhenStringLengthIsWithinLen(t *testing.T) {
  for _, str := range []string{"1234", "12345", "123456"} {
    validator := New().Len(str, 4, 6, "generic")
    assertGenericValid(t, validator)
  }
}

// MinLen
func TestInvalidWhenStringIsTooShort(t *testing.T) {
  for _, str := range []string{"1", "12", "123"} {
    validator := New().MinLen(str, 4, "generic")
    assertGenericInvalid(t, validator)
  }
}

func TestValidWhenStringIsLongerThanMinLen(t *testing.T) {
  for _, str := range []string{"1234", "12345", "1234567890"} {
    validator := New().MinLen(str, 4, "generic")
    assertGenericValid(t, validator)
  }
}

// MaxLen
func TestInvalidWhenStringIsTooLong(t *testing.T) {
  for _, str := range []string{"1234567", "12345678", "1234567890"} {
    validator := New().MaxLen(str, 6, "generic")
    assertGenericInvalid(t, validator)
  }
}

func TestValidWhenStringLengthIsShorterThanMaxLen(t *testing.T) {
  for _, str := range []string{"1", "12", "123", "1234", "12345", "123456"} {
    validator := New().MaxLen(str, 6, "generic")
    assertGenericValid(t, validator)
  }
}

// Same
func TestInvalidWhenTheTwoStringsArentEqual(t *testing.T) {
  validator := New().Same("password", "passwrod", "generic")
  assertGenericInvalid(t, validator)
}

func TestInvalidWhenTheTwoStringsAreEqual(t *testing.T) {
  validator := New().Same("it's over 9000", "it's over 9000", "generic")
  assertGenericValid(t, validator)
}

func assertGenericInvalid(t *testing.T, validator *Validator) {
  spec := gspec.New(t)
  spec.Expect(len(validator.Errors)).ToEqual(1)
  spec.Expect(validator.IsValid()).ToEqual(false)
  res, valid := validator.Response()
  spec.Expect(res).ToNotBeNil()
  spec.Expect(valid).ToEqual(false)
  spec.Expect(len(validator.Errors["x"])).ToEqual(1)
  spec.Expect(validator.Errors["x"][0].message).ToEqual("whyuno")
}

func assertGenericValid(t *testing.T, validator *Validator) {
  spec := gspec.New(t)
  spec.Expect(len(validator.Errors)).ToEqual(0)
  spec.Expect(validator.IsValid()).ToEqual(true)
  res, valid := validator.Response()
  spec.Expect(res).ToBeNil()
  spec.Expect(valid).ToEqual(true)
}
