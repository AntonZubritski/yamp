// Package tomlutil provides shared helpers for minimal TOML parsing
// used by the local and radio providers.
package tomlutil

import "strconv"

// Unquote strips surrounding double quotes from a TOML string value,
// handling escape sequences (written by Go's %q format verb).
func Unquote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		if u, err := strconv.Unquote(s); err == nil {
			return u
		}
		// Fall back to naive strip if Unquote fails.
		return s[1 : len(s)-1]
	}
	return s
}
