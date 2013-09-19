package auwfg

import (
  "net/http"
)

type ContextFactory func(*BaseContext) interface{}
func genericContextFactory(c *BaseContext) interface{} { return c }

type BaseContext struct {
  Route *Route
  Params *Params
  Body interface{}
  Req *http.Request
}

func newBaseContext(route *Route, params *Params, req *http.Request) *BaseContext{
  return &BaseContext{Route: route, Params: params, Req: req,}
}
