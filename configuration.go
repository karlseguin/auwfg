package auwfg

import (
  "strings"
)

type Configuration struct {
  address string
  notFound Response
  routes versions
  dispatcher Dispatcher
  contextFactory ContextFactory
}

type versions map[string]controllers
type controllers map[string]actions
type actions map[string]*Route

func Configure() *Configuration{
  return &Configuration{
    routes: make(versions),
    address: "127.0.0.1:4577",
    dispatcher: genericDispatcher,
    contextFactory: genericContextFactory,
    notFound: JsonResponse(`{"error":"not found","code":404}`, 404),
  }
}

func (c *Configuration) Address(address string) *Configuration {
  c.address = address
  return c
}

func (c *Configuration) NotFoundResponse(r Response) *Configuration {
  c.notFound = r
  return c
}

func (c *Configuration) NotFound(body string) *Configuration {
  return c.NotFoundResponse(JsonResponse(body, 404))
}

func (c *Configuration) Dispatcher(d Dispatcher) *Configuration {
  c.dispatcher = d
  return c
}

func (c *Configuration) ContextFactory(cf ContextFactory) *Configuration {
  c.contextFactory = cf
  return c
}

func (c *Configuration) Route(r *RouteBuilder) *Configuration {
  r.method = strings.ToUpper(r.method)
  r.version = strings.ToLower(r.version)
  r.resource = strings.ToLower(r.resource)
  if _, exists := c.routes[r.version]; !exists { c.routes[r.version] = make(controllers) }
  if _, exists := c.routes[r.version][r.resource]; !exists { c.routes[r.version][r.resource] = make(actions) }
  c.routes[r.version][r.resource][r.method] = &Route{
    Action: r.action,
    BodyFactory: r.bf,
  }
  return c
}
