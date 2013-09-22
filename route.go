package auwfg

type BodyFactory func() interface{}

type Route struct {
  Action interface{}
  BodyFactory BodyFactory
}

type RouteBuilder struct {
  method string
  version string
  resource string
  action interface{}
  parentResource string
  bf BodyFactory
}

func R(method, version, resource string, action interface{}) *RouteBuilder {
  return &RouteBuilder{
    action: action,
    method: method,
    version: version,
    resource: resource,
  }
}

func (r *RouteBuilder) BodyFactory(bf func() interface{}) *RouteBuilder {
  r.bf = bf
  return r
}

func (r *RouteBuilder) Parent(parentResource string) *RouteBuilder {
  r.parentResource = parentResource
  return r
}
