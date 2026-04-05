# Product Brief: `focalytics`

## 1. Product Vision & Summary
**`focalytics`** is a fast, local, privacy-first command-line utility built in Go. It recursively scans deeply nested directories containing decades of photographs, parses embedded metadata and sidecar files, and generates a comprehensive, beautiful, and completely self-contained HTML dashboard. 

The tool is designed to give photographers deep, data-driven insights into their shooting habits, gear usage, and technical progression over time. By adhering strictly to the Unix philosophy, it does one thing exceptionally well, requires zero external runtimes, operates entirely offline, and requires no configuration from the user.

## 2. Target Audience
* **Photography Enthusiasts & Hobbyists:** Looking to understand their shooting style and justify gear retention or acquisition.
* **Professional Photographers:** Reviewing their body of work across years of client shoots.
* **Archivists & Data Nerds:** Seeking a high-level, programmatic overview of massive digital photo collections via the terminal.

---

## 3. User Experience & CLI Design

The tool is designed for absolute simplicity. It is invoked directly against a target directory and requires zero additional flags.

**Usage:**
```bash
focalytics ~/Pictures/Archive
```
*Behavior:* Scans the provided directory and automatically generates a timestamped report in the current working directory with minute precision (e.g., `focalytics_report_20260405_1045.html`).

**Terminal UI (TUI):**
During execution, the terminal will display a high-quality progress interface showing:
* A smooth progress bar or spinner.
* Files scanned per second.
* Estimated time of completion (ETA).
* The current sub-directory being processed.

---

## 4. Core Features & Data Processing

* **Recursive Folder Scanning:** Seamlessly handles deeply nested directories with tens of thousands of files.
* **Multi-Tiered Metadata Extraction:**
    1. **Primary:** Embedded EXIF/IPTC data (using robust, fault-tolerant parsing).
    2. **Secondary:** XMP Sidecar files (reads `.xmp` files sharing the same base filename as the image).
    3. **Fallbacks:** OS-level file creation dates or regex-based directory parsing (e.g., extracting "2018" from `/Archive/2018/`) if all internal metadata is stripped.
* **Automated Sensor Size Normalization:** Automatically converts focal lengths to the 35mm full-frame equivalent. The tool prioritizes the `FocalLengthIn35mmFormat` tag, falling back to a hardcoded, internal camera model database to apply the correct crop factor transparently.
* **Pre-Aggregation Strategy:** To ensure the generated HTML remains lightweight, the Go backend aggregates all statistics in-memory (e.g., `{"50mm": 4500}`) rather than passing individual photo records to the frontend.
* **Graceful Degradation:** Corrupted files log a warning to `stderr` but **do not** halt the scanning process.

---

## 5. The Output: Static HTML Report

The output is a single, heavily optimized `.html` file. It relies on the `//go:embed` directive to inject CSS directly into the file. **It uses no external JavaScript charting libraries and makes no network calls to CDNs.**

### Visual Presentation
* **Aesthetics:** Sleek, modern "Dark Mode" UI mimicking professional photo editing software.
* **Handling Missing Data:** To prevent skewed visualizations, files missing metadata for a specific metric will NOT be charted. Instead, every chart will include a subtle footnote below it (e.g., *"Note: 1,402 photos were excluded from this chart due to missing lens metadata"*).

### Dashboard Contents

**A. The High-Level Overview (Hero Section)**
* **Total Photos Analyzed:** The sheer volume of the scanned archive.
* **Time Span:** "Photos taken between [First Date] and [Last Date]."
* **Top 3 Trophies:** Callouts for the most used camera, most used lens, and most used focal length.

**B. The Timeline (Shooting Habits)**
* **Photos per Year:** A CSS bar chart showing volume over the decades.
* **Activity Heatmap:** A GitHub-style daily contribution calendar (generated via inline SVG) showing active days/months.

**C. The Gear Locker**
* **Cameras Used:** A leaderboard of camera bodies based on total shutter actuations.
* **Lenses Used:** A breakdown of specific lenses utilized.
* **The "Prime vs. Zoom" Ratio:** The percentage of photos taken with fixed focal lengths versus variable zooms.

**D. Technical Analytics**
* **Focal Length Distribution (Normalized):** A histogram of the 35mm equivalent focal lengths, grouped into standard buckets (Ultra-wide < 24mm, 35mm, 50mm, 85mm, Telephoto > 200mm).
* **Aperture (f-stop) Usage:** A distribution chart showing how often lenses are shot wide open versus stopped down.
* **Shutter Speed:** Grouped into logical buckets (Long exposure > 1s, Handheld 1/60s–1/250s, Action > 1/1000s).
* **ISO Distribution:** A chart highlighting typical lighting environments.

---

## 6. Technical Architecture & Stack

* **Language:** Go (Golang)
* **CLI Framework:** `github.com/spf13/cobra` (Routing and standard CLI scaffolding)
* **Metadata Parser:** `github.com/dsoprea/go-exif` (plus standard XML parsing for `.xmp`)
* **Terminal UX:** `github.com/charmbracelet/bubbletea` and `bubbles` (For progress bars and TUI)
* **Templating:** Go Standard Library `html/template`
* **Asset Bundling:** Go Standard Library `embed` package (Zero-dependency binary)

---

## 7. Distribution & CI/CD Strategy

* **Compilation:** Go will be used to cross-compile single static binaries for macOS (Intel & Apple Silicon), Windows, and Linux.
* **Automation:** GitHub Actions will compile binaries, generate SHA256 hashes, create GitHub Releases, and automatically update package managers (Homebrew Tap and WinGet).
