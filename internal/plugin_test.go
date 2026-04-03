package internal_test

import (
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
