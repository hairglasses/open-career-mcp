package career

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func LoadFixtures(dir string) (Resume, []Opportunity, error) {
	resumePath := filepath.Join(dir, "resume.json")
	opportunityPath := filepath.Join(dir, "opportunities.json")

	resumeBytes, err := os.ReadFile(resumePath)
	if err != nil {
		return Resume{}, nil, fmt.Errorf("read resume fixture: %w", err)
	}

	var resume Resume
	if err := json.Unmarshal(resumeBytes, &resume); err != nil {
		return Resume{}, nil, fmt.Errorf("parse resume fixture: %w", err)
	}
	if !resume.Synthetic {
		return Resume{}, nil, fmt.Errorf("resume fixture must set synthetic=true")
	}

	opportunityBytes, err := os.ReadFile(opportunityPath)
	if err != nil {
		return Resume{}, nil, fmt.Errorf("read opportunities fixture: %w", err)
	}

	var opportunities []Opportunity
	if err := json.Unmarshal(opportunityBytes, &opportunities); err != nil {
		return Resume{}, nil, fmt.Errorf("parse opportunities fixture: %w", err)
	}
	for _, opportunity := range opportunities {
		if !opportunity.Synthetic {
			return Resume{}, nil, fmt.Errorf("opportunity %q must set synthetic=true", opportunity.ID)
		}
	}
	sort.Slice(opportunities, func(i, j int) bool {
		return opportunities[i].ID < opportunities[j].ID
	})

	return resume, opportunities, nil
}

func FindOpportunity(opportunities []Opportunity, id string) (Opportunity, error) {
	for _, opportunity := range opportunities {
		if opportunity.ID == id {
			return opportunity, nil
		}
	}
	return Opportunity{}, fmt.Errorf("opportunity %q not found", id)
}
