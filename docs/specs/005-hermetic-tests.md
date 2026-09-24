# 005: Make Network-Dependent Tests Hermetic

Status: Draft

## Summary

Make `runtime.TestInitHTTP` pass without internet access by serving its HTTP
request from a local `httptest.Server` instead of `https://example.com`. Also
record the full local verification checklist in one place.

## Motivation

`TestInitHTTP` (`runtime/httpcache_test.go`) loads `testdata/httpcache.star`,
which fetches `https://example.com`. The test fails in any sandboxed or offline
environment and whenever example.com is unreachable. This makes CI results less
trustworthy and hides real failures behind "that one always fails".

## Goals / Non-goals

Goals:
- `go test ./...` passes offline.
- The test still exercises the HTTP cache initialization path.

Non-goals:
- Rewriting the HTTP cache.

## Requirements

- **R1** `TestInitHTTP` MUST NOT make requests to external hosts.
- **R2** The test MUST still verify that the Starlark `http` module goes through the cache initialized by `InitHTTP`.
- **R3** No other test in `go test ./...` SHOULD require external network access; any found MUST be converted the same way or listed in this spec.
- **R4** The test MUST NOT be skipped or disabled to achieve R1.

## Design

- Start an `httptest.NewServer` in the test that returns a fixed body.
- Pass its URL into the applet through config (or by templating the `.star`
  source) instead of the hard-coded `example.com`.
- Keep the assertions on the response and cache behavior.

To find other offending tests, run `go test ./...` with outbound network blocked
(for example `HTTPS_PROXY=http://127.0.0.1:9` for proxy-respecting clients) and
fix anything that fails.

## Local verification checklist

(Replaces `docs/build-test-todo.md`.)

```bash
npm ci && npm run build          # frontend
make build                       # Go binary with embedded frontend
go test ./...                    # all tests, must pass offline
go vet ./...                     # see spec 006
./pixlet render examples/clock/clock.star
./pixlet serve examples/schema_hello_world/schema_hello_world.star  # UI smoke test
```

## Compatibility

Tests only. No effect on C1–C6.

## Acceptance criteria

- [ ] `go test ./...` passes with outbound network blocked.
- [ ] `TestInitHTTP` still fails if the cache is not initialized (verified by temporarily breaking it).

## Tasks

1. Rewrite `TestInitHTTP` and `testdata/httpcache.star` to use a local server.
2. Run the suite offline to find any other network-dependent tests.
3. Delete `docs/build-test-todo.md`.

## Open questions

- None.
