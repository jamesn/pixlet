# Project: Fork and Maintain `starlib`

Status: Draft
Owner: @jamesn
Last updated: 2026-09-24
Supersedes: `docs/specs/008-starlib-vendoring.md`

## 1. Summary

Create and maintain a fork of the archived `github.com/qri-io/starlib` as
`github.com/jamesn/starlib`, trimmed to the Starlark modules Pixlet exposes to
apps. Pixlet then depends on the fork, which can receive dependency updates,
security fixes, and (additive) improvements that the archived upstream never
will.

## 2. Background

Pixlet loads these starlib packages into every app (`runtime/applet.go`):

| Starlark `load()` path | starlib package | Third-party deps |
|------------------------|-----------------|------------------|
| `bsoup.star` | `bsoup` | `dustmop/soup`, `golang.org/x/net/html` |
| `compress/gzip.star` | `compress/gzip` | none |
| `encoding/base64.star` | `encoding/base64` | none |
| `encoding/csv.star` | `encoding/csv` | uses `util`, `util/replacecr` |
| `hash.star` | `hash` | none |
| `html.star` | `html` | `PuerkitoBio/goquery` |
| `re.star` | `re` | uses `util` |
| `zipfile.star` | `zipfile` | none |
| (Go helper) | `util` | `pkg/errors`; used by `runtime/modules/starlarkhttp` |

Facts:
- Upstream is archived. Pixlet pins pseudo-version `v0.5.1-0.20220611014110-7fb7ff9ec804`.
- About 2,900 lines of Go (including tests) across the packages above.
- License: MIT (© 2018 QRI, Inc.). Forking is permitted with the notice preserved.
- Upstream `go.mod` declares `go 1.12` and requires packages Pixlet never uses
  (`excelize`, `paulmach/orb`, `yaml.v2`, and others for `dataframe`, `geo`,
  `xlsx`), plus an old `pkg/errors`.
- `runtime/applet_test.go` and `runtime/modules/starlarkhttp/starlarkhttp_test.go`
  also import starlib (`encoding/base64`, `testdata`).

## 3. Goals / Non-goals

**Goals**
- G1: Pixlet depends on a maintained module under the owner's control.
- G2: App-visible behavior of every module is unchanged (contract C2 in `docs/specs/README.md`).
- G3: The fork has its own CI, vulnerability scanning, and dependency updates.
- G4: Pixlet's dependency tree shrinks by dropping unused starlib packages.

**Non-goals**
- Reviving the unused packages (`dataframe`, `geo`, `xlsx`, `math`, `time`, `http`, `asset`).
- Changing the behavior of any existing function. New behavior is additive only (§8).
- Publishing the fork for general community use (possible later; not a goal).

## 4. Key decisions

| # | Decision | Choice | Rationale |
|---|----------|--------|-----------|
| D1 | Fork vs. vendor into Pixlet | **Fork** (separate repo) | Owner preference; independent CI and history; reusable. |
| D2 | Module path | **Rename** to `github.com/jamesn/starlib` | A `replace` directive in Pixlet would work for building Pixlet, but is ignored by anyone importing `tidbyt.dev/pixlet` as a library. A real path is unambiguous. |
| D3 | Repo visibility | **Public** (decided 2026-09-24) | Private modules need `GOPRIVATE` plus credentials in every CI job and on every dev machine. MIT allows public. |
| D4 | Package scope | **Trim** to the 8 modules + `util` (+ `testdata` helper) | Removes unused deps and maintenance surface. |
| D5 | Versioning | Semver tags starting at `v0.6.0` | Continues after upstream's `v0.5.x`; Pixlet pins exact tags. |

## 5. Requirements

### Compatibility (hard gates)
- **R1** Every function in the 8 modules MUST keep its name, parameters (positional and keyword), defaults, return types, and return values.
- **R2** Error messages surfaced to Starlark (e.g. from `fail()` paths or Go errors converted to Starlark errors) MUST keep the same text, since apps may match on them.
- **R3** Load paths (`"re.star"`, `"encoding/csv.star"`, …) are defined in Pixlet and MUST NOT change.
- **R4** A conformance suite (§6, Phase 0) MUST pass against both the old upstream and the fork before Pixlet switches.

### Fork hygiene
- **R5** The fork MUST keep the MIT `LICENSE` and add a `NOTICE`/README section crediting QRI, Inc. and linking the upstream.
- **R6** `go.mod` MUST declare a Go version matching Pixlet's (currently `go 1.26.0`).
- **R7** The fork MUST have CI running `go test ./...`, `go vet ./...`, and `govulncheck ./...` on PRs and weekly.
- **R8** The fork MUST have Dependabot configured for `gomod` and `github-actions`.
- **R9** `github.com/pkg/errors` SHOULD be replaced with standard-library `errors`/`fmt.Errorf("%w")`, subject to R2.

### Pixlet integration
- **R10** Pixlet MUST import only `github.com/jamesn/starlib/...` after the switch; `go.mod` MUST NOT reference `qri-io/starlib`.
- **R11** `docs/modules.md` MUST credit the fork and link to it.

## 6. Plan

### Phase 0: Conformance suite in Pixlet (before any fork work)

Capture current behavior so the fork can prove it didn't change anything.

- Add `runtime/starlib_conformance_test.go` with a Starlark script per module
  that exercises every exported function, including edge cases and error paths,
  and asserts exact results (for error paths, exact message text).
  Examples:
  - `re`: `findall`, `split`, `sub` with groups, flags, no-match, invalid pattern error text.
  - `encoding/csv`: `read_all` with quotes, CRLF input (`replacecr`), `skip`, custom `comma`; `write_all`.
  - `bsoup` / `html`: parse, find, attrs, text on nested markup; malformed HTML.
  - `hash`: `md5`, `sha1`, `sha256` on empty and unicode strings.
  - `encoding/base64`: std/url/raw encodings round-trip; invalid input error text.
  - `compress/gzip`, `zipfile`: decompress fixtures; list/read entries; corrupt-archive errors.
- Runs against the current upstream pin and must pass on `main`.

**Exit:** conformance suite merged and green on `main`.

### Phase 1: Create the fork

1. Fork `qri-io/starlib` → `jamesn/starlib` on GitHub (keeps history).
2. On a branch:
   - Rename the module path to `github.com/jamesn/starlib` and rewrite internal imports.
   - Delete unused packages: `dataframe`, `geo`, `xlsx`, `math`, `time`, `http`, `asset`, plus any top-level aggregator (`starlib.go`) that imports them. Keep `testdata/` (used by tests and by Pixlet's `starlarkhttp` tests).
   - Set `go 1.26.0`; run `go mod tidy` (drops excelize, orb, yaml.v2, etc.).
   - Update README: purpose, credit to QRI, supported modules, compatibility policy (§8).
3. Add CI (R7) and Dependabot (R8).
4. Confirm the upstream tests for the kept packages pass unchanged.
5. Tag `v0.6.0`.

**Exit:** `jamesn/starlib@v0.6.0` exists, CI green, no dependency on unused packages.

### Phase 2: Switch Pixlet to the fork

1. Replace imports in `runtime/applet.go`, `runtime/applet_test.go`,
   `runtime/modules/starlarkhttp/starlarkhttp.go`, and `starlarkhttp_test.go`.
2. `go get github.com/jamesn/starlib@v0.6.0 && go mod tidy`.
3. Run the conformance suite, all tests, and govulncheck.
4. Render every example in `examples/` before and after; compare output hashes.
5. Update `docs/modules.md` (R11).

**Exit:** Pixlet PR merged; `go mod why github.com/qri-io/starlib` reports not needed.

### Phase 3: Modernize dependencies in the fork (behavior-preserving)

Each step is its own fork PR and tag, followed by a Pixlet bump PR that runs the conformance suite.

1. Replace `pkg/errors` (R9), keeping error text byte-identical.
2. Update `goquery`, `golang.org/x/net`, `go.starlark.net` to current versions.
3. Evaluate `dustmop/soup` (last updated 2019). Options: keep as-is, fork it, or reimplement the small subset `bsoup` uses on `x/net/html`. Decide based on conformance results; tracked as a sub-decision.

**Exit:** fork at current dependency versions, govulncheck clean, conformance green.

### Phase 4: Ongoing maintenance

- Dependabot PRs on the fork weekly; tag a patch release when dependencies change.
- Pixlet's Dependabot (spec 002) picks up new fork tags.
- Security advisories: fix in the fork, tag, bump Pixlet within one week for high/critical.

## 7. Testing strategy

| Layer | Where | What it proves |
|-------|-------|----------------|
| Upstream unit tests | fork | Package internals unchanged. |
| Conformance suite | Pixlet | App-visible behavior (R1–R2) unchanged across fork versions. |
| Example render hashes | Pixlet | End-to-end output unchanged (C1). |
| govulncheck | both | No known vulnerabilities. |

## 8. Compatibility policy for the fork

- **Bug fixes that change output are breaking.** If an existing function behaves
  incorrectly, fix it by adding a new function or an opt-in keyword argument, never
  by changing default behavior. Document the quirk in the module README.
- New functions and new optional kwargs are allowed (minor version bump).
- Removing or renaming anything is not allowed.

## 9. Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Dependency updates (goquery, x/net/html) change parsing of malformed HTML | Medium | Conformance suite includes malformed-HTML cases; pin the dependency if it differs. |
| Error text changes when removing `pkg/errors` | Medium | R2 plus error-text assertions in the conformance suite. |
| Maintenance fatigue (a second repo) | Medium | Dependabot plus grouped updates keep it to about one PR per week. |

## 10. Deliverables checklist

- [ ] Phase 0: `runtime/starlib_conformance_test.go` merged in Pixlet
- [ ] Phase 1: `jamesn/starlib` fork, trimmed, renamed, CI, Dependabot, tag `v0.6.0`
- [ ] Phase 2: Pixlet on the fork; `qri-io/starlib` gone from `go.mod`
- [ ] Phase 3: `pkg/errors` removed; dependencies current; `soup` decision recorded
- [ ] Phase 4: maintenance cadence running

## 11. Open questions

1. Should Pixlet's own `starlarkhttp` module (currently derived from starlib's
   `http`) move into the fork too? Proposed: **no**. It contains Pixlet-specific
   caching and belongs with Pixlet.
2. Any modules you'd like added later (e.g. `yaml`)? Out of scope for this project, but the fork makes it possible.
