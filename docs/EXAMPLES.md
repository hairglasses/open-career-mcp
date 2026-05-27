# Examples

All examples use synthetic fixtures from this repository.

## List Opportunities

```bash
go run ./cmd/open-career-mcp opportunities list
```

## Tailor A Resume

```bash
go run ./cmd/open-career-mcp resume tailor --opportunity syn-001
```

## Build Interview Prep

```bash
go run ./cmd/open-career-mcp interview prep --opportunity syn-001
```

## MCP-Style Manifest

```bash
go run ./cmd/open-career-mcp mcp manifest
```

## MCP-Style Tool Call

```bash
go run ./cmd/open-career-mcp mcp call open_career_tailor_resume --param opportunity_id=syn-001
```
