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
4. See [docs/MCP_TOOLS.md](docs/MCP_TOOLS.md) for the full MCP-style tool
   manifest (names, descriptions, required params).
5. Use [docs/PORTFOLIO_PROOF.md](docs/PORTFOLIO_PROOF.md) for the architecture
   diagram, walkthrough plan, tradeoffs, and interview prompts.
6. Check the public/private line in [PUBLIC_BOUNDARY.md](PUBLIC_BOUNDARY.md).

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

### Demo Transcript

```text
$ go run ./cmd/open-career-mcp opportunities list
syn-001 | Northstar Tools | Staff Engineer, Agent Infrastructure | remote | research
syn-002 | Cedar Analytics | Principal Platform Engineer | hybrid | screen

$ go run ./cmd/open-career-mcp mcp call open_career_tailor_resume --param opportunity_id=syn-001
{
  "opportunity_id": "syn-001",
  "company": "Northstar Tools",
  "role": "Staff Engineer, Agent Infrastructure",
  "summary": "Avery Example is a synthetic match for Staff Engineer, Agent Infrastructure because the role emphasizes Distributed systems, Go.",
  "matched_skills": [
    "Distributed systems",
    "Go"
  ],
  "gaps": [
    "Policy and approval gates",
    "Operational telemetry"
  ],
  "bullets": [
    "Position Avery Example as a Staff platform engineer focused on agent infrastructure for Staff Engineer, Agent Infrastructure using only synthetic portfolio examples.",
    "Lead with Agent Control Plane and Fleet Baseline Guard when discussing Northstar Tools.",
    "Keep all outreach, submission, and account-connector steps out of scope for this public sample."
  ],
  "approval_note": "Dry run only. This sample never sends messages, submits applications, or touches live accounts."
}
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
