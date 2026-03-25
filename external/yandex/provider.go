package yandex

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"yamp/config"
	"yamp/internal/browser"
	"yamp/playlist"
)

const radioPrefix = "radio:"

// Provider implements playlist.Provider for Yandex Music.
// Playlists() returns the user's playlists, liked tracks, and radio stations
// (including "Моя волна"). Tracks are returned with Yandex Music URLs;
// direct MP3 URLs are resolved lazily at playback time via ResolveDirectURL.
type Provider struct {
	client        *Client
	mu            sync.Mutex
	playlistCache []playlist.PlaylistInfo
	trackCache    map[string][]playlist.Track
}

const oauthURL = "https://oauth.yandex.ru/authorize?response_type=token&client_id=23cabbbdc6cd418abb4b39c32c41195d"

// New returns a Provider for the given OAuth token.
// If token is empty the provider still works but will request sign-in.
func New(token string) *Provider {
	p := &Provider{
		trackCache: make(map[string][]playlist.Track),
	}
	// Strip control chars that may have leaked from a previous paste.
	token = strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, strings.TrimSpace(token))
	if token != "" {
		p.client = NewClient(token)
	}
	return p
}

// Name returns the display name for the provider selector.
func (p *Provider) Name() string { return "Yandex Music" }

// NeedsAuth reports whether the provider has no token and requires sign-in.
func (p *Provider) NeedsAuth() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.client == nil
}

// SetToken configures the provider with a new OAuth token, saving it to config.
func (p *Provider) SetToken(token string) error {
	// Strip control characters that may leak from clipboard paste.
	clean := strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, strings.TrimSpace(token))
	if clean == "" {
		return fmt.Errorf("yandex: empty token")
	}
	if err := config.SaveYandexToken(clean); err != nil {
		return err
	}
	p.mu.Lock()
	p.client = NewClient(clean)
	p.playlistCache = nil
	p.mu.Unlock()
	return nil
}

// Authenticate opens the browser with the Yandex OAuth page.
// The actual token input is handled in the UI layer.
func (p *Provider) Authenticate() error {
	_ = browser.Open(oauthURL)
	return nil
}

// Playlists returns the user's playlists from Yandex Music.
// Radio stations (Моя волна, genre stations) are listed first,
// followed by "Liked Tracks", then user playlists.
func (p *Provider) Playlists() ([]playlist.PlaylistInfo, error) {
	if p.NeedsAuth() {
		return nil, playlist.ErrNeedsAuth
	}
	p.mu.Lock()
	if p.playlistCache != nil {
		cached := p.playlistCache
		p.mu.Unlock()
		return cached, nil
	}
	p.mu.Unlock()

	var lists []playlist.PlaylistInfo

	// Radio stations first (Моя волна etc.)
	stations, _ := p.client.RadioStations() // non-fatal if fails
	if len(stations) > 0 {
		for _, s := range stations {
			name := s.Name
			if s.ID == "user:onyourwave" {
				name = "\u266b \u041c\u043e\u044f \u0432\u043e\u043b\u043d\u0430" // ♫ Моя волна
			} else {
				name = "\u266a " + name // ♪ prefix for radio
			}
			lists = append(lists, playlist.PlaylistInfo{
				ID:   radioPrefix + s.ID,
				Name: name,
			})
		}
	} else {
		// Always show "Моя волна" even if dashboard fails.
		lists = append(lists, playlist.PlaylistInfo{
			ID:   radioPrefix + "user:onyourwave",
			Name: "\u266b \u041c\u043e\u044f \u0432\u043e\u043b\u043d\u0430",
		})
	}

	// Liked tracks
	lists = append(lists, playlist.PlaylistInfo{
		ID:   "likes",
		Name: "Liked Tracks",
	})

	// User playlists
	playlists, err := p.client.Playlists()
	if err != nil {
		return nil, err
	}
	uid, err := p.client.resolveUID()
	if err != nil {
		return nil, err
	}
	for _, pl := range playlists {
		ownerUID := fmt.Sprintf("%d", pl.Owner.UID)
		if ownerUID == "0" {
			ownerUID = uid
		}
		lists = append(lists, playlist.PlaylistInfo{
			ID:         fmt.Sprintf("%s:%d", ownerUID, pl.Kind),
			Name:       pl.Title,
			TrackCount: pl.TrackCount,
		})
	}

	p.mu.Lock()
	p.playlistCache = lists
	p.mu.Unlock()

	return lists, nil
}

// Tracks returns the tracks in the given playlist.
// playlistID is "likes", "ownerUID:kind", or "radio:type:tag".
func (p *Provider) Tracks(playlistID string) ([]playlist.Track, error) {
	// Radio stations are never cached — each call fetches a fresh batch.
	if strings.HasPrefix(playlistID, radioPrefix) {
		return p.loadRadioTracks(strings.TrimPrefix(playlistID, radioPrefix))
	}

	p.mu.Lock()
	if cached, ok := p.trackCache[playlistID]; ok {
		p.mu.Unlock()
		return cached, nil
	}
	p.mu.Unlock()

	var tracks []playlist.Track

	if playlistID == "likes" {
		t, err := p.loadLikedTracks()
		if err != nil {
			return nil, err
		}
		tracks = t
	} else {
		t, err := p.loadPlaylistTracks(playlistID)
		if err != nil {
			return nil, err
		}
		tracks = t
	}

	p.mu.Lock()
	p.trackCache[playlistID] = tracks
	p.mu.Unlock()

	return tracks, nil
}

func (p *Provider) loadLikedTracks() ([]playlist.Track, error) {
	ids, err := p.client.LikedTracks()
	if err != nil {
		return nil, err
	}
	if len(ids) > 200 {
		ids = ids[:200]
	}
	infos, err := p.client.TracksInfo(ids)
	if err != nil {
		return nil, err
	}
	return convertTracks(infos), nil
}

func (p *Provider) loadPlaylistTracks(playlistID string) ([]playlist.Track, error) {
	ownerUID, kindStr, ok := strings.Cut(playlistID, ":")
	if !ok {
		return nil, fmt.Errorf("yandex: invalid playlist ID %q", playlistID)
	}
	kind, err := strconv.Atoi(kindStr)
	if err != nil {
		return nil, fmt.Errorf("yandex: invalid playlist ID %q", playlistID)
	}
	infos, err := p.client.PlaylistTracks(ownerUID, kind)
	if err != nil {
		return nil, err
	}
	return convertTracks(infos), nil
}

func (p *Provider) loadRadioTracks(stationID string) ([]playlist.Track, error) {
	// Send "radioStarted" feedback to initialize the session.
	_ = p.client.RadioFeedback(stationID, "radioStarted", "")

	infos, err := p.client.RadioTracks(stationID)
	if err != nil {
		return nil, fmt.Errorf("yandex: radio %s: %w", stationID, err)
	}
	return convertTracks(infos), nil
}

// LoadRadioTracks fetches a batch of tracks from the given radio station.
// Exported for use as a callback in the UI model for auto-refill.
func (p *Provider) LoadRadioTracks(stationID string) ([]playlist.Track, error) {
	return p.loadRadioTracks(stationID)
}

// convertTracks turns Yandex API track info into playlist.Track values.
func convertTracks(infos []trackInfo) []playlist.Track {
	tracks := make([]playlist.Track, 0, len(infos))
	for _, t := range infos {
		trackID := fmt.Sprintf("%v", t.ID)
		if trackID == "" || trackID == "0" {
			continue
		}
		artist := ""
		if len(t.Artists) > 0 {
			artist = t.Artists[0].Name
		}
		album := ""
		year := 0
		if len(t.Albums) > 0 {
			album = t.Albums[0].Title
			year = t.Albums[0].Year
		}
		tracks = append(tracks, playlist.Track{
			Path:         "https://music.yandex.ru/track/" + trackID,
			Title:        t.Title,
			Artist:       artist,
			Album:        album,
			Year:         year,
			DurationSecs: t.DurationMs / 1000,
			Stream:       true,
		})
	}
	return tracks
}

// ResolveDirectURL resolves a direct MP3 stream URL for a track.
// trackPath should be like "https://music.yandex.ru/track/12345".
func (p *Provider) ResolveDirectURL(trackPath string) (string, error) {
	trackID := trackPath
	if idx := strings.LastIndex(trackPath, "/"); idx >= 0 {
		trackID = trackPath[idx+1:]
	}
	return p.client.DirectURL(trackID)
}
