# Research: Aggregate Insight Metrics
> Feature | 2026-04-05 | Purpose: inform deterministic archive aggregation design for E004

## Aggregate-Only Summary Modeling
- **Decision**: Use one-way typed accumulators that emit archive summary sections without carrying per-photo rendering payloads.
- **Rationale**: This keeps memory bound to aggregate cardinality and matches the report-layer contract in the PRD and SAD.
- **Rejected**: Passing raw fact lists into rendering because it would couple E004 to E005 and scale with file count instead of summary size.
- **Pitfalls**: Do not let map iteration define public ordering or summary shape.
- **Sources**: https://exiftool.org/TagNames/EXIF.html, https://go.dev/ref/spec

## Deterministic Bucket Strategy
- **Decision**: Bucket on canonical UTC date keys and numeric metric keys, then sort explicitly before publishing aggregate results.
- **Rationale**: Stable keys prevent formatting drift and repeated-run instability.
- **Rejected**: Float-string buckets and locale-derived labels because they make ordering and equality brittle.
- **Pitfalls**: Do not mix presentation labels with primary grouping keys.
- **Sources**: https://www.rfc-editor.org/rfc/rfc3339, https://exiftool.org/TagNames/EXIF.html

## Exclusion Transparency
- **Decision**: Preserve exclusions and warnings as first-class aggregate summaries with per-metric and per-reason counts.
- **Rationale**: Data-quality transparency must survive aggregation so later charts can explain omissions honestly.
- **Rejected**: Reconstructing exclusions from logs because it would lose denominator clarity and stable grouping.
- **Pitfalls**: Do not collapse unrelated reasons into generic unknown totals.
- **Sources**: https://pkg.go.dev/testing, https://go.dev/doc/security/fuzz/

## Testing Strategy
- **Decision**: Use table-driven unit tests, deterministic ordering checks, stage integration tests, and regression-friendly malformed-input coverage.
- **Rationale**: Aggregation is pure logic with clear boundary conditions, which fits Go's native testing model well.
- **Rejected**: Snapshotting unordered map output because it creates flaky tests and hides ordering bugs.
- **Pitfalls**: Do not depend on local timezone, wall clock, or shared mutable fixtures.
- **Sources**: https://pkg.go.dev/testing, https://go.dev/doc/security/fuzz/

## Summary
| Topic | Decision | Rationale |
|-------|----------|-----------|
| Aggregate Model | Typed accumulators | Preserves aggregate-only contract and bounded memory |
| Bucket Keys | Canonical numeric/time keys | Ensures deterministic counts and sort order |
| Exclusions | First-class summaries | Keeps omission reasons visible to later reports |
| Testing | Table-driven plus ordering checks | Validates bucket boundaries and reproducibility |

## Sources Index
| URL | Topic | Fetched |
|-----|-------|---------|
| https://exiftool.org/TagNames/EXIF.html | Aggregate Model | 2026-04-05 |
| https://go.dev/ref/spec | Aggregate Model | 2026-04-05 |
| https://www.rfc-editor.org/rfc/rfc3339 | Bucket Keys | 2026-04-05 |
| https://pkg.go.dev/testing | Exclusions, Testing | 2026-04-05 |
| https://go.dev/doc/security/fuzz/ | Exclusions, Testing | 2026-04-05 |