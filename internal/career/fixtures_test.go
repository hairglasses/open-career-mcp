package career

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFixtures(t *testing.T) {
	resume, opportunities, err := LoadFixtures("../../fixtures")
	if err != nil {
		t.Fatalf("LoadFixtures() error = %v", err)
	}
	if !resume.Synthetic {
		t.Fatal("resume fixture is not synthetic")
	}
	if len(opportunities) != 2 {
		t.Fatalf("expected 2 opportunities, got %d", len(opportunities))
	}
	for _, opportunity := range opportunities {
		if !opportunity.Synthetic {
			t.Fatalf("opportunity %s is not synthetic", opportunity.ID)
		}
	}
}

func TestBuildTailoringPlan(t *testing.T) {
	resume, opportunities, err := LoadFixtures("../../fixtures")
	if err != nil {
		t.Fatalf("LoadFixtures() error = %v", err)
	}
	opportunity, err := FindOpportunity(opportunities, "syn-001")
	if err != nil {
		t.Fatalf("FindOpportunity() error = %v", err)
	}
	plan := BuildTailoringPlan(resume, opportunity)
	if plan.OpportunityID != "syn-001" {
		t.Fatalf("unexpected opportunity id %q", plan.OpportunityID)
	}
	if len(plan.MatchedSkills) == 0 {
		t.Fatal("expected matched skills")
	}
	if !strings.Contains(plan.ApprovalNote, "Dry run only") {
		t.Fatalf("approval note does not preserve dry-run boundary: %q", plan.ApprovalNote)
	}
}

func TestManifest(t *testing.T) {
	manifest := Manifest()
	if len(manifest.Tools) != 3 {
		t.Fatalf("expected 3 tools, got %d", len(manifest.Tools))
	}
	for _, tool := range manifest.Tools {
		if !strings.HasPrefix(tool.Name, "open_career_") {
			t.Fatalf("tool name %q does not use open_career_ prefix", tool.Name)
		}
	}
}

func TestFixturesStaySynthetic(t *testing.T) {
	fixtureDir := "../../fixtures"
	entries, err := os.ReadDir(fixtureDir)
	if err != nil {
		t.Fatalf("ReadDir() error = %v", err)
	}
	forbidden := []string{
		"gmail",
		"linkedin",
		"oauth",
		"cookie",
		"tenant",
		"simplify",
		"calendar",
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(fixtureDir, entry.Name()))
		if err != nil {
			t.Fatalf("ReadFile(%s) error = %v", entry.Name(), err)
		}
		lower := strings.ToLower(string(data))
		if !strings.Contains(lower, `"synthetic": true`) {
			t.Fatalf("%s does not contain synthetic=true marker", entry.Name())
		}
		for _, word := range forbidden {
			if strings.Contains(lower, word) {
				t.Fatalf("%s contains forbidden fixture marker %q", entry.Name(), word)
			}
		}
	}
}
