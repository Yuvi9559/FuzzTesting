package llm

import (
	"context"
	"encoding/json"
	"fmt"
)

// TriageRequest contains crash info for LLM-assisted triage
type TriageRequest struct {
	CrashType  string
	StackTrace string
	ASANOutput string
	Function   string
	File       string
}

// TriageResult is the structured output from crash triage
type TriageResult struct {
	Severity    string `json:"severity"`    // CRITICAL | HIGH | MEDIUM | LOW
	Exploitable bool   `json:"exploitable"`
	RootCause   string `json:"root_cause"`
	Summary     string `json:"summary"`
	CVSSScore   string `json:"cvss_score"`
}

// TriageCrash uses Claude to automatically triage a crash
func (c *Client) TriageCrash(ctx context.Context, req TriageRequest) (*TriageResult, error) {
	system := `You are a vulnerability triage expert at a security operations center.
Analyze the crash output and return ONLY a valid JSON object.
No markdown, no explanation, no backticks — raw JSON only.
JSON fields: severity (string), exploitable (bool), root_cause (string), summary (string), cvss_score (string)`

	user := fmt.Sprintf(`Triage this crash and classify its severity:

Crash Type:  %s
Function:    %s
File:        %s

ASAN Output:
%s

Stack Trace:
%s

Classify severity as:
- CRITICAL: Memory corruption that is likely exploitable (RCE, heap overflow with control)
- HIGH:     Memory safety issue, exploitability unclear (UAF, stack overflow)
- MEDIUM:   Crash but limited exploit potential (null deref, assertion)
- LOW:      Informational or DoS only`,
		req.CrashType, req.Function, req.File,
		req.ASANOutput, req.StackTrace)

	raw, err := c.Complete(ctx, system, user)
	if err != nil {
		return nil, fmt.Errorf("triage failed: %w", err)
	}

	var result TriageResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		// If JSON parse fails, return a best-effort result
		return &TriageResult{
			Severity:    "MEDIUM",
			Exploitable: false,
			RootCause:   req.CrashType,
			Summary:     raw,
			CVSSScore:   "5.0",
		}, nil
	}
	return &result, nil
}