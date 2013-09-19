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
  parentResource string
  action interface{}
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
