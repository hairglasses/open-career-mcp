# Getting Started

## Requirements

- Go 1.26 or newer
- `make`
- Optional: `gitleaks` and `actionlint` for local parity with release checks

## Five-Minute Review

```bash
git clone https://github.com/hairglasses/open-career-mcp.git
cd open-career-mcp
make ci
go run ./cmd/open-career-mcp opportunities list
go run ./cmd/open-career-mcp resume tailor --opportunity syn-001
go run ./cmd/open-career-mcp interview prep --opportunity syn-001
go run ./cmd/open-career-mcp mcp manifest
```

All commands are local and deterministic. They read only committed synthetic
fixtures and write no account state.

## Review Path

1. Start with `README.md` for the public proof summary.
2. Inspect `fixtures/` to confirm the examples are synthetic.
3. Read `docs/EXAMPLES.md` for expected output shapes.
4. Read `PUBLIC_BOUNDARY.md` before adding fixtures, examples, or tools.
5. Run `make ci` before opening a pull request or publishing a release.
