package auwfg
// Another useless web framework for Go

import (
  "log"
  "net"
  "time"
  "net/http"
)

func Run(config *Configuration) {
  s := &http.Server{
    MaxHeaderBytes: 8192,
    Handler: &Router{config,},
    ReadTimeout: 10 * time.Second,
  }
  l, err := net.Listen("tcp", config.address)
  if err != nil { log.Fatal(err) }
  log.Fatal(s.Serve(l))
}
