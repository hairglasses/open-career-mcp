package main

import "testing"

func TestRunMCPManifest(t *testing.T) {
	if err := run([]string{"mcp", "manifest"}); err != nil {
		t.Fatalf("run() error = %v", err)
	}
}

func TestRunTailorJSON(t *testing.T) {
	if err := run([]string{"resume", "tailor", "--fixtures", "../../fixtures", "--opportunity", "syn-001", "--json"}); err != nil {
		t.Fatalf("run() error = %v", err)
	}
}
