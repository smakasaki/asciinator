package handlers

import (
	"errors"

	"github.com/smakasaki/asciinator/internal/art_processor"
)

type ImageHandler interface {
	Handle(request string, flags art_processor.Flags) (string, error)
	SetNext(handler ImageHandler) ImageHandler
}

type BaseHandler struct {
	next ImageHandler
}

func (b *BaseHandler) Handle(request string, flags art_processor.Flags) (string, error) {
	if b.next != nil {
		return b.next.Handle(request, flags)
	}
	return "", errors.New("no handler found for request")
}

func (b *BaseHandler) SetNext(handler ImageHandler) ImageHandler {
	b.next = handler
	return handler
}
