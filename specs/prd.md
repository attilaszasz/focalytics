# Product Requirements Document: focalytics

> Date: 2026-04-05 | Status: Draft

## Product Overview

focalytics is a privacy-first local command-line product that turns large personal photo archives into a self-contained HTML analytics report. It serves photographers who want fast, offline insight into how they shoot, what gear they actually use, and how their habits evolve over time, without importing their library into a cloud service or a heavyweight catalog system.

## Vision and Why Now

Digital photo archives now span decades, multiple cameras, and inconsistent metadata practices. Many photographers have rich historical data trapped in folder trees, but the tools that can surface it are often oriented toward editing workflows, catalog management, or connected ecosystems rather than instant archive-level insight.

focalytics exists to make archive intelligence immediate: point the tool at a directory, let it scan locally, and get a polished report that explains the collection in minutes. The timing is favorable because large personal archives are common, metadata standards remain widespread across photo workflows, and users increasingly value local-first tools that preserve privacy and ownership.

## Problem Statement

People with large photo archives struggle to answer simple, high-value questions about their own work: when they were most active, which cameras and lenses mattered most, how their focal-length preferences changed, or whether a piece of gear still earns its place. Existing approaches are often too manual, too catalog-centric, too network-dependent, or too fragile when metadata is incomplete.

Without a fast and trustworthy archive-insight product, users either ignore their historical body of work or rely on ad hoc scripts and scattered app views that do not provide a cohesive picture.

## Background and Evidence

The current product direction is grounded in founder observation and domain intuition rather than formal market validation. The strongest supporting signals are:

- Long-lived photo archives commonly accumulate across years, storage devices, and folder conventions.
- Photography workflows routinely depend on embedded EXIF/IPTC metadata and XMP sidecars, which makes metadata-driven analysis viable but inconsistent.
- Local-first, offline workflows matter for users who do not want private libraries copied to hosted services.
- Metadata quality varies across devices and tools, which creates a product need for graceful degradation and explicit transparency about excluded data.

Research-backed framing incorporated into this PRD:

- Photo-management ecosystems already rely on metadata and sidecar conventions as a practical source of organization and retrieval.
- Local-first products are valued for privacy, ownership, and longevity.
- Early validation should focus on proving the core user value, not raw download volume.

## Target Users, Stakeholders, and Core Personas

### Target Users

- Photography enthusiasts and hobbyists with years of local photos and limited patience for setup.
- Professional photographers who want a retrospective view of shooting patterns across client and personal work.
- Archivists and data-oriented users who want a concise analytical overview of large image collections.

### Stakeholders

- End users who need trustworthy, low-friction insight from local archives.
- Project maintainers responsible for product direction and distribution.
- Early contributors and testers who help validate which analytics are genuinely useful.

### Core Personas

- **Archive-curious hobbyist** — Has a large local library, wants quick answers about gear usage and shooting habits, and will abandon products that require setup or import workflows.
- **Working photographer** — Wants high-level retrospective insight across years of work and expects the product to cope with messy real-world metadata.
- **Metadata-minded archivist** — Values completeness, transparency, and the ability to analyze a deep folder tree without giving up local control.

## User Needs / Jobs To Be Done

- Help me understand what is inside my photo archive without reorganizing it first.
- Show me meaningful patterns in time, gear, and technical choices from the metadata I already have.
- Give me trustworthy results even when parts of the archive are incomplete or inconsistent.
- Let me run the analysis privately on my machine without accounts, uploads, or ongoing services.
- Produce an output I can revisit later without rerunning the product immediately.

## Product Principles or UX Principles

- **Local-first by default**: The product must preserve privacy, ownership, and offline use by keeping analysis and output on the user's machine.
- **Zero-config first value**: The first successful run should require only a target directory, not workflow setup, indexing, or account creation.
- **Insight over asset management**: The product exists to explain an archive, not to replace editing software, DAM systems, or catalog tools.
- **Honest about data quality**: Missing, conflicting, or corrupted metadata must be surfaced clearly so users can trust what the report does and does not represent.
- **One-command simplicity**: The core experience should remain understandable as a direct CLI invocation that produces a finished report.

## Scope Summary

The initial product scope is a single-user local workflow that scans a photo archive, reconciles available metadata, and generates a static HTML dashboard that provides useful archive-level insight on the first run. The first release is successful if target users can run focalytics against real, large archives and find the resulting report meaningfully useful.

### In-Scope Capabilities

- Recursive analysis of deeply nested local photo directories.
- Metadata extraction from embedded image metadata and matching XMP sidecars.
- Fallback strategies that preserve partial insight when canonical metadata is missing.
- Archive-level visual summaries for time, activity, gear usage, focal-length behavior, and exposure patterns.
- Clear reporting of excluded files or missing dimensions so charts are not misleading.
- Self-contained offline report output that does not depend on network access after generation.

### Out-of-Scope Items

- Cloud sync, hosted dashboards, user accounts, or remote processing.
- Photo editing, asset curation, keywording, or catalog-management workflows.
- Metadata write-back, file mutation, or automatic cleanup of user archives.
- Team collaboration, shared review flows, or multi-user permissions.
- General-purpose plugin ecosystems or highly customizable analytics pipelines in the initial release.

## Product Capability Map

Project-level execution anchors used by `specs/project-plan.md`. Keep these as capability clusters, not feature-level user stories.

| Capability ID | Capability | Priority | Outcome |
|---------------|------------|----------|---------|
| CAP-001 | Archive Discovery | P1 | Users can point focalytics at a large local photo directory and have the product traverse it reliably with minimal input. |
| CAP-002 | Metadata Recovery and Normalization | P1 | The product can derive usable photo facts from embedded metadata, sidecars, and defined fallbacks so analysis remains broadly useful. |
| CAP-003 | Offline Report Generation | P1 | Users receive a single self-contained HTML report that can be opened locally without accounts, services, or network dependencies. |
| CAP-004 | Timeline and Activity Insight | P1 | Users can understand when they shot, how active they were over time, and where archive density is concentrated. |
| CAP-005 | Gear and Technical Insight | P1 | Users can see meaningful patterns in camera, lens, focal-length, aperture, shutter-speed, and ISO usage. |
| CAP-006 | Data Quality Transparency | P1 | Users can trust the report because missing, conflicting, or excluded data is surfaced explicitly rather than hidden. |
| CAP-007 | Accessible Cross-Platform Adoption | P2 | Target users can discover, install, and run the product on major desktop operating systems without specialist tooling knowledge. |

## Success Metrics / KPIs / Desired Outcomes

| Metric | Target | Why It Matters | Measurement Window |
|--------|--------|----------------|--------------------|
| Successful report completion on representative large archives | At least 90% of validation runs complete and generate a usable report on archives sized for target users | Confirms the product can deliver its core outcome on real-world collections | Initial validation period |
| First-run usefulness rating | At least 80% of validation users rate the generated report as useful or very useful after a successful run | Tests the central product thesis: the output creates meaningful insight | Initial validation period |
| Time to first insight | At least 80% of validation users can install, run, and open their first report within 15 minutes without guided setup | Verifies zero-config, low-friction product value | Initial validation period |
| Rerun intent after first report | At least 60% of validation users say they would run focalytics again on the same or another archive | Indicates that the report provides durable value beyond first-run novelty | Initial validation period |

## Assumptions

- Target users already possess local photo archives large enough for archive-level analytics to feel valuable.
- Enough embedded or sidecar metadata will exist for the initial report to provide useful patterns, even when some dimensions are sparse.
- A zero-config command-line workflow is acceptable for the primary early-adopter audience.
- Static offline reporting can validate demand before the product needs deeper interactivity or collaboration features.

## Constraints

- The product must operate locally and preserve a privacy-first posture.
- The generated report must be self-contained and usable offline.
- The initial user experience should remain centered on a single obvious command against a target directory.
- Product scope must stay focused on archive insight rather than becoming a general photo-management suite.
- The initial release should remain broadly usable across major desktop operating systems.

## Dependencies

- Availability and quality of embedded metadata and sidecar files in user archives.
- Reliable filesystem access to large nested directories on user machines.
- Release distribution channels that make local CLI installation approachable for the target audience.
- Early testers willing to validate usefulness on real archives and share feedback.

## Risks

- Real-world metadata inconsistency may reduce perceived report accuracy if exclusions are not communicated clearly.
- Large archives may expose performance or reliability issues that undermine trust during first use.
- The value proposition may be diluted if messaging tries to serve hobbyists, professionals, and archivists equally at launch.
- A static report may not satisfy users who expect interactive exploration or archive remediation workflows.

## Open Questions

- Which report sections most strongly drive usefulness for hobbyist users after the first run?
- What archive size and file-type diversity should define the representative validation set?
- When, if ever, should optional flags be introduced without weakening the one-command product promise?
- How much demand exists for sharing or publishing generated reports versus keeping them strictly personal?

## Release or Validation Approach

Validation should focus on real archive runs, not broad feature breadth. The first release should be evaluated with representative local libraries from target users, beginning with photography enthusiasts and hobbyists. Success is defined by whether focalytics can be installed with low friction, complete a scan on real archives, and produce a report that users consider meaningfully insightful.

Public open-source distribution can support discovery and feedback collection, but downloads, stars, and package-manager reach are secondary signals. The primary release-learning loop should emphasize completion success, perceived usefulness, and rerun intent.

## Domain Glossary / Terminology

- **Archive insight**: High-level understanding of a photo collection's patterns, trends, and composition.
- **Embedded metadata**: Information stored inside the image file itself, such as EXIF or IPTC fields.
- **XMP sidecar**: A separate metadata file associated with an image and commonly used in photo workflows.
- **35mm equivalent focal length**: A normalized focal-length representation used to compare framing across different sensor sizes.
- **Self-contained report**: A generated HTML file that includes the assets it needs and does not depend on network-hosted resources.

## Handoff Guidance

Context that downstream architecture design or governance work must preserve.

- **Product intent to preserve**: focalytics must feel like a direct, privacy-first archive-insight tool rather than a general photo-management platform.
- **Scope boundaries to respect**: Do not expand the initial product into editing, cloud sync, collaboration, or metadata write-back workflows.
- **Critical constraints**: Preserve the local-first posture, offline report output, low-friction first run, and transparency around imperfect metadata.
- **Open decisions needing technical input**: Representative performance targets for large archives, practical definition of supported file sets, and how to preserve simplicity if optional controls become necessary later.

## Project Context Baseline Updates

- focalytics is a project-level product initiative aimed first at photography enthusiasts and hobbyists with large local archives.
- The canonical product promise is privacy-first, offline archive insight from a single local CLI invocation.
- Initial release validation should prioritize successful report generation on real archives and perceived usefulness of the generated report.