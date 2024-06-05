package handlers

func CreateChain() ImageHandler {
	// gifHandler := &GIFHandler{}
	// urlHandler := &URLHandler{}
	localHandler := &LocalHandler{}

	// gifHandler.SetNext(urlHandler)
	// urlHandler.SetNext(fileHandler)

	return localHandler
}
