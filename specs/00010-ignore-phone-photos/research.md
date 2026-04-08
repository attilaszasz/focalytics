# Research: Ignore Phone Photos
> Feature | 2026-04-08 | Purpose: inform runtime filter design, classification rules, and disclosure contracts

## CLI Filter Contract
- **Decision**: Add one explicit boolean opt-in flag that defaults to disabled and leaves the current stdout success contract unchanged.
- **Rationale**: Predictable CLI behavior requires additive options without changing existing output semantics.
- **Rejected**: Silent default exclusion and overloaded multi-mode flags because they would weaken scriptability and backward compatibility.
- **Pitfalls**: Do not append filter summaries to stdout; keep them on stderr or the interactive UI channel.
- **Sources**: https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html, https://semver.org/

## Device Classification Evidence
- **Decision**: Classify photos as phone-made only from trusted device identity fields such as make or model recovered from embedded, sidecar, or platform metadata.
- **Rationale**: Device identity tags are the strongest available evidence while software and provenance tags often describe edits rather than capture origin.
- **Rejected**: Software-tag inference, focal-length heuristics, and provenance-only hints because they over-classify ambiguous files.
- **Pitfalls**: Treat missing or conflicting metadata as unknown and keep those files included.
- **Sources**: https://exiftool.org/TagNames/EXIF.html, https://developer.adobe.com/xmp/docs/xmp-namespaces/xmp-mm/

## Filtered Scope Disclosure
- **Decision**: Show an always-visible scope note near the report overview, per-section notes for affected analytics, and terminal completion feedback outside stdout.
- **Rationale**: Users need scope changes in text and structure wherever interpretation changes, not hidden behind optional disclosure.
- **Rejected**: Details-only disclosure and color-only badges because they bury or weaken essential context.
- **Pitfalls**: Keep timeline and total-photo figures explicitly whole-archive in this increment so mixed-scope reporting stays understandable.
- **Sources**: https://www.w3.org/WAI/WCAG22/Understanding/info-and-relationships.html, https://design-system.service.gov.uk/components/details/

## Summary
| Topic | Decision | Rationale |
|-------|----------|-----------|
| CLI Filter Contract | One boolean opt-in flag, stdout unchanged | Preserves backward-compatible scriptability |
| Device Classification Evidence | Trusted make/model metadata only | Avoids speculative exclusion |
| Filtered Scope Disclosure | Always-visible overview note plus affected-section notes | Keeps filtered analytics interpretable |

## Sources Index
| URL | Topic | Fetched |
|-----|-------|---------|
| https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html | CLI Filter Contract | 2026-04-08 |
| https://semver.org/ | CLI Filter Contract | 2026-04-08 |
| https://exiftool.org/TagNames/EXIF.html | Device Classification Evidence | 2026-04-08 |
| https://developer.adobe.com/xmp/docs/xmp-namespaces/xmp-mm/ | Device Classification Evidence | 2026-04-08 |
| https://www.w3.org/WAI/WCAG22/Understanding/info-and-relationships.html | Filtered Scope Disclosure | 2026-04-08 |
| https://design-system.service.gov.uk/components/details/ | Filtered Scope Disclosure | 2026-04-08 |
