// Package internal implements the workflow-plugin-crm plugin.
package internal

import (
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// Version is set at build time via -ldflags
// "-X github.com/GoCodeAlone/workflow-plugin-crm/internal.Version=X.Y.Z"
var Version = "0.0.0"

// crmPlugin implements sdk.PluginProvider, sdk.ModuleProvider, sdk.StepProvider,
// and sdk.SchemaProvider.
type crmPlugin struct{}

// NewCRMPlugin returns a new crmPlugin instance.
func NewCRMPlugin() sdk.PluginProvider {
	return &crmPlugin{}
}

// Manifest returns plugin metadata.
func (p *crmPlugin) Manifest() sdk.PluginManifest {
	return sdk.PluginManifest{
		Name:        "workflow-plugin-crm",
		Version:     Version,
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

// ModuleSchemas implements sdk.SchemaProvider, returning UI/contract schema
// definitions for all module types provided by this plugin.
func (p *crmPlugin) ModuleSchemas() []sdk.ModuleSchemaData {
	return []sdk.ModuleSchemaData{
		{
			Type:        "crm.provider",
			Label:       "CRM Provider",
			Category:    "crm",
			Description: "Vendor-neutral CRM connection provider. Manages authentication and connection lifecycle for a CRM backend (currently supports Salesforce).",
			ConfigFields: []sdk.ConfigField{
				{Name: "provider", Type: "string", Description: "CRM backend type. Currently only 'salesforce' is supported.", DefaultValue: "salesforce"},
				{Name: "authType", Type: "string", Description: "Authentication type: 'oauth2' (recommended) or 'password'."},
				{Name: "clientId", Type: "string", Description: "OAuth2 connected-app client ID."},
				{Name: "clientSecret", Type: "string", Description: "OAuth2 connected-app client secret."},
				{Name: "refreshToken", Type: "string", Description: "OAuth2 refresh token for token-refresh flows."},
				{Name: "username", Type: "string", Description: "Salesforce username for password-flow or SOAP login."},
				{Name: "password", Type: "string", Description: "Salesforce password (password-flow authentication)."},
				{Name: "security_token", Type: "string", Description: "Salesforce security token appended to password during login."},
				{Name: "accessToken", Type: "string", Description: "Pre-obtained OAuth2 access token (skips token exchange)."},
				{Name: "instanceUrl", Type: "string", Description: "Salesforce instance URL (e.g. https://na1.salesforce.com)."},
				{Name: "apiVersion", Type: "string", Description: "Salesforce API version to use (e.g. 60.0)."},
				{Name: "loginUrl", Type: "string", Description: "Salesforce login URL override (defaults to https://login.salesforce.com)."},
				{Name: "sandbox", Type: "bool", Description: "When true, connects to a sandbox org (uses https://test.salesforce.com)."},
			},
		},
	}
}
