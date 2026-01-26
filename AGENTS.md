# Repository Guidelines

## Project Structure & Module Organization
- `builder/`: public builders and PDF generation entry points.
- `cmd/generate-samples/`: helper CLI to generate local PDF outputs under `samples/`.
- `core/`: shared document models and calculation helpers.
- `docs/`: contributor-facing notes (see `docs/layouts.md` for built-in layouts).
- `i18n/`: localization loader; translations live in `i18n/locales/*.toml` (supported: `en`, `ja`, `zh_cn`, `zh_tw`).
- `samples/`: YAML inputs used by examples/tests.
- `samples/`: generated sample PDFs (ignored by `.gitignore`).
- `fonts/`: bundled fonts for CJK demos (referenced by config).
- `sample-*.pdf`: legacy generated artifacts in repo root (also ignored).

## Build, Test, and Development Commands
- `go test ./...`: run all tests (note: this generates PDFs in the repo root).
- `go test ./builder -run TestGenerateInvoice -v`: run a single generator test.
- `go run ./cmd/generate-samples`: generate PDFs for all built-in layouts into `samples/`.
- `go vet ./...`: extra static checks beyond compilation.
- `gofmt -w .`: format Go code.
- `go mod tidy`: clean up dependencies after changing imports.

## Coding Style & Naming Conventions
- Always run `gofmt` on Go changes; keep imports gofmt-sorted.
- Exported identifiers use `PascalCase`; unexported use `camelCase`.
- Keep the public API in `builder/` stable; prefer additive options over breaking signatures.
- Layout selection uses stable names (see `builder/layouts.go`), e.g. `classic`, `modern`, `spotlight`.

## Testing Guidelines
- Tests live next to code as `*_test.go` (currently `builder/builder_test.go`).
- Avoid committing generated PDFs (`*.pdf` is ignored). If you need fixtures, prefer small text/YAML under `samples/`.

## Commit & Pull Request Guidelines
- Match the commit style in history: `feat: ...`, `fix: ...`, `refactor: ...`, `chore: ...` (optionally add a scope).
- PRs should include: a short description, how you validated (`go test ./...`), and (if PDF output changes) a screenshot or generated `sample-*.pdf` attached as an artifact, not committed.
