package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// --- step.crm_query ---

type queryStep struct {
	name       string
	moduleName string
}

func newQueryStep(name string, config map[string]any) (*queryStep, error) {
	return &queryStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *queryStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	provider, ok := GetProvider(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "crm provider not found: " + s.moduleName}}, nil
	}
	query := resolveValue("query", current, config)
	if query == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "query is required"}}, nil
	}
	result, err := provider.Query(ctx, query)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	records := make([]any, len(result.Records))
	for i, r := range result.Records {
		records[i] = r
	}
	return &sdk.StepResult{Output: map[string]any{
		"records":    records,
		"total_size": result.TotalSize,
		"done":       result.Done,
	}}, nil
}

// --- step.crm_search ---

type searchStep struct {
	name       string
	moduleName string
}

func newSearchStep(name string, config map[string]any) (*searchStep, error) {
	return &searchStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *searchStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	provider, ok := GetProvider(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "crm provider not found: " + s.moduleName}}, nil
	}
	query := resolveValue("query", current, config)
	if query == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "query is required"}}, nil
	}
	result, err := provider.Search(ctx, query)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	results := make([]any, len(result.Results))
	for i, r := range result.Results {
		results[i] = r
	}
	return &sdk.StepResult{Output: map[string]any{
		"results": results,
		"count":   len(result.Results),
	}}, nil
}
