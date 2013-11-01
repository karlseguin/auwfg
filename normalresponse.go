package auwfg

import (
  "strconv"
  "net/http"
)

type NormalResponse struct {
  S int
  B []byte
  H http.Header
}

func (r *NormalResponse) SetStatus(status int) {
  r.S = status
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


func newNormalResponse(b []byte, status int) Response {
  return &NormalResponse{
    S: status,
    B: b,
    H: http.Header{"Content-Type": JsonHeader, "Content-Length": []string{strconv.Itoa(len(b))}},
  }
}
