package handlers

import (
	"fmt"

	"github.com/smakasaki/asciinator/internal/art_processor"
)

type LocalHandler struct {
	BaseHandler
}

func (l *LocalHandler) Handle(request string, flags art_processor.Flags) (string, error) {
	processor := art_processor.Processor{}
	converter := art_processor.Converter{}
	img, err := processor.Load(request)
	if err != nil {
		return l.BaseHandler.Handle(request, flags)
	}

	termWidth, termHeight := GetTerminalSize()
	img = processor.Resize(img, termWidth, termHeight)

	asciiArt, err := converter.Convert(img, flags)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return asciiArt, nil
}
