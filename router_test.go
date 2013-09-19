package auwfg

import (
  "testing"
  "net/http/httptest"
  "github.com/viki-org/gspec"
)

func TestNotFoundOnShortUrls(t *testing.T) {
  req := gspec.Request().Url("/x").Req
  res := httptest.NewRecorder()
  newRouter(Configure()).ServeHTTP(res, req)
  assertResponse(t, res, 404, `{"error":"not found","code":404}`)
}

func TestNotFoundOnUnknownVersions(t *testing.T) {
  req := gspec.Request().Url("/v3/sessions.json").Req
  res := httptest.NewRecorder()
  newRouter(Configure()).ServeHTTP(res, req)
  assertResponse(t, res, 404, `{"error":"not found","code":404}`)
}

func TestNotFoundOnUnknownControllers(t *testing.T) {
  req := gspec.Request().Url("/v4/cats.json").Req
  res := httptest.NewRecorder()
  newRouter(Configure()).ServeHTTP(res, req)
  assertResponse(t, res, 404, `{"error":"not found","code":404}`)
}

func TestNotFoundOnUnknownActions(t *testing.T) {
  req := gspec.Request().Url("/v4/sessions.json").Req
  res := httptest.NewRecorder()
  newRouter(Configure()).ServeHTTP(res, req)
  assertResponse(t, res, 404, `{"error":"not found","code":404}`)
}

func TestExecutesAListAction(t *testing.T) {
  f := func(context interface{}) Response {return JsonResponse(`{"spice":"mustflow"}`, 200)}
  req := gspec.Request().Url("/v4/sessions.json").Req
  res := httptest.NewRecorder()
  router := newRouter(Configure().Route(R("LIST", "v4", "sessions", f)))
  router.ServeHTTP(res, req)
  assertResponse(t, res, 200, `{"spice":"mustflow"}`)
}

func TestExecutesAGetAction(t *testing.T) {
  f := func(context interface{}) Response { return JsonResponse(`{"name":"duncan"}`, 200) }
  req := gspec.Request().Url("/v2/gholas/123g.json").Req
  res := httptest.NewRecorder()
  router := newRouter(Configure().Route(R("GET", "v2", "gholas", f)))
  router.ServeHTTP(res, req)
  assertResponse(t, res, 200, `{"name":"duncan"}`)
}

func TestExecutesAPostAction(t *testing.T) {
  f := func(context interface{}) Response { return JsonResponse(`{"name":"shaihulud"}`, 200) }
  req := gspec.Request().Url("/v1/worms.json").Method("POST").Req
  res := httptest.NewRecorder()
  router := newRouter(Configure().Route(R("POST", "v1", "worms", f)))
  router.ServeHTTP(res, req)
  assertResponse(t, res, 200, `{"name":"shaihulud"}`)
}

func TestExecutesAPutAction(t *testing.T) {
  f := func(context interface{}) Response { return JsonResponse(`{"name":"shaihulud"}`, 200) }
  req := gspec.Request().Url("/v1/worms/22w.json").Method("PUT").Req
  res := httptest.NewRecorder()
  router := newRouter(Configure().Route(R("PUT", "v1", "worms", f)))
  router.ServeHTTP(res, req)
  assertResponse(t, res, 200, `{"name":"shaihulud"}`)
}

func TestCreatingAndDispatchingThroughCustomTypes(t *testing.T) {
  spec := gspec.New(t)
  f := func(context *TestContext) Response {
    spec.Expect(context.Name).ToEqual("leto")
    spec.Expect(context.Params.Version).ToEqual("v1")
    spec.Expect(context.Params.Resource).ToEqual("worms")
    spec.Expect(context.Params.Id).ToEqual("22w")
    return JsonResponse(`{"name":"bigjohn"}`, 203)
  }
  c := Configure().Route(R("GET", "v1", "worms", f)).ContextFactory(testContextFactory).Dispatcher(testDispatcher)
  req := gspec.Request().Url("/v1/worms/22w.json").Method("GET").Req
  res := httptest.NewRecorder()
  newRouter(c).ServeHTTP(res, req)
  assertResponse(t, res, 203, `{"name":"bigjohn"}`)
}

func TestDefaulsTheBodyValueWhenNoBodyIsPresent(t *testing.T) {
  spec := gspec.New(t)
  f := func(context *TestContext) Response {
    input := context.Body.(*TestBody)
    spec.Expect(input.Hello).ToEqual("")
    return JsonResponse("", 200)
  }
  c := Configure().Route(R("GET", "v1", "worms", f).BodyFactory(testBodyFactory)).ContextFactory(testContextFactory).Dispatcher(testDispatcher)
  req := gspec.Request().Url("/v1/worms/22w.json").Method("GET").Req
  newRouter(c).ServeHTTP(httptest.NewRecorder(), req)
}

func TestParsesABody(t *testing.T) {
  spec := gspec.New(t)
  f := func(context *TestContext) Response {
    input := context.Body.(*TestBody)
    spec.Expect(input.Hello).ToEqual("World")
    return JsonResponse("", 200)
  }
  c := Configure().Route(R("GET", "v1", "worms", f).BodyFactory(testBodyFactory)).ContextFactory(testContextFactory).Dispatcher(testDispatcher)
  req := gspec.Request().Url("/v1/worms/22w.json").Method("GET").BodyString(`{"hello":"World"}`).Req
  newRouter(c).ServeHTTP(httptest.NewRecorder(), req)
}

func TestHandlesANilResponse(t *testing.T) {
  f := func(context interface{}) Response { return nil }
  c := Configure().Route(R("GET", "v1", "worms", f))
  req := gspec.Request().Url("/v1/worms/22w.json").Method("GET").Req
  res := httptest.NewRecorder()
  newRouter(c).ServeHTTP(res, req)
  assertResponse(t, res, 500, `{"error":"internal server error","code":500}`)
}

func TestHandlesBodiesLargerThanAllowed(t *testing.T) {
  f := func(context *TestContext) Response { return JsonResponse("", 200) }
  c := Configure().Route(R("GET", "v1", "worms", f).BodyFactory(testBodyFactory)).ContextFactory(testContextFactory).Dispatcher(testDispatcher).BodyPool(3, 1)
  req := gspec.Request().Url("/v1/worms/22w.json").Method("GET").BodyString(`{"hello":"World"}`).Req
  res := httptest.NewRecorder()
  newRouter(c).ServeHTTP(res, req)
  assertResponse(t, res, 413, `{"error":"body too large","code":413}`)
}

func assertResponse(t *testing.T, res *httptest.ResponseRecorder, status int, raw string) {
  spec := gspec.New(t)
  spec.Expect(res.Code).ToEqual(status)
  spec.Expect(string(res.Body.Bytes())).ToEqual(raw)
}

type TestContext struct {
  Name string
  *BaseContext
}

type TestBody struct {
  Hello string
}

func testContextFactory(c *BaseContext) interface{} { return &TestContext{"leto",c} }
func testDispatcher(route *Route, context interface{}) Response {
  return route.Action.(func(*TestContext) Response)(context.(*TestContext))
}
func testBodyFactory() interface{} {
  return new(TestBody)
}
