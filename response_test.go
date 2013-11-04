package auwfg

import (
  "testing"
  "github.com/viki-org/gspec"
)

func TestCreatesAJsonResponse(t *testing.T) {
  spec := gspec.New(t)
  r := Json("this is the body").Status(9001).Response
  spec.Expect(string(r.GetBody())).ToEqual("this is the body")
  spec.Expect(r.GetStatus()).ToEqual(9001)
  spec.Expect(len(r.GetHeader())).ToEqual(2)
  spec.Expect(r.GetHeader().Get("Content-Type")).ToEqual("application/json; charset=utf-8")
  spec.Expect(r.GetHeader().Get("Content-Length")).ToEqual("16")
}
