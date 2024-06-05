package image_processor

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/buger/goterm"
	"github.com/muesli/termenv"
	"github.com/nfnt/resize"
)

// Convert function to handle the main conversion logic
func Convert(filePath string, flags Flags) (string, error) {
	localImage, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Error opening file: %v", err)
	}
	defer localImage.Close()

	return convertImage(localImage, flags)
}

// convertImage function to handle image decoding and ASCII conversion
func convertImage(localImage *os.File, flags Flags) (string, error) {
	img, _, err := image.Decode(localImage)
	if err != nil {
		return "", fmt.Errorf("Error decoding image: %v", err)
	}

	// Get terminal size
	termWidth := uint(goterm.Width())
	termHeight := uint(goterm.Height()) - 2 // Adjust for terminal borders

	// Resize the image to fit within the terminal size
	img = resizeToFit(img, termWidth, termHeight)

	asciiArt, err := imageToASCII(img, flags)
	if err != nil {
		return "", err
	}

	return asciiArt, nil
}

// resizeToFit resizes the image to fit within the specified width and height
func resizeToFit(img image.Image, maxWidth, maxHeight uint) image.Image {
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

// imageToASCII function to convert the image to ASCII art with optional colors
func imageToASCII(img image.Image, flags Flags) (string, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	asciiArt := ""
	profile := termenv.ColorProfile()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			grayColor := colorToGray(c)
			char := grayToASCII(grayColor, flags.CustomMap)

			if flags.Colored {
				// Convert the color to termenv format and set it
				r, g, b, _ := c.RGBA()
				color := profile.Color(fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8))
				asciiArt += termenv.String(char).Foreground(color).String()
			} else {
				asciiArt += char
			}
		}
		asciiArt += "\n"
	}

	return asciiArt, nil
}

// colorToGray function to convert a color to grayscale value
func colorToGray(c color.Color) uint8 {
	r, g, b, _ := c.RGBA()
	gray := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	return uint8(gray / 256)
}

// grayToASCII function to map grayscale value to ASCII character
func grayToASCII(gray uint8, customMap string) string {
	const asciiChars = ` .'` + "`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

	if customMap == "" {
		customMap = asciiChars
	}

	// Map the gray value to the corresponding ASCII character
	scale := float64(len(customMap)) / 256.0
	index := int(float64(gray) * scale)
	if index >= len(customMap) {
		index = len(customMap) - 1
	}
	return string(customMap[index])
}
