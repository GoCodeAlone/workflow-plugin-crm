package internal

import "sync"

var (
	providerMu       sync.RWMutex
	providerRegistry = make(map[string]CRMProvider)
)

// RegisterProvider adds a CRM provider to the global registry.
func RegisterProvider(name string, p CRMProvider) {
	providerMu.Lock()
	defer providerMu.Unlock()
	providerRegistry[name] = p
}

// GetProvider looks up a CRM provider by name.
func GetProvider(name string) (CRMProvider, bool) {
	providerMu.RLock()
	defer providerMu.RUnlock()
	p, ok := providerRegistry[name]
	return p, ok
}

// UnregisterProvider removes a provider from the registry.
func UnregisterProvider(name string) {
	providerMu.Lock()
	defer providerMu.Unlock()
	delete(providerRegistry, name)
}
