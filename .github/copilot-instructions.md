# Copilot instructions for Pixlet

**Purpose**: Help AI coding agents become productive quickly in the Pixlet codebase by documenting the architecture, developer workflows, and repository-specific conventions.

**Big Picture**:
- **What Pixlet is**: a Go-based CLI and runtime that executes Starlark applets to render pixel graphics (WebP/GIF) for Tidbyt devices. See `README.md` for user-facing details.
- **Core responsibilities**: parse/execute Starlark apps (`runtime/`), provide UI primitives (`render/`), encode frames (`encode/`), and expose CLI commands in `cmd/`.
- **Why it's structured this way**: separation of concerns — `runtime` handles Starlark integration and app lifecycle, `render` is a declarative widget toolkit, and `encode`/`bundle` handle output/publishing.

**Key files & directories** (most important first):
- `cmd/` : CLI commands. `main.go` wires `cmd.RenderCmd`, `cmd.PushCmd`, etc.
- `runtime/` : Starlark runtime integration and app lifecycle.
- `render/` : Widgets and layout primitives (Text, Box, Animation, Root).
- `encode/` : WebP/GIF encoding and benchmarks.
- `examples/` : Real `.star` example applets (useful for reproducing bugs and tests).
- `docs/` : Tutorials and widget references — mirror of expected runtime behavior.
- `Makefile` : canonical developer commands (build, test, codegen targets).

**Build / Test / Dev workflows** (exact commands)
- Build the CLI binary: `make build` (runs `go build -o pixlet tidbyt.dev/pixlet`).
- Run the full test suite: `make test` -> `go test -v -cover ./...`.
- Run specific package tests: `go test ./render -v` or `go test ./runtime -run TestName -v`.
- Generate widgets (codegen step used by repo): `make widgets` (runs `go run runtime/gen/main.go`).
- Embed fonts: `make embedfonts` (runs `go run render/gen/embedfonts.go`).
- Lint/format (buildifier for Starlark files): `make format` / `make lint`.

**Common developer tasks & examples**
- Render an example locally: `pixlet render examples/hello_world/hello_world.star`.
- Serve a remote script: `curl https://.../hello_world.star | pixlet serve /dev/stdin` (README example).
- Push a generated WebP: `pixlet push <DEVICE_ID> examples/bitcoin/bitcoin.webp`.
- CLI wiring example: `main.go` registers commands with Cobra: `rootCmd.AddCommand(cmd.RenderCmd)` — look in `cmd/` for per-command behavior and flags.

**Conventions & patterns to follow (project-specific)**
- Starlark files use the `.star` extension and load `render.star`/`time.star` etc. Inspect `examples/` for idiomatic usage.
- Code generation is tracked as canonical: run `make widgets` after changing widget sources or templates.
- Tests often live next to implementations and frequently use `testutil.go` helpers (see `render/testutil.go`).
- Module path is `tidbyt.dev/pixlet` (see `go.mod`) — imports use that path across packages.

**Integration points & external dependencies**
- Starlark integration: `go.starlark.net` and `github.com/qri-io/starlib` provide core language and libs.
- Encoding: `github.com/tidbyt/go-libwebp` and internal `encode/` package.
- CLI framework: `github.com/spf13/cobra` (see `main.go` + `cmd/`).

**Investigation tips for AI agents**
- To understand rendering flow, trace from a `.star` example -> `runtime` evaluation -> `render` widget graph -> `encode` frame generation.
- Examine `examples/` when proposing UI or starlark changes — they are executable and help reproduce behavior.
- If changing widgets, run `make widgets` then `make test` to validate generated code + behavior.

**What not to assume**
- Don’t assume a web server is part of the core binary — serving is a convenience for development (`pixlet serve`), the core is the render/runtime pipeline.
- Don’t change generated files without updating the corresponding generator (`runtime/gen/` or `render/gen/`).

If anything here is unclear or you'd like more detail in a specific area (e.g., testing patterns, generator internals, or Starlark runtime primitives), tell me which area and I will expand or add examples from the codebase.
