package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		return serveStdio(os.Stdin, os.Stdout, resolveFixtureDir("fixtures"))
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
	case "serve", "stdio":
		return mcpServe(args[1:])
	case "help", "-h", "--help":
		return usage()
	default:
		return fmt.Errorf("unknown command %q\n\n%w", args[0], usage())
	}
}

func resolveFixtureDir(requested string) string {
	if requested != "" && requested != "fixtures" {
		return requested
	}
	candidates := []string{
		os.Getenv("OPEN_CAREER_FIXTURES"),
		"fixtures",
		filepath.Join("..", "..", "fixtures"),
		filepath.Join("open-career-mcp", "fixtures"),
	}
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}
	return "fixtures"
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
	_, opportunities, err := career.LoadFixtures(resolveFixtureDir(*fixtureDir))
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
	resume, opportunities, err := career.LoadFixtures(resolveFixtureDir(*fixtureDir))
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
	resume, opportunities, err := career.LoadFixtures(resolveFixtureDir(*fixtureDir))
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
		return errors.New("usage: open-career-mcp mcp manifest|call|serve")
	}
	switch args[0] {
	case "manifest":
		return printJSON(career.Manifest())
	case "call":
		return mcpCall(args[1:])
	case "serve", "stdio":
		return mcpServe(args[1:])
	default:
		return fmt.Errorf("unknown mcp command %q", args[0])
	}
}

func mcpServe(args []string) error {
	fs := flag.NewFlagSet("mcp serve", flag.ContinueOnError)
	fixtureDir := fs.String("fixtures", "fixtures", "fixture directory")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return serveStdio(os.Stdin, os.Stdout, resolveFixtureDir(*fixtureDir))
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
	result, err := executeTool(resolveFixtureDir(*fixtureDir), tool, map[string]string{
		"stage":          params.Get("stage"),
		"opportunity_id": params.Get("opportunity_id"),
	})
	if err != nil {
		return err
	}
	return printJSON(result)
}

func executeTool(fixtureDir, tool string, args map[string]string) (any, error) {
	resumeData, opportunitiesData, err := career.LoadFixtures(fixtureDir)
	if err != nil {
		return nil, err
	}
	switch tool {
	case "open_career_list_opportunities":
		stage := args["stage"]
		if stage != "" {
			var filtered []career.Opportunity
			for _, opp := range opportunitiesData {
				if strings.EqualFold(opp.Stage, stage) {
					filtered = append(filtered, opp)
				}
			}
			return filtered, nil
		}
		return opportunitiesData, nil
	case "open_career_tailor_resume":
		opportunity, err := career.FindOpportunity(opportunitiesData, args["opportunity_id"])
		if err != nil {
			return nil, err
		}
		return career.BuildTailoringPlan(resumeData, opportunity), nil
	case "open_career_interview_prep":
		opportunity, err := career.FindOpportunity(opportunitiesData, args["opportunity_id"])
		if err != nil {
			return nil, err
		}
		return career.BuildInterviewPacket(resumeData, opportunity), nil
	default:
		return nil, fmt.Errorf("unknown MCP-style tool %q", tool)
	}
}

type jsonrpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type jsonrpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *jsonrpcError   `json:"error,omitempty"`
}

func normalizeJSONRPCID(raw json.RawMessage) json.RawMessage {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		return json.RawMessage("null")
	}
	return json.RawMessage(trimmed)
}

func writeJSONRPCResponse(out io.Writer, resp jsonrpcResponse) error {
	resp.JSONRPC = "2.0"
	resp.ID = normalizeJSONRPCID(resp.ID)
	payload, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "%s\n", payload)
	return err
}

func serveStdio(in io.Reader, out io.Writer, fixtureDir string) error {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var raw map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			if writeErr := writeJSONRPCResponse(out, jsonrpcResponse{
				ID:    json.RawMessage("null"),
				Error: &jsonrpcError{Code: -32700, Message: "Parse error"},
			}); writeErr != nil {
				return writeErr
			}
			continue
		}
		rawID, hasID := raw["id"]
		var method string
		if rawMethod, ok := raw["method"]; ok {
			_ = json.Unmarshal(rawMethod, &method)
		}
		if !hasID {
			continue
		}
		if method == "" {
			if err := writeJSONRPCResponse(out, jsonrpcResponse{
				ID:    rawID,
				Error: &jsonrpcError{Code: -32600, Message: "Invalid Request"},
			}); err != nil {
				return err
			}
			continue
		}
		switch method {
		case "initialize":
			if err := writeJSONRPCResponse(out, jsonrpcResponse{
				ID: rawID,
				Result: map[string]any{
					"protocolVersion": "2025-06-18",
					"capabilities": map[string]any{
						"tools": map[string]any{},
					},
					"serverInfo": map[string]any{
						"name":    "open-career-mcp",
						"version": "0.1.0",
					},
				},
			}); err != nil {
				return err
			}
		case "ping":
			if err := writeJSONRPCResponse(out, jsonrpcResponse{
				ID:     rawID,
				Result: map[string]any{},
			}); err != nil {
				return err
			}
		case "tools/list":
			manifest := career.Manifest()
			tools := make([]map[string]any, 0, len(manifest.Tools))
			for _, tool := range manifest.Tools {
				tools = append(tools, map[string]any{
					"name":        tool.Name,
					"description": tool.Description,
					"inputSchema": tool.InputSchema,
				})
			}
			if err := writeJSONRPCResponse(out, jsonrpcResponse{
				ID: rawID,
				Result: map[string]any{
					"tools": tools,
				},
			}); err != nil {
				return err
			}
		case "tools/call":
			var callParams struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments"`
			}
			if rawParams, ok := raw["params"]; ok && len(rawParams) > 0 {
				if err := json.Unmarshal(rawParams, &callParams); err != nil {
					if writeErr := writeJSONRPCResponse(out, jsonrpcResponse{
						ID:    rawID,
						Error: &jsonrpcError{Code: -32602, Message: "Invalid params"},
					}); writeErr != nil {
						return writeErr
					}
					continue
				}
			}
			strArgs := map[string]string{}
			for k, v := range callParams.Arguments {
				if s, ok := v.(string); ok {
					strArgs[k] = s
				}
			}
			result, err := executeTool(fixtureDir, callParams.Name, strArgs)
			if err != nil {
				if writeErr := writeJSONRPCResponse(out, jsonrpcResponse{
					ID: rawID,
					Result: map[string]any{
						"content": []map[string]any{
							{"type": "text", "text": err.Error()},
						},
						"isError": true,
					},
				}); writeErr != nil {
					return writeErr
				}
				continue
			}
			encoded, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return err
			}
			if err := writeJSONRPCResponse(out, jsonrpcResponse{
				ID: rawID,
				Result: map[string]any{
					"content": []map[string]any{
						{"type": "text", "text": string(encoded)},
					},
					"isError": false,
				},
			}); err != nil {
				return err
			}
		default:
			if err := writeJSONRPCResponse(out, jsonrpcResponse{
				ID:    rawID,
				Error: &jsonrpcError{Code: -32601, Message: fmt.Sprintf("Method not found: %s", method)},
			}); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}

func usage() error {
	return errors.New(`usage:
  open-career-mcp opportunities list [--json]
  open-career-mcp resume tailor --opportunity syn-001 [--json]
  open-career-mcp interview prep --opportunity syn-001 [--json]
  open-career-mcp mcp manifest
  open-career-mcp mcp call <tool> --param opportunity_id=syn-001
  open-career-mcp mcp serve [--fixtures fixtures]`)
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
