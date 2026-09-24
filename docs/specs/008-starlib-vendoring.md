# 008: Vendor the Archived `qri-io/starlib` Modules

Status: Draft

## Summary

Copy the `qri-io/starlib` packages that Pixlet uses into this repository (with
license attribution), so the fork controls and can patch the Starlark modules
that apps depend on.

## Motivation

`github.com/qri-io/starlib` is archived and pinned to a 2022 pseudo-version.
Pixlet exposes these modules to every app:

| Starlark load path | starlib package |
|--------------------|-----------------|
| `bsoup.star` | `bsoup` |
| `compress/gzip.star` | `compress/gzip` |
| `encoding/base64.star` | `encoding/base64` |
| `encoding/csv.star` | `encoding/csv` |
| `hash.star` | `hash` |
| `html.star` | `html` |
| `re.star` | `re` |
| `zipfile.star` | `zipfile` |
| (helper) | `util`, used by `runtime/modules/starlarkhttp` |

Because it's archived, any bug or security issue in these modules can never be
fixed upstream. It also pulls in extra transitive dependencies
(`dustmop/soup`, `PuerkitoBio/goquery`).

## Goals / Non-goals

Goals:
- Keep identical Starlark behavior with code the fork owns.
- Reduce transitive dependencies where possible.

Non-goals:
- Adding new functions to these modules (can follow later).

## Requirements

- **R1** Each package listed above MUST be copied under `runtime/modules/starlib/<pkg>` with its original license (MIT) and a NOTICE.
- **R2** Load paths, function names, signatures, and return values MUST be unchanged (C2).
- **R3** Upstream starlib tests for the vendored packages SHOULD be copied and pass.
- **R4** `go.mod` MUST no longer require `github.com/qri-io/starlib` after this change.
- **R5** Only packages actually imported MAY be vendored. Unused starlib packages MUST NOT be copied.

## Design

- Copy the source at the currently pinned commit (`7fb7ff9ec804`), then change
  only the import paths.
- Update the imports in `runtime/applet.go` and
  `runtime/modules/starlarkhttp/starlarkhttp.go`.
- Run the copied tests plus the existing runtime tests.

## Compatibility

- **C2 at risk**, guarded by copying the code as-is (no logic changes) and running both test suites.

## Acceptance criteria

- [ ] `go mod why github.com/qri-io/starlib` reports it is not needed.
- [ ] All copied upstream tests pass.
- [ ] All existing runtime tests pass.
- [ ] License and attribution are present.

## Tasks

1. Copy the packages and their tests at the pinned commit.
2. Rewrite import paths.
3. `go mod tidy`; run tests and govulncheck.

## Open questions

- Alternatively, switch to the actively maintained `github.com/1set/starlet`
  modules? That would change behavior (C2 risk), so vendoring is proposed first.
