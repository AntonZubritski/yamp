package ui

import (
	"math"
	"strings"
)

// zubritskiGlyphs holds 5×7 pixel bitmaps for each letter in "ZUBRITSKI".
// Each row is 5 bits wide; bit 4 (0x10) is the leftmost pixel.
var zubritskiGlyphs = [9][7]uint8{
	{0x1F, 0x01, 0x02, 0x04, 0x08, 0x10, 0x1F}, // Z
	{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E}, // U
	{0x1E, 0x11, 0x11, 0x1E, 0x11, 0x11, 0x1E}, // B
	{0x1E, 0x11, 0x11, 0x1E, 0x14, 0x12, 0x11}, // R
	{0x0E, 0x04, 0x04, 0x04, 0x04, 0x04, 0x0E}, // I
	{0x1F, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04}, // T
	{0x0E, 0x11, 0x10, 0x0E, 0x01, 0x11, 0x0E}, // S
	{0x11, 0x11, 0x0A, 0x04, 0x0A, 0x11, 0x11}, // K
	{0x0E, 0x04, 0x04, 0x04, 0x04, 0x04, 0x0E}, // I
}

const (
	zLetterW    = 5
	zLetterH    = 7
	zNumLetters = 9
	zGap        = 1
	zTotalW     = zNumLetters*zLetterW + (zNumLetters-1)*zGap // 53
)

// renderZubritski draws "ZUBRITSKI" in pixel art using Braille dots.
// Inside each letter, equalizer bars rise and fall with the music.
// The bars fill the letter shapes from bottom to top based on band energy,
// creating a visual effect of the equalizer pumping through the text.
func (v *Visualizer) renderZubritski(bands [numBands]float64) string {
	height := v.Rows
	dotRows := height * 4
	dotCols := panelWidth * 2

	grid := make([]bool, dotRows*dotCols)

	// Scale letters to fill the panel.
	scaleX := dotCols / zTotalW
	scaleY := (dotRows * 3 / 4) / zLetterH
	if scaleX < 1 {
		scaleX = 1
	}
	if scaleY < 1 {
		scaleY = 1
	}

	renderedW := zTotalW * scaleX
	renderedH := zLetterH * scaleY
	offsetX := (dotCols - renderedW) / 2
	baseOffsetY := (dotRows - renderedH) / 2

	// Map 9 letters across the 10 frequency bands with slight overlap.
	letterBand := [9]int{0, 1, 2, 3, 4, 5, 6, 7, 9}

	// Overall energy for global pulse.
	var totalEnergy float64
	for _, b := range bands {
		totalEnergy += b
	}
	totalEnergy /= numBands

	for li := range zNumLetters {
		energy := bands[letterBand[li]]

		// Subtle bounce from energy + traveling wave.
		wave := math.Sin(float64(v.frame)*0.08+float64(li)*0.7) * 1.2
		bounce := int(energy*float64(baseOffsetY)*0.25 + wave)

		letterX := offsetX + li*(zLetterW+zGap)*scaleX
		letterY := baseOffsetY - bounce

		// The equalizer bar fill level for this letter (0..1 maps to bottom..top of letter).
		barLevel := energy*0.85 + totalEnergy*0.15

		for py := range zLetterH {
			row := zubritskiGlyphs[li][py]
			for px := range zLetterW {
				if row&(1<<(zLetterW-1-px)) == 0 {
					continue
				}

				for sy := range scaleY {
					for sx := range scaleX {
						dx := letterX + px*scaleX + sx
						dy := letterY + py*scaleY + sy

						if dx < 0 || dx >= dotCols || dy < 0 || dy >= dotRows {
							continue
						}

						// Compute how high in the letter this pixel is (0=bottom, 1=top).
						pixelY := float64(zLetterH*scaleY-1-(py*scaleY+sy)) / float64(zLetterH*scaleY)

						// Equalizer fill: pixels below barLevel are solid, above dissolve.
						if pixelY <= barLevel {
							grid[dy*dotCols+dx] = true
						} else {
							// Sparse dissolve above the bar line — a few dots sparkle.
							sparkle := scatterHash(li, py*scaleY+sy, px*scaleX+sx, v.frame)
							threshold := (pixelY - barLevel) * 3.0
							if sparkle > threshold {
								grid[dy*dotCols+dx] = true
							}
						}
					}
				}
			}
		}
	}

	// Convert dot grid to Braille characters with spectrum coloring.
	lines := make([]string, height)
	for row := range height {
		var content strings.Builder
		for ch := range panelWidth {
			var braille rune = '\u2800'
			for dr := range 4 {
				for dc := range 2 {
					if grid[(row*4+dr)*dotCols+ch*2+dc] {
						braille |= brailleBit[dr][dc]
					}
				}
			}
			content.WriteRune(braille)
		}
		lines[row] = specStyle(float64(height-1-row) / float64(height)).Render(content.String())
	}

	return strings.Join(lines, "\n")
}
