# Build & Test Verification TODO

Purpose: Commands and checks to validate changes and ensure CI parity.

Last updated: 2025-11-16

Tasks

- [ ] **Local build**:

```bash
# clean, install, build
rm -rf node_modules dist
npm ci
npm run build
```

- [ ] **Run frontend tests** (if present):

```bash
npm test
```

- [ ] **Run Go tests / full test suite**:

```bash
# repository root
make test
# or
go test ./... -v
```

- [ ] **Generate widgets / embed fonts** (if making rendering changes):

```bash
make widgets
make embedfonts
```

- [ ] **CI validation**: Ensure the CI workflow runs the same commands and that the `dist` artifacts are validated; add bundle-size check if needed.

- [ ] **Record results**: After running builds/tests, copy key outputs (bundle sizes, failing tests, audit results) into `docs/build-test-notes.md` for traceability.

Notes

- If any build or test step fails locally, collect logs and open an issue or a branch to fix before merging.
- Consider adding a lightweight bundle-size CI job that fails on regressions.
