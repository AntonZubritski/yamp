package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"yamp/config"
)

// openCustomVisOverlay opens the custom visualizer editor overlay.
func (m *Model) openCustomVisOverlay() {
	m.customVis.visible = true
	m.customVis.screen = 0 // list
	m.customVis.configs = loadCustomVis()
	m.customVis.cursor = 0
	m.customVis.textBuf = ""
	m.customVis.effectIdx = 0
}

// handleCustomVisKey processes keys in the custom visualizer overlay.
func (m *Model) handleCustomVisKey(msg tea.KeyMsg) tea.Cmd {
	switch m.customVis.screen {
	case 0:
		return m.handleCustomVisListKey(msg)
	case 1:
		return m.handleCustomVisTextKey(msg)
	case 2:
		return m.handleCustomVisEffectKey(msg)
	}
	return nil
}

// screen 0: list of custom visualizers + "New..." + "Delete"
func (m *Model) handleCustomVisListKey(msg tea.KeyMsg) tea.Cmd {
	count := len(m.customVis.configs) + 1 // +1 for "+ New..."
	switch msg.String() {
	case "esc", "ctrl+e":
		m.customVis.visible = false
	case "ctrl+c":
		m.customVis.visible = false
		return m.quit()
	case "up", "k":
		if m.customVis.cursor > 0 {
			m.customVis.cursor--
		}
	case "down", "j":
		if m.customVis.cursor < count-1 {
			m.customVis.cursor++
		}
	case "enter":
		if m.customVis.cursor == len(m.customVis.configs) {
			// "+ New..." selected
			m.customVis.screen = 1
			m.customVis.textBuf = ""
		} else {
			// Activate selected custom visualizer
			cfg := m.customVis.configs[m.customVis.cursor]
			m.vis.CustomConfigs = m.customVis.configs
			m.vis.CustomIdx = m.customVis.cursor
			m.vis.Mode = VisCustom
			_ = config.Save("visualizer", fmt.Sprintf("%q", "Custom"))
			m.status.text = fmt.Sprintf("Visualizer: %s", cfg.Text)
			m.status.ttl = 60
			m.customVis.visible = false
		}
	case "d", "D", "delete":
		// Delete selected custom visualizer
		if m.customVis.cursor < len(m.customVis.configs) {
			m.customVis.configs = append(
				m.customVis.configs[:m.customVis.cursor],
				m.customVis.configs[m.customVis.cursor+1:]...,
			)
			_ = saveCustomVis(m.customVis.configs)
			m.vis.CustomConfigs = m.customVis.configs
			if m.customVis.cursor > 0 && m.customVis.cursor >= len(m.customVis.configs) {
				m.customVis.cursor--
			}
		}
	}
	return nil
}

// screen 1: text input
func (m *Model) handleCustomVisTextKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyEscape:
		m.customVis.screen = 0
	case tea.KeyEnter:
		text := strings.TrimSpace(m.customVis.textBuf)
		if text != "" {
			m.customVis.textBuf = text
			m.customVis.screen = 2
			m.customVis.effectIdx = 0
			// Start preview
			m.vis.Mode = VisCustom
		}
	case tea.KeyBackspace:
		m.customVis.textBuf = removeLastRune(m.customVis.textBuf)
	case tea.KeySpace:
		m.customVis.textBuf += " "
	case tea.KeyRunes:
		for _, r := range msg.Runes {
			if r >= 0x20 && r != 0x7f {
				m.customVis.textBuf += string(r)
			}
		}
	}
	return nil
}

// screen 2: effect selection with live preview
func (m *Model) handleCustomVisEffectKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		m.customVis.screen = 1
	case "up", "k":
		if m.customVis.effectIdx > 0 {
			m.customVis.effectIdx--
		}
	case "down", "j":
		if m.customVis.effectIdx < int(effectCount)-1 {
			m.customVis.effectIdx++
		}
	case "enter":
		// Save the new custom visualizer
		cfg := CustomVisConfig{
			Text:   strings.ToUpper(strings.TrimSpace(m.customVis.textBuf)),
			Effect: effectNames[m.customVis.effectIdx],
		}
		m.customVis.configs = append(m.customVis.configs, cfg)
		_ = saveCustomVis(m.customVis.configs)
		m.vis.CustomConfigs = m.customVis.configs
		m.vis.CustomIdx = len(m.customVis.configs) - 1
		m.vis.Mode = VisCustom
		_ = config.Save("visualizer", fmt.Sprintf("%q", "Custom"))
		m.status.text = fmt.Sprintf("Saved visualizer: %s", cfg.Text)
		m.status.ttl = 60
		m.customVis.visible = false
	}
	return nil
}

// renderCustomVisOverlay renders the custom visualizer editor overlay.
func (m Model) renderCustomVisOverlay() string {
	var lines []string

	// Live preview on screens 1 and 2.
	if text, eff, ok := m.customVisPreview(); ok {
		preview := m.vis.renderCustomText(text, eff, m.vis.Analyze(nil))
		if preview != "" {
			lines = append(lines, preview)
			lines = append(lines, "")
		}
	}

	switch m.customVis.screen {
	case 0: // List
		lines = append(lines, titleStyle.Render(" Custom Visualizers / Свои визуализаторы "))
		lines = append(lines, "")
		if len(m.customVis.configs) == 0 {
			lines = append(lines, dimStyle.Render("  Пока нет / No custom visualizers yet"))
			lines = append(lines, "")
		}
		for i, cfg := range m.customVis.configs {
			prefix, style := "  ", dimStyle
			if i == m.customVis.cursor {
				prefix = "> "
				style = playlistSelectedStyle
			}
			label := fmt.Sprintf("%s%s [%s]", prefix, cfg.Text, cfg.Effect)
			lines = append(lines, style.Render(label))
		}
		// "+ New..." option
		prefix, style := "  ", dimStyle
		if m.customVis.cursor == len(m.customVis.configs) {
			prefix = "> "
			style = playlistSelectedStyle
		}
		lines = append(lines, style.Render(prefix+"+ New... / Создать..."))
		lines = append(lines, "")
		lines = append(lines, dimStyle.Render("  Enter — выбрать/select | D — удалить/delete | Esc — закрыть/close"))

	case 1: // Text input
		lines = append(lines, titleStyle.Render(" New Custom Visualizer / Новый визуализатор "))
		lines = append(lines, "")
		lines = append(lines, dimStyle.Render("  Введите текст (будет заглавными):"))
		lines = append(lines, dimStyle.Render("  Enter text (will be uppercased):"))
		lines = append(lines, "")
		display := m.customVis.textBuf
		if display == "" {
			display = "_"
		} else {
			display += "_"
		}
		lines = append(lines, playlistSelectedStyle.Render("  > "+display))
		lines = append(lines, "")
		lines = append(lines, dimStyle.Render("  Enter — далее/next | Esc — назад/back"))

	case 2: // Effect selection with live preview
		lines = append(lines, titleStyle.Render(fmt.Sprintf(" Эффект для: %s ", strings.ToUpper(m.customVis.textBuf))))
		lines = append(lines, "")
		for i := range int(effectCount) {
			prefix, style := "  ", dimStyle
			if i == m.customVis.effectIdx {
				prefix = "> "
				style = playlistSelectedStyle
			}
			lines = append(lines, style.Render(fmt.Sprintf("%s%-12s %s", prefix, effectNames[i], effectDescs[i])))
		}
		lines = append(lines, "")
		lines = append(lines, dimStyle.Render("  Enter — сохранить/save | Esc — назад/back"))
	}

	return m.centerOverlay(strings.Join(lines, "\n"))
}

// customVisPreviewBands returns the text and effect to use for live preview
// while the custom vis overlay is open.
func (m Model) customVisPreview() (string, CustomVisEffect, bool) {
	if !m.customVis.visible {
		return "", 0, false
	}
	if m.customVis.screen == 2 {
		return strings.ToUpper(m.customVis.textBuf), CustomVisEffect(m.customVis.effectIdx), true
	}
	if m.customVis.screen == 1 && m.customVis.textBuf != "" {
		return strings.ToUpper(m.customVis.textBuf), EffectDissolve, true
	}
	return "", 0, false
}
