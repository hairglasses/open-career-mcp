# open-career-mcp - Agent Instructions

> Canonical instructions: AGENTS.md

Synthetic-data career workflow MCP sample. This repository is intended to be
public-safe: every fixture, example, and commit must preserve the synthetic
boundary.

## Public Boundary

- Use only synthetic people, companies, jobs, resumes, messages, and dates.
- Do not copy live `jobb` data, tenant files, OAuth state, browser state,
  recruiter messages, application records, or personal resume artifacts.
- Do not add Gmail, LinkedIn, Simplify, calendar, or browser automation
  connectors to the public sample.
- Keep examples generic and local-first.

## Verification

- Run `make ci` before committing meaningful changes.
- Run `gitleaks detect --source . --redact` before any tagged release or major
  public example expansion.
- Keep README, examples, and fixtures aligned with `PUBLIC_BOUNDARY.md`.

<!-- BEGIN GENERATED LANDING POLICY -->
Use `hg-land --repo <checkout>` (`--no-verify` while the fleet gate switch
`~/.config/hg/gates-disabled` or `HG_GATES_DISABLED=1` is active; `--local-gate <log>`
under `HG_GATES_DISABLED=0`, or `hg-gate-then-land --require-gate`) to land through a
pull request. Owned by `claude-harness/.agents/policies/landing.md`.

- While fleet gates are disabled (operator 2026-11-08), `hg-land` lands without a
  local gate log; when gated, bind the log to `head=<full SHA>` with final `rc=0`.
- Observe the repository's current remote requirements. Required checks and
  branch protections remain binding; hosted CI is not assumed required everywhere.
- Never bypass a failed gate, force-push, or replace unrelated work. Do not use
  direct pushes to main as a substitute for the configured landing route.
- Optoinvest repositories remain local-only: no remote publication.
- Prove landing with the tool receipt and exact remote main SHA. Source updates
  and runtime/deployment verification are separate claims.
<!-- END GENERATED LANDING POLICY -->
