package auwfg
// Another useless web framework for Go

import (
  "log"
  "net"
  "time"
  "net/http"
  "github.com/viki-org/auwfg/validation"
)

var NotFound = Json(`{"error":"not found","code":404}`, 404)
var InternalServerError = Json(`{"error":"internal server error","code":500}`, 500)
var Deleted = &NormalResponse{204, []byte{}, http.Header{"Content-Length": []string{"0"}}}

func Run(config *Configuration) {
  NotFound = config.notFound
  InternalServerError = config.internalServerError
  validation.InitInvalidPool(config.invalidPoolSize, config.maxInvalidSize)

  s := &http.Server{
    MaxHeaderBytes: 8192,
    Handler: newRouter(config),
    ReadTimeout: 10 * time.Second,
  }
  l, err := net.Listen("tcp", config.address)
  if err != nil { log.Fatal(err) }
  log.Fatal(s.Serve(l))
}
