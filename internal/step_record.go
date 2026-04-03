package internal

import (
	"context"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// --- step.crm_create_record ---

type createRecordStep struct {
	name       string
	moduleName string
}

func newCreateRecordStep(name string, config map[string]any) (*createRecordStep, error) {
	return &createRecordStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createRecordStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	provider, ok := GetProvider(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "crm provider not found: " + s.moduleName}}, nil
	}
	objectType := resolveValue("object_type", current, config)
	if objectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "object_type is required"}}, nil
	}
	fields := resolveMap("fields", current, config)
	if fields == nil {
		fields = map[string]any{}
	}
	result, err := provider.CreateRecord(ctx, objectType, fields)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":      result.ID,
		"success": result.Success,
		"errors":  result.Errors,
	}}, nil
}

// --- step.crm_get_record ---

type getRecordStep struct {
	name       string
	moduleName string
}

func newGetRecordStep(name string, config map[string]any) (*getRecordStep, error) {
	return &getRecordStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *getRecordStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	provider, ok := GetProvider(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "crm provider not found: " + s.moduleName}}, nil
	}
	objectType := resolveValue("object_type", current, config)
	recordID := resolveValue("record_id", current, config)
	if objectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "object_type is required"}}, nil
	}
	if recordID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "record_id is required"}}, nil
	}
	result, err := provider.GetRecord(ctx, objectType, recordID)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: result}, nil
}

// --- step.crm_update_record ---

type updateRecordStep struct {
	name       string
	moduleName string
}

func newUpdateRecordStep(name string, config map[string]any) (*updateRecordStep, error) {
	return &updateRecordStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *updateRecordStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	provider, ok := GetProvider(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "crm provider not found: " + s.moduleName}}, nil
	}
	objectType := resolveValue("object_type", current, config)
	recordID := resolveValue("record_id", current, config)
	if objectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "object_type is required"}}, nil
	}
	if recordID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "record_id is required"}}, nil
	}
	fields := resolveMap("fields", current, config)
	if fields == nil {
		fields = map[string]any{}
	}
	if err := provider.UpdateRecord(ctx, objectType, recordID, fields); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"success": true}}, nil
}

// --- step.crm_upsert_record ---

type upsertRecordStep struct {
	name       string
	moduleName string
}

func newUpsertRecordStep(name string, config map[string]any) (*upsertRecordStep, error) {
	return &upsertRecordStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *upsertRecordStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	provider, ok := GetProvider(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "crm provider not found: " + s.moduleName}}, nil
	}
	objectType := resolveValue("object_type", current, config)
	extField := resolveValue("external_id_field", current, config)
	extValue := resolveValue("external_id_value", current, config)
	if objectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "object_type is required"}}, nil
	}
	if extField == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "external_id_field is required"}}, nil
	}
	if extValue == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "external_id_value is required"}}, nil
	}
	fields := resolveMap("fields", current, config)
	if fields == nil {
		fields = map[string]any{}
	}
	result, err := provider.UpsertRecord(ctx, objectType, extField, extValue, fields)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{
		"id":      result.ID,
		"success": result.Success,
	}}, nil
}

// --- step.crm_delete_record ---

type deleteRecordStep struct {
	name       string
	moduleName string
}

func newDeleteRecordStep(name string, config map[string]any) (*deleteRecordStep, error) {
	return &deleteRecordStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *deleteRecordStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	provider, ok := GetProvider(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "crm provider not found: " + s.moduleName}}, nil
	}
	objectType := resolveValue("object_type", current, config)
	recordID := resolveValue("record_id", current, config)
	if objectType == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "object_type is required"}}, nil
	}
	if recordID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "record_id is required"}}, nil
	}
	if err := provider.DeleteRecord(ctx, objectType, recordID); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"success": true}}, nil
}

// suppress unused import
var _ = fmt.Sprintf
