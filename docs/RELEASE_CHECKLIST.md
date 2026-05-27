# Release Checklist

`open-career-mcp` stays private until this checklist is complete.

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
- Enable push and pull request CI triggers before public release.

## GitHub Visibility

The repository is currently private. Do not switch to public until the local
gates, full-history scan, and manual public-boundary review all pass.
