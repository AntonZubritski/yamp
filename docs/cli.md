# CLI Flags

Override any config option for a single session without editing `~/.config/yamp/config.toml`. Flags can appear before or after file/URL arguments.

## Playback

```sh
yamp --volume -5 track.mp3          # volume in dB [-30, +6]
yamp --shuffle ~/Music              # enable shuffle
yamp --repeat all ~/Music           # repeat mode: off, all, one
yamp --mono track.mp3               # downmix to mono
yamp --no-mono track.mp3            # force stereo
yamp --auto-play ~/Music            # start playback immediately
```

## Audio engine

```sh
yamp --sample-rate 48000 track.mp3      # output sample rate (22050, 44100, 48000, 96000, 192000)
yamp --buffer-ms 200 track.mp3          # speaker buffer in ms (50–500)
yamp --resample-quality 1 track.mp3     # resample quality factor (1–4)
yamp --bit-depth 32 track.m4a           # PCM bit depth: 16 (default) or 32 (lossless)
```

## Appearance

```sh
yamp --compact ~/Music                     # cap width at 80 columns
yamp --eq-preset "Bass Boost" ~/Music
```

## Search

Search and play a track directly from the command line (requires [yt-dlp](https://github.com/yt-dlp/yt-dlp)):

```sh
yamp search "never gonna give you up"       # search YouTube
yamp search-sc "lofi beats"                  # search SoundCloud
```

Press `f` in the player to search YouTube interactively, or `F` (Shift+F) to search SoundCloud.

## General

| Flag | Short | Description |
|------|-------|-------------|
| `--help` | `-h` | Show help and exit |
| `--version` | `-v` | Print version and exit |
| `--upgrade` | | Update to the latest release |

## Mixing flags and files

Flags can appear anywhere — before, after, or between positional arguments:

```sh
yamp --shuffle track.mp3 --volume -5
yamp track.mp3 --repeat all --mono ~/Music
```

## Flag reference

| Flag | Type | Default | Range / Values |
|------|------|---------|----------------|
| `--volume` | float | 0 | -30 to +6 dB |
| `--shuffle` | bool | false | |
| `--repeat` | string | off | off, all, one |
| `--mono` / `--no-mono` | bool | false | |
| `--auto-play` | bool | false | |
| `--compact` | bool | false | |
| `--theme` | string | | theme name |
| `--eq-preset` | string | | preset name |
| `--sample-rate` | int | 44100 | 22050, 44100, 48000, 96000, 192000 |
| `--buffer-ms` | int | 100 | 50–500 |
| `--resample-quality` | int | 4 | 1–4 |
| `--bit-depth` | int | 16 | 16, 32 |

CLI flags override config file values for the current session only. They are not persisted.
