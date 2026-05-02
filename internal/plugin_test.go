package internal_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/GoCodeAlone/workflow-plugin-crm/internal"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

func TestNewPlugin_ImplementsPluginProvider(t *testing.T) {
	var _ sdk.PluginProvider = internal.NewCRMPlugin()
}

func TestPlugin_ImplementsModuleProvider(t *testing.T) {
	p := internal.NewCRMPlugin()
	mp, ok := p.(sdk.ModuleProvider)
	if !ok {
		t.Fatal("plugin does not implement ModuleProvider")
	}
	types := mp.ModuleTypes()
	if len(types) != 1 || types[0] != "crm.provider" {
		t.Errorf("unexpected module types: %v", types)
	}
}

func TestPlugin_ImplementsStepProvider(t *testing.T) {
	p := internal.NewCRMPlugin()
	sp, ok := p.(sdk.StepProvider)
	if !ok {
		t.Fatal("plugin does not implement StepProvider")
	}
	types := sp.StepTypes()
	if len(types) != 10 {
		t.Errorf("expected 10 step types, got %d: %v", len(types), types)
	}
}

func TestManifest_HasRequiredFields(t *testing.T) {
	p := internal.NewCRMPlugin()
	m := p.Manifest()
	if m.Name == "" {
		t.Error("manifest Name is empty")
	}
	if m.Version == "" {
		t.Error("manifest Version is empty")
	}
	if m.Description == "" {
		t.Error("manifest Description is empty")
	}
}

// TestPlugin_ImplementsSchemaProvider verifies the plugin implements sdk.SchemaProvider
// and returns a non-empty module schema for crm.provider.
func TestPlugin_ImplementsSchemaProvider(t *testing.T) {
	p := internal.NewCRMPlugin()
	sp, ok := p.(sdk.SchemaProvider)
	if !ok {
		t.Fatal("plugin does not implement SchemaProvider")
	}
	schemas := sp.ModuleSchemas()
	if len(schemas) == 0 {
		t.Fatal("ModuleSchemas returned empty slice")
	}
	found := false
	for _, s := range schemas {
		if s.Type == "crm.provider" {
			found = true
			if s.Label == "" {
				t.Error("crm.provider schema has empty Label")
			}
			if s.Description == "" {
				t.Error("crm.provider schema has empty Description")
			}
			if len(s.ConfigFields) == 0 {
				t.Error("crm.provider schema has no ConfigFields")
			}
		}
	}
	if !found {
		t.Error("crm.provider not found in ModuleSchemas")
	}
}

// TestModuleSchema_CRMProviderFields checks that all expected config field names
// are declared in the crm.provider module schema contract.
func TestModuleSchema_CRMProviderFields(t *testing.T) {
	p := internal.NewCRMPlugin()
	sp, ok := p.(sdk.SchemaProvider)
	if !ok {
		t.Fatal("plugin does not implement SchemaProvider")
	}
	var crmSchema *sdk.ModuleSchemaData
	for _, s := range sp.ModuleSchemas() {
		s := s
		if s.Type == "crm.provider" {
			crmSchema = &s
			break
		}
	}
	if crmSchema == nil {
		t.Fatal("crm.provider schema not found")
	}

	expectedFields := []string{
		"provider", "authType", "clientId", "clientSecret",
		"refreshToken", "username", "password", "security_token",
		"accessToken", "instanceUrl", "apiVersion", "loginUrl", "sandbox",
	}
	fieldSet := make(map[string]bool, len(crmSchema.ConfigFields))
	for _, f := range crmSchema.ConfigFields {
		fieldSet[f.Name] = true
	}
	for _, name := range expectedFields {
		if !fieldSet[name] {
			t.Errorf("expected config field %q not found in crm.provider schema", name)
		}
	}
}

// TestPluginJSON_HasStepSchemas verifies that plugin.json exists and declares
// stepSchemas for every step type advertised in stepTypes.
func TestPluginJSON_HasStepSchemas(t *testing.T) {
	data, err := os.ReadFile("../plugin.json")
	if err != nil {
		t.Fatalf("cannot read plugin.json: %v", err)
	}

	var manifest struct {
		StepTypes  []string `json:"stepTypes"`
		StepSchemas []struct {
			Type string `json:"type"`
		} `json:"stepSchemas"`
		// handle legacy capabilities object shape
		Capabilities *struct {
			StepTypes []string `json:"stepTypes"`
		} `json:"capabilities"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("cannot parse plugin.json: %v", err)
	}

	// Collect all declared step types (canonical or legacy)
	allStepTypes := manifest.StepTypes
	if manifest.Capabilities != nil {
		allStepTypes = append(allStepTypes, manifest.Capabilities.StepTypes...)
	}

	if len(allStepTypes) == 0 {
		t.Fatal("plugin.json declares no step types")
	}
	if len(manifest.StepSchemas) == 0 {
		t.Fatal("plugin.json has no stepSchemas")
	}

	schemaByType := make(map[string]bool, len(manifest.StepSchemas))
	for _, s := range manifest.StepSchemas {
		schemaByType[s.Type] = true
	}
	for _, st := range allStepTypes {
		if !schemaByType[st] {
			t.Errorf("step type %q has no entry in plugin.json stepSchemas", st)
		}
	}
}

// TestPluginJSON_CanonicalFormat verifies that the plugin.json uses the canonical
// manifest shape with moduleTypes/stepTypes declared inside a top-level capabilities object.
func TestPluginJSON_CanonicalFormat(t *testing.T) {
	data, err := os.ReadFile("../plugin.json")
	if err != nil {
		t.Fatalf("cannot read plugin.json: %v", err)
	}

	var manifest struct {
		Capabilities *struct {
			ModuleTypes []string `json:"moduleTypes"`
			StepTypes   []string `json:"stepTypes"`
		} `json:"capabilities"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("cannot parse plugin.json: %v", err)
	}
	if manifest.Capabilities == nil {
		t.Fatal("plugin.json: top-level capabilities object is missing (canonical format required)")
	}
	if len(manifest.Capabilities.ModuleTypes) == 0 {
		t.Error("plugin.json: capabilities.moduleTypes is empty")
	}
	if len(manifest.Capabilities.StepTypes) == 0 {
		t.Error("plugin.json: capabilities.stepTypes is empty")
	}
}
