# 009: Replace Archived `nfnt/resize`

Status: Superseded by [`docs/projects/image-resize/PLAN.md`](../projects/image-resize/PLAN.md)

## Summary

Replace `github.com/nfnt/resize` (archived since 2018) with
`golang.org/x/image/draw`, only if output is proven pixel-identical.

## Motivation

`render/image.go:192` uses `resize.Resize(w, h, img, resize.NearestNeighbor)`
to scale images in the `render.Image` widget. The library is archived. The code
is small and stable, so the risk of leaving it is low, but it is an unmaintained
dependency in the rendering path.

## Goals / Non-goals

Goals:
- Remove an archived dependency from the render path.

Non-goals:
- Changing scaling quality or adding interpolation modes.

## Requirements

- **R1** Output MUST be pixel-identical to the current implementation for all tested inputs (C1).
- **R2** If R1 cannot be met, this spec MUST be abandoned (mark Status: Rejected) rather than accepting visual changes.
- **R3** A golden-image test MUST cover upscale, downscale, non-integer ratios, odd dimensions, and animated GIF frames.

## Design

1. Write a golden test first. Generate scaled outputs with `nfnt/resize` for a
   matrix of source images and target sizes, and store checksums.
2. Implement with `draw.NearestNeighbor.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Src, nil)`.
3. Compare. Nearest-neighbor sampling-point rounding may differ between the
   libraries; if it does, R2 applies.

## Compatibility

- **C1 at risk**, strictly guarded by R1/R2.

## Acceptance criteria

- [ ] Golden test exists and passes on the current implementation.
- [ ] Either the replacement passes the golden test and `nfnt/resize` is removed from `go.mod`, or the spec is marked Rejected with the diff evidence.

## Tasks

1. Golden test PR.
2. Replacement PR (or rejection note).

## Open questions

- None.
