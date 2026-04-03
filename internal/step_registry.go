package internal

import (
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type stepConstructor func(name string, config map[string]any) (sdk.StepInstance, error)

var stepRegistry = map[string]stepConstructor{
	"step.crm_create_record":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newCreateRecordStep(n, c) },
	"step.crm_get_record":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newGetRecordStep(n, c) },
	"step.crm_update_record":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newUpdateRecordStep(n, c) },
	"step.crm_upsert_record":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newUpsertRecordStep(n, c) },
	"step.crm_delete_record":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newDeleteRecordStep(n, c) },
	"step.crm_query":           func(n string, c map[string]any) (sdk.StepInstance, error) { return newQueryStep(n, c) },
	"step.crm_search":          func(n string, c map[string]any) (sdk.StepInstance, error) { return newSearchStep(n, c) },
	"step.crm_bulk_import":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newBulkImportStep(n, c) },
	"step.crm_describe_object": func(n string, c map[string]any) (sdk.StepInstance, error) { return newDescribeObjectStep(n, c) },
	"step.crm_get_limits":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newGetLimitsStep(n, c) },
}

func createStep(typeName, name string, config map[string]any) (sdk.StepInstance, error) {
	constructor, ok := stepRegistry[typeName]
	if !ok {
		return nil, fmt.Errorf("crm plugin: unknown step type %q", typeName)
	}
	return constructor(name, config)
}

func allStepTypes() []string {
	types := make([]string, 0, len(stepRegistry))
	for k := range stepRegistry {
		types = append(types, k)
	}
	return types
}
