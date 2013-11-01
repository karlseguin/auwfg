package auwfg

import (
  "errors"
  "net/http"
  "github.com/viki-org/bytepool"
)

var (
  JsonHeader = []string{"application/json; charset=utf-8"}
)

type Response interface {
  SetStatus(status int)

  Status() int
  Body() []byte
  Header() http.Header
  Close()
}

type ResponseBuilder struct {
  Response Response
}

func (b *ResponseBuilder) Status(status int) *ResponseBuilder {
  b.Response.SetStatus(status)
  return b
}

func Json(body interface{}) *ResponseBuilder {
  switch b := body.(type) {
  case string:
    return &ResponseBuilder{newNormalResponse([]byte(b), 200)}
  case []byte:
    return &ResponseBuilder{newNormalResponse(b, 200)}
  case *bytepool.Item:
    return &ResponseBuilder{newClosableResponse(b, 200)}
  default:
    return &ResponseBuilder{Fatal(errors.New("unknown body type"))}
  }
}

type FatalResponse struct {
  err error
  Response
}

func Fatal(err error) Response {
  return &FatalResponse{err, InternalServerError,}
}

