# 002: Consolidate Dependency Automation on Dependabot

Status: In progress (Phase 1 PR)

## Summary

Use one tool to propose dependency updates: Dependabot. Remove the leftover
`renovate.json` and add a `.github/dependabot.yml` that groups updates for Go
modules, npm, and GitHub Actions.

## Motivation

- Dependabot already opens security PRs on this repo (e.g. #5, #6).
- `renovate.json` came with the upstream fork. It sets `automerge: true` and
  extends `config:base`, which Renovate has deprecated. If Renovate were ever
  enabled, the automerge setting could merge updates that change app behavior
  (see the `go-humanize` note in #11).
- Two bots proposing the same updates create duplicate PRs.

## Goals / Non-goals

Goals:
- Receive grouped, reviewable update PRs on a predictable schedule.
- Keep security updates flowing without delay.

Non-goals:
- Auto-merging. Every update is reviewed, because behavior-affecting upgrades
  (C1/C2) are possible.

## Requirements

- **R1** `renovate.json` MUST be removed.
- **R2** `.github/dependabot.yml` MUST cover the `gomod`, `npm`, and `github-actions` ecosystems.
- **R3** Version updates SHOULD run weekly and be grouped per ecosystem (minor + patch together), limiting PR noise to about 3 per week.
- **R4** Major version updates SHOULD be separate PRs so they can be reviewed individually.
- **R5** `github.com/dustin/go-humanize` MUST be ignored for versions ≥ 1.1.0 until a compatibility decision is made (C2).
- **R6** Security updates MUST remain enabled (repo setting).

## Design

```yaml
version: 2
updates:
  - package-ecosystem: gomod
    directory: /
    schedule: { interval: weekly }
    groups:
      go-minor-patch: { update-types: [minor, patch] }
    ignore:
      - dependency-name: github.com/dustin/go-humanize
        versions: [">= 1.1.0"]   # relative_time output change; see spec README C2
  - package-ecosystem: npm
    directory: /
    schedule: { interval: weekly }
    groups:
      npm-minor-patch: { update-types: [minor, patch] }
  - package-ecosystem: github-actions
    directory: /
    schedule: { interval: weekly }
    groups:
      actions: { patterns: ["*"] }
```

## Compatibility

Indirect. Update PRs may threaten C1/C2, which is why nothing auto-merges and
spec 001 plus the existing tests gate every update.

## Acceptance criteria

- [ ] `renovate.json` is deleted.
- [ ] Dependabot opens grouped PRs for each ecosystem after the first weekly run.
- [ ] No Dependabot PR bumps `go-humanize` past 1.0.1.

## Tasks

1. Delete `renovate.json`.
2. Add `.github/dependabot.yml`.
3. Confirm in GitHub settings that Dependabot security updates are enabled.

## Open questions

- Is Renovate installed as a GitHub App on this repo? If so, uninstall it too.
