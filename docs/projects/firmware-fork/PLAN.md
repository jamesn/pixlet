# Project: Firmware Fork (Device Independence)

Status: Draft (Phase 0 not started)
Owner: @jamesn
Last updated: 2026-09-24
Track: **Independent**. This project does not block, and is not blocked by, the
other work in [`docs/ROADMAP.md`](../../ROADMAP.md). It supersedes
[spec 011 Phase B](../../specs/011-device-push-targets.md).

## 1. Summary

Run the house signs on firmware and a server the owner controls, so the signs keep
working regardless of the Tidbyt cloud. "Control" means having the source,
building it ourselves, and flashing it ourselves. The preferred route is to adopt
an existing open-source firmware and self-hosted server, validate them on one
sign, then fork them (the same pattern as the
[starlib fork](../starlib-fork/PLAN.md)).

## 2. Background: why pixlet alone isn't enough

```
TODAY
  .star app ──pixlet render──► WebP ──pixlet push──► Tidbyt cloud ◄──polls/streams── Tidbyt firmware ──► LED panel
                                                     (decides what shows, when)      (talks only to Tidbyt)

TARGET
  .star app ──render──► WebP ──► self-hosted server ◄──fetches── firmware we build ──► LED panel
                                 (rotation, schedule)          (points at our server)
```

- Pixlet's job ends at producing a WebP image. Delivery and scheduling happen in
  the Tidbyt cloud, and the stock firmware only accepts content from there.
- [Spec 011 Phase A](../../specs/011-device-push-targets.md) (a configurable push
  target) is useful groundwork for the server side, but **does not by itself make
  the signs independent**: stock firmware will never fetch from another server.
- Therefore independence requires both **(a) a server we run** and **(b) firmware
  we control**.

### What we know / don't know yet

| Topic | Status |
|-------|--------|
| Sign hardware model(s) and count | **Unknown.** Owner to confirm (Phase 0). Tidbyt Gen 1 is believed to be ESP32-based with a 64×32 HUB75 panel. Verify, and check Gen 2 separately. |
| Open-source replacement firmware for Tidbyt hardware | Believed to exist in the community (firmware + self-hosted server). **Not yet researched.** Current state, license, and supported models to be verified in Phase 0. |
| Whether stock Tidbyt firmware can be restored after re-flashing | **Unknown.** This is a critical Phase 0 question. |
| Device IDs and installations in the Tidbyt cloud | Retrievable today via `pixlet devices` / `pixlet list`. **Back up now.** |

## 3. Goals / Non-goals

**Goals**
- G1: Every sign can operate with **no dependency on Tidbyt infrastructure**.
- G2: The owner can build firmware and server from source, reproducibly, and flash or update every sign.
- G3: Existing pixlet apps display the same on the new stack (contracts C1/C2 in [`docs/specs/README.md`](../../specs/README.md)).
- G4: Migration is **reversible per sign** until the owner decides otherwise.

**Non-goals**
- Designing new hardware.
- Supporting devices the owner doesn't own.
- Cloud hosting. The server runs on the home network.
- Replacing pixlet as the renderer. The server uses pixlet (this fork) to render apps.

## 4. Principles and constraints

1. **Adopt, then fork. Don't write from scratch** unless Phase 0 finds nothing usable.
2. **One sign first.** All other signs stay on stock firmware until the pilot passes (Phase 2 exit).
3. **Know the way back before flashing.** Either a verified restore-to-stock procedure exists, or the owner explicitly accepts one-way migration for the pilot sign.
4. **Local-first and secure:** the server listens only on the home network; firmware stores Wi-Fi credentials securely; OTA updates are authenticated.
5. **Same apps, same output:** the server renders with this pixlet fork, so contracts C1/C2 carry over automatically.

## 5. Requirements

### Independence
- **R1** A sign on the new stack MUST display apps with Tidbyt's cloud unreachable (verified by blocking `*.tidbyt.com` at the router).
- **R2** Firmware and server MUST build from source in CI under the owner's GitHub account.
- **R3** Licenses of adopted projects MUST permit forking, modifying, and private use; the license and attribution MUST be preserved.

### Safety
- **R4** A documented, tested recovery procedure MUST exist for a sign that fails to boot (e.g. USB re-flash).
- **R5** A restore-to-stock procedure SHOULD exist and be tested on the pilot sign before Phase 3. If it can't exist, the owner MUST sign off on one-way migration.
- **R6** Firmware updates over the network (OTA) MUST be authenticated. Unsigned images MUST be rejected.

### Functionality parity (minimum)
- **R7** App rotation across multiple installed apps with per-app refresh intervals.
- **R8** Brightness control (at minimum manual; scheduled/auto if the hardware supports it).
- **R9** Recovery from Wi-Fi and server outages without manual intervention.
- **R10** The server MUST render apps with this pixlet fork (as a library or CLI) and MUST support app config (schema values) and secrets.

### Operability
- **R11** Server health and per-sign last-seen status MUST be visible (web page or log).
- **R12** Backups: server config and app installations MUST be exportable and restorable.

## 6. Key decisions (to be made in-project)

| # | Decision | Options | Decide in |
|---|----------|---------|-----------|
| FD1 | Firmware base | Adopt community firmware → fork / write own | Phase 0 |
| FD2 | Server base | Adopt the community server that pairs with FD1 → fork / build own around pixlet | Phase 0 |
| FD3 | Device ↔ server protocol | Whatever FD1/FD2 use (preferred, less work) / define our own | Phase 0 |
| FD4 | Where the server runs | Existing home machine / small dedicated box (e.g. Pi) / container on NAS | Phase 1 |
| FD5 | Repo layout | Separate forks (`jamesn/<firmware>`, `jamesn/<server>`) / monorepo | Phase 3 |
| FD6 | Relationship to spec 011 Phase A | Keep pixlet push targets pointing at the server / server pulls and renders itself (push unnecessary) | Phase 1 |

## 7. Plan

### Phase 0: Research and inventory (no hardware changes)

1. **Inventory the signs:** model/generation, count, and where each is mounted
   (for USB access during flashing). Photograph labels and board revisions.
2. **Back up Tidbyt data:** `pixlet devices` and `pixlet list <device>` for every
   sign; export app configs; store privately. (Also listed as a task in the roadmap's Next horizon.)
3. **Survey candidate firmware and servers.** For each candidate, score:

   | Criterion | Why it matters |
   |-----------|----------------|
   | Supports the owner's hardware model(s) | Hard requirement |
   | License (forkable, R3) | Hard requirement |
   | Maintenance activity (commits, releases, issue response) | Longevity of the upstream you fork from |
   | Build toolchain (e.g. ESP-IDF / PlatformIO) reproducible in CI | R2 |
   | OTA with authentication | R6 |
   | Protocol simplicity and documentation | FD3; ease of owning it |
   | Server uses pixlet (which version?) | R10; compatibility with this fork |
   | Restore-to-stock documented | R5 |
   | Feature parity (rotation, brightness) | R7–R8 |

4. **Answer the reversibility question:** can stock Tidbyt firmware be dumped
   before flashing and written back? Document the exact tools and steps.
5. **Write the findings** into this plan (§2 table, FD1–FD3) and update the roadmap.

**Exit:** hardware confirmed; a chosen firmware + server (or a decision to build our
own); a restore/recovery procedure documented; Tidbyt data backed up.

### Phase 1: Bench spike (one sign, off the wall)

1. Dump the stock firmware from the pilot sign (if possible) and verify the dump.
2. Stand up the chosen server on the home network (FD4). Point it at this pixlet fork.
3. Flash the pilot sign with the chosen firmware; configure Wi-Fi and server URL.
4. Install 2–3 of the owner's real apps; verify output matches the Tidbyt-rendered version (screenshots/WebP hashes).
5. Test R1 (block `*.tidbyt.com`), R9 (pull power, kill Wi-Fi, stop server), R4 (USB recovery).
6. **Restore the pilot sign to stock** and confirm it works with Tidbyt again (R5). Then re-flash for the pilot.

**Exit:** all of the above pass, or findings recorded and the owner decides whether to proceed.

### Phase 2: Pilot (one sign, in daily use)

- Run the pilot sign on the new stack for **at least 2–4 weeks** alongside the stock signs.
- Track: uptime, missed refreshes, crashes/reboots, any app rendering differences.
- Fix or file issues upstream; keep notes of every local change needed.

**Exit criteria:** no unrecovered failures during the pilot; owner is satisfied with day-to-day behavior.

### Phase 3: Fork and own

1. Fork the firmware and server repos to the owner's account (FD5); keep license and attribution (R3).
2. Add CI: reproducible firmware builds (artifacts attached to tags), server tests,
   dependency and vulnerability scanning, Dependabot (mirroring pixlet's specs 001/002).
3. Set up **OTA signing keys** (R6); document key storage and rotation.
4. Pin the server to this pixlet fork's releases (e.g. `v0.4x`) so app compatibility is governed by pixlet's contracts.
5. Re-flash the pilot sign from **our** build to prove the chain end to end.

**Exit:** the pilot sign runs firmware built by our CI and talks to our server build.

### Phase 4: Rollout

- Migrate the remaining signs one at a time, confirming each before the next.
- Keep the Tidbyt data backup until all signs have run cleanly for a month.
- Update pixlet docs: how the house signs are programmed now.

**Exit:** all signs independent (G1). The roadmap's P4 pillar is marked achieved.

### Phase 5: Ongoing maintenance

- Merge upstream fixes into the forks (monthly check) while upstream is alive.
- Security: patch firmware/server within 7 days for high/critical issues (same policy as pixlet).
- OTA rollouts: pilot sign first, then the rest.

## 8. Relationship to other work

| Item | Relationship |
|------|--------------|
| Pixlet contracts C1/C2 | Server renders with this fork, so app output guarantees carry over. |
| Spec 011 Phase A (push targets) | Optional. Useful if the chosen server accepts pushes; unnecessary if the server renders apps itself (FD6). Neither blocks the other. |
| starlib fork / resize / Tink projects | Improve the pixlet the server uses; not prerequisites. |
| Roadmap | Independent track; appears in `docs/ROADMAP.md` alongside the horizons, not inside them. |

## 9. Risks

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Re-flashing bricks a sign | Lost sign | Low–Med | USB recovery procedure (R4) tested in Phase 1; pilot on one sign |
| Stock firmware can't be restored | One-way migration | Unknown | Answer in Phase 0; owner sign-off required (R5) |
| Community project goes dormant | Fork carries all maintenance | Medium | That's why Phase 3 forks; choose the simplest codebase in Phase 0 |
| Hardware revision differences between signs | Firmware works on some signs only | Medium | Inventory in Phase 0; test each revision before rollout |
| Server machine failure | All signs stale | Medium | R9 (signs keep last content / retry), R12 backups, simple redeploy |
| Unauthenticated OTA or open server port | Someone else controls the signs | Low | R6 signing; local-only server; no port forwarding |
| Server's bundled pixlet differs from this fork | App output differences | Medium | Phase 3 step 4: pin to this fork |

## 10. Deliverables checklist

- [ ] Phase 0: hardware inventory; Tidbyt data backup; candidate evaluation; restore/recovery answer; FD1–FD3 decided
- [ ] Phase 1: bench spike passes; restore-to-stock verified (or signed off)
- [ ] Phase 2: pilot runs 2–4 weeks cleanly
- [ ] Phase 3: forks under owner's account with CI, signed OTA, pinned to this pixlet
- [ ] Phase 4: all signs migrated
- [ ] Phase 5: maintenance cadence running

## 11. Open questions

1. Which hardware are the signs, and how many? *(Blocks Phase 0 evaluation.)*
2. Is a one-way migration acceptable if restore-to-stock turns out to be impossible?
3. Where should the server run (FD4)? Is there an always-on machine at home already?
4. Must any sign keep working with the Tidbyt mobile app during the transition (e.g. for other household members)?
