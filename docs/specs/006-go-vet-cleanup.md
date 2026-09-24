# 006: Fix `go vet` Findings and Enforce Vet in CI

Status: Draft

## Summary

Fix the existing `go vet` findings (surfaced by the Go 1.26 toolchain), then add
`go vet ./...` to CI so new ones are caught.

## Motivation

Findings on `main` as of #12:

| Location | Finding |
|----------|---------|
| `server/loader/loader.go:180` | `context.WithTimeoutCause` cancel func discarded (context leak) |
| `cmd/render.go:139` | same |
| `render/vector.go:20-21` | malformed struct tags `starlark: "main_align"` (space after colon) |
| `render/animation/curve.go:56` | unreachable code |
| `runtime/modules/animation_runtime/percentage.go:14` | unkeyed fields in struct literal |

The two context leaks are real bugs: timers stay alive until they fire. The
struct tags are malformed, so `reflect.StructTag.Get("starlark")` returns an empty
string for those fields.

## Goals / Non-goals

Goals:
- `go vet ./...` is clean and stays clean.

Non-goals:
- Adopting a broader linter (golangci-lint) in this spec.

## Requirements

- **R1** Each finding above MUST be fixed.
- **R2** Fixing the `render/vector.go` tags MUST NOT change how `main_align` / `cross_align` behave from Starlark (C2). Verify how the tags are consumed first.
- **R3** CI MUST run `go vet ./...` and fail on findings.
- **R4** Render output MUST stay pixel-identical (C1).

## Design

- Context leaks: keep the cancel func and `defer cancel()` in the owning scope.
- Struct tags: correct to `starlark:"main_align"`. Before changing, grep for how
  `starlark` tags are read. If Starlark attribute lookup currently works through
  another path (e.g. field-name mapping in `runtime/gen`), confirm the fix does
  not introduce duplicate or renamed attributes.
- Unreachable code / unkeyed fields: mechanical fixes.

## Compatibility

- C1/C2 at low risk (`vector.go` tags). Guarded by R2 and the existing widget tests.

## Acceptance criteria

- [ ] `go vet ./...` reports nothing.
- [ ] `go test ./...` passes.
- [ ] `examples/` render byte-identical WebP output before and after (script compares hashes).
- [ ] CI fails on a deliberately introduced vet finding.

## Tasks

1. Capture golden hashes of `pixlet render` for every example.
2. Fix each finding.
3. Re-render and compare hashes.
4. Add a `go vet ./...` step to `pull-request.yml` (lint job).

## Open questions

- Should `staticcheck` be added at the same time? Proposed: separate spec if wanted.
