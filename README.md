# open-career-mcp

[![ci](https://github.com/hairglasses/open-career-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/hairglasses/open-career-mcp/actions/workflows/ci.yml)

Synthetic-data career workflow MCP sample for resume tailoring, opportunity
tracking, and interview preparation.

This repository is intentionally separate from private career automation
systems. It contains only public-safe, synthetic fixtures and local dry-run
examples.

## What Works Now

- Track synthetic opportunities and review state.
- Tailor synthetic resume records against synthetic job descriptions.
- Generate interview prep packets from fake company and role data.
- Expose the workflow through small CLI and MCP-style tool examples.
- Demonstrate dry-run approval boundaries without live account connectors.

## Why It Exists

`open-career-mcp` is a compact reference implementation for career-workflow
automation patterns that are safe to inspect publicly: fixture validation,
dry-run planning, MCP-style manifests, in-process tool calls, and explicit
connector boundaries.

## Usage

```bash
go run ./cmd/open-career-mcp opportunities list
go run ./cmd/open-career-mcp resume tailor --opportunity syn-001
go run ./cmd/open-career-mcp interview prep --opportunity syn-001
go run ./cmd/open-career-mcp mcp manifest
```

See [docs/EXAMPLES.md](docs/EXAMPLES.md) for more commands.

## Not In Scope

- Gmail, LinkedIn, Simplify, calendar, browser automation, OAuth, cookies, or
  live application submission.
- Personal resumes, recruiter threads, real companies, tenant databases, saved
  job IDs, or application history.

See [PUBLIC_BOUNDARY.md](PUBLIC_BOUNDARY.md) before adding examples or fixtures.

## Verification

```bash
make ci
gitleaks detect --source . --redact
```
