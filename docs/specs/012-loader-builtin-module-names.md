# 012: App Files Named After Built-in Modules Fail to Load

Status: Draft

## Summary

An app file whose name matches a built-in module (for example `humanize.star`,
`qrcode.star`, `sunrise.star`, `time.star`) cannot `load()` that built-in module.
The loader resolves the name to the app's own file first, then reports a circular
dependency. Fix the loader so that a `load()` that would cycle back to a file
already being loaded falls back to the built-in module, and keep every other
resolution exactly as it is today.

## Motivation

Three bundled examples fail to render, on `main` and on the v0.41 / v0.42
releases alike:

```
$ pixlet render examples/humanize/
Error: failed to load applet: starlark.ExecFile: cannot load humanize.star:
circular dependency detected: humanize.star -> humanize.star
```

Same for `examples/qrcode/` and `examples/sunrise/`. Any of the owner's apps
named after a module it loads will fail the same way.

### Root cause

`runtime/applet.go`, `ensureLoaded()` installs a `thread.Load` that resolves a
module name in this order:

1. If a file with that name exists in the app's filesystem, load it as a
   **local** module (via `ensureLoaded`, which tracks `currentlyLoading` to
   detect cycles).
2. Otherwise fall back to `loadModule()`: first the optional custom
   `a.loader`, then the built-in modules.

When `humanize.star` executes `load("humanize.star", "humanize")`, step 1 finds
the app's own `humanize.star`. That file is already in `currentlyLoading`, so
the cycle check fires. The built-in module is never consulted.

This ordering arrived upstream with multi-file applet support (`ce0d274`,
"runtime: Support applets with multiple Starlark files (#1041)", 2024-04-18).
Before that change, `load()` went straight to `loadModule`, so these examples
would have resolved to the built-ins (inferred from the pre-change loader, not
re-run). **This is a regression**, not a design choice about single-file apps.

### Built-in module names (current `loadModule` switch)

`render.star`, `animation.star`, `schema.star`, `cache.star`, `secret.star`,
`xpath.star`, `bsoup.star`, `compress/gzip.star`, `compress/zipfile.star`,
`encoding/base64.star`, `encoding/csv.star`, `encoding/json.star`, `hash.star`,
`hmac.star`, `http.star`, `html.star`, `humanize.star`, `math.star`, `re.star`,
`sunrise.star`, `time.star`, `random.star`, `qrcode.star`, `assert.star`, plus
anything a caller-supplied `WithModuleLoader` resolves.

## Goals / Non-goals

**Goals**
- A file named after a built-in module can load that built-in module.
- The three bundled examples render again.
- No change to how any currently-working app resolves its `load()` calls (contract C2).

**Non-goals**
- Changing precedence in general (local files continue to win when there is no cycle).
- Rendering a *single file* that depends on sibling data files (e.g.
  `pixlet render examples/bitcoin/bitcoin.star` failing on `icon.png`). That is
  expected; render the directory instead. Possibly a docs note, not this spec.

## Resolution behavior: current vs proposed

| # | Situation | Today | Proposed |
|---|-----------|-------|----------|
| 1 | `X.star` loads `"X.star"`, `X` is a built-in | ❌ circular-dependency error | ✅ built-in module |
| 2 | `A.star` loads `"X.star"`; a local `X.star` helper exists; `X` is a built-in | Local helper | **Unchanged**: local helper. `pixlet lint`/`check` warns that it shadows a built-in |
| 3 | `A.star` loads `"X.star"`; local file exists; `X` not a built-in | Local file | **Unchanged** |
| 4 | Local cycle `A → B → A`, neither a built-in | Circular-dependency error | **Unchanged**: same error text |
| 5 | `X.star` loads `"util.star"`, which loads `"X.star"`; `X` is a built-in | ❌ circular-dependency error | ✅ built-in module (the local candidate would be a cycle) |
| 6 | Root files load alphabetically; `humanize.star` has loaded, then `zz.star` loads `"humanize.star"` | Gets the app file's globals | **Unchanged** (see Open questions) |

## Requirements

- **R1** In cases 1 and 5, `load()` MUST resolve to the built-in module (`loadModule`, including any `WithModuleLoader` loader).
- **R2** In cases 2, 3 and 6, resolution MUST be unchanged (C2).
- **R3** In case 4 (and in case 1/5 when *no* built-in or custom module exists for the name), the existing error `circular dependency detected: a -> b -> a` MUST be returned with unchanged text.
- **R4** The fix MUST NOT require a hard-coded list of built-in names in the loader. Resolution falls back to `loadModule()`, so custom loaders keep working.
- **R5** `pixlet lint` / `pixlet check` SHOULD warn when an app's `.star` file name equals a built-in module name. Case 2 is legal but surprising, and case 1 is now legal but worth a note.
- **R6** `examples/humanize/`, `examples/qrcode/`, and `examples/sunrise/` MUST render (as directories).
- **R7** Rendering output for every currently-working example MUST be byte-identical before and after (C1).

## Design

In `thread.Load` inside `ensureLoaded`, keep the local-first check, but detect
the cycle **before** recursing:

```go
modulePath := path.Clean(module)
if _, err := fs.Stat(fsys, modulePath); err == nil {
    if slices.Contains(currentlyLoading, modulePath) {
        // Loading this local file would be a cycle. If a built-in (or
        // custom-loader) module has this name, the author almost certainly
        // meant that, e.g. humanize.star loading "humanize.star".
        if mod, berr := a.loadModule(thread, module); berr == nil {
            return mod, nil
        }
        // No such module: report the cycle exactly as before (R3).
    }
    if err := a.ensureLoaded(fsys, modulePath, currentlyLoading...); err != nil {
        return nil, err
    }
    ...
}
return a.loadModule(thread, module)
```

- Only the "would cycle" path changes, so cases 2, 3 and 6 never reach the new code (R2).
- `loadModule` already tries the custom loader and then the built-in switch, so no name list is needed (R4).
- The error path is the existing `ensureLoaded` call, so the message is unchanged (R3).

**Lint warning (R5):** factor the built-in names out of the `loadModule` switch
into a single exported set (e.g. `runtime.BuiltinModules()`), so the lint check
and the loader can't drift apart. The lint rule then flags any root or nested
`.star` file whose cleaned path is in that set.

### Alternative considered: built-ins always win

Reserve built-in names: resolve to the built-in first, local file second. This is
simpler to explain, but it **changes case 2**. An app with a local `math.star`
helper loaded as `"math.star"` would silently get the built-in instead and break
(C2). Rejected for now. Could be revisited behind a deprecation period once the
lint warning (R5) has been in place for a while.

## Compatibility

- C2 (module behavior): protected by R2/R3; only cases that error today change.
- C1 (render output): R7, verified by example render hashes.
- C3 (CLI): unchanged, aside from a new lint warning.

## Acceptance criteria

- [ ] New tests in `runtime/applet_test.go` (MapFS-based, offline):
  - case 1: an app file `humanize.star` that loads `"humanize.star"` runs and calls `humanize.ordinal(1)`;
  - case 5: a two-file app where the built-in-named file loads a helper that loads the built-in name;
  - case 2: a local `math.star` helper loaded as `"math.star"` still returns the helper's globals;
  - case 4: `a.star ↔ b.star` still returns the exact cycle error;
  - case 1 with a name that has no built-in or custom module (e.g. `foo.star` loading `"foo.star"`) still returns the cycle error;
  - custom loader: a `WithModuleLoader` module named like the file resolves via the fallback.
- [ ] `pixlet render examples/{humanize,qrcode,sunrise}/` succeed.
- [ ] Every other example renders byte-identical output to `main` (hash comparison).
- [ ] `pixlet lint` warns on a file named `time.star`, and on nothing else in `examples/`.
- [ ] `go test ./...` passes offline; `go vet` adds no new findings.

## Tasks

1. Capture current render hashes for all examples (directories) on `main`.
2. Write the failing tests for cases 1 and 5, plus the guard tests for 2, 4, custom loader, and no-builtin cycle.
3. Implement the loader change.
4. Extract the built-in module name set; add the lint/check warning.
5. Re-render examples; compare hashes; confirm the three fixed examples render.
6. Optionally add an `examples` render test (offline) so this class of regression is caught in CI.

## Open questions

1. **Case 6 (load-order dependence):** after a root file named like a built-in
   has loaded, a *different* file loading that name gets the app file's globals,
   not the built-in. Keeping this unchanged is the compatible choice, but it is
   surprising. Should the lint warning (R5) also flag a load of a built-in
   name that resolves to a local file? (Proposed: yes, as a warning.)
2. Should the fix be released as a patch (`v0.42.1`), since it repairs a
   regression, or ride the next minor (`v0.43.0`) with spec 006? (Proposed:
   patch, as it is small and self-contained.)
