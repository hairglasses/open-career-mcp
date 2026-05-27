package career

func Manifest() ToolManifest {
	return ToolManifest{Tools: []Tool{
		{
			Name:        "open_career_list_opportunities",
			Description: "List synthetic career opportunities from local fixtures.",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			Name:        "open_career_tailor_resume",
			Description: "Build a dry-run resume tailoring plan for a synthetic opportunity.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"opportunity_id": map[string]any{"type": "string"},
				},
				"required": []string{"opportunity_id"},
			},
		},
		{
			Name:        "open_career_interview_prep",
			Description: "Build a synthetic interview prep packet for a synthetic opportunity.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"opportunity_id": map[string]any{"type": "string"},
				},
				"required": []string{"opportunity_id"},
			},
		},
	}}
}
