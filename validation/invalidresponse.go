package validation

import (
  "net/http"
)

var header = http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}

type InvalidResponse struct {
  Errors map[string][]*Definition
}

func (r *InvalidResponse) Status() int {
  return 400
}

func (r *InvalidResponse) Body() []byte {
  return []byte("Need closable responses to implement this, I think")
}

func (r *InvalidResponse) Header() http.Header {
  return header
}
