package auwfg

import (
	"errors"
	"net/http"
	"strconv"
)

var (
	JsonHeader = []string{"application/json; charset=utf-8"}
)

type Response interface {
	SetStatus(status int)

	GetStatus() int
	GetBody() []byte
	GetHeader() http.Header
	Close() error
}

type ResponseBuilder struct {
	Response Response
}

func (b *ResponseBuilder) Status(status int) *ResponseBuilder {
	b.Response.SetStatus(status)
	return b
}

func (b *ResponseBuilder) Cache(duration int) *ResponseBuilder {
	b.Response.GetHeader().Set("Cache-Control", "private,max-age="+strconv.Itoa(duration))
	return b
}

func (b *ResponseBuilder) Header(key, value string) *ResponseBuilder {
	b.Response.GetHeader().Set(key, value)
	return b
}

func (b *ResponseBuilder) SetStatus(status int) {
	b.Response.SetStatus(status)
}

func (b *ResponseBuilder) GetStatus() int {
	return b.Response.GetStatus()
}

func (b *ResponseBuilder) GetBody() []byte {
	return b.Response.GetBody()
}

func (b *ResponseBuilder) GetHeader() http.Header {
	return b.Response.GetHeader()
}

func (b *ResponseBuilder) Close() error {
	return b.Response.Close()
}

func Json(body interface{}) *ResponseBuilder {
	switch b := body.(type) {
	case string:
		return &ResponseBuilder{newNormalResponse([]byte(b), 200)}
	case []byte:
		return &ResponseBuilder{newNormalResponse(b, 200)}
	case ByteCloser:
		return &ResponseBuilder{newClosableResponse(b, 200)}
	default:
		return &ResponseBuilder{Fatal(errors.New("unknown body type"))}
	}
}

type FatalResponse struct {
	err error
	Response
}

func Fatal(err error) Response {
	return &FatalResponse{err, InternalServerError}
}
