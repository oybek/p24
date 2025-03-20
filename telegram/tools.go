package telegram

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func DrawTextToImage(text string, fontBytes []byte, bgColor color.RGBA) ([]byte, error) {
	fontSize := float64(60)
	DPI := float64(36)

	fontWidth, fontHeight, err := calculateFontSize(fontBytes, fontSize, DPI)
	if err != nil {
		return nil, err
	}

	textWidth, textHeight := getTextDimensions(text)

	width := fontWidth*textWidth + fontWidth*2
	height := fontHeight*textHeight + fontHeight*2

	// Set starting position
	x := fontWidth
	y := fontHeight * 3 / 2

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Load the font with Cyrillic support
	parsedFont, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     DPI,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, err
	}

	// Draw text
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  fixed.P(x, y),
	}

	for _, line := range splitLines(text) {
		d.Dot = fixed.P(x, y)
		d.DrawString(line)
		y += (fontHeight + 1)
	}

	// Encode image to PNG format in memory
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func splitLines(text string) []string {
	var lines []string
	start := 0
	for i, ch := range text {
		if ch == '\n' {
			lines = append(lines, text[start:i])
			start = i + 1
		}
	}
	lines = append(lines, text[start:])
	return lines
}

func getTextDimensions(text string) (int, int) {
	maxLength, curLen, numLines := 0, 0, 1

	// Iterate over runes (Unicode characters)
	for _, ch := range text {
		if ch == '\n' {
			if curLen > maxLength {
				maxLength = curLen
			}
			curLen = 0
			numLines++
		} else {
			curLen++
		}
	}

	// Handle the last line (if it doesn't end with '\n')
	if curLen > maxLength {
		maxLength = curLen
	}

	return maxLength, numLines
}

// Load the font and calculate character width and text height
func calculateFontSize(fontBytes []byte, fontSize float64, dpi float64) (int, int, error) {
	// Parse the font
	parsedFont, err := opentype.Parse(fontBytes)
	if err != nil {
		return 0, 0, err
	}

	// Create a font face
	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return 0, 0, err
	}
	defer face.Close()

	// Measure font metrics
	metrics := face.Metrics()
	ascent := metrics.Ascent.Round()   // Baseline to top
	descent := metrics.Descent.Round() // Baseline to bottom
	//lineGap := float64(metrics.Height.Round())  // Total line height

	// Calculate total height per line
	fontHeight := ascent + descent // If needed, add lineGap

	// Measure average character width (using '0' as a representative character)
	advance, _ := face.GlyphAdvance('0') // Works well for monospace fonts
	charWidth := advance.Round()         // Convert to pixels

	return charWidth, fontHeight, nil
}
