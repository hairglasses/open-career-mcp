.PHONY: build test vet smoke public-boundary gitleaks actionlint ci

GO ?= env GOWORK=off go

build:
	tmp_dir="$$(mktemp -d)"; \
	trap 'rm -rf "$$tmp_dir"' EXIT; \
	$(GO) build -o "$$tmp_dir/open-career-mcp" ./cmd/open-career-mcp

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

smoke:
	$(GO) run ./cmd/open-career-mcp opportunities list
	$(GO) run ./cmd/open-career-mcp resume tailor --opportunity syn-001
	$(GO) run ./cmd/open-career-mcp interview prep --opportunity syn-001
	$(GO) run ./cmd/open-career-mcp mcp manifest
	$(GO) run ./cmd/open-career-mcp mcp call open_career_tailor_resume --param opportunity_id=syn-001

public-boundary:
	GO_CMD="$(GO)" scripts/check-public-boundary.sh

gitleaks:
	gitleaks detect --source . --no-git --redact

actionlint:
	actionlint .github/workflows/*.yml

ci: test vet build smoke public-boundary
