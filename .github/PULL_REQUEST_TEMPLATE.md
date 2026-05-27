## Summary

## Public Boundary Checklist

- [ ] Uses only synthetic fixtures, examples, and output.
- [ ] Does not add live account connectors or account-mutating actions.
- [ ] Does not include credentials, OAuth state, cookies, local paths, personal
      resume data, recruiter messages, tenant data, or application records.
- [ ] Updates tests for command or fixture changes.
- [ ] `go test ./...` passes.
- [ ] `gitleaks detect --source . --no-git --redact` passes.
