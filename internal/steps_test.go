package internal

import (
	"context"
	"testing"
)

// mockProvider implements CRMProvider for unit testing.
type mockProvider struct {
	createRecordFn   func(ctx context.Context, objectType string, fields map[string]any) (*RecordResult, error)
	getRecordFn      func(ctx context.Context, objectType, id string) (map[string]any, error)
	updateRecordFn   func(ctx context.Context, objectType, id string, fields map[string]any) error
	upsertRecordFn   func(ctx context.Context, objectType, extField, extValue string, fields map[string]any) (*RecordResult, error)
	deleteRecordFn   func(ctx context.Context, objectType, id string) error
	queryFn          func(ctx context.Context, query string) (*QueryResult, error)
	searchFn         func(ctx context.Context, query string) (*SearchResult, error)
	bulkOperationFn  func(ctx context.Context, op BulkOp) (*BulkResult, error)
	describeObjectFn func(ctx context.Context, objectType string) (map[string]any, error)
	getLimitsFn      func(ctx context.Context) (map[string]any, error)
}

func (m *mockProvider) Connect(_ context.Context, _ ProviderConfig) error { return nil }
func (m *mockProvider) Close() error                                      { return nil }

func (m *mockProvider) CreateRecord(ctx context.Context, objectType string, fields map[string]any) (*RecordResult, error) {
	if m.createRecordFn != nil {
		return m.createRecordFn(ctx, objectType, fields)
	}
	return &RecordResult{ID: "mock-id", Success: true}, nil
}

func (m *mockProvider) GetRecord(ctx context.Context, objectType, id string) (map[string]any, error) {
	if m.getRecordFn != nil {
		return m.getRecordFn(ctx, objectType, id)
	}
	return map[string]any{"Id": id, "Name": "Mock Record"}, nil
}

func (m *mockProvider) UpdateRecord(ctx context.Context, objectType, id string, fields map[string]any) error {
	if m.updateRecordFn != nil {
		return m.updateRecordFn(ctx, objectType, id, fields)
	}
	return nil
}

func (m *mockProvider) UpsertRecord(ctx context.Context, objectType, extField, extValue string, fields map[string]any) (*RecordResult, error) {
	if m.upsertRecordFn != nil {
		return m.upsertRecordFn(ctx, objectType, extField, extValue, fields)
	}
	return &RecordResult{ID: "upsert-id", Success: true}, nil
}

func (m *mockProvider) DeleteRecord(ctx context.Context, objectType, id string) error {
	if m.deleteRecordFn != nil {
		return m.deleteRecordFn(ctx, objectType, id)
	}
	return nil
}

func (m *mockProvider) Query(ctx context.Context, query string) (*QueryResult, error) {
	if m.queryFn != nil {
		return m.queryFn(ctx, query)
	}
	return &QueryResult{
		Records:   []map[string]any{{"Id": "001a"}, {"Id": "001b"}},
		TotalSize: 2,
		Done:      true,
	}, nil
}

func (m *mockProvider) Search(ctx context.Context, query string) (*SearchResult, error) {
	if m.searchFn != nil {
		return m.searchFn(ctx, query)
	}
	return &SearchResult{Results: []map[string]any{{"Id": "001a"}}}, nil
}

func (m *mockProvider) BulkOperation(ctx context.Context, op BulkOp) (*BulkResult, error) {
	if m.bulkOperationFn != nil {
		return m.bulkOperationFn(ctx, op)
	}
	return &BulkResult{JobID: "job-1", State: "Open", RecordsProcessed: 0, RecordsFailed: 0}, nil
}

func (m *mockProvider) DescribeObject(ctx context.Context, objectType string) (map[string]any, error) {
	if m.describeObjectFn != nil {
		return m.describeObjectFn(ctx, objectType)
	}
	return map[string]any{"name": objectType, "fields": []any{}}, nil
}

func (m *mockProvider) GetLimits(ctx context.Context) (map[string]any, error) {
	if m.getLimitsFn != nil {
		return m.getLimitsFn(ctx)
	}
	return map[string]any{"DailyApiRequests": map[string]any{"Max": 15000, "Remaining": 14500}}, nil
}

// registerMock is a test helper that registers a mock and returns a cleanup func.
func registerMock(name string, mock *mockProvider) func() {
	RegisterProvider(name, mock)
	return func() { UnregisterProvider(name) }
}

// --- Step tests using mock provider ---

func TestCreateRecordStep_MissingProvider(t *testing.T) {
	step, _ := newCreateRecordStep("test", map[string]any{"module": "nonexistent"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{"object_type": "Account"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing provider")
	}
}

func TestCreateRecordStep_MissingObjectType(t *testing.T) {
	cleanup := registerMock("test-create-no-type", &mockProvider{})
	defer cleanup()

	step, _ := newCreateRecordStep("test", map[string]any{"module": "test-create-no-type"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing object_type")
	}
}

func TestCreateRecordStep_Success(t *testing.T) {
	cleanup := registerMock("test-create-ok", &mockProvider{})
	defer cleanup()

	step, _ := newCreateRecordStep("test", map[string]any{"module": "test-create-ok"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{
		"object_type": "Account",
		"fields":      map[string]any{"Name": "Acme Corp"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] != nil {
		t.Fatalf("unexpected error: %v", result.Output["error"])
	}
	if result.Output["id"] != "mock-id" {
		t.Errorf("expected id=mock-id, got %v", result.Output["id"])
	}
	if result.Output["success"] != true {
		t.Errorf("expected success=true, got %v", result.Output["success"])
	}
}

func TestGetRecordStep_MissingRecordID(t *testing.T) {
	cleanup := registerMock("test-get-noid", &mockProvider{})
	defer cleanup()

	step, _ := newGetRecordStep("test", map[string]any{"module": "test-get-noid"})
	result, _ := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{"object_type": "Account"})
	if result.Output["error"] == nil {
		t.Error("expected error for missing record_id")
	}
}

func TestGetRecordStep_Success(t *testing.T) {
	cleanup := registerMock("test-get-ok", &mockProvider{})
	defer cleanup()

	step, _ := newGetRecordStep("test", map[string]any{"module": "test-get-ok"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{
		"object_type": "Account",
		"record_id":   "001xx",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["Id"] != "001xx" {
		t.Errorf("expected Id=001xx, got %v", result.Output["Id"])
	}
}

func TestUpdateRecordStep_Success(t *testing.T) {
	cleanup := registerMock("test-update-ok", &mockProvider{})
	defer cleanup()

	step, _ := newUpdateRecordStep("test", map[string]any{"module": "test-update-ok"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{
		"object_type": "Account",
		"record_id":   "001xx",
		"fields":      map[string]any{"Name": "Updated"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["success"] != true {
		t.Errorf("expected success=true, got %v", result.Output["success"])
	}
}

func TestUpsertRecordStep_MissingExtField(t *testing.T) {
	cleanup := registerMock("test-upsert-nofield", &mockProvider{})
	defer cleanup()

	step, _ := newUpsertRecordStep("test", map[string]any{"module": "test-upsert-nofield"})
	result, _ := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{
		"object_type": "Account",
	})
	if result.Output["error"] == nil {
		t.Error("expected error for missing external_id_field")
	}
}

func TestUpsertRecordStep_Success(t *testing.T) {
	cleanup := registerMock("test-upsert-ok", &mockProvider{})
	defer cleanup()

	step, _ := newUpsertRecordStep("test", map[string]any{"module": "test-upsert-ok"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{
		"object_type":       "Account",
		"external_id_field": "ExtId__c",
		"external_id_value": "ext-123",
		"fields":            map[string]any{"Name": "Upserted"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["id"] != "upsert-id" {
		t.Errorf("expected id=upsert-id, got %v", result.Output["id"])
	}
}

func TestDeleteRecordStep_Success(t *testing.T) {
	cleanup := registerMock("test-delete-ok", &mockProvider{})
	defer cleanup()

	step, _ := newDeleteRecordStep("test", map[string]any{"module": "test-delete-ok"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{
		"object_type": "Account",
		"record_id":   "001xx",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["success"] != true {
		t.Errorf("expected success=true, got %v", result.Output["success"])
	}
}

func TestQueryStep_MissingQuery(t *testing.T) {
	cleanup := registerMock("test-query-noq", &mockProvider{})
	defer cleanup()

	step, _ := newQueryStep("test", map[string]any{"module": "test-query-noq"})
	result, _ := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if result.Output["error"] == nil {
		t.Error("expected error for missing query")
	}
}

func TestQueryStep_Success(t *testing.T) {
	cleanup := registerMock("test-query-ok", &mockProvider{})
	defer cleanup()

	step, _ := newQueryStep("test", map[string]any{"module": "test-query-ok"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{
		"query": "SELECT Id FROM Account",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["total_size"].(int) != 2 {
		t.Errorf("expected total_size=2, got %v", result.Output["total_size"])
	}
	if result.Output["done"].(bool) != true {
		t.Error("expected done=true")
	}
}

func TestSearchStep_Success(t *testing.T) {
	cleanup := registerMock("test-search-ok", &mockProvider{})
	defer cleanup()

	step, _ := newSearchStep("test", map[string]any{"module": "test-search-ok"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{
		"query": "FIND {Acme}",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["count"].(int) != 1 {
		t.Errorf("expected count=1, got %v", result.Output["count"])
	}
}

func TestBulkImportStep_MissingObjectType(t *testing.T) {
	cleanup := registerMock("test-bulk-noobj", &mockProvider{})
	defer cleanup()

	step, _ := newBulkImportStep("test", map[string]any{"module": "test-bulk-noobj"})
	result, _ := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if result.Output["error"] == nil {
		t.Error("expected error for missing object_type")
	}
}

func TestBulkImportStep_Success(t *testing.T) {
	cleanup := registerMock("test-bulk-ok", &mockProvider{})
	defer cleanup()

	step, _ := newBulkImportStep("test", map[string]any{"module": "test-bulk-ok"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{
		"object_type": "Account",
		"operation":   "insert",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["job_id"] != "job-1" {
		t.Errorf("expected job_id=job-1, got %v", result.Output["job_id"])
	}
}

func TestDescribeObjectStep_Success(t *testing.T) {
	cleanup := registerMock("test-describe-ok", &mockProvider{})
	defer cleanup()

	step, _ := newDescribeObjectStep("test", map[string]any{"module": "test-describe-ok"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{
		"object_type": "Account",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["name"] != "Account" {
		t.Errorf("expected name=Account, got %v", result.Output["name"])
	}
}

func TestGetLimitsStep_Success(t *testing.T) {
	cleanup := registerMock("test-limits-ok", &mockProvider{})
	defer cleanup()

	step, _ := newGetLimitsStep("test", map[string]any{"module": "test-limits-ok"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["DailyApiRequests"] == nil {
		t.Error("expected DailyApiRequests in output")
	}
}

func TestStepRegistry_AllTypesConstructible(t *testing.T) {
	for typeName := range stepRegistry {
		_, err := createStep(typeName, "test-"+typeName, map[string]any{})
		if err != nil {
			t.Errorf("createStep(%q): unexpected error: %v", typeName, err)
		}
	}
}

func TestStepRegistry_UnknownType(t *testing.T) {
	_, err := createStep("step.crm_unknown_xyz", "test", map[string]any{})
	if err == nil {
		t.Error("expected error for unknown step type")
	}
}

func TestModuleStopUnregistersProvider(t *testing.T) {
	mock := &mockProvider{}
	RegisterProvider("mod-stop-test", mock)

	m := &crmModule{name: "mod-stop-test"}
	if err := m.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok := GetProvider("mod-stop-test"); ok {
		t.Error("expected provider to be unregistered after Stop")
	}
}

func TestHelpers_ResolveValue(t *testing.T) {
	current := map[string]any{"key": "from_current"}
	config := map[string]any{"key": "from_config", "other": "cfg_only"}

	if v := resolveValue("key", current, config); v != "from_current" {
		t.Errorf("expected from_current, got %q", v)
	}
	if v := resolveValue("other", current, config); v != "cfg_only" {
		t.Errorf("expected cfg_only, got %q", v)
	}
	if v := resolveValue("missing", current, config); v != "" {
		t.Errorf("expected empty, got %q", v)
	}
}

func TestHelpers_ResolveMap(t *testing.T) {
	fields := map[string]any{"Name": "Test"}
	current := map[string]any{"fields": fields}
	config := map[string]any{}

	result := resolveMap("fields", current, config)
	if result["Name"] != "Test" {
		t.Errorf("expected Name=Test, got %v", result["Name"])
	}

	result = resolveMap("missing", current, config)
	if result != nil {
		t.Errorf("expected nil for missing key, got %v", result)
	}
}

func TestHelpers_GetModuleName(t *testing.T) {
	if name := getModuleName(map[string]any{"module": "custom"}); name != "custom" {
		t.Errorf("expected custom, got %q", name)
	}
	if name := getModuleName(map[string]any{}); name != "crm" {
		t.Errorf("expected crm, got %q", name)
	}
}
