# YAMP — Changelog

## Session 2026-03-25

### Yandex Music Integration
- Yandex Music provider with playlists, liked tracks, radio stations
- "My Wave" (Моя волна) — infinite personalized radio with auto-refill
- Direct MP3 URL resolution bypassing yt-dlp for Yandex tracks
- In-app OAuth flow: press Enter → browser opens → paste token → done
- Token sanitization (strips control characters from clipboard paste)
- Token saved to `~/.config/yamp/config.toml` automatically
- `config.SaveYandexToken()` — writes token into `[yandex]` section

### All Providers Always Visible
- Navidrome, Plex, Spotify, YouTube shown even without config
- `setupProvider` placeholder — shows setup instructions in TUI when selected
- Each provider displays bilingual config instructions (RU/EN)
- Yandex Music interactive OAuth, others show config.toml snippets

### Custom Text Visualizers
- `Ctrl+E` — opens custom visualizer editor overlay
- 3-screen flow: list → text input → effect selection with live preview
- 12 effects: Dissolve, Equalizer, Bounce, Rain, Matrix, Flame, Pulse, Glitch, Wave, Binary, Scatter, Lightning
- 5×7 bitmap font: A-Z, 0-9, space, punctuation (`vis_font.go`)
- `renderCustomText()` — universal renderer for any text + any effect
- Saved to `~/.config/yamp/custom_vis.json`
- Appears in `v` cycle as "Custom" mode
- Minimum 20% pixel visibility so text is always readable

### Zubritski Signature Visualizer
- Built-in "Zubritski" mode — ZUBRITSKI text with EQ bars inside letters
- Each letter mapped to a frequency band
- Energy-driven fill from bottom + sparkle dissolve above bar line
- Gentle bounce and traveling wave animation

### Clipboard Support
- `Ctrl+V` in token input reads system clipboard
- Windows: PowerShell `Get-Clipboard`
- macOS: `pbpaste`
- Linux: `xclip` / `xsel`
- `ui/clipboard.go` — cross-platform clipboard reader

### Web UI for Mobile (PWA)
- `yamp --web` — starts HTTP server on port 8080
- Beautiful dark PWA interface (Apple Music style)
- Glassmorphism, animated visualizer bars, smooth controls
- Play/Pause, Next, Prev, Seek, Volume
- Provider selector tabs
- Auto-refreshes state every 500ms
- Add to iPhone Home Screen = native app experience
- Prints local IP addresses for phone connection
- `web/server.go` — Go HTTP server with REST API
- `web/index.html` — single-file PWA with embedded CSS/JS

### Global Media Keys
- Windows: `RegisterHotKey` API — Play/Pause, Next, Prev, Stop
- macOS: CoreGraphics event tap for media key interception
- Linux: already handled by MPRIS D-Bus
- Works in background, no configuration needed
- `mpris/mediakeys_windows.go`, `mpris/mediakeys_darwin.go`

### UI Fixes
- Added space between pause icon and "Paused" text
- Fixed Yandex token input capturing Ctrl+V control characters
- Fixed empty token `""` in config still creating API client
- Sanitized tokens on read and write (strip `\x00` etc.)

### GitHub & CI/CD
- Repository: github.com/AntonZubritski/yamp
- Single-commit clean history, sole contributor
- GitHub Actions release workflow (`release.yml`):
  - macOS arm64 on `macos-latest`
  - Linux amd64 on `ubuntu-latest`
  - Windows amd64 on `windows-latest`
  - Auto-creates release with binaries + checksums
- Removed old workflows (pages, aur)
- Removed CNAME, site/ from original project

### Install Scripts
- `install.sh` — one-command install for macOS/Linux (curl | sh)
- `install.ps1` — one-command install for Windows (irm | iex)
- Auto-detects OS and architecture
- Downloads binary from GitHub Releases
- Installs to PATH

### README
- Full bilingual documentation (English + Russian)
- Step-by-step install for Windows, macOS, Linux (Debian, Fedora, Arch)
- "Download binary" section — no Go needed
- "One command install" section
- Keybindings table
- All providers with config examples
- Custom visualizer editor documentation
- Mobile (iPhone/Android) section with PWA instructions
- Termux instructions for Android
- Configuration reference with all TOML options
- Troubleshooting (PipeWire/PulseAudio audio fix)

### Files Created
```
external/yandex/client.go          — Yandex Music API client
external/yandex/provider.go        — Yandex Music provider
ui/clipboard.go                    — Cross-platform clipboard
ui/custom_vis_overlay.go           — Custom visualizer editor UI
ui/setup_provider.go               — Placeholder provider for unconfigured services
ui/vis_custom.go                   — Custom text visualizer renderer (12 effects)
ui/vis_font.go                     — 5×7 bitmap font (A-Z, 0-9, symbols)
ui/vis_zubritski.go                — Zubritski signature visualizer
web/index.html                     — PWA mobile remote control
web/server.go                      — HTTP API server
install.sh                         — macOS/Linux installer
install.ps1                        — Windows installer
mpris/mediakeys_windows.go         — Windows global media keys
mpris/mediakeys_darwin.go          — macOS global media keys
mpris/mediakeys_linux.go           — Linux stub (MPRIS handles it)
.github/workflows/release.yml      — CI/CD release pipeline
CHANGELOG.md                       — This file
```

### Files Modified
```
main.go                            — Yandex provider, --web flag, media keys
config/config.go                   — SaveYandexToken()
config/flags.go                    — --web flag
ui/model.go                        — Yandex auth state, custom vis, imports
ui/keys.go                         — Yandex token input, Ctrl+E, media key handlers
ui/view.go                         — Yandex auth screen, setup instructions, custom vis overlay
ui/state.go                        — customVisOverlay struct
ui/visualizer.go                   — VisZubritski, VisCustom modes, custom config loading
ui/keymap.go                       — Ctrl+E keybinding entry
README.md                          — Complete rewrite with bilingual docs
install.sh                         — Updated repo URL
.gitignore                         — Added *.exe
```
