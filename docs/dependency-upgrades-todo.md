# Dependency Upgrades TODO

Purpose: Track remaining frontend dependency work and verification after recent upgrades (color picker, webpack-dev-server, etc.).

Last updated: 2025-11-16

Tasks

- [ ] **Pin upgraded packages**: Update `package.json` to pin versions for packages that were upgraded (e.g., `webpack-dev-server`, `react-colorful`, any transitive dependencies you care about).

- [ ] **Resolve peer dependency warnings**: Re-install and address any peer-dependency warnings. If unavoidable, add a short note to `docs/` explaining why.

- [ ] **Install and audit**:

```bash
npm ci
npm audit fix --force
npm install
```

- [ ] **Manual verification**: Run the app locally (`npm run start` / `npm run build`) and check for console warnings and runtime errors.

- [ ] **Update lockfile**: Commit updated `package-lock.json` or `yarn.lock` after confirming tests/build succeed.

- [ ] **Record security status**: Add a short note with `npm audit` output (vulnerabilities fixed / remaining) into this file.

Notes

- Some upgrades (e.g., `webpack-dev-server`) were already applied earlier in the session; this checklist ensures we finish validation and commit the resulting lockfile.
