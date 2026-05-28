#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

go_cmd="${GO_CMD:-go}"
forbidden_fixture_pattern='gmail|linkedin|oauth|cookie|tenant|simplify|calendar|mitch|mitchell|/home/hg|hairglasses-studio|jobb|@[A-Za-z0-9._%+-]+\.[A-Za-z]{2,}'
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

echo "== fixture boundary =="
if rg --ignore-case --line-number "$forbidden_fixture_pattern" fixtures; then
  echo "fixtures contain private-boundary markers" >&2
  exit 1
fi

echo "== sample output boundary =="
$go_cmd run ./cmd/open-career-mcp opportunities list > "$tmp_dir/opportunities.txt"
$go_cmd run ./cmd/open-career-mcp resume tailor --opportunity syn-001 > "$tmp_dir/tailor.txt"
$go_cmd run ./cmd/open-career-mcp interview prep --opportunity syn-001 > "$tmp_dir/interview.txt"
$go_cmd run ./cmd/open-career-mcp mcp manifest > "$tmp_dir/manifest.json"
$go_cmd run ./cmd/open-career-mcp mcp call open_career_tailor_resume --param opportunity_id=syn-001 > "$tmp_dir/mcp-tailor.json"

if rg --ignore-case --line-number "$forbidden_fixture_pattern" "$tmp_dir"; then
  echo "sample output contains private-boundary markers" >&2
  exit 1
fi

echo "== gitleaks =="
if command -v gitleaks >/dev/null 2>&1; then
  gitleaks detect --source . --no-git --redact
else
  echo "gitleaks not installed; skipping local secret scan"
fi

echo "== actionlint =="
if command -v actionlint >/dev/null 2>&1; then
  actionlint .github/workflows/*.yml
else
  echo "actionlint not installed; skipping workflow lint"
fi

echo "public boundary checks passed"
