# 003: Refresh GitHub Actions Versions

Status: In progress (Phase 1 PR)

## Summary

Update the GitHub Actions used by the workflows to current major versions, and
pin third-party actions to commit SHAs.

## Motivation

The workflows use older action majors:

| Action | Current in repo |
|--------|-----------------|
| `actions/checkout` | v4 |
| `actions/setup-go` | v5 |
| `actions/setup-node` | v4 |
| `actions/cache` | v4 |
| `actions/upload-artifact` / `download-artifact` | v4 / v5 |
| `github/codeql-action/*` | v3 |
| `msys2/setup-msys2` | v2 |
| `softprops/action-gh-release` | **v1** |

Older majors run on deprecated Node runtimes and eventually stop working. The
release action `softprops/action-gh-release@v1` is the oldest, and it publishes
the release binaries.

## Goals / Non-goals

Goals:
- Current majors for all actions.
- Third-party actions (`softprops/*`, `msys2/*`) pinned by SHA for supply-chain safety.

Non-goals:
- Restructuring the workflows.

## Requirements

- **R1** Every `uses:` MUST reference a currently supported major version.
- **R2** Third-party (non-`actions/*`, non-`github/*`) actions MUST be pinned to a full commit SHA, with the version in a trailing comment.
- **R3** The redundant `actions/cache` steps for Go modules SHOULD be removed, because `setup-go` already caches via `cache-dependency-path`.
- **R4** A tag push MUST still produce release artifacts for Linux, macOS, and Windows.

## Design

Update each `uses:` line, checking each action's changelog for breaking input
changes (notably `upload-artifact`/`download-artifact`, whose artifact names must
be unique per job, which the matrix already does).

Spec 002's `github-actions` ecosystem keeps these current afterward.

## Compatibility

CI only. R4 protects the release process.

## Acceptance criteria

- [ ] All workflows pass on a PR.
- [ ] A test tag on a fork or throwaway repo produces a draft release with all three platform binaries.
- [ ] No Node-runtime deprecation warnings in workflow logs.

## Tasks

1. Update action versions in `main.yml`, `pull-request.yml`, and `codeql-analysis.yml`.
2. Pin `softprops/action-gh-release` and `msys2/setup-msys2` by SHA.
3. Remove the duplicate Go module cache steps.
4. Dry-run the release path.

## Open questions

- None.
