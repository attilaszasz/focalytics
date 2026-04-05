# Research: Render Offline Report
> Feature | 2026-04-05 | Purpose: inform self-contained HTML report generation for E005

## Self-Contained Templating
- **Decision**: Use Go `html/template` with `embed.FS` for templates and inline CSS generation.
- **Rationale**: This keeps the report safe, offline, and bundled into the binary without runtime asset dependencies.
- **Rejected**: `text/template` and manual escaping because they weaken HTML safety and increase rendering risk.
- **Pitfalls**: Do not pass untrusted values as trusted template HTML or JavaScript fragments.
- **Sources**: https://pkg.go.dev/html/template, https://pkg.go.dev/embed

## Offline Visualization Strategy
- **Decision**: Combine semantic HTML sections with inline SVG charts and script-free disclosure elements for exclusion notes.
- **Rationale**: Native browser primitives cover bar charts, heatmaps, and summary marks without external libraries or network access.
- **Rejected**: Canvas or JS chart libraries because they add runtime complexity and violate the no-CDN goal.
- **Pitfalls**: Do not hide essential counts only inside graphics; keep labels and notes visible in text.
- **Sources**: https://developer.mozilla.org/en-US/docs/Web/SVG/Element/svg, https://developer.mozilla.org/en-US/docs/Web/HTML/Element/details

## Output Path Rules
- **Decision**: Write one report file into the current working directory using a deterministic timestamped name with `.html` suffix.
- **Rationale**: This preserves the one-command CLI promise from the product brief while keeping the artifact easy to locate and revisit.
- **Rejected**: Hidden output directories or multi-file report assets because they weaken discoverability and portability.
- **Pitfalls**: Do not let path generation depend on platform-specific separators or ambiguous relative traversal.
- **Sources**: https://pkg.go.dev/path/filepath, https://clig.dev/#arguments-and-flags

## Testing Strategy
- **Decision**: Use golden HTML tests, structural assertions, and integration coverage for file creation plus CLI output.
- **Rationale**: Rendering failures can come from templating, formatting, or file I/O, so tests should isolate those boundaries.
- **Rejected**: String-fragment-only assertions because they miss escaping, section-order, and file-writing regressions.
- **Pitfalls**: Do not snapshot nondeterministic timestamps directly; normalize them in tests.
- **Sources**: https://pkg.go.dev/testing, https://pkg.go.dev/html/template

## Summary
| Topic | Decision | Rationale |
|-------|----------|-----------|
| Templates | `html/template` + `embed.FS` | Safe, offline, single-binary rendering |
| Charts | Semantic HTML + inline SVG | Self-contained visuals without JS libraries |
| Output | Timestamped `.html` in cwd | Matches the one-command CLI flow |
| Testing | Golden plus structural tests | Covers template, layout, and file output seams |

## Sources Index
| URL | Topic | Fetched |
|-----|-------|---------|
| https://pkg.go.dev/html/template | Templates, Testing | 2026-04-05 |
| https://pkg.go.dev/embed | Templates | 2026-04-05 |
| https://developer.mozilla.org/en-US/docs/Web/SVG/Element/svg | Charts | 2026-04-05 |
| https://developer.mozilla.org/en-US/docs/Web/HTML/Element/details | Charts | 2026-04-05 |
| https://pkg.go.dev/path/filepath | Output | 2026-04-05 |
| https://clig.dev/#arguments-and-flags | Output | 2026-04-05 |
| https://pkg.go.dev/testing | Testing | 2026-04-05 |