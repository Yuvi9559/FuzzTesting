package llm

import (
	"context"
	"fmt"
)

// HarnessRequest contains everything needed to generate a fuzzing harness
type HarnessRequest struct {
	TargetType  string // "c_cpp" | "python" | "rest_api"
	TargetName  string // function name or endpoint
	InputVector string // "file" | "stdin" | "network" | "http_request"
	Description string // context about what the target does
}

// HarnessResponse contains the generated harness code
type HarnessResponse struct {
	Code     string
	Language string
}

// GenerateHarness uses Claude to generate a fuzzing harness for the given target
func (c *Client) GenerateHarness(ctx context.Context, req HarnessRequest) (*HarnessResponse, error) {
	system := `You are an expert fuzzing engineer at a security research firm.
Generate a complete, production-grade fuzzing harness.
Return ONLY the code with inline comments explaining key decisions.
No markdown fences, no explanation outside the code itself.`

	user := fmt.Sprintf(`Generate a fuzzing harness for the following target:

Target type:      %s
Function/endpoint: %s
Input vector:     %s
Description:      %s

Requirements based on target type:
- c_cpp:    Use libFuzzer format with LLVMFuzzerTestOneInput().
            Include AddressSanitizer hooks and corpus seed hints in comments.
- python:   Use Atheris format with atheris.Setup() and atheris.Fuzz().
            Include proper exception handling for fuzzing context.
- rest_api: Use Python requests library with structured mutation strategies.
            Include headers, auth handling, and response validation.

Always include:
1. Boundary condition tests
2. Null/empty input handling
3. Mutation strategy hints in comments
4. How to compile/run at the top as a comment`,
		req.TargetType, req.TargetName, req.InputVector, req.Description)

	code, err := c.Complete(ctx, system, user)
	if err != nil {
		return nil, fmt.Errorf("harness generation failed: %w", err)
	}

	lang := "c"
	if req.TargetType == "python" || req.TargetType == "rest_api" {
		lang = "python"
	}

	return &HarnessResponse{Code: code, Language: lang}, nil
}