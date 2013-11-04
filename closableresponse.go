package auwfg

import (
  "strconv"
  "net/http"
)

type ByteCloser interface {
  Len() int
  Close()
  Bytes() []byte
}

type ClosableResponse struct {
  S int
  H http.Header
  B ByteCloser
}

func (r *ClosableResponse) SetStatus(status int) {
  r.S = status
}

func (r *ClosableResponse) GetStatus() int {
  return r.S
}

func (r *ClosableResponse) GetBody() []byte {
  return r.B.Bytes()
}

func (r *ClosableResponse) GetHeader() http.Header {
  return r.H
}

func (r *ClosableResponse) Length() int {
  return r.B.Len()
}

func (r *ClosableResponse) Close() {
  r.B.Close()
}

func newClosableResponse(b ByteCloser, s int) Response {
  return &ClosableResponse{
    S: s,
    B: b,
    H: http.Header{"Content-Type": JsonHeader, "Content-Length": []string{strconv.Itoa(b.Len())}},
  }
}
