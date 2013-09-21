package auwfg

import (
  "testing"
  "github.com/viki-org/gspec"
)

func TestCreatesAJsonResponse(t *testing.T) {
  spec := gspec.New(t)
  r := Json("this is the body", 9001)
  spec.Expect(string(r.Body())).ToEqual("this is the body")
  spec.Expect(r.Status()).ToEqual(9001)
  spec.Expect(len(r.Header())).ToEqual(2)
  spec.Expect(r.Header().Get("Content-Type")).ToEqual("application/json; charset=utf-8")
  spec.Expect(r.Header().Get("Content-Length")).ToEqual("16")
}
