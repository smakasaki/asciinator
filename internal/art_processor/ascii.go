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
	Braille   bool
}

// ASCIIConverter interface for ASCII conversion
type ASCIIConverter interface {
	Convert(img image.Image, flags Flags) (string, error)
}

// Converter struct implementing ASCIIConverter
type Converter struct{}

// Convert converts the image to ASCII art with optional colors
func (c Converter) Convert(img image.Image, flags Flags) (string, error) {
	if flags.Braille {
		return c.convertToBraille(img, flags)
	}
	return c.convertToASCII(img, flags)
}

func (c Converter) convertToASCII(img image.Image, flags Flags) (string, error) {
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

// convertToBraille converts the image to Braille art with optional colors
func (c Converter) convertToBraille(img image.Image, flags Flags) (string, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	brailleArt := ""
	profile := termenv.ColorProfile()

	for y := 0; y < height; y += 1 {
		for x := 0; x < width; x += 2 {
			char := brailleChar(img, x, y)

			if flags.Colored {
				r, g, b, _ := img.At(x, y).RGBA()
				color := profile.Color(fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8))
				brailleArt += termenv.String(char).Foreground(color).String()
			} else {
				brailleArt += char
			}
		}
		brailleArt += "\n"
	}

	return brailleArt, nil
}

// brailleChar converts a 2x4 block of pixels to a single Braille character
func brailleChar(img image.Image, x, y int) string {
	pixels := [8]bool{}
	offsets := []struct{ dx, dy int }{
		{0, 0}, {1, 0}, {0, 1}, {1, 1},
		{0, 2}, {1, 2}, {0, 3}, {1, 3},
	}

	for i, offset := range offsets {
		if x+offset.dx < img.Bounds().Max.X && y+offset.dy < img.Bounds().Max.Y {
			gray := colorToGray(img.At(x+offset.dx, y+offset.dy))
			pixels[i] = gray < 128
		}
	}

	brailleValue := 0x2800
	for i, pixel := range pixels {
		if pixel {
			brailleValue |= 1 << uint(i)
		}
	}

	return string(rune(brailleValue))
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
