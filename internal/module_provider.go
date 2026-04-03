package internal

import (
	"context"
	"fmt"
)

// crmModule creates and manages a CRM provider connection.
type crmModule struct {
	name   string
	config map[string]any
}

func newCRMModule(name string, config map[string]any) (*crmModule, error) {
	return &crmModule{name: name, config: config}, nil
}

// Init creates the appropriate CRM provider based on config and registers it.
func (m *crmModule) Init() error {
	providerType := "salesforce"
	if v, ok := m.config["provider"].(string); ok && v != "" {
		providerType = v
	}

	var provider CRMProvider
	switch providerType {
	case "salesforce":
		provider = &salesforceAdapter{}
	default:
		return fmt.Errorf("crm.provider %q: unsupported provider type %q", m.name, providerType)
	}

	cfg := ProviderConfig{
		Provider:      providerType,
		AuthType:      strCfg(m.config, "authType"),
		ClientID:      strCfg(m.config, "clientId"),
		ClientSecret:  strCfg(m.config, "clientSecret"),
		RefreshToken:  strCfg(m.config, "refreshToken"),
		Username:      strCfg(m.config, "username"),
		Password:      strCfg(m.config, "password"),
		SecurityToken: strCfg(m.config, "security_token"),
		AccessToken:   strCfg(m.config, "accessToken"),
		InstanceURL:   strCfg(m.config, "instanceUrl"),
		APIVersion:    strCfg(m.config, "apiVersion"),
		LoginURL:      strCfg(m.config, "loginUrl"),
		Sandbox:       boolCfg(m.config, "sandbox"),
	}

	if err := provider.Connect(context.Background(), cfg); err != nil {
		return fmt.Errorf("crm.provider %q: %w", m.name, err)
	}

	RegisterProvider(m.name, provider)
	return nil
}

// Start is a no-op for this module.
func (m *crmModule) Start(_ context.Context) error { return nil }

// Stop unregisters the provider and closes it.
func (m *crmModule) Stop(_ context.Context) error {
	if p, ok := GetProvider(m.name); ok {
		_ = p.Close()
	}
	UnregisterProvider(m.name)
	return nil
}

func strCfg(config map[string]any, key string) string {
	if v, ok := config[key].(string); ok {
		return v
	}
	return ""
}

func boolCfg(config map[string]any, key string) bool {
	switch v := config[key].(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "1" || v == "yes"
	}
	return false
}
