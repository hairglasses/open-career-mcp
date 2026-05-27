.PHONY: build test vet smoke public-boundary gitleaks actionlint ci

build:
	tmp_dir="$$(mktemp -d)"; \
	trap 'rm -rf "$$tmp_dir"' EXIT; \
	go build -o "$$tmp_dir/open-career-mcp" ./cmd/open-career-mcp

test:
	go test ./...

vet:
	go vet ./...

smoke:
	go run ./cmd/open-career-mcp opportunities list
	go run ./cmd/open-career-mcp resume tailor --opportunity syn-001
	go run ./cmd/open-career-mcp interview prep --opportunity syn-001
	go run ./cmd/open-career-mcp mcp manifest
	go run ./cmd/open-career-mcp mcp call open_career_tailor_resume --param opportunity_id=syn-001

public-boundary:
	scripts/check-public-boundary.sh

gitleaks:
	gitleaks detect --source . --no-git --redact

actionlint:
	actionlint .github/workflows/*.yml

ci: test vet build smoke public-boundary
