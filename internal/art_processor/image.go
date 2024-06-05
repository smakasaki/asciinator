package art_processor

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

// ImageProcessor interface for image processing
type ImageProcessor interface {
	Load(filePath string) (image.Image, error)
	LoadFromURL(url string) (image.Image, error)
	Resize(img image.Image, maxWidth, maxHeight uint) image.Image
}

// Processor struct implementing ImageProcessor
type Processor struct{}

// Load loads an image from a file
func (p Processor) Load(filePath string) (image.Image, error) {
	localImage, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer localImage.Close()

	img, _, err := image.Decode(localImage)
	return img, err
}

// LoadFromURL loads an image from a URL
func (p Processor) LoadFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	return img, err
}

// Resize resizes the image to fit within the specified width and height
func (p Processor) Resize(img image.Image, maxWidth, maxHeight uint) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Adjust for terminal character aspect ratio (approx 2:1)
	charAspectRatio := 2.0
	termAspectRatio := float64(maxWidth) / (float64(maxHeight) * charAspectRatio)

	// Calculate new dimensions based on the maximum dimensions while maintaining aspect ratio
	var newWidth, newHeight uint
	if float64(width)/float64(height) > termAspectRatio {
		newWidth = maxWidth
		newHeight = uint(float64(maxWidth) / float64(width) * float64(height) / charAspectRatio)
	} else {
		newHeight = maxHeight
		newWidth = uint(float64(maxHeight) / float64(height) * float64(width) * charAspectRatio)
	}

	return resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
}
