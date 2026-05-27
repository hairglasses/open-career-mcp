package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hairglasses/open-career-mcp/internal/career"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return usage()
	}

	switch args[0] {
	case "opportunities":
		return opportunities(args[1:])
	case "resume":
		return resume(args[1:])
	case "interview":
		return interview(args[1:])
	case "mcp":
		return mcp(args[1:])
	case "help", "-h", "--help":
		return usage()
	default:
		return fmt.Errorf("unknown command %q\n\n%w", args[0], usage())
	}
}

func opportunities(args []string) error {
	if len(args) == 0 || args[0] != "list" {
		return errors.New("usage: open-career-mcp opportunities list [--fixtures fixtures] [--json]")
	}
	fs := flag.NewFlagSet("opportunities list", flag.ContinueOnError)
	fixtureDir := fs.String("fixtures", "fixtures", "fixture directory")
	jsonOut := fs.Bool("json", false, "print JSON")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	_, opportunities, err := career.LoadFixtures(*fixtureDir)
	if err != nil {
		return err
	}
	if *jsonOut {
		return printJSON(opportunities)
	}
	for _, opportunity := range opportunities {
		fmt.Printf("%s | %s | %s | %s | %s\n", opportunity.ID, opportunity.Company, opportunity.Role, opportunity.WorkMode, opportunity.Stage)
	}
	return nil
}

func resume(args []string) error {
	if len(args) == 0 || args[0] != "tailor" {
		return errors.New("usage: open-career-mcp resume tailor --opportunity syn-001 [--fixtures fixtures] [--json]")
	}
	fs := flag.NewFlagSet("resume tailor", flag.ContinueOnError)
	fixtureDir := fs.String("fixtures", "fixtures", "fixture directory")
	opportunityID := fs.String("opportunity", "", "synthetic opportunity id")
	jsonOut := fs.Bool("json", false, "print JSON")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if *opportunityID == "" {
		return errors.New("--opportunity is required")
	}
	resume, opportunities, err := career.LoadFixtures(*fixtureDir)
	if err != nil {
		return err
	}
	opportunity, err := career.FindOpportunity(opportunities, *opportunityID)
	if err != nil {
		return err
	}
	plan := career.BuildTailoringPlan(resume, opportunity)
	if *jsonOut {
		return printJSON(plan)
	}
	fmt.Printf("# Tailoring Plan: %s\n\n%s\n\nMatched skills: %s\n\nGaps: %s\n\n", plan.Role, plan.Summary, strings.Join(plan.MatchedSkills, ", "), strings.Join(plan.Gaps, ", "))
	for _, bullet := range plan.Bullets {
		fmt.Printf("- %s\n", bullet)
	}
	fmt.Printf("\n%s\n", plan.ApprovalNote)
	return nil
}

func interview(args []string) error {
	if len(args) == 0 || args[0] != "prep" {
		return errors.New("usage: open-career-mcp interview prep --opportunity syn-001 [--fixtures fixtures] [--json]")
	}
	fs := flag.NewFlagSet("interview prep", flag.ContinueOnError)
	fixtureDir := fs.String("fixtures", "fixtures", "fixture directory")
	opportunityID := fs.String("opportunity", "", "synthetic opportunity id")
	jsonOut := fs.Bool("json", false, "print JSON")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if *opportunityID == "" {
		return errors.New("--opportunity is required")
	}
	resume, opportunities, err := career.LoadFixtures(*fixtureDir)
	if err != nil {
		return err
	}
	opportunity, err := career.FindOpportunity(opportunities, *opportunityID)
	if err != nil {
		return err
	}
	packet := career.BuildInterviewPacket(resume, opportunity)
	if *jsonOut {
		return printJSON(packet)
	}
	fmt.Printf("# Interview Prep: %s at %s\n\n", packet.Role, packet.Company)
	fmt.Println("Talking points:")
	for _, point := range packet.TalkingPoints {
		fmt.Printf("- %s\n", point)
	}
	fmt.Println("\nPractice questions:")
	for _, question := range packet.Questions {
		fmt.Printf("- %s\n", question)
	}
	fmt.Printf("\n%s\n", packet.BoundaryNote)
	return nil
}

func mcp(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: open-career-mcp mcp manifest|call")
	}
	switch args[0] {
	case "manifest":
		return printJSON(career.Manifest())
	case "call":
		return mcpCall(args[1:])
	default:
		return fmt.Errorf("unknown mcp command %q", args[0])
	}
}

func mcpCall(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: open-career-mcp mcp call <tool> [--param opportunity_id=syn-001] [--fixtures fixtures]")
	}
	tool := args[0]
	fs := flag.NewFlagSet("mcp call", flag.ContinueOnError)
	fixtureDir := fs.String("fixtures", "fixtures", "fixture directory")
	params := multiFlag{}
	fs.Var(&params, "param", "key=value parameter")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	resume, opportunities, err := career.LoadFixtures(*fixtureDir)
	if err != nil {
		return err
	}
	switch tool {
	case "open_career_list_opportunities":
		return printJSON(opportunities)
	case "open_career_tailor_resume":
		opportunity, err := career.FindOpportunity(opportunities, params.Get("opportunity_id"))
		if err != nil {
			return err
		}
		return printJSON(career.BuildTailoringPlan(resume, opportunity))
	case "open_career_interview_prep":
		opportunity, err := career.FindOpportunity(opportunities, params.Get("opportunity_id"))
		if err != nil {
			return err
		}
		return printJSON(career.BuildInterviewPacket(resume, opportunity))
	default:
		return fmt.Errorf("unknown MCP-style tool %q", tool)
	}
}

func usage() error {
	return errors.New(`usage:
  open-career-mcp opportunities list [--json]
  open-career-mcp resume tailor --opportunity syn-001 [--json]
  open-career-mcp interview prep --opportunity syn-001 [--json]
  open-career-mcp mcp manifest
  open-career-mcp mcp call <tool> --param opportunity_id=syn-001`)
}

func printJSON(value any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

type multiFlag []string

func (m *multiFlag) String() string {
	return strings.Join(*m, ",")
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func (m multiFlag) Get(key string) string {
	prefix := key + "="
	for _, value := range m {
		if strings.HasPrefix(value, prefix) {
			return strings.TrimPrefix(value, prefix)
		}
	}
	return ""
}
