# 011: Pluggable Push Targets (Tidbyt API + Self-Hosted)

Status: Draft

## Summary

Introduce a push-target abstraction so `pixlet push` (and the related device
commands) can send rendered images either to the Tidbyt cloud API (the current,
default behavior) or to a self-hosted target. Tidbyt behavior stays exactly as
it is; self-hosted becomes opt-in.

## Motivation

Today the signs are programmed through the Tidbyt cloud:

| Command | Endpoint |
|---------|----------|
| `push` | `POST https://api.tidbyt.com/v0/devices/{id}/push` |
| `devices` | `GET https://api.tidbyt.com/v0/devices` |
| `list` | `GET https://api.tidbyt.com/v0/devices/{id}/installations` |
| `delete` | `DELETE https://api.tidbyt.com/v0/devices/{id}/installations/{iid}` |
| `login` | OAuth at `login.tidbyt.com` |
| `serve` UI push button | same push endpoint (`server/browser/push.go`) |

These URLs are hard-coded in several places (`cmd/push.go`, `cmd/list.go`,
`cmd/delete.go`, `cmd/devices.go`, `server/browser/push.go`, `cmd/config/config.go`).
Pixlet itself is no longer maintained upstream, so relying on a single vendor
cloud is a long-term risk. A self-hosted path lets the signs keep working if the
cloud changes, and allows local-network programming without an internet round trip.

## Goals / Non-goals

Goals:
- One code path for "send a WebP to a device", with interchangeable backends.
- Zero behavior change for current Tidbyt usage.
- A documented, simple protocol for a self-hosted target.

Non-goals (this spec):
- Writing replacement device firmware.
- Building a full scheduling / rotation server. (Possible follow-up spec; see Phase B.)

## Requirements

### Compatibility
- **R1** With no new configuration, every command MUST behave exactly as today: same endpoints, payloads, auth sources (`--api-token`, `$TIDBYT_API_TOKEN`, `pixlet login` token), positional args (including the legacy 3rd arg for installation ID), and error messages (C3, C6).
- **R2** Existing config files MUST load unchanged (C4).

### Abstraction
- **R3** A Go interface MUST define device operations:
  ```go
  type Target interface {
      Push(ctx context.Context, deviceID string, webp []byte, opts PushOptions) error
      ListDevices(ctx context.Context) ([]Device, error)                // MAY return ErrUnsupported
      ListInstallations(ctx context.Context, deviceID string) ([]Installation, error) // MAY return ErrUnsupported
      DeleteInstallation(ctx context.Context, deviceID, installationID string) error  // MAY return ErrUnsupported
  }
  type PushOptions struct { InstallationID string; Background bool }
  ```
- **R4** The Tidbyt implementation MUST be a straight refactor of the existing HTTP code, and MUST be the default.
- **R5** All hard-coded `api.tidbyt.com` URLs in `cmd/` and `server/browser/` MUST go through the target (the `cmd/private` commands already accept `--url` and MAY be left as-is).

### Target selection
- **R6** Target selection MUST follow this precedence: `--target` flag > `$PIXLET_TARGET` > `target` key in the pixlet config file > default `tidbyt`.
- **R7** A target is named and configured in the pixlet config file, e.g.:
  ```yaml
  targets:
    livingroom:
      type: http
      url: http://signs.local:8000
      token_env: SIGNS_TOKEN   # optional; reads token from this env var
  target: livingroom          # optional default
  ```
- **R8** Secrets (tokens) MUST NOT be stored in plain text in the config file; reference env vars or the existing token store.

### Self-hosted HTTP target (Phase A)
- **R9** An `http` target MUST implement `Push` as:
  `POST {url}/v0/devices/{deviceID}/push` with the **same JSON body** as the Tidbyt API
  (`deviceID`, `image` (base64 WebP), `installationID`, `background`) and an optional
  `Authorization: Bearer` header. Reusing the Tidbyt shape means any server that
  emulates the Tidbyt API works unchanged.
- **R10** The `http` target SHOULD implement the list and delete operations with the same path shapes as the Tidbyt API.
- **R11** The protocol MUST be documented in `docs/push-protocol.md` so a server can be written against it.

### UI
- **R12** The `pixlet serve` push button MUST use the selected target.

## Design

### Phase A: Abstraction + HTTP target (this spec)

```
cmd/push.go ─┐
cmd/list.go ─┤
cmd/delete.go┼─► target.Resolve(cmd) ─► target.Target
cmd/devices.go┤                            ├─ tidbyt (default; current code moved here)
server/browser/push.go┘                    └─ http   (Tidbyt-compatible API at a custom base URL)
```

- New package `target/` with the interface, `tidbyt.go`, `http.go`, and `resolve.go`.
- The Tidbyt target is the `http` target with base URL `https://api.tidbyt.com`
  plus the existing token lookup. This keeps one HTTP implementation, and the
  default case is provably identical.
- Contract tests: an `httptest.Server` records requests; the same test table runs
  against both targets and asserts identical method, path, headers, and body.

### Phase B: Device side (moved)

> **Superseded:** Phase B is now the independent [firmware fork project](../projects/firmware-fork/PLAN.md).
> Independence requires firmware we control, not only a push target. The options below are kept for history.

#### Original Phase B options

Once pixlet can push to any URL, the remaining question is what receives the
image and drives the sign. Options:

| Option | What it is | Trade-offs |
|--------|------------|------------|
| **B1. Stay on Tidbyt cloud** | Default target only. | No work; depends on the vendor cloud. |
| **B2. Community self-hosted server + replacement firmware** | Open-source projects exist that re-flash Tidbyt hardware to pull from a self-hosted server that renders pixlet apps. (Verify the current state and license before choosing.) | Full independence; requires re-flashing each sign; must track that project. |
| **B3. Own minimal server** | A small service in this repo (`pixlet serve-devices`?) that accepts R9 pushes, stores the latest WebP per device and installation, and serves a rotation feed the firmware polls. | Most control; still needs device firmware that polls it (B2's firmware or custom). |

Recommendation: do Phase A now (valuable on its own, zero risk to current usage),
then evaluate B2 vs B3 in a follow-up spec once the hardware model(s) of the
signs are confirmed.

## Compatibility

- **C3 / C6 at risk** from the refactor, guarded by R1 and the contract tests,
  which assert the default target sends byte-identical requests to what the
  current code sends.
- C4 guarded by R2 (new config keys are optional).

## Acceptance criteria

- [ ] Contract tests prove the default target's requests match the pre-refactor requests exactly (method, URL, headers, JSON body) for push (with/without installation ID, background), devices, list, and delete.
- [ ] `pixlet push <device> <file>` with no new config works against the real Tidbyt API (manual check by owner).
- [ ] `pixlet push --target livingroom <device> <file>` sends a Tidbyt-shaped request to the configured URL.
- [ ] The `serve` UI push button honors the selected target.
- [ ] `docs/push-protocol.md` documents the protocol.

## Tasks

1. Capture current request shapes in contract tests (before refactor).
2. Add the `target/` package; move the Tidbyt HTTP code into it.
3. Route `push`, `list`, `delete`, `devices`, and the UI push handler through `target.Resolve`.
4. Add the `http` target and config/flag/env selection.
5. Write `docs/push-protocol.md`.
6. Manual check against the Tidbyt API.

## Open questions

- Which hardware are the signs (Tidbyt Gen 1 / Gen 2 / other HUB75 panels)? This drives the Phase B choice.
- Does `pixlet login` (OAuth to `login.tidbyt.com`) need a self-hosted equivalent, or is a static bearer token enough for a home network?
- Should the self-hosted target support pushing to multiple signs at once (fan-out), e.g. `--device all`?
