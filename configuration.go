package auwfg

import (
  "strings"
)

type Configuration struct {
  address string
  maxBodySize int64
  bodyPoolSize int
  notFound Response
  bodyTooLarge Response
  invalidFormat Response
  internalServerError Response
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
    maxBodySize: 32769,
    bodyPoolSize: 1024,
    dispatcher: genericDispatcher,
    contextFactory: genericContextFactory,
    notFound: JsonResponse(`{"error":"not found","code":404}`, 404),
    bodyTooLarge: JsonResponse(`{"error":"body too large","code":413}`, 413),
    invalidFormat: JsonResponse(`{"error":"invalid input format","code":400}`, 400),
    internalServerError: JsonResponse(`{"error":"internal server error","code":500}`, 500),
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

func (c *Configuration) BodyTooLargeResponse(r Response) *Configuration {
  c.bodyTooLarge = r
  return c
}

func (c *Configuration) BodyTooLarge(body string) *Configuration {
  return c.BodyTooLargeResponse(JsonResponse(body, 413))
}

func (c *Configuration) InvalidFormatResponse(r Response) *Configuration {
  c.invalidFormat = r
  return c
}

func (c *Configuration) InvalidFormat(body string) *Configuration {
  return c.InvalidFormatResponse(JsonResponse(body, 400))
}

func (c *Configuration) InternalServerErrorResponse(r Response) *Configuration {
  c.internalServerError = r
  return c
}

func (c *Configuration) InternalServerError(body string) *Configuration {
  return c.InternalServerErrorResponse(JsonResponse(body, 500))
}

func (c *Configuration) Dispatcher(d Dispatcher) *Configuration {
  c.dispatcher = d
  return c
}

func (c *Configuration) ContextFactory(cf ContextFactory) *Configuration {
  c.contextFactory = cf
  return c
}

func (c *Configuration) BodyPool(maxBodySize int, poolSize int) *Configuration {
  c.bodyPoolSize = poolSize
  c.maxBodySize = int64(maxBodySize + 1)
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
