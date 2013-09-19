package auwfg
// Another useless web framework for Go

import (
  "log"
  "net"
  "time"
  "net/http"
)

var NotFound Response

func Run(config *Configuration) {
  NotFound = config.notFound
  s := &http.Server{
    MaxHeaderBytes: 8192,
    Handler: newRouter(config),
    ReadTimeout: 10 * time.Second,
  }
  l, err := net.Listen("tcp", config.address)
  if err != nil { log.Fatal(err) }
  log.Fatal(s.Serve(l))
}
