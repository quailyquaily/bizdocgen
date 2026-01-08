# Repository Guidelines

## Project Structure & Module Organization
- `builder/`: public builders and PDF generation entry points.
- `core/`: shared document models and calculation helpers.
- `i18n/`: localization loader; translations live in `i18n/locales/*.toml`.
- `sample-params/`: YAML inputs used by examples/tests.
- `fonts/`: bundled fonts for CJK demos (referenced by config).
- `sample-*.pdf`: generated artifacts (ignored by `.gitignore`).

## Build, Test, and Development Commands
- `go test ./...`: run all tests (note: this generates PDFs in the repo root).
- `go test ./builder -run TestGenerateInvoice -v`: run a single generator test.
- `go vet ./...`: extra static checks beyond compilation.
- `gofmt -w .`: format Go code.
- `go mod tidy`: clean up dependencies after changing imports.

## Coding Style & Naming Conventions
- Always run `gofmt` on Go changes; keep imports gofmt-sorted.
- Exported identifiers use `PascalCase`; unexported use `camelCase`.
- Keep the public API in `builder/` stable; prefer additive options over breaking signatures.

## Testing Guidelines
- Tests live next to code as `*_test.go` (currently `builder/builder_test.go`).
- Avoid committing generated PDFs (`*.pdf` is ignored). If you need fixtures, prefer small text/YAML under `sample-params/`.

## Commit & Pull Request Guidelines
- Match the commit style in history: `feat: ...`, `fix: ...`, `refactor: ...`, `chore: ...` (optionally add a scope).
- PRs should include: a short description, how you validated (`go test ./...`), and (if PDF output changes) a screenshot or generated `sample-*.pdf` attached as an artifact, not committed.
