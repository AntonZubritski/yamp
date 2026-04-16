//go:build linux

package mpris

// StartMediaKeyListener is a no-op on Linux — MPRIS D-Bus handles media keys.
func StartMediaKeyListener(send func(interface{})) {}
