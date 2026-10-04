# 004: Remove Dead Frontend Build Dependencies

Status: In progress (Phase 1 PR)

## Summary

Remove unused and deprecated Babel packages from `package.json`, which also drops
the unmaintained `core-js@2` from the install tree.

## Motivation

- `babel-preset-react@6` (Babel 6 era) is in `devDependencies` but is not
  referenced. `.babelrc` uses `@babel/preset-react`. It pulls in `core-js@2`,
  which `npm ci` warns is unmaintained.
- `@babel/plugin-proposal-class-properties` is deprecated. Class properties are
  standard JavaScript and are handled by `@babel/plugin-transform-class-properties`
  or `@babel/preset-env`.

## Goals / Non-goals

Goals:
- Smaller, warning-free `npm ci`.
- Identical built output behavior.

Non-goals:
- Migrating the bundler or changing browser targets.

## Requirements

- **R1** `babel-preset-react` MUST be removed.
- **R2** `@babel/plugin-proposal-class-properties` MUST be replaced with `@babel/plugin-transform-class-properties` in `package.json` and `.babelrc`.
- **R3** `npm ci` MUST complete without the `core-js@2` deprecation warning.
- **R4** `npm run build` MUST succeed, and the UI MUST work as before (see acceptance criteria).
- **R5** `copy-webpack-plugin` SHOULD move from `dependencies` to `devDependencies`, since it is build-only.

## Design

Package changes only. Check whether any source file uses class-property syntax;
if none does, the plugin can be dropped entirely rather than replaced.

## Compatibility

Build only. No effect on C1–C6. The frontend is embedded in the Go binary, so the
UI smoke test below is the guard.

## Acceptance criteria

- [ ] `npm ci` shows no `core-js` deprecation warning.
- [ ] `npm run build` succeeds.
- [ ] `pixlet serve examples/schema_hello_world/schema_hello_world.star` loads with no console errors, and config edits re-render the preview.

## Tasks

1. `npm uninstall babel-preset-react @babel/plugin-proposal-class-properties`.
2. Add `@babel/plugin-transform-class-properties` if class properties are used; update `.babelrc`.
3. Move `copy-webpack-plugin` to devDependencies.
4. Build, then run the UI smoke test.

## Open questions

- None.
