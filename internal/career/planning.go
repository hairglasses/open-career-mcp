package career

import (
	"fmt"
	"slices"
	"strings"
)

func BuildTailoringPlan(resume Resume, opportunity Opportunity) TailoringPlan {
	matched := matchSkills(resume.Skills, append(opportunity.Tags, opportunity.Requirements...))
	gaps := missingRequirements(opportunity.Requirements, matched)
	bullets := []string{
		fmt.Sprintf("Position %s as a %s for %s using only synthetic portfolio examples.", resume.Candidate, resume.Headline, opportunity.Role),
		fmt.Sprintf("Lead with %s when discussing %s.", projectNames(resume.Projects), opportunity.Company),
		"Keep all outreach, submission, and account-connector steps out of scope for this public sample.",
	}

	return TailoringPlan{
		OpportunityID: opportunity.ID,
		Company:       opportunity.Company,
		Role:          opportunity.Role,
		Summary:       fmt.Sprintf("%s is a synthetic match for %s because the role emphasizes %s.", resume.Candidate, opportunity.Role, strings.Join(matched, ", ")),
		MatchedSkills: matched,
		Gaps:          gaps,
		Bullets:       bullets,
		ApprovalNote:  "Dry run only. This sample never sends messages, submits applications, or touches live accounts.",
	}
}

func BuildInterviewPacket(resume Resume, opportunity Opportunity) InterviewPacket {
	talkingPoints := []string{
		fmt.Sprintf("Explain how %s maps to %s.", resume.Headline, opportunity.Role),
		"Describe safe execution boundaries before describing automation depth.",
		"Use synthetic examples and avoid private workspace or account-specific details.",
	}
	for _, signal := range opportunity.Signals {
		talkingPoints = append(talkingPoints, signal)
	}

	return InterviewPacket{
		OpportunityID: opportunity.ID,
		Company:       opportunity.Company,
		Role:          opportunity.Role,
		TalkingPoints: talkingPoints,
		Questions: []string{
			"Which workflow boundaries must remain dry-run only?",
			"How should a public sample prove value without live account connectors?",
			"What reliability signal would make this workflow trustworthy in production?",
		},
		BoundaryNote: "Use synthetic company, role, resume, and interview data only.",
	}
}

func matchSkills(skills []string, needs []string) []string {
	seen := map[string]bool{}
	var matched []string
	for _, skill := range skills {
		skillKey := normalize(skill)
		for _, need := range needs {
			if strings.Contains(normalize(need), skillKey) || strings.Contains(skillKey, normalize(need)) {
				if !seen[skill] {
					seen[skill] = true
					matched = append(matched, skill)
				}
			}
		}
	}
	if len(matched) == 0 && len(skills) > 0 {
		matched = append(matched, skills[0])
	}
	slices.Sort(matched)
	return matched
}

func missingRequirements(requirements, matched []string) []string {
	var gaps []string
	for _, requirement := range requirements {
		found := false
		for _, skill := range matched {
			if strings.Contains(normalize(requirement), normalize(skill)) || strings.Contains(normalize(skill), normalize(requirement)) {
				found = true
				break
			}
		}
		if !found {
			gaps = append(gaps, requirement)
		}
	}
	if len(gaps) == 0 {
		return []string{"No synthetic gap identified; validate against the real role before using this pattern."}
	}
	return gaps
}

func projectNames(projects []Project) string {
	var names []string
	for _, project := range projects {
		names = append(names, project.Name)
	}
	if len(names) == 0 {
		return "synthetic projects"
	}
	return strings.Join(names, " and ")
}

func normalize(value string) string {
	value = strings.ToLower(value)
	value = strings.ReplaceAll(value, "-", " ")
	return strings.TrimSpace(value)
}
