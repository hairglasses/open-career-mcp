# Public Boundary

`open-career-mcp` is a planned synthetic-data career workflow MCP sample. It is
private staging until the public boundary is proven by tests and scans.

## Included

- Synthetic opportunity tracking.
- Synthetic resume tailoring examples.
- Synthetic interview prep packet generation.
- MCP-style tool manifest and local demo calls.
- Dry-run approval-boundary examples.

## Excluded

- Real resumes, recruiter messages, Gmail, LinkedIn, Simplify, browser
  automation, OAuth files, cookies, calendar data, tenant databases, saved jobs,
  application history, and personal profile exports.
- Local workstation paths, account names, emails, hostnames, private repo
  manifests, secrets, or generated private artifacts.
- Live send, live submit, calendar decline, or any account-mutating connector.

## Release Gate

Before public visibility:

```bash
go test ./...
gitleaks detect --source . --no-git --redact
```

The first public release also needs CI and a full history review.
