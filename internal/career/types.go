package career

type Resume struct {
	Synthetic bool      `json:"synthetic"`
	Candidate string    `json:"candidate"`
	Headline  string    `json:"headline"`
	Skills    []string  `json:"skills"`
	Projects  []Project `json:"projects"`
}

type Project struct {
	Name    string   `json:"name"`
	Summary string   `json:"summary"`
	Skills  []string `json:"skills"`
}

type Opportunity struct {
	Synthetic    bool     `json:"synthetic"`
	ID           string   `json:"id"`
	Company      string   `json:"company"`
	Role         string   `json:"role"`
	Location     string   `json:"location"`
	WorkMode     string   `json:"work_mode"`
	Stage        string   `json:"stage"`
	Tags         []string `json:"tags"`
	Description  string   `json:"description"`
	Requirements []string `json:"requirements"`
	Signals      []string `json:"signals"`
}

type TailoringPlan struct {
	OpportunityID string   `json:"opportunity_id"`
	Company       string   `json:"company"`
	Role          string   `json:"role"`
	Summary       string   `json:"summary"`
	MatchedSkills []string `json:"matched_skills"`
	Gaps          []string `json:"gaps"`
	Bullets       []string `json:"bullets"`
	ApprovalNote  string   `json:"approval_note"`
}

type InterviewPacket struct {
	OpportunityID string   `json:"opportunity_id"`
	Company       string   `json:"company"`
	Role          string   `json:"role"`
	TalkingPoints []string `json:"talking_points"`
	Questions     []string `json:"questions"`
	BoundaryNote  string   `json:"boundary_note"`
}

type ToolManifest struct {
	Tools []Tool `json:"tools"`
}

type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"input_schema"`
}
