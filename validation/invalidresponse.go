package validation

import (
  "fmt"
  "strconv"
  "net/http"
  "github.com/viki-org/bytepool"
)

var invalidBytePool = bytepool.New(1, 1024)
// more than a little ugly...
func InitInvalidPool(poolSize, bufferSize int) {
  invalidBytePool = bytepool.New(poolSize, bufferSize)
}

type InvalidResponse struct {
  buffer *bytepool.Item
}

func NewResponse(errors map[string][]*Definition) *InvalidResponse {
  buffer := invalidBytePool.Checkout()
  //am I really doing this in a public repo?!
  buffer.WriteByte(byte('{'))
  for field, definitions := range errors {
    buffer.WriteString(fmt.Sprintf("%q:[", field))
    for _, definition := range definitions {
      buffer.WriteString(fmt.Sprintf("%q,", definition.message))
    }
    buffer.Position(buffer.Len() -1) //strip trailing comma
    buffer.WriteString("],")
  }
  buffer.Position(buffer.Len() -1) //strip trailing comma
  buffer.WriteByte(byte('}'))
  return &InvalidResponse{buffer}
}

func (r *InvalidResponse) SetStatus(status int) {}

func (r *InvalidResponse) GetStatus() int {
  return 400
}

func (r *InvalidResponse) GetBody() []byte {
  return r.buffer.Bytes()
}

func (r *InvalidResponse) GetHeader() http.Header {
  return http.Header{"Content-Type": []string{"application/json; charset=utf-8"}, "Content-Length": []string{strconv.Itoa(r.buffer.Len())}}
}

func (r *InvalidResponse) Close() error {
  return r.buffer.Close()
}
