package handlers

import "github.com/buger/goterm"

func CreateChain() ImageHandler {
	// gifHandler := &GIFHandler{}
	urlHandler := &UrlHandler{}
	localHandler := &LocalHandler{}

	// gifHandler.SetNext(urlHandler)
	urlHandler.SetNext(localHandler)

	return urlHandler
}

func GetTerminalSize() (uint, uint) {
	termWidth := uint(goterm.Width())
	termHeight := uint(goterm.Height()) - 2 // Adjust for terminal borders
	return termWidth, termHeight
}
