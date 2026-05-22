package llm

import (
	"context"
	"fmt"
)

// PatchRequest describes the crash to generate a fix for
type PatchRequest struct {
	CrashID    string
	CrashType  string // e.g. "Heap Buffer Overflow"
	Function   string // e.g. "parse_input()"
	File       string // e.g. "parser.c:142"
	StackTrace string
	ASANOutput string
}

// PatchResponse contains the generated patch and metadata
type PatchResponse struct {
	PatchCode   string
	RootCause   string
	Confidence  string // "HIGH" | "MEDIUM" | "LOW"
}

// GeneratePatch uses Claude to suggest a security patch for a crash
func (c *Client) GeneratePatch(ctx context.Context, req PatchRequest) (*PatchResponse, error) {
	system := `You are a senior security engineer specializing in memory safety vulnerabilities.
Generate a minimal, correct security patch.
Return ONLY the patched code with inline comments.
No markdown fences, no preamble outside the code itself.`

	user := fmt.Sprintf(`Generate a security patch for this vulnerability:

Crash ID:    %s
Type:        %s
Function:    %s
Source file: %s

ASAN Output:
%s

Stack Trace:
%s

Requirements:
1. Write the corrected function with the vulnerability fixed
2. Add a comment block at the top explaining:
   - Root cause of the vulnerability
   - What the fix does and why it is correct
   - Any edge cases the fix handles
3. Add a comment with a regression test recommendation
4. Keep changes minimal — only fix what is broken`,
		req.CrashID, req.CrashType, req.Function, req.File,
		req.ASANOutput, req.StackTrace)

	code, err := c.Complete(ctx, system, user)
	if err != nil {
		return nil, fmt.Errorf("patch generation failed: %w", err)
	}

	return &PatchResponse{
		PatchCode:  code,
		RootCause:  req.CrashType,
		Confidence: "MEDIUM",
	}, nil
}