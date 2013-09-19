package auwfg

import (
  "net/http"
)
type ContextFactory func(*BaseContext) interface{}

type BaseContext struct {
  Route *Route
  Params *Params
  Req *http.Request
}

func genericContextFactory(c *BaseContext) interface{} { return c }
