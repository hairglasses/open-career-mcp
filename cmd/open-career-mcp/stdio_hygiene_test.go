package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func validateStdioHygieneFrame(line string, expectedRawID string) (map[string]json.RawMessage, error) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return nil, errors.New("stdiohygiene: empty frame")
	}
	if !strings.HasPrefix(trimmed, "{") || !strings.HasSuffix(trimmed, "}") {
		return nil, fmt.Errorf("stdiohygiene: non-JSON stdout noise %q", trimmed)
	}
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal([]byte(trimmed), &envelope); err != nil {
		return nil, fmt.Errorf("stdiohygiene: malformed JSON frame: %w", err)
	}
	rawVersion, ok := envelope["jsonrpc"]
	if !ok {
		return nil, errors.New("stdiohygiene: missing jsonrpc field")
	}
	var version string
	if err := json.Unmarshal(rawVersion, &version); err != nil || version != "2.0" {
		return nil, fmt.Errorf("stdiohygiene: invalid jsonrpc version %s", string(rawVersion))
	}
	rawID, hasID := envelope["id"]
	if !hasID {
		return nil, errors.New("stdiohygiene: missing id field in response")
	}
	gotID := strings.TrimSpace(string(rawID))
	if gotID != expectedRawID {
		return nil, fmt.Errorf("stdiohygiene: id mismatch got=%s want=%s", gotID, expectedRawID)
	}
	_, hasResult := envelope["result"]
	_, hasError := envelope["error"]
	if !hasResult && !hasError {
		return nil, errors.New("stdiohygiene: response missing both result and error")
	}
	return envelope, nil
}

func validateMCPToolSchema(toolName string, schema map[string]any, path string) []string {
	var findings []string
	kind, _ := schema["type"].(string)
	switch kind {
	case "array":
		if _, hasItems := schema["items"]; !hasItems {
			if _, hasPrefix := schema["prefixItems"]; !hasPrefix {
				findings = append(findings, fmt.Sprintf("ARRAY_MISSING_ITEMS:%s%s", toolName, path))
			}
		}
	case "object":
		props, ok := schema["properties"].(map[string]any)
		if !ok || len(props) == 0 {
			findings = append(findings, fmt.Sprintf("OBJECT_MISSING_PROPERTIES:%s%s", toolName, path))
		}
	case "":
		findings = append(findings, fmt.Sprintf("NODE_MISSING_TYPE:%s%s", toolName, path))
	}
	if props, ok := schema["properties"].(map[string]any); ok {
		if reqList, ok := schema["required"].([]any); ok {
			for _, rawReq := range reqList {
				if reqName, ok := rawReq.(string); ok {
					if _, exists := props[reqName]; !exists {
						findings = append(findings, fmt.Sprintf("REQUIRED_NOT_IN_PROPERTIES:%s%s.%s", toolName, path, reqName))
					}
				}
			}
		}
		for propName, rawSub := range props {
			if subMap, ok := rawSub.(map[string]any); ok {
				findings = append(findings, validateMCPToolSchema(toolName, subMap, path+"."+propName)...)
			}
		}
	}
	if itemsMap, ok := schema["items"].(map[string]any); ok {
		findings = append(findings, validateMCPToolSchema(toolName, itemsMap, path+"[]")...)
	}
	return findings
}

func TestStdioHygieneAndJSONRPCFraming(t *testing.T) {
	t.Run("stdiohygiene_live_binary", func(t *testing.T) {
		repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
		if err != nil {
			t.Fatalf("resolve repo root: %v", err)
		}
		binDir := t.TempDir()
		binPath := filepath.Join(binDir, "open-career-mcp")

		buildCmd := exec.Command("go", "build", "-o", binPath, "./cmd/open-career-mcp")
		buildCmd.Dir = repoRoot
		buildCmd.Env = append(os.Environ(), "GOWORK=off")
		if out, err := buildCmd.CombinedOutput(); err != nil {
			t.Fatalf("go build failed: %v\n%s", err, string(out))
		}

		isolatedHome := t.TempDir()
		cmd := exec.Command(binPath)
		cmd.Dir = repoRoot
		cmd.Env = []string{
			"HOME=" + isolatedHome,
			"PATH=" + os.Getenv("PATH"),
			"GOWORK=off",
		}

		stdinPipe, err := cmd.StdinPipe()
		if err != nil {
			t.Fatalf("StdinPipe: %v", err)
		}
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			t.Fatalf("StdoutPipe: %v", err)
		}
		var stderrBuf bytes.Buffer
		cmd.Stderr = &stderrBuf

		if err := cmd.Start(); err != nil {
			t.Fatalf("Start binary: %v", err)
		}
		defer func() {
			_ = stdinPipe.Close()
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
		}()

		linesCh := make(chan string, 16)
		readErrCh := make(chan error, 1)
		go func() {
			scanner := bufio.NewScanner(stdoutPipe)
			scanner.Buffer(make([]byte, 64*1024), 1024*1024)
			for scanner.Scan() {
				linesCh <- scanner.Text()
			}
			readErrCh <- scanner.Err()
			close(linesCh)
		}()

		select {
		case unexpected := <-linesCh:
			t.Fatalf("stdiohygiene violation: unexpected startup stdout bytes before request: %q", unexpected)
		case <-time.After(80 * time.Millisecond):
		}

		sendLine := func(payload string) {
			t.Helper()
			if _, err := fmt.Fprintf(stdinPipe, "%s\n", payload); err != nil {
				t.Fatalf("write stdin: %v", err)
			}
		}

		readNextFrame := func(expectedID string) map[string]json.RawMessage {
			t.Helper()
			select {
			case line, ok := <-linesCh:
				if !ok {
					t.Fatalf("stdout closed while waiting for id=%s (stderr=%s)", expectedID, stderrBuf.String())
				}
				env, err := validateStdioHygieneFrame(line, expectedID)
				if err != nil {
					t.Fatalf("frame validation failed for id=%s: %v (line=%q)", expectedID, err, line)
				}
				return env
			case <-time.After(5 * time.Second):
				t.Fatalf("timeout waiting for JSON-RPC response id=%s (stderr=%s)", expectedID, stderrBuf.String())
				return nil
			}
		}

		sendLine(`{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"stdiohygiene-test","version":"1.0.0"}}}`)
		initEnv := readNextFrame("0")
		var initResult struct {
			ProtocolVersion string `json:"protocolVersion"`
			ServerInfo      struct {
				Name string `json:"name"`
			} `json:"serverInfo"`
		}
		if err := json.Unmarshal(initEnv["result"], &initResult); err != nil {
			t.Fatalf("unmarshal initialize result: %v", err)
		}
		if initResult.ServerInfo.Name != "open-career-mcp" {
			t.Fatalf("unexpected serverInfo.name=%q", initResult.ServerInfo.Name)
		}

		sendLine(`{"jsonrpc":"2.0","method":"notifications/initialized"}`)
		sendLine(`{"jsonrpc":"2.0","id":null,"method":"ping"}`)
		_ = readNextFrame("null")

		sendLine(`{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}`)
		listEnv := readNextFrame("1")

		var listResult struct {
			Tools []struct {
				Name        string         `json:"name"`
				Description string         `json:"description"`
				InputSchema map[string]any `json:"inputSchema"`
			} `json:"tools"`
		}
		if err := json.Unmarshal(listEnv["result"], &listResult); err != nil {
			t.Fatalf("unmarshal tools/list result: %v", err)
		}
		if len(listResult.Tools) != 3 {
			t.Fatalf("expected 3 tools from tools/list, got %d", len(listResult.Tools))
		}
		expectedTools := map[string]bool{
			"open_career_list_opportunities": false,
			"open_career_tailor_resume":      false,
			"open_career_interview_prep":     false,
		}
		for _, tool := range listResult.Tools {
			if _, ok := expectedTools[tool.Name]; !ok {
				t.Fatalf("unexpected tool %q", tool.Name)
			}
			expectedTools[tool.Name] = true
			if findings := validateMCPToolSchema(tool.Name, tool.InputSchema, ""); len(findings) != 0 {
				t.Fatalf("tool %s schema findings: %v", tool.Name, findings)
			}
		}
		for name, seen := range expectedTools {
			if !seen {
				t.Fatalf("missing expected tool %s", name)
			}
		}

		sendLine(`{"jsonrpc":"2.0","id":"call-2","method":"tools/call","params":{"name":"open_career_tailor_resume","arguments":{"opportunity_id":"syn-001"}}}`)
		callEnv := readNextFrame(`"call-2"`)
		var callResult struct {
			IsError bool `json:"isError"`
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		}
		if err := json.Unmarshal(callEnv["result"], &callResult); err != nil {
			t.Fatalf("unmarshal tools/call result: %v", err)
		}
		if callResult.IsError || len(callResult.Content) == 0 || !strings.Contains(callResult.Content[0].Text, "syn-001") {
			t.Fatalf("unexpected tools/call payload: %+v", callResult)
		}

		if err := stdinPipe.Close(); err != nil {
			t.Fatalf("close stdin: %v", err)
		}

		waitCh := make(chan error, 1)
		go func() {
			waitCh <- cmd.Wait()
		}()

		select {
		case waitErr := <-waitCh:
			if waitErr != nil {
				t.Fatalf("process did not exit cleanly on stdin EOF: %v (stderr=%s)", waitErr, stderrBuf.String())
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("process did not exit within 5s after stdin EOF")
		}

		for trailing := range linesCh {
			t.Fatalf("stdiohygiene violation: trailing stdout line after final response: %q", trailing)
		}
		if readErr := <-readErrCh; readErr != nil {
			t.Fatalf("stdout scanner error: %v", readErr)
		}
	})

	t.Run("negative_controls", func(t *testing.T) {
		cases := []struct {
			name       string
			line       string
			expectedID string
		}{
			{
				name:       "reject_stdout_banner_noise",
				line:       "open-career-mcp starting on stdio...",
				expectedID: "0",
			},
			{
				name:       "reject_malformed_json_line",
				line:       `{"jsonrpc":"2.0","id":0,"result":`,
				expectedID: "0",
			},
			{
				name:       "reject_null_id_when_zero_requested",
				line:       `{"jsonrpc":"2.0","id":null,"result":{"protocolVersion":"2025-06-18"}}`,
				expectedID: "0",
			},
			{
				name:       "reject_mismatched_id",
				line:       `{"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2025-06-18"}}`,
				expectedID: "0",
			},
			{
				name:       "reject_invalid_jsonrpc_version",
				line:       `{"jsonrpc":"1.0","id":0,"result":{}}`,
				expectedID: "0",
			},
			{
				name:       "reject_missing_result_and_error",
				line:       `{"jsonrpc":"2.0","id":0}`,
				expectedID: "0",
			},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if _, err := validateStdioHygieneFrame(tc.line, tc.expectedID); err == nil {
					t.Fatalf("expected stdiohygiene validator to reject %s (%q), but it passed", tc.name, tc.line)
				}
			})
		}

		badSchemas := []struct {
			name   string
			schema map[string]any
		}{
			{
				name:   "reject_array_missing_items",
				schema: map[string]any{"type": "object", "properties": map[string]any{"tags": map[string]any{"type": "array"}}},
			},
			{
				name:   "reject_object_missing_properties",
				schema: map[string]any{"type": "object", "properties": map[string]any{}},
			},
			{
				name:   "reject_required_not_in_properties",
				schema: map[string]any{"type": "object", "properties": map[string]any{"stage": map[string]any{"type": "string"}}, "required": []any{"missing_field"}},
			},
		}
		for _, bs := range badSchemas {
			t.Run(bs.name, func(t *testing.T) {
				if findings := validateMCPToolSchema("bad_tool", bs.schema, ""); len(findings) == 0 {
					t.Fatalf("expected schema validator to reject %s", bs.name)
				}
			})
		}
	})
}
