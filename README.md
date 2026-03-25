```
██    ██  █████  ███    ███ ██████
 ██  ██  ██   ██ ████  ████ ██   ██
  ████   ███████ ██ ████ ██ ██████
   ██    ██   ██ ██  ██  ██ ██
   ██    ██   ██ ██      ██ ██
```

A retro terminal music player. Play local files, streams, podcasts, YouTube, YouTube Music, SoundCloud, Bilibili, Spotify, Yandex Music, Navidrome, and Plex with 22+ spectrum visualizers, custom text visualizers, parametric EQ, and playlist management.

## Features

- Local files (MP3, FLAC, OGG, WAV, AAC, ALAC, Opus, WMA via ffmpeg)
- Internet radio — 30,000+ stations from [Radio Browser](https://www.radio-browser.info/)
- YouTube, YouTube Music, SoundCloud, Bandcamp, Bilibili (via yt-dlp)
- Spotify Connect (Premium required)
- Yandex Music — playlists and infinite radio ("My Wave") with in-app OAuth
- Navidrome and Plex media servers
- 22+ spectrum visualizers including Zubritski signature style
- Custom text visualizers — type any word, pick from 12 effects (Ctrl+E)
- Parametric 10-band EQ with presets
- Lyrics display
- Themes and custom keybindings
- All providers visible with setup instructions for unconfigured ones

## Install

### Windows

```powershell
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp.exe .
```

Copy `yamp.exe` to a directory in your PATH, or run directly:

```powershell
.\yamp.exe ~\Music
```

### macOS

```sh
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
sudo mv yamp /usr/local/bin/
```

No extra dependencies — CoreAudio is used.

### Linux (Debian/Ubuntu)

```sh
sudo apt install libasound2-dev
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
sudo mv yamp /usr/local/bin/
```

### Linux (Fedora)

```sh
sudo dnf install alsa-lib-devel
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
sudo mv yamp /usr/local/bin/
```

### Linux (Arch)

```sh
sudo pacman -S alsa-lib
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
sudo mv yamp /usr/local/bin/
```

### Prerequisites

- [Go](https://go.dev/dl/) 1.25.5 or later

### Optional runtime dependencies

- [ffmpeg](https://ffmpeg.org/) — for AAC, ALAC, Opus, and WMA playback
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) — for YouTube, SoundCloud, Bandcamp, and Bilibili

## Quick Start

```sh
yamp ~/Music                     # play a directory
yamp *.mp3 *.flac               # play files
yamp https://example.com/stream  # play a URL
```

## Keybindings

| Key | Action |
|-----|--------|
| `Space` | Play / Pause |
| `> .` | Next track |
| `< ,` | Previous track |
| `← →` | Seek ±5s |
| `Shift+← →` | Seek ±30s |
| `+ -` | Volume up/down |
| `v` | Cycle visualizer |
| `V` | Full-screen visualizer |
| `Ctrl+E` | Custom visualizer editor |
| `e` | Cycle EQ preset |
| `t` | Choose theme |
| `z` | Toggle shuffle |
| `r` | Cycle repeat |
| `/` | Search playlist |
| `R` | Radio catalog (30,000+ stations) |
| `N` | Navidrome browser |
| `o` | Open file browser |
| `p` | Playlist manager |
| `y` | Show lyrics |
| `J` | Jump to time |
| `S` | Save/download track |
| `Ctrl+K` | All keybindings |
| `q` | Quit |

## Providers

All providers are always visible in the player. Unconfigured ones show setup instructions when selected.

### Yandex Music

Select **Yandex Music** in the player and press Enter. The app will:
1. Open your browser with the Yandex OAuth page
2. Ask you to paste the token from the URL
3. Save it to `~/.config/yamp/config.toml`

"My Wave" provides infinite personalized radio.

### Spotify

Requires Spotify Premium. Create an app at [developer.spotify.com/dashboard](https://developer.spotify.com/dashboard), then add to config:

```toml
[spotify]
client_id = "your_client_id"
```

### Navidrome

```toml
[navidrome]
url = "https://your-server.com"
user = "username"
password = "password"
```

### Plex

```toml
[plex]
url = "http://192.168.1.10:32400"
token = "your-plex-token"
```

### YouTube / YouTube Music

Requires [yt-dlp](https://github.com/yt-dlp/yt-dlp). Enable in config:

```toml
[ytmusic]
enabled = true
```

## Custom Visualizers

Press `Ctrl+E` to open the custom visualizer editor:

1. Select **+ New...**
2. Type any text (e.g. your name, band name)
3. Choose from 12 effects:
   - **Dissolve** — dots appear/disappear with energy
   - **Equalizer** — EQ bars fill letters from bottom
   - **Bounce** — letters bounce with the beat
   - **Rain** — drops fall through letter shapes
   - **Matrix** — falling matrix characters
   - **Flame** — fire rising through letters
   - **Pulse** — letters pulse with beat
   - **Glitch** — random block corruption
   - **Wave** — sine wave through letters
   - **Binary** — 0/1 streaming through letters
   - **Scatter** — sparkle particles inside
   - **Lightning** — electric bolts inside letters
4. Live preview shown while selecting
5. Saved to `~/.config/yamp/custom_vis.json`

## Configuration

Config file location: `~/.config/yamp/config.toml`

```toml
# Default provider on startup
provider = "radio"

# Theme
theme = "ethereal"

# Visualizer
visualizer = "Zubritski"

# EQ preset
eq_preset = "Electronic"

# Seek step for Shift+Arrow (seconds)
seek_large_step_sec = 30

[yandex]
token = "your_token"

[spotify]
client_id = "your_client_id"

[navidrome]
url = "https://your-server.com"
user = "username"
password = "password"

[plex]
url = "http://192.168.1.10:32400"
token = "your-plex-token"

[ytmusic]
enabled = true
```

## Troubleshooting

**No audio output (silence with no errors)**

On Linux with PipeWire or PulseAudio, install the ALSA bridge:

```sh
# PipeWire
sudo pacman -S pipewire-alsa        # Arch
sudo apt install pipewire-alsa       # Debian/Ubuntu

# PulseAudio
sudo pacman -S pulseaudio-alsa       # Arch
```

## Author

[Anton Zubritski](https://github.com/AntonZubritski)

## License

MIT

---

# README (Русский)

```
██    ██  █████  ███    ███ ██████
 ██  ██  ██   ██ ████  ████ ██   ██
  ████   ███████ ██ ████ ██ ██████
   ██    ██   ██ ██  ██  ██ ██
   ██    ██   ██ ██      ██ ██
```

Ретро-музыкальный плеер для терминала. Воспроизводит локальные файлы, интернет-радио, подкасты, YouTube, YouTube Music, SoundCloud, Bilibili, Spotify, Яндекс Музыку, Navidrome и Plex. 22+ визуализатора спектра, кастомные текстовые визуализаторы, параметрический эквалайзер и управление плейлистами.

## Возможности

- Локальные файлы (MP3, FLAC, OGG, WAV, AAC, ALAC, Opus, WMA через ffmpeg)
- Интернет-радио — более 30 000 станций из [Radio Browser](https://www.radio-browser.info/)
- YouTube, YouTube Music, SoundCloud, Bandcamp, Bilibili (через yt-dlp)
- Spotify Connect (нужен Premium)
- Яндекс Музыка — плейлисты и бесконечное радио («Моя волна») с OAuth в приложении
- Медиасерверы Navidrome и Plex
- 22+ визуализатора спектра, включая фирменный стиль Zubritski
- Кастомные текстовые визуализаторы — введите любое слово, выберите из 12 эффектов (Ctrl+E)
- Параметрический 10-полосный эквалайзер с пресетами
- Отображение текстов песен
- Темы оформления и настройка горячих клавиш
- Все провайдеры видны с инструкциями по настройке

## Установка

### Windows

```powershell
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp.exe .
```

Скопируйте `yamp.exe` в директорию из PATH, или запускайте напрямую:

```powershell
.\yamp.exe ~\Music
```

### macOS

```sh
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
sudo mv yamp /usr/local/bin/
```

Дополнительные зависимости не нужны — используется CoreAudio.

### Linux (Debian/Ubuntu)

```sh
sudo apt install libasound2-dev
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
sudo mv yamp /usr/local/bin/
```

### Linux (Fedora)

```sh
sudo dnf install alsa-lib-devel
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
sudo mv yamp /usr/local/bin/
```

### Linux (Arch)

```sh
sudo pacman -S alsa-lib
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
sudo mv yamp /usr/local/bin/
```

### Требования

- [Go](https://go.dev/dl/) 1.25.5 или новее

### Опциональные зависимости

- [ffmpeg](https://ffmpeg.org/) — для воспроизведения AAC, ALAC, Opus, WMA
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) — для YouTube, SoundCloud, Bandcamp, Bilibili

## Быстрый старт

```sh
yamp ~/Music                     # играть директорию
yamp *.mp3 *.flac               # играть файлы
yamp https://example.com/stream  # играть URL
```

## Горячие клавиши

| Клавиша | Действие |
|---------|----------|
| `Space` | Воспроизведение / Пауза |
| `> .` | Следующий трек |
| `< ,` | Предыдущий трек |
| `← →` | Перемотка ±5с |
| `Shift+← →` | Перемотка ±30с |
| `+ -` | Громкость |
| `v` | Переключить визуализатор |
| `V` | Полноэкранный визуализатор |
| `Ctrl+E` | Редактор кастомных визуализаторов |
| `e` | Переключить пресет EQ |
| `t` | Выбрать тему |
| `z` | Перемешать |
| `r` | Повтор |
| `/` | Поиск в плейлисте |
| `R` | Каталог радиостанций |
| `N` | Браузер Navidrome |
| `o` | Файловый браузер |
| `p` | Менеджер плейлистов |
| `y` | Показать текст песни |
| `J` | Перейти к времени |
| `S` | Сохранить/скачать трек |
| `Ctrl+K` | Все горячие клавиши |
| `q` | Выход |

## Провайдеры

Все провайдеры всегда видны в плеере. Ненастроенные показывают инструкции при выборе.

### Яндекс Музыка

Выберите **Yandex Music** и нажмите Enter. Приложение:
1. Откроет браузер с OAuth-страницей Яндекса
2. Попросит вставить токен из URL
3. Сохранит его в `~/.config/yamp/config.toml`

«Моя волна» — бесконечное персонализированное радио.

### Spotify

Нужен Spotify Premium. Создайте приложение на [developer.spotify.com/dashboard](https://developer.spotify.com/dashboard):

```toml
[spotify]
client_id = "ваш_client_id"
```

### Navidrome

```toml
[navidrome]
url = "https://ваш-сервер.com"
user = "имя_пользователя"
password = "пароль"
```

### Plex

```toml
[plex]
url = "http://192.168.1.10:32400"
token = "ваш-plex-токен"
```

### YouTube / YouTube Music

Нужен [yt-dlp](https://github.com/yt-dlp/yt-dlp):

```toml
[ytmusic]
enabled = true
```

## Кастомные визуализаторы

Нажмите `Ctrl+E`:

1. Выберите **+ New...**
2. Введите текст (например, ваше имя)
3. Выберите из 12 эффектов:
   - **Dissolve** — точки появляются/исчезают
   - **Equalizer** — EQ столбцы заполняют буквы
   - **Bounce** — буквы прыгают в такт
   - **Rain** — капли падают через буквы
   - **Matrix** — падающие символы матрицы
   - **Flame** — огонь поднимается через буквы
   - **Pulse** — буквы пульсируют с битом
   - **Glitch** — случайные блоки
   - **Wave** — синусоида через буквы
   - **Binary** — потоки 0/1
   - **Scatter** — мерцающие частицы
   - **Lightning** — молнии внутри букв
4. Живое превью при выборе
5. Сохраняется в `~/.config/yamp/custom_vis.json`

## Автор

[Anton Zubritski](https://github.com/AntonZubritski)

## Лицензия

MIT
