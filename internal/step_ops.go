package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// --- step.crm_bulk_import ---

type bulkImportStep struct {
	name       string
	moduleName string
}

func newBulkImportStep(name string, config map[string]any) (*bulkImportStep, error) {
	return &bulkImportStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *bulkImportStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	provider, ok := GetProvider(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "crm provider not found: " + s.moduleName}}, nil
	}
	objectType := resolveValue("object_type", current, config)
	operation := resolveValue("operation", current, config)
	if objectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "object_type is required"}}, nil
	}
	if operation == "" {
		operation = "insert"
	}
	records := resolveAnySlice("records", current, config)

	op := BulkOp{
		Operation:  operation,
		ObjectType: objectType,
		Records:    records,
	}
	result, err := provider.BulkOperation(ctx, op)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"job_id":            result.JobID,
		"state":             result.State,
		"records_processed": result.RecordsProcessed,
		"records_failed":    result.RecordsFailed,
	}}, nil
}

// --- step.crm_describe_object ---

type describeObjectStep struct {
	name       string
	moduleName string
}

func newDescribeObjectStep(name string, config map[string]any) (*describeObjectStep, error) {
	return &describeObjectStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *describeObjectStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	provider, ok := GetProvider(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "crm provider not found: " + s.moduleName}}, nil
	}
	objectType := resolveValue("object_type", current, config)
	if objectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "object_type is required"}}, nil
	}
	result, err := provider.DescribeObject(ctx, objectType)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// --- step.crm_get_limits ---

type getLimitsStep struct {
	name       string
	moduleName string
}

func newGetLimitsStep(name string, config map[string]any) (*getLimitsStep, error) {
	return &getLimitsStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *getLimitsStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	provider, ok := GetProvider(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "crm provider not found: " + s.moduleName}}, nil
	}
	result, err := provider.GetLimits(ctx)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}
