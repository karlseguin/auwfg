package auwfg

import (
  "errors"
  "strconv"
  "net/http"
  "github.com/viki-org/bytepool"
)

var (
  JsonHeader = []string{"application/json; charset=utf-8"}
)

type Response interface {
  Status() int
  Body() []byte
  Header() http.Header
  Close()
}

type NormalResponse struct {
  S int
  B []byte
  H http.Header
}

func (r *NormalResponse) Status() int {
  return r.S
}

func (r *NormalResponse) Body() []byte {
  return r.B
}

func (r *NormalResponse) Header() http.Header {
  return r.H
}

func (r *NormalResponse) Close() {}

func Json(body interface{}, status int) Response {
  switch b := body.(type) {
  case string:
    return normalResponse([]byte(b), status)
  case []byte:
    return normalResponse(b, status)
  case *bytepool.Item:
    return closableResponse(b, status)
  default:
    return Fatal(errors.New("unknown body type"))
  }
}

func normalResponse(b []byte, status int) Response {
  return &NormalResponse{
    S: status,
    B: b,
    H: http.Header{"Content-Type": JsonHeader, "Content-Length": []string{strconv.Itoa(len(b))}},
  }
}

func closableResponse(b *bytepool.Item, s int) Response {
  return &ClosableResponse{
    S: s,
    B: b,
    H: http.Header{"Content-Type": JsonHeader, "Content-Length": []string{strconv.Itoa(b.Len())}},
  }
}

type FatalResponse struct {
  err error
  Response
}

func Fatal(err error) Response {
  return &FatalResponse{err, InternalServerError,}
}
