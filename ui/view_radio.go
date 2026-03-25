package ui

import (
	"fmt"
	"strings"
)

// renderRadioCatalog renders the radio catalog browser overlay.
func (m Model) renderRadioCatalog() string {
	title := "R A D I O   C A T A L O G"
	if m.radioCatalog.showFavorites {
		title = "R A D I O   F A V O R I T E S"
	}
	lines := []string{titleStyle.Render(title), ""}

	// Search bar (hidden in favorites view)
	if m.radioCatalog.showFavorites {
		count := len(m.radioCatalog.stations)
		lines = append(lines, dimStyle.Render(fmt.Sprintf("  %d favorite stations", count)), "")
	} else if m.radioCatalog.searching {
		prompt := "  Search: " + m.radioCatalog.query + "_"
		lines = append(lines, playlistSelectedStyle.Render(prompt), "")
	} else if m.radioCatalog.query != "" {
		lines = append(lines, dimStyle.Render("  Search: "+m.radioCatalog.query), "")
	} else {
		lines = append(lines, dimStyle.Render("  Top stations by votes"), "")
	}

	// Loading state
	if m.radioCatalog.loading {
		lines = append(lines, dimStyle.Render("  Loading stations..."))
		lines = append(lines, "", m.radioCatalogHelp())
		return m.centerOverlay(strings.Join(lines, "\n"))
	}

	// Error state
	if m.radioCatalog.err != "" {
		lines = append(lines, errorStyle.Render("  "+m.radioCatalog.err))
		lines = append(lines, "", m.radioCatalogHelp())
		return m.centerOverlay(strings.Join(lines, "\n"))
	}

	// Empty state
	if len(m.radioCatalog.stations) == 0 {
		lines = append(lines, dimStyle.Render("  No stations found."))
		lines = append(lines, "", m.radioCatalogHelp())
		return m.centerOverlay(strings.Join(lines, "\n"))
	}

	// Station list
	maxVisible := m.plVisible
	if maxVisible < 5 {
		maxVisible = 5
	}

	scroll := m.radioCatalog.scroll
	rendered := 0
	for i := scroll; i < len(m.radioCatalog.stations) && rendered < maxVisible; i++ {
		s := m.radioCatalog.stations[i]
		star := ""
		if m.radioCatalog.favorites != nil && m.radioCatalog.favorites.Contains(s.URL) {
			star = "* "
		}
		label := star + s.Name
		if s.Bitrate > 0 {
			label += fmt.Sprintf(" [%dk]", s.Bitrate)
		}
		if s.Country != "" {
			label += " · " + s.Country
		}
		label = truncate(label, panelWidth-6)
		lines = append(lines, cursorLine(label, i == m.radioCatalog.cursor))
		rendered++
	}
	lines = padLines(lines, maxVisible, rendered)

	// Footer
	lines = append(lines, "",
		dimStyle.Render(fmt.Sprintf("  %d/%d stations",
			m.radioCatalog.cursor+1, len(m.radioCatalog.stations))))
	lines = append(lines, "", m.radioCatalogHelp())

	return m.centerOverlay(strings.Join(lines, "\n"))
}

// radioCatalogHelp returns the help line for the radio catalog.
func (m Model) radioCatalogHelp() string {
	if m.radioCatalog.searching {
		return helpKey("Enter", "Search ") + helpKey("Esc", "Cancel")
	}
	if m.radioCatalog.showFavorites {
		return helpKey("↑↓", "Navigate ") +
			helpKey("Enter", "Play ") +
			helpKey("f", "Unfav ") +
			helpKey("a", "Append ") +
			helpKey("F", "Catalog ") +
			helpKey("Esc", "Close")
	}
	return helpKey("↑↓", "Navigate ") +
		helpKey("Enter", "Play ") +
		helpKey("f", "Fav ") +
		helpKey("/", "Search ") +
		helpKey("a", "Append ") +
		helpKey("F", "Favs ") +
		helpKey("Esc", "Close")
}
