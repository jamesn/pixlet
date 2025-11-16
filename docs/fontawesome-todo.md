# Font Awesome: Selective Registration TODO

This file records the current plan and progress for replacing full Font Awesome pack registration
(`library.add(fas, fab)`) with selective, per-icon registration to reduce bundle size.

Last updated: 2025-11-16

## Plan

- [ ] **Enumerate used FA icons**: Collect Font Awesome icon names used at runtime and in schema/tests (scan `src/features/schema`, `examples/`, and tests) and produce a canonical list of solid/brand icons to import.

- [ ] **Create icons registration module**: Add `src/features/icons/registerFontAwesome.js` that imports only the required icons from `@fortawesome/free-solid-svg-icons` and `@fortawesome/free-brands-svg-icons` and calls `library.add(...)`. Replace `library.add(fas, fab)` in `src/features/theme/DevToolsTheme.jsx` to import this module instead.

- [ ] **Update FieldIcon resolution**: Ensure `src/features/schema/FieldIcon.jsx` resolves to registered icons. Either keep `findIconDefinition` (works if icons are added to library) or map icon string to imported icon objects for direct usage if needed.

- [ ] **Build and measure**: Run `npm run build` and compare vendor bundle size before/after. Inspect `dist` outputs and source maps to confirm FA packs are no longer fully bundled.

- [ ] **Iterate / fallback plan**: If bundle size still large, consider migrating some icons to MUI or inline SVGs and/or lazy-load Font Awesome icons. Prepare PR and tests.

## Notes

- Current status: TODO list created and saved. No code changes have been applied yet.
- Suggested next action: I can start step 1 (enumerate icons) and create `src/features/icons/registerFontAwesome.js` once we have the list. Reply "Yes" to continue.

## References

- `src/features/theme/DevToolsTheme.jsx` — current location of `library.add(fas, fab)` (needs change)
- `src/features/schema/FieldIcon.jsx` — resolves icon names and currently uses `findIconDefinition` and `FontAwesomeIcon` component


