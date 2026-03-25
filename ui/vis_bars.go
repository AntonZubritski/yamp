package ui

import "strings"

// renderBars is the default smooth spectrum with fractional Unicode blocks.
func (v *Visualizer) renderBars(bands [numBands]float64) string {
	height := v.Rows
	lines := make([]string, height)

	for row := range height {
		var content strings.Builder
		rowBottom := float64(height-1-row) / float64(height)
		rowTop := float64(height-row) / float64(height)

		for i, level := range bands {
			bw := visBandWidth(i)
			block := fracBlock(level, rowBottom, rowTop)
			for range bw {
				content.WriteString(block)
			}
			if i < numBands-1 {
				content.WriteByte(' ')
			}
		}
		lines[row] = specStyle(rowBottom).Render(content.String())
	}

	return strings.Join(lines, "\n")
}
