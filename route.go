package auwfg

import (
  "strings"
)

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
  parentResource := ""
  if index := strings.Index(resource, "/"); index > -1 {
    parentResource = resource[:index]
    resource = resource[index+1:]
  }

  return &RouteBuilder{
    action: action,
    method: method,
    version: version,
    resource: resource,
    parentResource: parentResource,
  }
}

func (r *RouteBuilder) BodyFactory(bf func() interface{}) *RouteBuilder {
  r.bf = bf
  return r
}
