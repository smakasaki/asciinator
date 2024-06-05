package handlers

import (
	"fmt"

	"github.com/buger/goterm"
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

	// Get terminal size
	termWidth := uint(goterm.Width())
	termHeight := uint(goterm.Height()) - 2 // Adjust for terminal borders

	// Resize the image to fit within the terminal size
	img = processor.Resize(img, termWidth, termHeight)

	asciiArt, err := converter.Convert(img, flags)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return asciiArt, nil
}
