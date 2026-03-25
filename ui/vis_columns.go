package ui

import "strings"

// renderColumns draws many thin single-character-wide columns, interpolating
// between bands so adjacent columns vary slightly for a dense, organic look.
func (v *Visualizer) renderColumns(bands [numBands]float64) string {
	height := v.Rows

	// Compute per-band column counts and flat-array offsets.
	var bandCols [numBands]int
	var offsets [numBands]int
	totalCols := 0
	for b := range numBands {
		offsets[b] = totalCols
		bandCols[b] = visBandWidth(b)
		totalCols += bandCols[b]
	}

	// Build per-column levels by interpolating between neighboring bands.
	cols := make([]float64, totalCols)
	for b, level := range bands {
		nextLevel := level
		if b+1 < numBands {
			nextLevel = bands[b+1]
		}
		for c := range bandCols[b] {
			t := float64(c) / float64(bandCols[b])
			cols[offsets[b]+c] = level*(1-t) + nextLevel*t
		}
	}

	lines := make([]string, height)

	for row := range height {
		var content strings.Builder
		rowBottom := float64(height-1-row) / float64(height)
		rowTop := float64(height-row) / float64(height)

		for b := range numBands {
			for c := range bandCols[b] {
				level := cols[offsets[b]+c]
				content.WriteString(fracBlock(level, rowBottom, rowTop))
			}
			if b < numBands-1 {
				content.WriteByte(' ')
			}
		}
		lines[row] = specStyle(rowBottom).Render(content.String())
	}

	return strings.Join(lines, "\n")
}
