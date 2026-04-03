package internal

import "context"

// CRMProvider defines a vendor-neutral interface for CRM operations.
// Each backend (Salesforce, HubSpot, Dynamics, etc.) implements this
// interface so that pipeline steps remain provider-agnostic.
type CRMProvider interface {
	Connect(ctx context.Context, config ProviderConfig) error
	CreateRecord(ctx context.Context, objectType string, fields map[string]any) (*RecordResult, error)
	GetRecord(ctx context.Context, objectType, id string) (map[string]any, error)
	UpdateRecord(ctx context.Context, objectType, id string, fields map[string]any) error
	UpsertRecord(ctx context.Context, objectType, extField, extValue string, fields map[string]any) (*RecordResult, error)
	DeleteRecord(ctx context.Context, objectType, id string) error
	Query(ctx context.Context, query string) (*QueryResult, error)
	Search(ctx context.Context, query string) (*SearchResult, error)
	BulkOperation(ctx context.Context, op BulkOp) (*BulkResult, error)
	DescribeObject(ctx context.Context, objectType string) (map[string]any, error)
	GetLimits(ctx context.Context) (map[string]any, error)
	Close() error
}

// RecordResult is the outcome of a create or upsert operation.
type RecordResult struct {
	ID      string
	Success bool
	Errors  []string
}

// QueryResult holds the results of a structured query (e.g. SOQL).
type QueryResult struct {
	Records   []map[string]any
	TotalSize int
	Done      bool
}

// SearchResult holds the results of a free-text search (e.g. SOSL).
type SearchResult struct {
	Results []map[string]any
}

// BulkOp describes a bulk CRM operation.
type BulkOp struct {
	Operation  string // "insert", "update", "upsert", "delete"
	ObjectType string
	Records    []map[string]any
}

// BulkResult is the outcome of a bulk operation.
type BulkResult struct {
	JobID            string
	State            string
	RecordsProcessed int
	RecordsFailed    int
}

// ProviderConfig holds authentication and connection settings.
type ProviderConfig struct {
	Provider     string // "salesforce"
	AuthType     string
	ClientID     string
	ClientSecret string
	RefreshToken string
	Username     string
	Password     string
	AccessToken  string
	InstanceURL  string
	APIVersion   string
	Sandbox      bool
	LoginURL     string
}
