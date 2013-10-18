package auwfg

import (
  "strconv"
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

func ClosableJson(b *bytepool.Item, s int) Response {
  return &ClosableResponse{
    S: s,
    B: b,
    H: http.Header{"Content-Type": JsonHeader, "Content-Length": []string{strconv.Itoa(b.Len())}},
  }
}
