package ui

import "yamp/playlist"

// setupProvider is a placeholder provider shown in the UI for services
// that are not yet configured. Selecting it triggers the sign-in / setup
// flow which displays configuration instructions.
type setupProvider struct {
	name string
}

// NewSetupProvider creates a placeholder provider for services not yet configured.
func NewSetupProvider(name string) *setupProvider {
	return &setupProvider{name: name}
}

func (p *setupProvider) Name() string { return p.name }

func (p *setupProvider) Playlists() ([]playlist.PlaylistInfo, error) {
	return nil, playlist.ErrNeedsAuth
}

func (p *setupProvider) Tracks(string) ([]playlist.Track, error) {
	return nil, playlist.ErrNeedsAuth
}
