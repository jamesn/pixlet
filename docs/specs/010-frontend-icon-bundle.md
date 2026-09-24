# 010: Reduce Font Awesome Bundle Size

Status: Draft

## Summary

Shrink the main frontend bundle (about 1.35 MB, mostly Font Awesome) by
**lazy-loading** the icon packs, not by registering a fixed subset of icons.
The subset approach would break apps.

## Motivation

`src/features/theme/DevToolsTheme.jsx` calls `library.add(fas, fab)`, which
bundles every solid and brand icon. The build output shows
`free-solid-svg-icons` (995 KiB) and `free-brands-svg-icons` (533 KiB) as the
largest modules.

## Why not selective registration

The earlier plan (`docs/fontawesome-todo.md`) proposed registering only the icons
in use. However, schema field icons come from **each app's `.star` file**
(`schema.Text(icon = "user")`), resolved at runtime in
`src/features/schema/FieldIcon.jsx` via `findIconDefinition`. Any icon name
not in the registered subset would silently render blank for existing and
future apps. That breaks the implicit contract that any Font Awesome free icon
name works.

## Goals / Non-goals

Goals:
- Smaller initial bundle; faster first paint of `pixlet serve`.
- Every Font Awesome free icon name keeps working.

Non-goals:
- Changing icon naming or the schema API.

## Requirements

- **R1** Any icon name resolvable today MUST still resolve (compatibility).
- **R2** Font Awesome packs MUST be loaded through dynamic `import()` in a separate chunk, not the main bundle.
- **R3** While icons load, fields MUST render (with a placeholder or no icon) and MUST NOT error.
- **R4** The main entry chunk size SHOULD drop by at least 1 MB.
- **R5** `pixlet community validate-icons` / `list-icons` behavior MUST be unchanged (they are Go-side and unaffected).

## Design

- Replace the static `library.add(fas, fab)` with an async loader:
  `Promise.all([import('@fortawesome/free-solid-svg-icons'), import('@fortawesome/free-brands-svg-icons')])`,
  then `library.add(...)` and a state update so `FieldIcon` re-renders.
- `FieldIcon` subscribes to an "icons ready" flag (context or Redux) and renders
  nothing until it is set.
- Measure with `npx source-map-explorer` before and after.

## Compatibility

- No Go-side contract affected. R1 protects app schema icons.

## Acceptance criteria

- [ ] Main chunk shrinks by ≥ 1 MB.
- [ ] An example using several icons (`schema_hello_world`: `user`, `compress`) shows its icons after load.
- [ ] No console errors during load.

## Tasks

1. Measure the baseline bundle.
2. Implement the async loader and ready flag.
3. Measure; run the UI smoke test.
4. Delete `docs/fontawesome-todo.md` and `docs/code-splitting-todo.md`.

## Open questions

- None.
