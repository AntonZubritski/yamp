package ui

import "strings"

// renderTerrain draws a scrolling side-view landscape where terrain height
// is the current spectrum energy. New data enters from the right and scrolls
// left, creating a moving mountain range silhouette. Braille dots give smooth
// sub-cell edges; spectrum coloring paints green valleys, yellow slopes, red peaks.
func (v *Visualizer) renderTerrain(bands [numBands]float64) string {
	height := v.Rows
	dotRows := height * 4
	dotCols := panelWidth * 2

	// Ensure terrain buffer matches current panel width.
	if len(v.terrainBuf) != dotCols {
		v.terrainBuf = make([]float64, dotCols)
	}

	// Scroll left by 2 dot columns per frame for visible movement.
	copy(v.terrainBuf, v.terrainBuf[2:])

	// Compute new rightmost height from average spectrum energy.
	var totalEnergy float64
	for _, e := range bands {
		totalEnergy += e
	}
	avg := totalEnergy / float64(numBands)

	// Two new columns with slight noise for organic ridge edges.
	v.terrainBuf[dotCols-2] = min(1.0, avg+scatterHash(0, 0, 0, v.frame)*0.12)
	v.terrainBuf[dotCols-1] = min(1.0, avg+scatterHash(0, 0, 1, v.frame)*0.12)

	// Render: each dot column is filled from its terrain height down to the bottom.
	lines := make([]string, height)
	for row := range height {
		var content strings.Builder
		for ch := range panelWidth {
			var braille rune = '\u2800'
			for dc := range 2 {
				x := ch*2 + dc
				terrainH := v.terrainBuf[x]
				// Top dot position — invert so 0 is bottom.
				topDot := dotRows - 1 - int(terrainH*float64(dotRows-1))
				for dr := range 4 {
					dotY := row*4 + dr
					if dotY >= topDot {
						braille |= brailleBit[dr][dc]
					}
				}
			}
			content.WriteRune(braille)
		}
		// Color by row height: green base, yellow middle, red peaks.
		lines[row] = specStyle(float64(height-1-row) / float64(height)).Render(content.String())
	}

	return strings.Join(lines, "\n")
}
