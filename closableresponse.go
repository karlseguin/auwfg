package auwfg

import (
  "net/http"
  "github.com/viki-org/bytepool"
)

type ClosableResponse struct {
  S int
  H http.Header
  B *bytepool.Item
}

func (r *ClosableResponse) Status() int {
  return r.S
}

func (r *ClosableResponse) Body() []byte {
  return r.B.Bytes()
}

func (r *ClosableResponse) Header() http.Header {
  return r.H
}

func (r *ClosableResponse) Length() int {
  return r.B.Len()
}

func (r *ClosableResponse) Close() {
  r.B.Close()
}
