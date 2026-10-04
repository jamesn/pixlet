# Pixlet Fork: Master Roadmap

Status: Active
Owner: @jamesn
Last updated: 2026-09-24
Current release: [v0.41.0](https://github.com/jamesn/pixlet/releases/tag/v0.41.0)

This is the top-level guide for development on this fork. It sets the strategy,
the order of work, and the rules every change follows. Detail lives in the
linked specs (`docs/specs/`) and projects (`docs/projects/`). This document
says *what* comes next and *why*.

---

## 1. Mission

**Keep the house sign boards running, securely, for the long term, without
depending on anyone else to maintain the software.**

Pixlet (upstream `tidbyt/pixlet`) is abandoned. This fork is the only
maintained copy the owner relies on. Every decision is weighed against one
question: *does this make the signs more likely to keep working, and easier to
keep working?*

## 2. Strategic pillars

| # | Pillar | What it means | How we measure it |
|---|--------|---------------|-------------------|
| P1 | **Secure by default** | No known vulnerabilities ship; new ones are caught automatically. | `govulncheck` = 0, `npm audit` high/critical = 0, on every PR and weekly. |
| P2 | **Never break existing apps** | Existing `.star` apps, CLI usage, config, and secrets keep working unchanged. | Compatibility contracts C1–C6 (§3) enforced by tests. |
| P3 | **Own the critical path** | Replace archived or abandoned dependencies with code we control. | Count of archived/unmaintained deps in the render and runtime path → 0. |
| P4 | **Reduce vendor lock-in** | The signs keep working without the Tidbyt cloud. This requires firmware we control, not just a push path. | At least one sign runs on firmware + server the owner builds; ultimately all signs. |
| P5 | **Cheap to maintain** | Routine upkeep is small, automated, and predictable. | About 1 grouped dependency PR per ecosystem per week; CI green on `main`. |

When pillars conflict, **P2 wins over everything except a critical
security fix (P1)**. Even then, the fix must be the least behavior-changing
option available.

## 3. Compatibility contracts (non-negotiable)

These are defined in [`docs/specs/README.md`](specs/README.md) and apply to every change:

| ID | Contract |
|----|----------|
| C1 | Existing apps render pixel-identical output. Geometry (dimensions, frame count) is a hard gate. |
| C2 | Starlark module names, signatures, return values, and error text are unchanged. |
| C3 | CLI commands, positional args, flags, and env vars keep working. |
| C4 | Existing config files and stored OAuth tokens remain valid. |
| C5 | Secrets encrypted by earlier versions still decrypt. |
| C6 | `pixlet push` to the Tidbyt API remains the default. |

**Rule for behavior changes:** if something is wrong but apps may depend on it,
fix it *additively* (new function, new optional argument, opt-in flag), never
by changing the default. Example: `go-humanize` is pinned at v1.0.1 because
v1.1.0 changes `humanize.relative_time` output.

## 4. Roadmap by horizon

```
 DONE            NOW               NEXT                    LATER                      FUTURE
 v0.41.0  ──►   v0.42.0   ──►     v0.43.0      ──►        v0.44+          ──►        (decision-gated)
 security       maintenance       code health +           own the critical           device
 baseline       baseline          push abstraction        path                       independence
```

### ✅ Done: Security baseline (v0.41.0)

| Item | Result |
|------|--------|
| #11: Go 1.26, all Go modules updated, npm fixes, CI on Node 22 | govulncheck 54 → 0; npm 33 → 2 |
| #12: React 18 + react-router 7.18 | npm 2 → 0 |
| v0.41.0 released | Linux, macOS, Windows binaries published |
| Spec plan + project plans written | `docs/specs/`, `docs/projects/` |

### 🔨 Now: Maintenance baseline (target v0.42.0)

One low-risk PR. No runtime behavior changes.

| Spec | Item | Pillar |
|------|------|--------|
| [001](specs/001-ci-security-scanning.md) | govulncheck + npm audit in CI, weekly scheduled scan | P1 |
| [002](specs/002-dependency-automation.md) | Dependabot only (remove `renovate.json`), grouped updates, humanize pinned | P1, P5 |
| [003](specs/003-github-actions-refresh.md) | Current Actions majors; SHA-pin third-party actions | P5 |
| [004](specs/004-frontend-build-cleanup.md) | Remove dead Babel packages (drops `core-js@2`) | P5 |
| [005](specs/005-hermetic-tests.md) | `TestInitHTTP` works offline; local verification checklist | P2, P5 |

**Why first:** 001 and 005 make every later PR safer to review, and 002 keeps
the security baseline from decaying while the bigger work happens.

### ⏭️ Next: Code health + push abstraction (target v0.43.0)

| Item | Summary | Pillar |
|------|---------|--------|
| [006](specs/006-go-vet-cleanup.md) | Fix `go vet` findings (2 real context leaks, malformed struct tags); enforce vet in CI | P2, P5 |
| [011 Phase A](specs/011-device-push-targets.md) | `Target` interface; Tidbyt stays default; opt-in Tidbyt-compatible HTTP target | P4 |
| **Tidbyt data backup** (task, no spec) | Export device IDs and installations (`pixlet devices`, `pixlet list`) to a private location while the API works | P4 |

**Why here:** the push abstraction has zero impact on current usage and turns
any future move off the Tidbyt cloud into a config change. The backup is cheap
insurance that only works while the cloud is up.

### 🔭 Later: Own the critical path (v0.44.0 onward)

Independent tracks. Each ships as its own PR(s) and minor release, in any order.

| Item | Summary | Contract guarded | Pillar |
|------|---------|------------------|--------|
| [Project: starlib fork](projects/starlib-fork/PLAN.md) | Public `jamesn/starlib`, trimmed to 8 modules; Pixlet-side conformance suite first | C2 | P3 |
| [Project: image resize](projects/image-resize/PLAN.md) | Bring nfnt's resize in-house behind a golden corpus (`x/image/draw` verified to differ) | C1 | P3 |
| [007](specs/007-tink-migration.md) | `google/tink/go` → `tink-crypto/tink-go/v2`, with decrypt fixtures first | C5 | P1, P3 |
| [010](specs/010-frontend-icon-bundle.md) | Lazy-load Font Awesome (not selective registration) | App schema icons | P5 |

**Suggested order:** starlib fork first (largest app-facing surface, archived,
has transitive deps), then Tink (security library), then resize (stable code,
low urgency), with 010 whenever convenient.

### 🧭 Independent track: Firmware fork (device independence)

[Project: firmware fork](projects/firmware-fork/PLAN.md). Runs in parallel with the
horizons above and neither blocks nor is blocked by them.

Pixlet only produces images; the Tidbyt firmware only talks to the Tidbyt cloud.
Real independence requires **firmware we control** plus **a server we run**.
Approach: adopt open-source firmware and server → bench spike on one sign →
multi-week pilot → fork under the owner's account with CI and signed OTA → roll
out sign by sign. Phase 0 (research + hardware inventory) needs no hardware changes.

The server runs on the owner's existing Kubernetes cluster:
[Project: k8s hosting](projects/k8s-hosting/PLAN.md) (~2–5 days). Its Phase 1, a
multi-arch **pixlet container image** published to `ghcr.io/jamesn/pixlet`, is
useful on its own and can ship in any pixlet release.

### 🌅 Future: Decision-gated work

These need a decision or a new spec before work starts.

| Item | Gate | Notes |
|------|------|-------|
| **Frontend platform modernization** | New spec | MUI 5 → current, `@mui/x-date-pickers` 5 → current, Redux Toolkit 1 → 2, react-redux 8 → 9, React 19 (unlocks react-router 8). Large UI-regression surface; needs a UI test harness first. |
| **Dependency watch list review** | New spec | Assess the remaining forks and old pins in the render/runtime path: `tidbyt/go-libwebp` (cgo), `tidbyt/gg`, `srwiley/oksvg` + `rasterx`, `ericpauley/go-quantize`, `zachomedia/go-bdf`, `skip2/go-qrcode`, `manifoldco/promptui`, `newm4n/go-dfe`, `dustmop/soup` (via starlib). Classify each: fine / fork / replace. |
| **Additive features** | Owner demand | e.g. opt-in `render.Image` scaling modes, new Starlark modules via the starlib fork, multi-sign push fan-out. |

## 5. Dependency graph

```
001 CI scanning ─┐
005 hermetic ────┼──► 006 vet ──► (all Later items benefit from strict CI)
002/003/004 ─────┘
                 └──► 011 Phase A   (optional input to the firmware track's server)

firmware fork (independent): Phase 0 research ─► bench spike ─► pilot ─► fork ─► rollout

starlib fork:  Phase 0 conformance tests ─► fork ─► switch ─► modernize (soup decision)
image resize:  Phase 0 golden corpus ─► internalize ─► remove dep
tink:          decrypt fixtures ─► migrate
frontend modernization:  UI test harness ─► MUI/Redux/React upgrades
```

Rule: **characterization tests land before the change they protect**, always
in a separate, earlier PR.

## 6. How work flows

1. **Plan:** every non-trivial change has a spec (`docs/specs/NNN-*.md`) or a
   project plan (`docs/projects/*/PLAN.md`) with numbered requirements and
   acceptance criteria. Status goes Draft → Approved → In progress → Done.
2. **Protect:** if a compatibility contract is at risk, first merge a test that
   captures current behavior (golden images, conformance scripts, decrypt fixtures).
3. **Change:** small PRs, one spec or one project phase each. Branch from `main`.
4. **Verify locally** (spec 005 checklist): `npm ci && npm run build`,
   `make build`, `go test ./...` offline, `go vet ./...`, render examples, and
   drive `pixlet serve` in a browser for UI changes.
5. **CI green** on Linux, macOS, and Windows, plus security scans (after 001).
6. **Merge** with a merge commit (repo convention).
7. **Release** when a horizon or a meaningful batch completes (§7).
8. **Update this roadmap**: move items between horizons and add to the decision log.

### Definition of done (any item)

- [ ] Acceptance criteria in its spec are all checked
- [ ] No compatibility contract regressed (tests prove it)
- [ ] CI green on all three platforms; govulncheck and npm audit clean
- [ ] Docs updated (`docs/modules.md`, specs status, this roadmap)
- [ ] Released, or queued for the next release

## 7. Release strategy

- Versions come from **git tags** (`vMAJOR.MINOR.PATCH`); pushing a tag runs the
  release build in `.github/workflows/main.yml`. No version is stored in files.
- **Minor** (`v0.42.0`): each completed horizon, or any toolchain, dependency-major, or feature change.
- **Patch** (`v0.42.1`): security fixes and dependency bumps with no behavior change.
- **Major (`v1.0.0`)**: reserved. Candidate milestone is when the fork owns its
  critical path (Later horizon complete) and 011 Phase A has shipped.
- Release notes: summarize every PR since the last tag (the auto-generated
  notes only show the last merge commit; edit them after publishing).
- High/critical security advisories: patch release within **7 days**.

## 8. Decision log

| # | Date | Decision | Rationale | Where |
|---|------|----------|-----------|-------|
| D-001 | 2026-09-24 | Go toolchain 1.26 (`toolchain go1.26.8`); CI reads the version from `go.mod` | Go 1.24 out of support; stdlib CVEs | #11 |
| D-002 | 2026-09-24 | Pin `go-humanize` v1.0.1 | v1.1.0 changes `relative_time` output (C2) | #11 |
| D-003 | 2026-09-24 | Node 22 in CI | Node 16 EOL; required by webpack-dev-server 6 | #11 |
| D-004 | 2026-09-24 | React 18 + react-router 7.18, not 8 | 8 requires React 19, a large migration for no security gain | #12 |
| D-005 | 2026-09-24 | Fork starlib as **public** `jamesn/starlib` with a renamed module path | Avoids private-module friction; `replace` does not propagate | [starlib plan](projects/starlib-fork/PLAN.md) |
| D-006 | 2026-09-24 | Internalize nfnt resize; **geometry is a hard gate** | `x/image/draw` changes downscale pixels; owner requires geometric compatibility | [resize plan](projects/image-resize/PLAN.md) |
| D-007 | 2026-09-24 | Font Awesome: lazy-load, not selective registration | Selective registration breaks app-chosen schema icons | [010](specs/010-frontend-icon-bundle.md) |
| D-008 | 2026-09-24 | Plan now for leaving the Tidbyt cloud; build only the abstraction (Phase A) | Low-cost insurance; device side depends on hardware | [011](specs/011-device-push-targets.md) |
| D-010 | 2026-09-24 | Device independence = control the firmware; run it as an independent project (adopt → pilot → fork) | Stock firmware only talks to the Tidbyt cloud | [firmware plan](projects/firmware-fork/PLAN.md) |
| D-011 | 2026-09-24 | Host the sign server on the existing Kubernetes cluster; publish a pixlet container image | Owner already runs a cluster; container-first server keeps deploys declarative | [k8s plan](projects/k8s-hosting/PLAN.md) |
| D-009 | proposed | Dependabot as the only update bot | Renovate config is stale and set to automerge | [002](specs/002-dependency-automation.md) |

## 9. Risk register

| Risk | Impact | Likelihood | Mitigation | Owner item |
|------|--------|------------|------------|------------|
| Tidbyt cloud API changes or shuts down | Signs stop updating | Unknown | Firmware fork project; Tidbyt data backup | Independent track |
| New CVE in a dependency goes unnoticed | Security exposure | High without automation | 001 scheduled scans; 002 Dependabot | Now |
| Archived dependency has an unfixable bug | Broken or insecure apps | Medium | starlib fork, resize internalization, Tink migration | Later |
| A "harmless" upgrade changes app output | Signs show wrong content | Medium | Contracts C1–C6; characterization tests first; no automerge | All |
| Windows/macOS-only build breaks | Release fails | Low | 3-platform CI on every PR | Ongoing |
| Maintenance fatigue | Fork decays again | Medium | Grouped weekly updates; small PRs; this roadmap keeps focus | Ongoing |

## 10. Open decisions (owner input needed)

| # | Question | Blocks |
|---|----------|--------|
| Q1 | Which hardware are the signs (Tidbyt Gen 1 / Gen 2 / other HUB75), and how many? | Firmware fork Phase 0 |
| Q2 | For self-hosted push: is a static bearer token on the home network enough, or is OAuth needed? | 011 Phase A design detail |
| Q3 | Is Renovate installed as a GitHub App on the repo? (If yes, uninstall it with 002.) | 002 |
| Q4 | Target date or trigger for `v1.0.0`? | Release strategy |
| Q5 | Cluster details: distribution, node architectures, LoadBalancer, storage class, GitOps tool? | k8s hosting Phase 0 |

## 11. Document map

```
docs/
├── ROADMAP.md                  ← you are here (strategy, order, rules)
├── specs/
│   ├── README.md               ← contracts C1–C6, spec format, phase index
│   └── 001…011-*.md            ← individual specs
└── projects/
    ├── starlib-fork/PLAN.md    ← multi-repo project
    ├── firmware-fork/PLAN.md   ← independent track: device independence
    ├── k8s-hosting/PLAN.md     ← sign server on the home Kubernetes cluster
    └── image-resize/PLAN.md    ← golden-corpus-gated project
```

Superseded: `docs/dependency-upgrades-todo.md`, `docs/build-test-todo.md`,
`docs/code-splitting-todo.md`, `docs/fontawesome-todo.md` (removed as their
replacing specs complete).
