# Examples

All examples use synthetic fixtures from this repository.

## List Opportunities

```bash
go run ./cmd/open-career-mcp opportunities list
```

Expected output:

```text
syn-001 | Northstar Tools | Staff Engineer, Agent Infrastructure | remote | research
syn-002 | Cedar Analytics | Principal Platform Engineer | hybrid | screen
```

## Tailor A Resume

```bash
go run ./cmd/open-career-mcp resume tailor --opportunity syn-001
```

Expected output shape:

```text
# Tailoring Plan: Staff Engineer, Agent Infrastructure

Matched skills: Distributed systems, Go
Gaps: Policy and approval gates, Operational telemetry

- Position Avery Example as a Staff platform engineer focused on agent infrastructure.
- Keep all outreach, submission, and account-connector steps out of scope.

Dry run only.
```

## Build Interview Prep

```bash
go run ./cmd/open-career-mcp interview prep --opportunity syn-001
```

Expected output shape:

```text
# Interview Prep: Staff Engineer, Agent Infrastructure at Northstar Tools

Talking points:
- Explain how the synthetic background maps to the synthetic role.
- Describe safe execution boundaries before describing automation depth.

Practice questions:
- Which workflow boundaries must remain dry-run only?
```

## MCP-Style Manifest

```bash
go run ./cmd/open-career-mcp mcp manifest
```

## MCP-Style Tool Call

```bash
go run ./cmd/open-career-mcp mcp call open_career_tailor_resume --param opportunity_id=syn-001
```

The MCP-style call returns the same dry-run tailoring plan as the CLI command.
No connector, browser, message, calendar, or application state is touched.
