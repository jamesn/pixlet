# Project: Replace `nfnt/resize`

Status: Draft
Owner: @jamesn
Last updated: 2026-09-24
Supersedes: `docs/specs/009-image-resize-replacement.md`

## 1. Summary

Remove the archived `github.com/nfnt/resize` dependency from Pixlet's render
path without changing how any image looks on the signs. Characterization tests
showed that the obvious replacement (`golang.org/x/image/draw`) produces
**different pixels**, so the plan is to bring the exact algorithm in-house,
guarded by a golden-image corpus.

## 2. Background

`render.Image` scales images when an app sets `width` and/or `height`
(`render/image.go`, `Image.Init`):

```go
for i := 0; i < len(p.imgs); i++ {
    p.imgs[i] = resize.Resize(uint(nw), uint(nh), p.imgs[i], resize.NearestNeighbor)
}
```

- Inputs: WebP, GIF (animated: every frame), SVG (rasterized), and any
  `image.Decode` format (PNG, JPEG).
- Pixlet computes the missing dimension itself to keep the aspect ratio
  (`int()` truncation).
- Library: `nfnt/resize` @ `83c6a9932646` (2018), archived. About 2,700 lines
  including tests. License: ISC (permissive; copying allowed with the notice kept).

### Verified behavior of the current code (2026-09-24 probe)

| Case | `nfnt/resize` NearestNeighbor | `x/image/draw` NearestNeighbor |
|------|-------------------------------|--------------------------------|
| Downscale 4×1 `[0,255,0,255]` → 2×1 | `[127,127]`: **averages** (the kernel widens on downscale, so it acts like a box filter) | `[255,255]`: picks one source pixel |
| `*image.Paletted` input (GIF frames) | returns `*image.RGBA64` | returns whatever destination type you allocate |
| A requested dimension of `0` | treated as "preserve aspect ratio" | produces an empty image |

**Conclusion:** swapping in `x/image/draw` would visibly change any downscaled
image (sharper and more aliased) and would change edge-case behavior. That
breaks contract C1 (pixel-identical renders).

## 3. Goals / Non-goals

**Goals**
- G1: No dependency on an archived module in the render path.
- G2: Pixel-identical output for every input `render.Image` accepts (C1).
- G3: No performance regression in `BenchmarkRunAndRender`.

**Owner input (2026-09-24):** none of the owner's apps downscale images, but
geometric compatibility must be maintained for all apps. Geometry (output
dimensions, aspect-ratio handling, zero-dimension semantics, frame count and
per-frame alignment) is therefore a **hard requirement** (R0). Pixel-identical
output (G2) remains the target, and Option D delivers both.

**Non-goals**
- New scaling modes or quality improvements (could be a later additive
  `render.Image(scale_mode=...)` option; see §8).
- Changing how Pixlet computes the missing dimension.

## 4. Options considered

| Option | Description | Pixel-identical? | Maintenance | Verdict |
|--------|-------------|------------------|-------------|---------|
| A | Swap to `x/image/draw.NearestNeighbor` | **No** (verified) | Lowest | Rejected (breaks C1) |
| B | Swap to `x/image/draw` with a custom kernel mimicking nfnt | Possibly, but hard to prove across all image types | Medium | Fallback only |
| C | Fork `nfnt/resize` to `jamesn/resize` | Yes | A second external repo (after starlib) for about 2,700 frozen lines | Not preferred |
| **D** | **Copy the needed nfnt code into `render/internal/resize`** (ISC notice preserved), trimmed to what Pixlet calls | **Yes** (same code) | Small, lives with its only consumer | **Recommended** |

Why D over C: unlike starlib, this code has no third-party dependencies to keep
updated and one caller. A separate repo adds release overhead with no benefit.

## 5. Requirements

- **R0 (geometry, hard gate)** For every input and requested `width`/`height`, the output MUST have exactly the same dimensions and bounds origin as today, including: the aspect-ratio dimension Pixlet computes (`int()` truncation), nfnt's "0 = preserve aspect ratio" rule when a computed dimension is 0, and identical frame count for animated sources. This holds even if a future opt-in scaling mode (Phase 3) changes pixel values.
- **R1** Output MUST be pixel-identical to `nfnt/resize` for every case in the golden corpus (§6, Phase 0).
- **R2** The returned image type for each input type MUST match nfnt's (e.g. Paletted → RGBA64), because downstream drawing and encoding may depend on it.
- **R3** Zero-dimension semantics MUST match (0 = preserve aspect ratio).
- **R4** The internal package MUST carry nfnt's ISC license text and attribution.
- **R5** Only code reachable from `Resize(..., NearestNeighbor)` SHOULD be kept; other filters and `Thumbnail` MAY be dropped.
- **R6** `BenchmarkRunAndRender` MUST NOT regress by more than 5%.
- **R7** After the switch, `go.mod` MUST NOT reference `github.com/nfnt/resize`.

## 6. Plan

### Phase 0: Golden corpus and characterization tests

Build the safety net first, against the current library.

1. Create `render/testdata/resize/` with source images covering:
   - **Types:** `RGBA`, `NRGBA`, `RGBA64`, `NRGBA64`, `Gray`, `Gray16`, `Paletted` (GIF),
     `YCbCr` 4:4:4 / 4:2:2 / 4:2:0 (JPEG), and any type SVG rasterization or WebP decoding produces.
   - **Content:** hard edges, 1-px checkerboard (shows averaging vs picking), gradients, alpha edges, fully transparent regions.
   - **Sizes:** tiny (1×1, 2×2), odd (7×5), the display size (64×32), large (256×256).
2. Scale matrix per source: upscale ×2, ×3, non-integer up (×1.5), downscale ÷2, ÷3,
   non-integer down (×0.7), width-only, height-only, one dimension computing to 0 or 1.
3. Generate expected outputs **with the current `nfnt/resize`**; store SHA-256 of
   raw pixel data plus the concrete Go type in `golden.json`.
4. Test `TestResizeGolden` asserts every case against `golden.json`.
   Record output bounds separately from pixel hashes, so a geometry failure (R0) is reported distinctly from a pixel difference (R1).
5. Add an end-to-end check: render every `examples/*.star` that uses
   `render.Image(width=…/height=…)`, plus a dedicated test app that scales a GIF, and
   store output hashes.

**Exit:** golden tests merged and green on current `main`.

### Phase 1: Bring the algorithm in-house

1. Create `render/internal/resize/` by copying `resize.go`, `filters.go`,
   `nearest.go`, `converter.go`, `ycc.go` (and their tests) from the pinned commit.
2. Trim to the NearestNeighbor path (R5), keeping per-type fast paths, because they determine output types (R2).
3. Add the ISC `LICENSE` and a header noting the origin and commit (R4).
4. Switch `render/image.go` to the internal package.
5. Run the golden tests, the upstream nfnt tests that were kept, and the example hashes.

**Exit:** all golden cases pass with the internal package.

### Phase 2: Remove the dependency

1. `go mod tidy`; confirm `nfnt/resize` is gone (R7).
2. Run `make bench` before and after (R6).
3. govulncheck.

**Exit:** PR merged.

### Phase 3 (optional, later): Additive scaling modes

Only if wanted: add an optional `render.Image(..., scale="nearest"|"smooth"|...)`
parameter whose **default stays the current behavior**. This could use
`x/image/draw` for new modes. It's a separate spec because it changes the
widget API (additive).

## 7. Testing strategy

| Test | Proves |
|------|--------|
| `TestResizeGolden` (Phase 0) | Pixel and type equality for the full matrix (R1–R3). |
| Kept nfnt unit tests | Internal code still behaves like upstream. |
| Example render hashes | End-to-end C1. |
| `make bench` | R6. |

## 8. Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Corpus misses an input type seen in real apps | Medium | Include every type the WebP/GIF/SVG/image decoders return; add a debug assertion logging unknown types during Phase 0 test runs. |
| nfnt uses goroutines, so output might depend on scheduling | Low | nfnt splits work by rows deterministically; run golden tests with `-count=20` and `GOMAXPROCS=1` and `=8`. |
| Owner later *wants* crisper downscaling | n/a | Phase 3 adds it as opt-in without changing existing apps. |

## 9. Deliverables checklist

- [ ] Phase 0: golden corpus + `TestResizeGolden` + example hashes on `main`
- [ ] Phase 1: `render/internal/resize` with ISC notice; `render/image.go` switched
- [ ] Phase 2: `nfnt/resize` removed from `go.mod`; benchmark within 5%
- [ ] Phase 3: (optional) spec for additive scaling modes

## 10. Open questions

1. ~~Do any apps downscale?~~ Answered: no, but geometric compatibility is required (R0).
2. Is there interest in the Phase 3 "crisp" mode for pixel-art sources, where
   averaging blurs hard edges?
