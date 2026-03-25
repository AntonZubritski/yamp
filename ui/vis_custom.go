package ui

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"yamp/internal/appdir"
)

// CustomVisEffect identifies the rendering effect for a custom text visualizer.
type CustomVisEffect int

const (
	EffectDissolve  CustomVisEffect = iota // dots scatter/fill with energy
	EffectEqualizer                        // bars rise inside letters
	EffectBounce                           // letters bounce vertically
	EffectRain                             // drops fall through letters
	EffectMatrix                           // falling matrix chars through letters
	EffectFlame                            // fire rising through letters
	EffectPulse                            // letters pulse size with energy
	EffectGlitch                           // random block corruption
	EffectWave                             // sine wave fills letters
	EffectBinary                           // 0/1 streaming through letters
	EffectScatter                          // sparkle particles inside letters
	EffectLightning                        // electric bolts inside letters
	effectCount
)

var effectNames = [effectCount]string{
	"Dissolve", "Equalizer", "Bounce", "Rain",
	"Matrix", "Flame", "Pulse", "Glitch",
	"Wave", "Binary", "Scatter", "Lightning",
}
var effectDescs = [effectCount]string{
	"Dots appear/disappear with energy",
	"EQ bars fill letters from bottom",
	"Letters bounce with the beat",
	"Drops fall through letter shapes",
	"Falling matrix characters",
	"Fire rising through letters",
	"Letters pulse size with beat",
	"Random block corruption",
	"Sine wave fills letters",
	"0/1 streaming through letters",
	"Sparkle particles inside",
	"Electric bolts inside letters",
}

// CustomVisConfig describes a user-created text visualizer.
type CustomVisConfig struct {
	Text   string `json:"text"`
	Effect string `json:"effect"`
}

func effectFromString(s string) CustomVisEffect {
	for i, n := range effectNames {
		if strings.EqualFold(n, s) {
			return CustomVisEffect(i)
		}
	}
	return EffectDissolve
}

// customVisPath returns the path to the custom visualizers JSON file.
func customVisPath() string {
	dir, err := appdir.Dir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "custom_vis.json")
}

// loadCustomVis loads saved custom visualizers from disk.
func loadCustomVis() []CustomVisConfig {
	path := customVisPath()
	if path == "" {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var configs []CustomVisConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil
	}
	return configs
}

// saveCustomVis writes custom visualizers to disk.
func saveCustomVis(configs []CustomVisConfig) error {
	path := customVisPath()
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// renderCustomText renders arbitrary text with the given effect and band data.
// Used by both the custom visualizer mode and the preview in the editor overlay.
func (v *Visualizer) renderCustomText(text string, effect CustomVisEffect, bands [numBands]float64) string {
	glyphs := textGlyphs(text)
	if len(glyphs) == 0 {
		return ""
	}

	height := v.Rows
	dotRows := height * 4
	dotCols := panelWidth * 2

	grid := make([]bool, dotRows*dotCols)

	letterGap := 1
	totalTextW := len(glyphs)*fontW + (len(glyphs)-1)*letterGap

	// Scale to fill panel.
	scaleX := dotCols / totalTextW
	scaleY := (dotRows * 3 / 4) / fontH
	if scaleX < 1 {
		scaleX = 1
	}
	if scaleY < 1 {
		scaleY = 1
	}
	// Cap scale so text doesn't look too blocky.
	if scaleX > scaleY*2 {
		scaleX = scaleY * 2
	}

	renderedW := totalTextW * scaleX
	renderedH := fontH * scaleY
	offsetX := (dotCols - renderedW) / 2
	baseOffsetY := (dotRows - renderedH) / 2

	nLetters := utf8.RuneCountInString(text)

	// Overall energy.
	var totalEnergy float64
	for _, b := range bands {
		totalEnergy += b
	}
	totalEnergy /= numBands

	for li, glyph := range glyphs {
		// Map letter to a band.
		bandIdx := 0
		if nLetters > 1 {
			bandIdx = li * (numBands - 1) / (nLetters - 1)
		}
		if bandIdx >= numBands {
			bandIdx = numBands - 1
		}
		energy := bands[bandIdx]

		letterX := offsetX + li*(fontW+letterGap)*scaleX

		var letterY int
		switch effect {
		case EffectBounce:
			wave := math.Sin(float64(v.frame)*0.08+float64(li)*0.7) * 1.5
			bounce := int(energy*float64(baseOffsetY)*0.35 + wave)
			letterY = baseOffsetY - bounce
		default:
			letterY = baseOffsetY
		}

		barLevel := energy*0.85 + totalEnergy*0.15

		for py := range fontH {
			row := glyph[py]
			for px := range fontW {
				if row&(1<<(fontW-1-px)) == 0 {
					continue
				}

				for sy := range scaleY {
					for sx := range scaleX {
						dx := letterX + px*scaleX + sx
						dy := letterY + py*scaleY + sy

						if dx < 0 || dx >= dotCols || dy < 0 || dy >= dotRows {
							continue
						}

						visible := false
						pixelNormY := float64(fontH*scaleY-1-(py*scaleY+sy)) / float64(fontH*scaleY)

						switch effect {
						case EffectDissolve:
							fill := energy*energy*0.75 + 0.15
							if scatterHash(li, py*scaleY+sy, px*scaleX+sx, v.frame) <= fill {
								visible = true
							}

						case EffectEqualizer:
							if pixelNormY <= barLevel {
								visible = true
							} else {
								sparkle := scatterHash(li, py*scaleY+sy, px*scaleX+sx, v.frame)
								threshold := (pixelNormY - barLevel) * 3.0
								if sparkle > threshold {
									visible = true
								}
							}

						case EffectBounce:
							fill := energy*energy*0.6 + 0.3
							if scatterHash(li, py*scaleY+sy, px*scaleX+sx, v.frame) <= fill {
								visible = true
							}

						case EffectRain:
							dropPhase := (float64(v.frame)*0.15 + float64(px*scaleX+sx)*0.3)
							dropY := math.Mod(dropPhase, float64(fontH*scaleY))
							dist := math.Abs(float64(py*scaleY+sy) - dropY)
							if dist < float64(scaleY)*1.5*energy+0.5 {
								visible = true
							}
							if energy > 0.3 && scatterHash(li, py*scaleY+sy, px*scaleX+sx, v.frame) < energy*0.4 {
								visible = true
							}

						case EffectMatrix:
							// Falling streams inside letters.
							col := px*scaleX + sx
							speed := 0.1 + float64(col%5)*0.04
							fall := math.Mod(float64(v.frame)*speed+float64(col)*2.7, float64(fontH*scaleY))
							trail := energy*float64(scaleY)*3 + 1
							dist := float64(py*scaleY+sy) - fall
							if dist < 0 {
								dist += float64(fontH * scaleY)
							}
							if dist < trail {
								visible = true
							}

						case EffectFlame:
							// Fire rising from bottom.
							flicker := scatterHash(li, px*scaleX+sx, py*scaleY+sy, v.frame/2) * 0.3
							flameH := (1.0-pixelNormY)*energy + flicker
							if flameH > 0.25 {
								visible = true
							}

						case EffectPulse:
							// Pulsate: all pixels visible at energy-based threshold.
							centerY := float64(fontH*scaleY) / 2
							centerX := float64(fontW*scaleX) / 2
							dy2 := float64(py*scaleY+sy) - centerY
							dx2 := float64(px*scaleX+sx) - centerX
							dist := math.Sqrt(dx2*dx2+dy2*dy2) / (centerY + centerX)
							pulse := math.Sin(float64(v.frame)*0.1) * 0.3
							if dist < energy+pulse+0.2 {
								visible = true
							}

						case EffectGlitch:
							// Random blocks appear/disappear with energy.
							block := scatterHash(li, (py*scaleY+sy)/3, (px*scaleX+sx)/3, v.frame/4)
							if block < energy*0.9+0.1 {
								visible = true
							}
							// Horizontal shift glitch.
							if energy > 0.5 {
								shift := int(scatterHash(li, py, 0, v.frame/3) * 4)
								sdx := dx + shift
								if sdx >= 0 && sdx < dotCols {
									grid[dy*dotCols+sdx] = true
								}
							}

						case EffectWave:
							// Sine wave sweeps through letters.
							waveY := 0.5 + 0.4*math.Sin(float64(v.frame)*0.08+float64(dx)*0.05)
							waveH := energy*0.4 + 0.1
							if math.Abs(pixelNormY-waveY) < waveH {
								visible = true
							}

						case EffectBinary:
							// Streaming 0/1 pattern.
							col := px*scaleX + sx
							row2 := py*scaleY + sy
							speed := 0.12 + float64(col%3)*0.03
							phase := int(float64(v.frame)*speed+float64(col)*1.3) + row2
							bit := scatterHash(li, phase%20, col, uint64(phase/20))
							if bit < energy*0.8+0.15 {
								visible = true
							}

						case EffectScatter:
							// Sparkle: random dots twinkle.
							sparkle := scatterHash(li, py*scaleY+sy, px*scaleX+sx, v.frame)
							threshold := energy * energy * 0.7
							if sparkle < threshold+0.05 {
								visible = true
							}
							// Persistent glow at high energy.
							if energy > 0.6 {
								visible = true
							}

						case EffectLightning:
							// Electric bolts: vertical jagged lines.
							col := px*scaleX + sx
							boltX := int(scatterHash(li, 0, 0, v.frame/3) * float64(fontW*scaleX))
							dist := math.Abs(float64(col - boltX))
							if dist < energy*float64(scaleX)*1.5+0.5 {
								jitter := scatterHash(li, py*scaleY+sy, col, v.frame)
								if jitter < energy*0.8+0.1 {
									visible = true
								}
							}
							// Base visibility.
							if energy > 0.4 && scatterHash(li, py*scaleY+sy, col, v.frame+99) < energy*0.3 {
								visible = true
							}
						}

						if visible {
							grid[dy*dotCols+dx] = true
						}
					}
				}
			}
		}
	}

	// Convert to Braille.
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
