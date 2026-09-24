# Pixlet Fork: Maintenance & Roadmap Plan

Status: Draft
Last updated: 2026-09-24

This fork of Pixlet programs the sign boards in the house. Upstream (tidbyt/pixlet)
is no longer maintained, so this plan describes how the fork stays secure,
buildable, and useful. Each work item is a spec in this directory. Implementation
happens in small PRs, one spec at a time.

## Guiding principle: compatibility first

Every spec in this directory MUST preserve the following unless the spec says
otherwise and the owner has explicitly agreed:

| ID | Compatibility contract |
|----|------------------------|
| C1 | Existing `.star` apps render **pixel-identical** WebP/GIF output. |
| C2 | Starlark module names, function signatures, and return values are unchanged (e.g. `humanize.relative_time` keeps its trailing space). |
| C3 | Existing CLI commands, positional args, flags, and env vars (`TIDBYT_API_TOKEN`) keep working. |
| C4 | Existing config files and stored OAuth tokens remain valid. |
| C5 | Secrets encrypted by earlier versions still decrypt. |
| C6 | `pixlet push` to the Tidbyt API remains the default behavior. |

Specs that could affect a contract must include a verification step that proves it
holds (golden images, round-trip tests, etc.).

## Spec format

Each spec uses the same sections:

- **Status**: Draft → Approved → In progress → Done
- **Summary**: one paragraph
- **Motivation**: why it matters
- **Goals / Non-goals**
- **Requirements**: numbered, testable, using MUST / SHOULD / MAY
- **Design**: how we intend to meet the requirements
- **Compatibility**: which contracts (C1–C6) are at risk and how they are protected
- **Acceptance criteria**: checklist that defines "done"
- **Tasks**: implementation steps
- **Open questions**

## Work items

### Phase 1: Low-risk maintenance (one PR)

These items do not change runtime behavior, so they can ship together.

| Spec | Title | Risk |
|------|-------|------|
| [001](001-ci-security-scanning.md) | Vulnerability scanning in CI | None (CI only) |
| [002](002-dependency-automation.md) | Consolidate dependency automation on Dependabot | None (repo config) |
| [003](003-github-actions-refresh.md) | Refresh GitHub Actions versions | Low (CI only) |
| [004](004-frontend-build-cleanup.md) | Remove dead frontend build dependencies | Low (build only) |
| [005](005-hermetic-tests.md) | Make network-dependent tests hermetic | None (tests only) |

### Phase 2: Code health

| Spec | Title | Risk |
|------|-------|------|
| [006](006-go-vet-cleanup.md) | Fix `go vet` findings and enforce vet in CI | Low |

### Phase 3: Replace deprecated / archived libraries

Each of these could affect a compatibility contract, so each gets its own PR with
verification.

| Spec | Title | Contract at risk |
|------|-------|------------------|
| [007](007-tink-migration.md) | Migrate `google/tink/go` → `tink-crypto/tink-go/v2` | C5 |
| [008](008-starlib-vendoring.md) | ~~Vendor `qri-io/starlib`~~. Now its own project: [starlib fork](../projects/starlib-fork/PLAN.md) | C2 |
| [009](009-image-resize-replacement.md) | ~~Replace `nfnt/resize`~~. Now its own project: [image resize](../projects/image-resize/PLAN.md) | C1 |
| [010](010-frontend-icon-bundle.md) | Reduce Font Awesome bundle size | Schema icons in apps |

### Phase 4: Device programming independence

| Spec | Title | Contract at risk |
|------|-------|------------------|
| [011](011-device-push-targets.md) | Pluggable push targets (Tidbyt API + self-hosted) | C3, C6 |

The signs are currently programmed through the Tidbyt cloud API
(`api.tidbyt.com`). Spec 011 keeps that path as the default and adds an
abstraction so a self-hosted target can be added without breaking existing
commands.

## Projects

Larger efforts that span multiple PRs (and possibly repos) are planned as projects:

| Project | Summary |
|---------|---------|
| [Fork and maintain starlib](../projects/starlib-fork/PLAN.md) | Fork `qri-io/starlib` to `jamesn/starlib`, trim it to the 8 modules Pixlet uses, and maintain it with CI and conformance tests. |
| [Replace nfnt/resize](../projects/image-resize/PLAN.md) | Bring the exact resize algorithm in-house behind a golden-image corpus. `x/image/draw` was verified to change downscaled output. |

## Sequencing

```
Phase 1 (001–005) ──► Phase 2 (006) ──► Phase 3 (007, 008, 009, 010 in any order)
                                    └─► Phase 4 (011)  ← can start after Phase 1
```

Phase 1 comes first because 001 (CI scanning) and 005 (hermetic tests) make every
later PR safer to review. Phase 4 does not depend on Phase 3, so it can run in
parallel if priorities shift.

## Completed work

| PR | Summary |
|----|---------|
| #11 | Go 1.26 toolchain; all Go modules updated (govulncheck 54 → 0); npm audit fixes (33 → 2); CI moved to `go-version-file` and Node 22. |
| #12 | React 18 + react-router 7.18 (npm audit 2 → 0). |

## Superseded documents

The earlier planning notes are superseded by this directory:

- `docs/dependency-upgrades-todo.md`: completed by #11 and #12.
- `docs/build-test-todo.md`: folded into specs 001 and 005.
- `docs/code-splitting-todo.md` and `docs/fontawesome-todo.md`: folded into spec 010.
