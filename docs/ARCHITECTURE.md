# Architecture

`open-career-mcp` is a compact, synthetic workflow sample. It is intentionally
small enough for a reviewer to inspect end to end.

## Data Flow

```text
synthetic fixtures
  -> opportunity/resume/interview domain functions
  -> CLI commands
  -> MCP-style manifest and in-process calls
  -> dry-run JSON or text output
```

## Components

| Area | Purpose |
| --- | --- |
| `fixtures/` | Synthetic opportunity, resume, and interview data. |
| `internal/career` | Domain logic for listing, tailoring, and prep packets. |
| `cmd/open-career-mcp` | CLI routing and MCP-style local adapter. |
| `scripts/check-public-boundary.sh` | Fixture and output boundary checks. |

## Safety Model

- Inputs are committed synthetic fixtures or explicit command parameters.
- Outputs are dry-run review packets.
- The sample does not include account connectors, browser automation, live send,
  live submit, calendar actions, or external mutations.
- `make ci` validates tests, vet, smoke commands, public-boundary checks, and a
  temporary build.
