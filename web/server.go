// Package web provides an HTTP server that serves a mobile-friendly PWA
// remote control for the yamp player. Start with ListenAndServe.
package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"time"

	"yamp/player"
	"yamp/playlist"
)

//go:embed index.html
var staticFS embed.FS

// PlayerState is the JSON response for GET /api/state.
type PlayerState struct {
	Playing   bool    `json:"playing"`
	Paused    bool    `json:"paused"`
	Track     string  `json:"track"`
	Artist    string  `json:"artist"`
	Album     string  `json:"album"`
	Position  float64 `json:"position"`
	Duration  float64 `json:"duration"`
	Volume    float64 `json:"volume"`
	Streaming bool    `json:"streaming"`
}

// Server holds references needed by the HTTP handlers.
type Server struct {
	player   *player.Player
	playlist *playlist.Playlist
}

// ListenAndServe starts the web server on the given port.
// It blocks until the server is shut down.
func ListenAndServe(p *player.Player, pl *playlist.Playlist, port int) error {
	s := &Server{player: p, playlist: pl}
	mux := http.NewServeMux()

	// Serve the PWA.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, _ := staticFS.ReadFile("index.html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	// API endpoints.
	mux.HandleFunc("/api/state", s.handleState)
	mux.HandleFunc("/api/play", s.handlePlay)
	mux.HandleFunc("/api/pause", s.handlePause)
	mux.HandleFunc("/api/next", s.handleNext)
	mux.HandleFunc("/api/prev", s.handlePrev)
	mux.HandleFunc("/api/volume", s.handleVolume)
	mux.HandleFunc("/api/seek", s.handleSeek)

	addr := fmt.Sprintf(":%d", port)

	// Print all local IPs so user can connect from phone.
	fmt.Printf("\n  YAMP Web UI running at:\n")
	fmt.Printf("  http://localhost:%d\n", port)
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				fmt.Printf("  http://%s:%d\n", ipnet.IP.String(), port)
			}
		}
	}
	fmt.Printf("\n  Open on your phone to control playback.\n\n")

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return srv.ListenAndServe()
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	t, _ := s.playlist.Current()
	pos := s.player.Position().Seconds()
	dur := s.player.Duration().Seconds()

	st := PlayerState{
		Playing:   s.player.IsPlaying(),
		Paused:    s.player.IsPaused(),
		Track:     t.Title,
		Artist:    t.Artist,
		Album:     t.Album,
		Position:  pos,
		Duration:  dur,
		Volume:    dbToLinear(s.player.Volume()),
		Streaming: t.Stream,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(st)
}

func (s *Server) handlePlay(w http.ResponseWriter, r *http.Request) {
	if s.player.IsPaused() {
		s.player.TogglePause()
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handlePause(w http.ResponseWriter, r *http.Request) {
	if !s.player.IsPaused() {
		s.player.TogglePause()
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleNext(w http.ResponseWriter, r *http.Request) {
	// Next track is handled by the UI model, not directly here.
	// For now, just acknowledge.
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handlePrev(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleVolume(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Volume float64 `json:"volume"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	s.player.SetVolume(linearToDB(body.Volume))
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSeek(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Position float64 `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	s.player.Seek(time.Duration(body.Position) * time.Second)
	w.WriteHeader(http.StatusOK)
}

// dbToLinear converts dB (-30..+6) to linear (0..1).
func dbToLinear(db float64) float64 {
	v := (db + 30) / 36
	return math.Max(0, math.Min(1, v))
}

// linearToDB converts linear (0..1) to dB (-30..+6).
func linearToDB(v float64) float64 {
	return v*36 - 30
}
