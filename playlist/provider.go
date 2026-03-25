package playlist

import "errors"

// ErrNeedsAuth is returned by providers that require interactive sign-in
// before they can be used.
var ErrNeedsAuth = errors.New("sign-in required")

// PlaylistInfo describes a playlist with its name and track count.
type PlaylistInfo struct {
	ID         string
	Name       string
	TrackCount int
}

// Provider is the interface for playlist sources (radio, Navidrome, Spotify, etc.).
type Provider interface {
	// Name returns the display name of this provider.
	Name() string

	// Playlists returns the available playlists from this provider.
	Playlists() ([]PlaylistInfo, error)

	// Tracks returns the tracks in the given playlist.
	Tracks(playlistID string) ([]Track, error)
}

// Authenticator is optionally implemented by providers that require sign-in.
type Authenticator interface {
	Authenticate() error
}
