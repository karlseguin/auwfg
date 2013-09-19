package auwfg

import (
  "testing"
  "github.com/viki-org/gspec"
)

func TestSetsTheAddress(t *testing.T) {
  spec := gspec.New(t)
  spec.Expect(Configure().Address("invalid-and-we-dont-care").address).ToEqual("invalid-and-we-dont-care")
}

func TestSetsTheNotFoundResponse(t *testing.T) {
  spec := gspec.New(t)
  expected := JsonResponse("the res", 244)
  actual := Configure().NotFoundResponse(expected).notFound
  spec.Expect(actual.Status()).ToEqual(244)
  spec.Expect(string(actual.Body())).ToEqual("the res")
}

func TestSetsTheNotFoundBody(t *testing.T) {
  spec := gspec.New(t)
  actual := Configure().NotFound("try again").notFound
  spec.Expect(actual.Status()).ToEqual(404)
  spec.Expect(string(actual.Body())).ToEqual("try again")
}

func TestAddingASimpleRoute(t *testing.T) {
  spec := gspec.New(t)
  c := Configure().Route(R("GET", "v1", "gholas", "badAction"))
  spec.Expect(c.routes["v1"]["gholas"]["GET"].Action.(string)).ToEqual("badAction")
}

func TestAddingMultipleSimpleRoutes(t *testing.T) {
  spec := gspec.New(t)
  c := Configure().
        Route(R("GET", "v1", "gholas", "getgholas")).
        Route(R("LIST", "v1", "gholas", "listgholas"))
  spec.Expect(c.routes["v1"]["gholas"]["GET"].Action.(string)).ToEqual("getgholas")
  spec.Expect(c.routes["v1"]["gholas"]["LIST"].Action.(string)).ToEqual("listgholas")
}

func TestFixesRouteCasing(t *testing.T) {
  spec := gspec.New(t)
  c := Configure().
        Route(R("get", "V1", "GHOLAS", "getgholas"))
  spec.Expect(c.routes["v1"]["gholas"]["GET"].Action.(string)).ToEqual("getgholas")
}
