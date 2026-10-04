# 007: Migrate `google/tink/go` to `tink-crypto/tink-go/v2`

Status: Draft

## Summary

Replace the deprecated `github.com/google/tink/go` module with its maintained
successor `github.com/tink-crypto/tink-go/v2`, and prove that secrets encrypted
by the current code still decrypt.

## Motivation

`google/tink/go` (v1.7.0) is deprecated; development moved to
`tink-crypto/tink-go`. It stays pinned forever with no security fixes. It is used
in `runtime/secret.go` for app secret encryption (hybrid encryption with a JSON
keyset, and a KEK-wrapped private keyset for decryption).

## Goals / Non-goals

Goals:
- Use a maintained Tink module.
- Keep existing encrypted secrets and keysets working.

Non-goals:
- Changing the encryption scheme or keyset format.

## Requirements

- **R1** All imports of `github.com/google/tink/go/...` MUST be replaced with `github.com/tink-crypto/tink-go/v2/...`.
- **R2** A ciphertext produced by the **current** implementation MUST decrypt with the migrated implementation (C5).
- **R3** A ciphertext produced by the migrated implementation MUST decrypt with the current implementation (so rollback is safe).
- **R4** The public API of `runtime.SecretEncryptionKey` / `SecretDecryptionKey` MUST NOT change.
- **R5** `pixlet encrypt` output format MUST NOT change.

## Design

1. **Before migrating**, add a test fixture: a public keyset JSON, an encrypted
   private keyset + KEK (test-only), and a ciphertext produced by the current code.
   Commit these under `runtime/testdata/secret/`.
2. Add a test that decrypts the fixture ciphertext.
3. Swap imports (`hybrid`, `keyset`, `tink`) to the v2 paths. Tink v2 keeps
   keyset serialization wire-compatible, but API names may differ slightly (check
   `keyset.ReadWithNoSecrets` and `keyset.Read`).
4. Run the fixture test to confirm R2, and generate a new ciphertext to confirm R3
   against a build of the old code.

## Compatibility

- **C5 at risk**, guarded by R2/R3 fixture tests written before the change.

## Acceptance criteria

- [ ] Fixture test exists on `main` before the migration PR.
- [ ] After migration, the fixture still decrypts.
- [ ] `go.mod` no longer references `github.com/google/tink/go`.
- [ ] govulncheck clean.

## Tasks

1. PR A: add fixture and decrypt test against the current code.
2. PR B: migrate imports; tests pass unchanged.

## Open questions

- Do any real secrets exist for your sign apps (e.g. API keys in `.star` files)?
  If so, include one real-format (non-sensitive) example in the fixture set.
