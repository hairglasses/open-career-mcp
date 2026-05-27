# Contributing

Keep this repository synthetic, portable, and reviewable.

- Do not add real resumes, recruiter messages, application records, account
  identifiers, credentials, tokens, cookies, host-specific paths, or user data.
- Use synthetic values in fixtures, docs, examples, tests, and command output.
- Add or update tests for new command behavior.
- Run `go test ./...` and `gitleaks detect --source . --no-git --redact`
  before opening a pull request.
- Keep live send, live submit, Gmail, LinkedIn, Simplify, calendar, and browser
  automation connectors out of this public sample.
