# Release Checklist

Use this checklist before visibility changes, tagged releases, or major example
expansions.

## Local Gates

```bash
make ci
gitleaks detect --source . --redact
```

`make ci` runs tests, vet, build, smoke commands, fixture/output boundary
checks, gitleaks when available, and actionlint when available.

## Review

- Confirm every fixture is synthetic and marked with `synthetic: true`.
- Confirm sample output contains no account names, emails, local paths, real
  companies, recruiter messages, tenant data, or application records.
- Confirm docs describe excluded connectors without adding implementation hooks
  for Gmail, LinkedIn, Simplify, calendar, browser automation, OAuth, or cookies.
- Review full Git history with gitleaks before changing visibility.
- Keep push, pull request, and manual CI triggers enabled.

## GitHub Visibility

The repository may be public only when local gates, full-history scan, CI, and
manual public-boundary review all pass.

## Latest Verification Snapshot

Checked on 2026-05-27 20:29 PDT:

- Local `make ci`: passed from this release-prep working tree.
- Full-history `gitleaks detect --source . --redact`: passed, scanning the
  complete 5-commit local Git history with no leaks found.
- Boundary scan from `scripts/check-public-boundary.sh`: passed.
- Fixture review: every committed fixture is synthetic and marked with
  `synthetic: true`.
- Output review: smoke output contains no account names, emails, local paths,
  real companies, recruiter messages, tenant data, or application records.
- GitHub visibility: approved to become public after this release-prep commit
  and CI pass.
