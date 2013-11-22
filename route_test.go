package auwfg

import (
  "testing"
  "github.com/viki-org/gspec"
)

func TestParsesANormalRoute(t *testing.T) {
  spec := gspec.New(t)
  rb := R("GET", "v1", "users", nil)
  spec.Expect(rb.resource).ToEqual("users")
  spec.Expect(rb.parentResource).ToEqual("")
}

func TestParsesANestedRoute(t *testing.T) {
  spec := gspec.New(t)
  rb := R("GET", "v1", "users/likes", nil)
  spec.Expect(rb.resource).ToEqual("likes")
  spec.Expect(rb.parentResource).ToEqual("users")
}
