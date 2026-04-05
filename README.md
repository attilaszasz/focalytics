# focalytics

<!-- release-version: 0.0.6 -->

focalytics is a local-first command-line tool that scans a photo archive and generates a self-contained HTML report about how that archive was shot over time.

It is built for photographers, archivists, and anyone sitting on years of images spread across nested folders. Point it at an archive root, let it scan locally, and it will produce a report covering shooting activity, gear usage, focal lengths, exposure patterns, and other metadata-driven trends.

The tool is designed to stay simple:

- It runs entirely offline.
- It does not modify the source archive.
- It does not require a database, a server, or external services.
- It works on macOS, Linux, and Windows through prebuilt GitHub Release artifacts.

## What It Does

focalytics recursively scans a directory of photos, reads embedded metadata and related sidecar files where available, falls back gracefully when metadata is incomplete, and renders a standalone report you can open in any browser.

The generated report includes:

- An overview of archive size, date range, and most-used gear.
- A timeline of shooting activity across years and days.
- Camera and lens usage breakdowns.
- Technical analytics for focal length, aperture, shutter speed, and ISO.

## Quick Start

After installing the binary, run focalytics against the root of your archive:

```bash
focalytics ~/Pictures/Archive
```

You can also use the explicit subcommand form:

```bash
focalytics run ~/Pictures/Archive
```

The command scans the target directory and writes a timestamped HTML report into the current working directory, for example:

```text
focalytics_report_20260405_1045.html
```

When you run the command in a terminal, focalytics shows a live stage-based progress view while it scans and renders. When stdout or stderr is redirected, it stays quiet apart from warnings and the final report path so shell pipelines remain script-friendly.

## Install

Download the correct archive for your platform from the latest GitHub Release:

https://github.com/attilaszasz/focalytics/releases/latest

Choose the asset that matches your operating system and CPU architecture:

- `darwin arm64`: Apple Silicon Macs
- `darwin amd64`: Intel Macs
- `linux arm64`: ARM64 Linux systems
- `linux amd64`: x86_64 Linux systems
- `windows arm64`: Windows on ARM
- `windows amd64`: 64-bit Intel or AMD Windows systems

### macOS

Apple Silicon (arm64):

```bash
curl -fsSLO https://github.com/attilaszasz/focalytics/releases/download/v0.0.6/focalytics_v0.0.6_darwin_arm64.tar.gz
tar -xzf focalytics_v0.0.6_darwin_arm64.tar.gz
sudo mv focalytics /usr/local/bin/
```

For Intel Macs, replace `arm64` with `amd64` in the commands above.

If you prefer a user-local install instead of `/usr/local/bin`:

```bash
mkdir -p "$HOME/.local/bin"
mv focalytics "$HOME/.local/bin/"
```

Add `~/.local/bin` to your `PATH` in your shell profile if it is not already there.

### Linux

x86_64 (amd64):

```bash
curl -fsSLO https://github.com/attilaszasz/focalytics/releases/download/v0.0.6/focalytics_v0.0.6_linux_amd64.tar.gz
tar -xzf focalytics_v0.0.6_linux_amd64.tar.gz
sudo mv focalytics /usr/local/bin/
```

For ARM64 systems, replace `amd64` with `arm64` in the commands above.

User-local install:

```bash
mkdir -p "$HOME/.local/bin"
mv focalytics "$HOME/.local/bin/"
```

### Windows

64-bit Intel or AMD (amd64) via PowerShell:

```powershell
Invoke-WebRequest -Uri https://github.com/attilaszasz/focalytics/releases/download/v0.0.6/focalytics_v0.0.6_windows_amd64.zip -OutFile focalytics.zip
Expand-Archive .\focalytics.zip -DestinationPath .\focalytics
New-Item -ItemType Directory -Force "$HOME\bin" | Out-Null
Move-Item .\focalytics\focalytics.exe "$HOME\bin\focalytics.exe" -Force
```

For Windows on ARM, replace `amd64` with `arm64` in the URL above.

If `$HOME\bin` is not on your user `Path`, add it so `focalytics.exe` can be started from any terminal.

## Verify Your Download

Each release includes a checksum file named `focalytics_v0.0.6_checksums.txt`.

To verify a downloaded archive:

```bash
curl -fsSLO https://github.com/attilaszasz/focalytics/releases/download/v0.0.6/focalytics_v0.0.6_checksums.txt
sha256sum -c focalytics_v0.0.6_checksums.txt --ignore-missing
```

On macOS, use `shasum -a 256` instead of `sha256sum`.

## Build From Source

If you prefer to compile the tool yourself, install Go `1.25.8` or newer compatible with the module and build from the `src` directory:

```bash
cd src
go build .
```

That produces a local `focalytics` binary in `src` on macOS and Linux, or `focalytics.exe` on Windows.

## Notes

- focalytics is distributed today through GitHub Releases.
- The tool is intended for offline analysis of local archives, not cloud sync, metadata editing, or hosted reporting.
