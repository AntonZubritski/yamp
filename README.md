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

---

## Install

### Download ready-made binary (easiest)

Go to [Releases](https://github.com/AntonZubritski/yamp/releases/latest) and download the file for your OS:

| OS | File |
|----|------|
| Windows | `yamp-windows-amd64.exe` |
| macOS (Intel) | `yamp-darwin-amd64` |
| macOS (Apple Silicon M1/M2/M3) | `yamp-darwin-arm64` |
| Linux | `yamp-linux-amd64` |

**Windows:** rename to `yamp.exe`, put in any folder, open terminal in that folder and run `.\yamp.exe`

**macOS / Linux:** make executable and move to PATH:

```sh
chmod +x yamp-*
sudo mv yamp-* /usr/local/bin/yamp
```

Then run `yamp` from any terminal.

---

### Build from source

> Requires [Go](https://go.dev/dl/) installed. **If `go` is not found** — download and install Go first from https://go.dev/dl/ (click the big blue button, install, restart terminal).

### Windows (PowerShell)

**Step 1.** Install Go from https://go.dev/dl/ (download the `.msi` file, run it, restart PowerShell)

**Step 2.** Open PowerShell (or Windows Terminal) and clone the repo:

```powershell
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
```

**Step 3.** Build:

```powershell
go build -o yamp.exe .
```

**Step 3.** Make it available from anywhere. Choose one option:

**Option A** — Copy to a folder already in PATH:

```powershell
copy yamp.exe "$env:LOCALAPPDATA\Microsoft\WinGet\Links\yamp.exe"
```

**Option B** — Add the current folder to PATH (run as Admin):

```powershell
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$PWD", "User")
```

**Step 4.** Open a **new** terminal and run:

```powershell
yamp
```

Or play a specific folder:

```powershell
yamp C:\Users\YourName\Music
```

---

### macOS

**Step 1.** Open Terminal and install Go (if not installed):

```sh
brew install go
```

Or download from https://go.dev/dl/

**Step 2.** Clone and build:

```sh
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
```

**Step 3.** Install system-wide so you can run it from anywhere:

```sh
sudo mv yamp /usr/local/bin/
```

**Step 4.** Run:

```sh
yamp
```

Or play your music folder:

```sh
yamp ~/Music
```

> macOS uses CoreAudio — no extra audio dependencies needed.

---

### Linux (Debian / Ubuntu / Mint)

**Step 1.** Install Go and audio development headers:

```sh
sudo apt update
sudo apt install -y golang libasound2-dev git
```

> If your distro's Go is too old, install from https://go.dev/dl/ instead.

**Step 2.** Clone and build:

```sh
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
```

**Step 3.** Install system-wide:

```sh
sudo mv yamp /usr/local/bin/
```

**Step 4.** Run:

```sh
yamp
```

Or play your music:

```sh
yamp ~/Music
```

**If you hear no sound** (PipeWire/PulseAudio systems):

```sh
# PipeWire
sudo apt install pipewire-alsa

# PulseAudio
sudo apt install pulseaudio-alsa
```

---

### Linux (Fedora)

**Step 1.** Install Go and audio development headers:

```sh
sudo dnf install golang alsa-lib-devel git
```

**Step 2.** Clone and build:

```sh
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
```

**Step 3.** Install system-wide:

```sh
sudo mv yamp /usr/local/bin/
```

**Step 4.** Run:

```sh
yamp ~/Music
```

---

### Linux (Arch / Manjaro)

**Step 1.** Install Go and audio headers:

```sh
sudo pacman -S go alsa-lib git
```

**Step 2.** Clone and build:

```sh
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
```

**Step 3.** Install system-wide:

```sh
sudo mv yamp /usr/local/bin/
```

**Step 4.** Run:

```sh
yamp ~/Music
```

**If you hear no sound:**

```sh
# PipeWire
sudo pacman -S pipewire-alsa

# PulseAudio
sudo pacman -S pulseaudio-alsa
```

---

### Optional dependencies (all platforms)

These are not required but unlock extra features:

| Dependency | What it does | Install |
|-----------|-------------|---------|
| [ffmpeg](https://ffmpeg.org/) | AAC, ALAC, Opus, WMA playback | `brew install ffmpeg` / `sudo apt install ffmpeg` / `sudo pacman -S ffmpeg` |
| [yt-dlp](https://github.com/yt-dlp/yt-dlp) | YouTube, SoundCloud, Bandcamp, Bilibili | `pip install yt-dlp` or `brew install yt-dlp` |

---

## How to run

After installation, open any terminal and type:

```sh
# Launch with TUI (browse providers, radio, playlists)
yamp

# Play a folder
yamp ~/Music

# Play specific files
yamp song1.mp3 song2.flac

# Play a URL (stream, podcast, YouTube)
yamp https://example.com/stream

# Play a YouTube video
yamp https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

On **Windows** if `yamp` is not found, use the full path:

```powershell
C:\path\to\yamp.exe
```

Or add it to PATH (see install steps above).

---

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
| `Ctrl+K` | Show all keybindings |
| `q` | Quit |

---

## Providers

All providers are always visible in the player. If one is not configured, selecting it shows setup instructions right in the app.

### Yandex Music

No manual config needed. Just:

1. Select **Yandex Music** in the player
2. Press **Enter** — browser opens with Yandex login
3. Log in, copy the token from the URL after `#access_token=`
4. Paste it in the app and press **Enter**
5. Done — your playlists and "My Wave" infinite radio are ready

### Spotify

Requires Spotify Premium.

1. Go to [developer.spotify.com/dashboard](https://developer.spotify.com/dashboard)
2. Create an app, copy the Client ID
3. Edit `~/.config/yamp/config.toml` (on Windows: `%USERPROFILE%\.config\yamp\config.toml`):

```toml
[spotify]
client_id = "your_client_id"
```

4. Restart yamp

### Navidrome

Edit `~/.config/yamp/config.toml`:

```toml
[navidrome]
url = "https://your-server.com"
user = "username"
password = "password"
```

Restart yamp.

### Plex

Edit `~/.config/yamp/config.toml`:

```toml
[plex]
url = "http://192.168.1.10:32400"
token = "your-plex-token"
```

To get your Plex token: open Plex Web, go to any media, click "Get Info" > "View XML" — the token is in the URL.

Restart yamp.

### YouTube / YouTube Music

1. Install [yt-dlp](https://github.com/yt-dlp/yt-dlp): `pip install yt-dlp`
2. Edit `~/.config/yamp/config.toml`:

```toml
[ytmusic]
enabled = true
```

3. Restart yamp

---

## Custom Visualizers

Press `Ctrl+E` to open the editor:

1. Select **+ New...**
2. Type any text (your name, band name, anything)
3. Pick from 12 effects:

| Effect | Description |
|--------|------------|
| Dissolve | Dots appear/disappear with energy |
| Equalizer | EQ bars fill letters from bottom |
| Bounce | Letters bounce with the beat |
| Rain | Drops fall through letter shapes |
| Matrix | Falling matrix characters |
| Flame | Fire rising through letters |
| Pulse | Letters pulse with beat |
| Glitch | Random block corruption |
| Wave | Sine wave through letters |
| Binary | 0/1 streaming through letters |
| Scatter | Sparkle particles inside |
| Lightning | Electric bolts inside letters |

4. See live preview while picking
5. Press **Enter** to save

Saved visualizers appear in the `v` cycle and are stored in `~/.config/yamp/custom_vis.json`.

---

## Configuration

Config file: `~/.config/yamp/config.toml`

On Windows: `%USERPROFILE%\.config\yamp\config.toml`

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

---

## Author

[Anton Zubritski](https://github.com/AntonZubritski)

## License

MIT

---
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
- Яндекс Музыка — плейлисты и бесконечное радио («Моя волна») с OAuth прямо в приложении
- Медиасерверы Navidrome и Plex
- 22+ визуализатора спектра, включая фирменный стиль Zubritski
- Кастомные текстовые визуализаторы — введите любое слово, выберите из 12 эффектов (Ctrl+E)
- Параметрический 10-полосный эквалайзер с пресетами
- Отображение текстов песен
- Темы оформления и настройка горячих клавиш
- Все провайдеры видны с инструкциями по настройке

---

## Установка

### Скачать готовый файл (проще всего)

Перейдите в [Releases](https://github.com/AntonZubritski/yamp/releases/latest) и скачайте файл для вашей ОС:

| ОС | Файл |
|----|------|
| Windows | `yamp-windows-amd64.exe` |
| macOS (Intel) | `yamp-darwin-amd64` |
| macOS (Apple Silicon M1/M2/M3) | `yamp-darwin-arm64` |
| Linux | `yamp-linux-amd64` |

**Windows:** переименуйте в `yamp.exe`, положите в любую папку, откройте терминал в этой папке и запустите `.\yamp.exe`

**macOS / Linux:** сделайте исполняемым и переместите в PATH:

```sh
chmod +x yamp-*
sudo mv yamp-* /usr/local/bin/yamp
```

Затем запускайте `yamp` из любого терминала.

---

### Сборка из исходников

> Нужен установленный [Go](https://go.dev/dl/). **Если `go` не найден** — скачайте и установите Go с https://go.dev/dl/ (нажмите большую синюю кнопку, установите, перезапустите терминал).

### Windows (PowerShell)

**Шаг 1.** Установите Go с https://go.dev/dl/ (скачайте `.msi` файл, запустите его, перезапустите PowerShell)

**Шаг 2.** Откройте PowerShell (или Windows Terminal) и клонируйте репозиторий:

```powershell
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
```

**Шаг 3.** Соберите:

```powershell
go build -o yamp.exe .
```

**Шаг 3.** Сделайте доступным отовсюду. Выберите один вариант:

**Вариант А** — Скопировать в папку которая уже в PATH:

```powershell
copy yamp.exe "$env:LOCALAPPDATA\Microsoft\WinGet\Links\yamp.exe"
```

**Вариант Б** — Добавить текущую папку в PATH (от имени администратора):

```powershell
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$PWD", "User")
```

**Шаг 4.** Откройте **новый** терминал и запустите:

```powershell
yamp
```

Или укажите папку с музыкой:

```powershell
yamp C:\Users\ВашеИмя\Music
```

---

### macOS

**Шаг 1.** Откройте Terminal и установите Go (если не установлен):

```sh
brew install go
```

Или скачайте с https://go.dev/dl/

**Шаг 2.** Клонируйте и соберите:

```sh
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
```

**Шаг 3.** Установите в систему чтобы запускать из любого места:

```sh
sudo mv yamp /usr/local/bin/
```

**Шаг 4.** Запустите:

```sh
yamp
```

Или укажите папку:

```sh
yamp ~/Music
```

> macOS использует CoreAudio — дополнительные зависимости не нужны.

---

### Linux (Debian / Ubuntu / Mint)

**Шаг 1.** Установите Go и заголовки для аудио:

```sh
sudo apt update
sudo apt install -y golang libasound2-dev git
```

> Если версия Go в вашем дистрибутиве слишком старая, установите с https://go.dev/dl/

**Шаг 2.** Клонируйте и соберите:

```sh
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
```

**Шаг 3.** Установите в систему:

```sh
sudo mv yamp /usr/local/bin/
```

**Шаг 4.** Запустите:

```sh
yamp
```

Или укажите папку с музыкой:

```sh
yamp ~/Music
```

**Если нет звука** (системы с PipeWire/PulseAudio):

```sh
# PipeWire
sudo apt install pipewire-alsa

# PulseAudio
sudo apt install pulseaudio-alsa
```

---

### Linux (Fedora)

**Шаг 1.** Установите Go и заголовки для аудио:

```sh
sudo dnf install golang alsa-lib-devel git
```

**Шаг 2.** Клонируйте и соберите:

```sh
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
```

**Шаг 3.** Установите:

```sh
sudo mv yamp /usr/local/bin/
```

**Шаг 4.** Запустите:

```sh
yamp ~/Music
```

---

### Linux (Arch / Manjaro)

**Шаг 1.** Установите Go и заголовки для аудио:

```sh
sudo pacman -S go alsa-lib git
```

**Шаг 2.** Клонируйте и соберите:

```sh
git clone https://github.com/AntonZubritski/yamp.git
cd yamp
go build -o yamp .
```

**Шаг 3.** Установите:

```sh
sudo mv yamp /usr/local/bin/
```

**Шаг 4.** Запустите:

```sh
yamp ~/Music
```

**Если нет звука:**

```sh
# PipeWire
sudo pacman -S pipewire-alsa

# PulseAudio
sudo pacman -S pulseaudio-alsa
```

---

### Опциональные зависимости (все платформы)

Не обязательны, но открывают дополнительные возможности:

| Зависимость | Что даёт | Установка |
|------------|---------|-----------|
| [ffmpeg](https://ffmpeg.org/) | Воспроизведение AAC, ALAC, Opus, WMA | `brew install ffmpeg` / `sudo apt install ffmpeg` / `sudo pacman -S ffmpeg` |
| [yt-dlp](https://github.com/yt-dlp/yt-dlp) | YouTube, SoundCloud, Bandcamp, Bilibili | `pip install yt-dlp` или `brew install yt-dlp` |

---

## Как запускать

После установки откройте любой терминал и введите:

```sh
# Запустить TUI (провайдеры, радио, плейлисты)
yamp

# Играть папку с музыкой
yamp ~/Music

# Играть конкретные файлы
yamp song1.mp3 song2.flac

# Играть URL (стрим, подкаст, YouTube)
yamp https://example.com/stream

# Играть видео с YouTube
yamp https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

На **Windows** если `yamp` не найден, используйте полный путь:

```powershell
C:\путь\к\yamp.exe
```

Или добавьте в PATH (см. шаги установки выше).

---

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
| `R` | Каталог радиостанций (30 000+) |
| `N` | Браузер Navidrome |
| `o` | Файловый браузер |
| `p` | Менеджер плейлистов |
| `y` | Показать текст песни |
| `J` | Перейти к времени |
| `S` | Сохранить/скачать трек |
| `Ctrl+K` | Показать все горячие клавиши |
| `q` | Выход |

---

## Провайдеры

Все провайдеры всегда видны в плеере. Если какой-то не настроен — при выборе покажется инструкция прямо в приложении.

### Яндекс Музыка

Ничего заранее настраивать не нужно:

1. Выберите **Yandex Music** в плеере
2. Нажмите **Enter** — откроется браузер со страницей входа Яндекса
3. Войдите, скопируйте токен из URL после `#access_token=`
4. Вставьте его в приложение и нажмите **Enter**
5. Готово — ваши плейлисты и «Моя волна» доступны

### Spotify

Нужен Spotify Premium.

1. Зайдите на [developer.spotify.com/dashboard](https://developer.spotify.com/dashboard)
2. Создайте приложение, скопируйте Client ID
3. Отредактируйте `~/.config/yamp/config.toml` (на Windows: `%USERPROFILE%\.config\yamp\config.toml`):

```toml
[spotify]
client_id = "ваш_client_id"
```

4. Перезапустите yamp

### Navidrome

Отредактируйте `~/.config/yamp/config.toml`:

```toml
[navidrome]
url = "https://ваш-сервер.com"
user = "логин"
password = "пароль"
```

Перезапустите yamp.

### Plex

Отредактируйте `~/.config/yamp/config.toml`:

```toml
[plex]
url = "http://192.168.1.10:32400"
token = "ваш-plex-токен"
```

Как получить токен Plex: откройте Plex Web, зайдите в любой медиафайл, нажмите «Get Info» > «View XML» — токен будет в URL.

Перезапустите yamp.

### YouTube / YouTube Music

1. Установите [yt-dlp](https://github.com/yt-dlp/yt-dlp): `pip install yt-dlp`
2. Отредактируйте `~/.config/yamp/config.toml`:

```toml
[ytmusic]
enabled = true
```

3. Перезапустите yamp

---

## Кастомные визуализаторы

Нажмите `Ctrl+E` чтобы открыть редактор:

1. Выберите **+ New...**
2. Введите любой текст (ваше имя, название группы, что угодно)
3. Выберите из 12 эффектов:

| Эффект | Описание |
|--------|----------|
| Dissolve | Точки появляются/исчезают с энергией |
| Equalizer | EQ столбцы заполняют буквы снизу |
| Bounce | Буквы прыгают в такт |
| Rain | Капли падают через буквы |
| Matrix | Падающие символы матрицы |
| Flame | Огонь поднимается через буквы |
| Pulse | Буквы пульсируют с битом |
| Glitch | Случайные блоки |
| Wave | Синусоида через буквы |
| Binary | Потоки 0/1 |
| Scatter | Мерцающие частицы |
| Lightning | Молнии внутри букв |

4. Смотрите живое превью при выборе
5. Нажмите **Enter** чтобы сохранить

Сохранённые визуализаторы появляются в цикле `v` и хранятся в `~/.config/yamp/custom_vis.json`.

---

## Конфигурация

Файл конфига: `~/.config/yamp/config.toml`

На Windows: `%USERPROFILE%\.config\yamp\config.toml`

```toml
# Провайдер по умолчанию
provider = "radio"

# Тема
theme = "ethereal"

# Визуализатор
visualizer = "Zubritski"

# Пресет EQ
eq_preset = "Electronic"

# Шаг перемотки для Shift+Стрелка (секунды)
seek_large_step_sec = 30

[yandex]
token = "ваш_токен"

[spotify]
client_id = "ваш_client_id"

[navidrome]
url = "https://ваш-сервер.com"
user = "логин"
password = "пароль"

[plex]
url = "http://192.168.1.10:32400"
token = "ваш-plex-токен"

[ytmusic]
enabled = true
```

---

## Автор

[Anton Zubritski](https://github.com/AntonZubritski)

## Лицензия

MIT
