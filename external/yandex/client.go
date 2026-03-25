// Package yandex implements a playlist.Provider for Yandex Music.
// It uses the unofficial Yandex Music API for metadata (playlists, tracks)
// and relies on yt-dlp for actual audio playback.
package yandex

import (
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	apiBase         = "https://api.music.yandex.net"
	maxResponseBody = 10 << 20 // 10 MB
)

// apiClient is used for all Yandex Music API calls.
var apiClient = &http.Client{Timeout: 30 * time.Second}

// Client speaks to the Yandex Music API.
type Client struct {
	token string // OAuth token
	uid   string // user ID, resolved on first call
}

// NewClient returns a Client for the given OAuth token.
func NewClient(token string) *Client {
	return &Client{token: token}
}

// playlistInfo is the JSON shape for a playlist in the API response.
type playlistInfo struct {
	Kind       int    `json:"kind"`
	Title      string `json:"title"`
	TrackCount int    `json:"trackCount"`
	Owner      struct {
		UID  int    `json:"uid"`
		Name string `json:"name"`
	} `json:"owner"`
}

// trackInfo is the JSON shape for a track in the API response.
type trackInfo struct {
	ID        any    `json:"id"` // can be int or string
	Title     string `json:"title"`
	DurationMs int   `json:"durationMs"`
	Artists   []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Albums []struct {
		Title string `json:"title"`
		Year  int    `json:"year"`
	} `json:"albums"`
}

// get issues an authenticated GET and decodes the JSON "result" field.
func (c *Client) get(path string, result any) error {
	req, err := http.NewRequest(http.MethodGet, apiBase+path, nil)
	if err != nil {
		return fmt.Errorf("yandex: %s: %w", path, err)
	}
	req.Header.Set("Authorization", "OAuth "+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := apiClient.Do(req)
	if err != nil {
		return fmt.Errorf("yandex: %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("yandex: token invalid or expired (HTTP %d)", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("yandex: %s: HTTP %s", path, resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBody))
	if err != nil {
		return fmt.Errorf("yandex: %s: %w", path, err)
	}

	// The Yandex Music API wraps responses in {"result": ...}.
	var envelope struct {
		Result json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("yandex: %s: parsing response: %w", path, err)
	}
	if envelope.Result == nil {
		// Some endpoints return the data at top level.
		return json.Unmarshal(body, result)
	}
	return json.Unmarshal(envelope.Result, result)
}

// resolveUID fetches and caches the current user's UID.
func (c *Client) resolveUID() (string, error) {
	if c.uid != "" {
		return c.uid, nil
	}
	var status struct {
		Account struct {
			UID int `json:"uid"`
		} `json:"account"`
	}
	if err := c.get("/account/status", &status); err != nil {
		return "", err
	}
	if status.Account.UID == 0 {
		return "", fmt.Errorf("yandex: could not resolve user ID from account status")
	}
	c.uid = fmt.Sprintf("%d", status.Account.UID)
	return c.uid, nil
}

// Playlists returns the user's playlists.
func (c *Client) Playlists() ([]playlistInfo, error) {
	uid, err := c.resolveUID()
	if err != nil {
		return nil, err
	}
	var playlists []playlistInfo
	if err := c.get("/users/"+uid+"/playlists/list", &playlists); err != nil {
		return nil, err
	}
	return playlists, nil
}

// PlaylistTracks returns tracks for the given playlist kind.
func (c *Client) PlaylistTracks(ownerUID string, kind int) ([]trackInfo, error) {
	var pl struct {
		Tracks []struct {
			Track trackInfo `json:"track"`
		} `json:"tracks"`
	}
	path := fmt.Sprintf("/users/%s/playlists/%d", ownerUID, kind)
	if err := c.get(path, &pl); err != nil {
		return nil, err
	}
	tracks := make([]trackInfo, 0, len(pl.Tracks))
	for _, t := range pl.Tracks {
		tracks = append(tracks, t.Track)
	}
	return tracks, nil
}

// LikedTracks returns the user's liked tracks.
func (c *Client) LikedTracks() ([]string, error) {
	uid, err := c.resolveUID()
	if err != nil {
		return nil, err
	}
	var result struct {
		Library struct {
			Tracks []struct {
				ID string `json:"id"`
			} `json:"tracks"`
		} `json:"library"`
	}
	if err := c.get("/users/"+uid+"/likes/tracks", &result); err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(result.Library.Tracks))
	for _, t := range result.Library.Tracks {
		ids = append(ids, t.ID)
	}
	return ids, nil
}

// TracksInfo fetches metadata for a batch of track IDs.
func (c *Client) TracksInfo(ids []string) ([]trackInfo, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	// POST /tracks with track IDs. This is one of the few POST endpoints.
	// We'll batch in groups to avoid overly large requests.
	const batchSize = 100
	var all []trackInfo
	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		batch := ids[i:end]
		tracks, err := c.postTracksInfo(batch)
		if err != nil {
			return nil, err
		}
		all = append(all, tracks...)
	}
	return all, nil
}

func (c *Client) postTracksInfo(ids []string) ([]trackInfo, error) {
	// Build form body: track-ids=id1,id2,...
	body := "track-ids="
	for i, id := range ids {
		if i > 0 {
			body += ","
		}
		body += id
	}

	req, err := http.NewRequest(http.MethodPost, apiBase+"/tracks", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "OAuth "+c.token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Body = io.NopCloser(io.LimitReader(
		readerFromString(body), int64(len(body)),
	))
	req.ContentLength = int64(len(body))

	resp, err := apiClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("yandex: /tracks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yandex: /tracks: HTTP %s", resp.Status)
	}

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBody))
	if err != nil {
		return nil, err
	}

	var envelope struct {
		Result []trackInfo `json:"result"`
	}
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return nil, err
	}
	return envelope.Result, nil
}

// radioStation describes a Yandex Music radio station.
type radioStation struct {
	ID   string // "type:tag", e.g. "user:onyourwave"
	Name string
}

// RadioStations returns the user's dashboard radio stations.
func (c *Client) RadioStations() ([]radioStation, error) {
	var dashboard struct {
		Stations []struct {
			Station struct {
				ID struct {
					Type string `json:"type"`
					Tag  string `json:"tag"`
				} `json:"id"`
				Name string `json:"name"`
			} `json:"station"`
		} `json:"stations"`
	}
	if err := c.get("/rotor/stations/dashboard", &dashboard); err != nil {
		return nil, err
	}
	stations := make([]radioStation, 0, len(dashboard.Stations))
	for _, s := range dashboard.Stations {
		stations = append(stations, radioStation{
			ID:   s.Station.ID.Type + ":" + s.Station.ID.Tag,
			Name: s.Station.Name,
		})
	}
	return stations, nil
}

// RadioTracks fetches a batch of tracks from the given radio station.
// stationID is "type:tag", e.g. "user:onyourwave".
func (c *Client) RadioTracks(stationID string) ([]trackInfo, error) {
	var result struct {
		Sequence []struct {
			Track trackInfo `json:"track"`
		} `json:"sequence"`
	}
	path := "/rotor/station/" + stationID + "/tracks?settings2=true"
	if err := c.get(path, &result); err != nil {
		return nil, err
	}
	tracks := make([]trackInfo, 0, len(result.Sequence))
	for _, s := range result.Sequence {
		tracks = append(tracks, s.Track)
	}
	return tracks, nil
}

// RadioFeedback sends playback feedback to the radio station so the algorithm
// can refine recommendations. feedbackType is e.g. "radioStarted", "trackStarted", "skip".
func (c *Client) RadioFeedback(stationID, feedbackType, trackID string) error {
	body := fmt.Sprintf(`{"type":"%s","timestamp":"%s"`, feedbackType, time.Now().UTC().Format(time.RFC3339))
	if trackID != "" {
		body += fmt.Sprintf(`,"from":"","trackId":"%s"`, trackID)
	}
	body += "}"

	req, err := http.NewRequest(http.MethodPost, apiBase+"/rotor/station/"+stationID+"/feedback", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "OAuth "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(readerFromString(body))
	req.ContentLength = int64(len(body))

	resp, err := apiClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// DirectURL resolves a direct MP3 stream URL for the given track ID.
// It picks the highest bitrate MP3 available.
func (c *Client) DirectURL(trackID string) (string, error) {
	// Step 1: get download info (list of available codecs/bitrates).
	var infos []struct {
		Codec           string `json:"codec"`
		BitrateInKbps   int    `json:"bitrateInKbps"`
		DownloadInfoURL string `json:"downloadInfoUrl"`
	}
	if err := c.get("/tracks/"+trackID+"/download-info", &infos); err != nil {
		return "", err
	}

	// Pick best MP3 (highest bitrate).
	bestIdx := -1
	bestBitrate := 0
	for i, info := range infos {
		if info.Codec == "mp3" && info.BitrateInKbps > bestBitrate {
			bestIdx = i
			bestBitrate = info.BitrateInKbps
		}
	}
	if bestIdx == -1 {
		// Fallback: pick any codec with highest bitrate.
		for i, info := range infos {
			if info.BitrateInKbps > bestBitrate {
				bestIdx = i
				bestBitrate = info.BitrateInKbps
			}
		}
	}
	if bestIdx == -1 {
		return "", fmt.Errorf("yandex: no download info for track %s", trackID)
	}

	// Step 2: fetch the download-info XML to get host/path/ts/s.
	req, err := http.NewRequest(http.MethodGet, infos[bestIdx].DownloadInfoURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := apiClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("yandex: download-info fetch: %w", err)
	}
	defer resp.Body.Close()

	var dlInfo struct {
		Host string `xml:"host"`
		Path string `xml:"path"`
		TS   string `xml:"ts"`
		S    string `xml:"s"`
	}
	if err := xml.NewDecoder(resp.Body).Decode(&dlInfo); err != nil {
		return "", fmt.Errorf("yandex: parsing download-info XML: %w", err)
	}

	// Step 3: build the signed URL.
	// Sign = md5("XGRlBW9FXlekgbPrRHuSiA" + path[1:] + s)
	sign := fmt.Sprintf("%x", md5.Sum([]byte("XGRlBW9FXlekgbPrRHuSiA"+strings.TrimPrefix(dlInfo.Path, "/")+dlInfo.S)))

	return fmt.Sprintf("https://%s/get-mp3/%s/%s%s", dlInfo.Host, sign, dlInfo.TS, dlInfo.Path), nil
}

type stringReader struct {
	s string
	i int
}

func readerFromString(s string) io.Reader {
	return &stringReader{s: s}
}

func (r *stringReader) Read(p []byte) (int, error) {
	if r.i >= len(r.s) {
		return 0, io.EOF
	}
	n := copy(p, r.s[r.i:])
	r.i += n
	return n, nil
}
