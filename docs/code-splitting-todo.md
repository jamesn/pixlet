# Code-Splitting & Lazy-Loading TODO

Purpose: Improvements to reduce initial bundle size and defer heavy code to runtime.

Last updated: 2025-11-16

Tasks

- [ ] **Audit large entry points**: Use `webpack --profile --json` or `source-map-explorer` on `dist` to identify large modules.

- [ ] **Apply `React.lazy` + `Suspense`**: Wrap heavy components (e.g., field editors, large form widgets, complex UIs) with dynamic imports.

- [ ] **Split large vendor libraries**: Ensure libraries that are not needed on first load are dynamically imported. Example: load certain charting or editor libs only when their view is opened.

- [ ] **Selective Font Awesome registration**: See `docs/fontawesome-todo.md` — import and register only used FA icons instead of `library.add(fas, fab)`.

- [ ] **Measure effects**:

```bash
npm run build
# then inspect bundle sizes (dist/*)
# or use source-map-explorer
npx source-map-explorer dist/*.js
```

- [ ] **Edge-case testing**: Test navigation flows to ensure lazy-loaded bundles load correctly and show appropriate fallbacks (e.g., `Suspense` loading spinners).

Notes

- Aim for reducing the largest `vendor` chunk first (commonly caused by large icon packs, UI libraries, or charting libs).
- Keep user experience smooth by deferring non-critical code and showing placeholders while loading.
