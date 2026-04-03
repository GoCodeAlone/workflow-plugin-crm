// Package internal implements the workflow-plugin-crm plugin.
package internal

import (
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// crmPlugin implements sdk.PluginProvider, sdk.ModuleProvider, and sdk.StepProvider.
type crmPlugin struct{}

// NewCRMPlugin returns a new crmPlugin instance.
func NewCRMPlugin() sdk.PluginProvider {
	return &crmPlugin{}
}

// Manifest returns plugin metadata.
func (p *crmPlugin) Manifest() sdk.PluginManifest {
	return sdk.PluginManifest{
		Name:        "workflow-plugin-crm",
		Version:     "0.1.0",
		Author:      "GoCodeAlone",
		Description: "Vendor-neutral CRM plugin with pluggable provider architecture (Salesforce adapter)",
	}
}

// ModuleTypes returns the module type names this plugin provides.
func (p *crmPlugin) ModuleTypes() []string {
	return []string{"crm.provider"}
}

// CreateModule creates a module instance of the given type.
func (p *crmPlugin) CreateModule(typeName, name string, config map[string]any) (sdk.ModuleInstance, error) {
	switch typeName {
	case "crm.provider":
		return newCRMModule(name, config)
	default:
		return nil, fmt.Errorf("crm plugin: unknown module type %q", typeName)
	}
}

// StepTypes returns the step type names this plugin provides.
func (p *crmPlugin) StepTypes() []string {
	return allStepTypes()
}

// CreateStep creates a step instance of the given type.
func (p *crmPlugin) CreateStep(typeName, name string, config map[string]any) (sdk.StepInstance, error) {
	return createStep(typeName, name, config)
}
