package validation

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func init() {
	Define("fail", "failfield", "failmessage", &FakeRule{false})
	Define("pass", "passfield", "passmessage", &FakeRule{true})
	InitInvalidPool(1, 1024)
}

func TestIsNotValidIfOneRuleFails(t *testing.T) {
	spec := gspec.New(t)
	validator := New()
	spec.Expect(validator.Validate("value", "fail")).ToEqual(false)
	spec.Expect(validator.IsValid()).ToEqual(false)
	spec.Expect(len(validator.Errors["failfield"])).ToEqual(1)
	spec.Expect(validator.Errors["failfield"][0].message).ToEqual("failmessage")
	res, valid := validator.Response()
	spec.Expect(res).ToNotBeNil()
	spec.Expect(valid).ToEqual(false)
}

func TestVali(t *testing.T) {
	spec := gspec.New(t)
	validator := New()
	spec.Expect(validator.Validate("value", "pass")).ToEqual(true)
	spec.Expect(validator.IsValid()).ToEqual(true)
	spec.Expect(len(validator.Errors)).ToEqual(0)
	res, valid := validator.Response()
	spec.Expect(res).ToBeNil()
	spec.Expect(valid).ToEqual(true)
}

type FakeRule struct {
	valid bool
}

func (f *FakeRule) Verify(value interface{}) bool {
	return f.valid
}
