package art_processor

import (
	"fmt"
	"image"
	"image/color"

	"github.com/muesli/termenv"
)

// Flags struct for command-line flags
type Flags struct {
	CustomMap string
	Colored   bool
}

// ASCIIConverter interface for ASCII conversion
type ASCIIConverter interface {
	Convert(img image.Image, flags Flags) (string, error)
}

// Converter struct implementing ASCIIConverter
type Converter struct{}

// Convert converts the image to ASCII art with optional colors
func (c Converter) Convert(img image.Image, flags Flags) (string, error) {
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
