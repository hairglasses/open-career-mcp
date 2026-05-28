# open-career-mcp

[![ci](https://github.com/hairglasses/open-career-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/hairglasses/open-career-mcp/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Synthetic-data career workflow MCP sample for resume tailoring, opportunity
tracking, and interview preparation.

This repository is intentionally separate from private career automation
systems. It contains only public-safe, synthetic fixtures and local dry-run
examples.

## Start Here

For a quick review path:

1. Run the five-minute commands in [docs/GETTING_STARTED.md](docs/GETTING_STARTED.md).
2. Compare output shapes in [docs/EXAMPLES.md](docs/EXAMPLES.md).
3. Review the data flow in [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).
4. Check the public/private line in [PUBLIC_BOUNDARY.md](PUBLIC_BOUNDARY.md).

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
make ci
go run ./cmd/open-career-mcp opportunities list
go run ./cmd/open-career-mcp resume tailor --opportunity syn-001
go run ./cmd/open-career-mcp interview prep --opportunity syn-001
go run ./cmd/open-career-mcp mcp manifest
go run ./cmd/open-career-mcp mcp call open_career_tailor_resume --param opportunity_id=syn-001
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
gitleaks detect --source . --no-git --redact
```

`make ci` runs tests, vet, a temporary build, smoke commands, public-boundary
checks, and optional local `gitleaks` / `actionlint` checks when those tools are
installed.
