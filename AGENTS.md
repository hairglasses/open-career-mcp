# open-career-mcp - Agent Instructions

> Canonical instructions: AGENTS.md

Synthetic-data career workflow MCP sample. This repository is private staging
until every fixture, example, and commit is safe for public release.

## Public Boundary

- Use only synthetic people, companies, jobs, resumes, messages, and dates.
- Do not copy live `jobb` data, tenant files, OAuth state, browser state,
  recruiter messages, application records, or personal resume artifacts.
- Do not add Gmail, LinkedIn, Simplify, calendar, or browser automation
  connectors to the public sample.
- Keep examples generic and local-first.

## Verification

- Run `gitleaks detect --source . --no-git --redact` before any visibility
  change.
- Keep README, examples, and fixtures aligned with `PUBLIC_BOUNDARY.md`.
