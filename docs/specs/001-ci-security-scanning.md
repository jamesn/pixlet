# 001: Vulnerability Scanning in CI

Status: In progress (Phase 1 PR)

## Summary

Run `govulncheck` and `npm audit` on every pull request and on a weekly schedule,
so new security advisories surface automatically instead of through a manual audit.

## Motivation

The fork accumulated 54 reachable Go vulnerabilities and 33 npm advisories before
they were noticed (fixed in #11 and #12). CI currently builds and tests but never
checks for known vulnerabilities, so new ones will pile up the same way.

## Goals / Non-goals

Goals:
- Flag known-vulnerable Go code paths and npm packages on every PR.
- Catch advisories published after merge through a scheduled run.

Non-goals:
- Automatically fixing vulnerabilities (spec 002 covers update PRs).
- Replacing CodeQL (it stays as-is).

## Requirements

- **R1** CI MUST run `govulncheck ./...` on pull requests to `main`.
- **R2** CI MUST run `npm audit --audit-level=high` on pull requests to `main`.
- **R3** A scheduled workflow MUST run both checks at least weekly against `main`.
- **R4** A `govulncheck` finding in reachable code MUST fail the job.
- **R5** `npm audit` MUST fail on high/critical findings. Moderate/low findings SHOULD be reported without failing, so a moderate advisory with no available fix does not block unrelated PRs.
- **R6** The job MUST install `libwebp` headers (cgo), matching the existing build jobs.
- **R7** The govulncheck version SHOULD be pinned for reproducible results.

## Design

One workflow, `.github/workflows/security.yml`, triggered by `pull_request`,
a weekly `schedule` (Mondays 13:00 UTC), and `workflow_dispatch`. Keeping PR
and scheduled scans in one file avoids duplicating the steps.

- **govulncheck** job: `setup-go` from `go.mod`, `libwebp-dev` from apt (cgo
  headers), then `go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...` (R1, R4, R6, R7).
- **npm-audit** job: `setup-node` 22, `npm ci`, `npm audit --audit-level=high` (R2, R5).
- `permissions: contents: read`.

Verified locally before merge: both pass on the Phase 1 branch, and
govulncheck exits non-zero (3) after downgrading to known-vulnerable versions.

## Compatibility

No runtime impact. None of C1–C6 are affected.

## Acceptance criteria

- [ ] PRs show a `security` check that runs govulncheck and npm audit.
- [ ] Introducing a known-vulnerable module version in a test branch fails the check.
- [ ] The scheduled workflow runs, and its results are visible in the Actions tab.
- [ ] Current `main` passes (verified 0 findings after #11 and #12).

## Tasks

1. Add the `security` job to `pull-request.yml`.
2. Add `security-scheduled.yml`.
3. Verify on a throwaway branch that a vulnerable dependency fails the job.

## Open questions

- Should a scheduled-run failure open an issue automatically, or is the Actions
  failure email enough?
