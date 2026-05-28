# Public Boundary

`open-career-mcp` is a synthetic-data career workflow MCP sample. Its public
boundary is deliberately narrow: fixtures, examples, tests, and command output
must remain synthetic, local-first, and dry-run only.

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

Before each release or major example expansion:

```bash
make ci
gitleaks detect --source . --redact
```

The first public release also required CI and a full-history review.
