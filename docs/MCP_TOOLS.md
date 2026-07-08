# MCP-Style Tool Surface

`open-career-mcp` ships a small, stable MCP-style tool manifest plus an
in-process call adapter. Every tool reads only committed synthetic fixtures
under `fixtures/` and returns a dry-run review packet — none of them touch
Gmail, LinkedIn, a browser, a calendar, or any account connector.

## Tools

| Tool | Description | Required params |
| --- | --- | --- |
| `open_career_list_opportunities` | List synthetic career opportunities from local fixtures. | none |
| `open_career_tailor_resume` | Build a dry-run resume tailoring plan for a synthetic opportunity. | `opportunity_id` |
| `open_career_interview_prep` | Build a synthetic interview prep packet for a synthetic opportunity. | `opportunity_id` |

## Manifest

```bash
go run ./cmd/open-career-mcp mcp manifest
```

Returns the tool list above as JSON, including each tool's JSON-schema
`InputSchema` (see `internal/career/mcp.go`).

## Calling A Tool

```bash
go run ./cmd/open-career-mcp mcp call open_career_list_opportunities
go run ./cmd/open-career-mcp mcp call open_career_tailor_resume --param opportunity_id=syn-001
go run ./cmd/open-career-mcp mcp call open_career_interview_prep --param opportunity_id=syn-001
```

Each call loads `fixtures/opportunities.json` and `fixtures/resume.json`,
resolves the requested `opportunity_id`, and prints the same JSON a CLI
subcommand would print. There is no separate network transport in this
sample — `mcp call` is an in-process adapter so the tool surface can be
exercised without standing up a server.

## Adding A Tool

1. Add a synthetic-only fixture case if the new tool needs new data (`fixtures/`
   entries must keep `"synthetic": true` and pass
   `scripts/check-public-boundary.sh`).
2. Add domain logic in `internal/career`.
3. Register the tool in `internal/career/mcp.go` (`Manifest()`), keeping the
   `open_career_` name prefix (enforced by
   `internal/career/fixtures_test.go`).
4. Wire the `mcp call` case in `cmd/open-career-mcp/main.go`.
5. Add an example to `docs/EXAMPLES.md` and re-run `make ci`.

See `docs/ARCHITECTURE.md` for how these pieces fit together and
`PUBLIC_BOUNDARY.md` for what may never appear in a fixture or tool output.
