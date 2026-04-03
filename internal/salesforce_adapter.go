package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	sfprovider "github.com/GoCodeAlone/workflow-plugin-salesforce/salesforce"
)

const defaultAPIVersion = "v63.0"

// salesforceAdapter implements CRMProvider by delegating to the
// Salesforce plugin's exported Provider for authentication and then
// making REST API calls with a standard HTTP client.
type salesforceAdapter struct {
	instanceURL string
	accessToken string
	apiVersion  string
	httpClient  *http.Client
}

func (a *salesforceAdapter) Connect(ctx context.Context, cfg ProviderConfig) error {
	apiVersion := cfg.APIVersion
	if apiVersion == "" {
		apiVersion = defaultAPIVersion
	}
	if !strings.HasPrefix(apiVersion, "v") {
		apiVersion = "v" + apiVersion
	}
	a.apiVersion = apiVersion
	a.httpClient = &http.Client{Timeout: 30 * time.Second}

	// Direct access token — skip OAuth flow.
	if cfg.AccessToken != "" {
		if cfg.InstanceURL == "" {
			return fmt.Errorf("crm/salesforce: instanceURL is required when using accessToken")
		}
		a.instanceURL = strings.TrimRight(cfg.InstanceURL, "/")
		a.accessToken = cfg.AccessToken
		return nil
	}

	// Delegate authentication to the salesforce plugin's provider.
	sfCfg := sfprovider.Config{
		AuthType:      cfg.AuthType,
		ClientID:      cfg.ClientID,
		ClientSecret:  cfg.ClientSecret,
		RefreshToken:  cfg.RefreshToken,
		Username:      cfg.Username,
		Password:      cfg.Password,
		SecurityToken: cfg.SecurityToken,
		AccessToken:   cfg.AccessToken,
		InstanceURL:   cfg.InstanceURL,
		LoginURL:      cfg.LoginURL,
		APIVersion:    cfg.APIVersion,
		Sandbox:       cfg.Sandbox,
	}

	provider, err := sfprovider.NewProvider(ctx, sfCfg)
	if err != nil {
		return fmt.Errorf("crm/salesforce: auth failed: %w", err)
	}

	a.instanceURL = strings.TrimRight(provider.Client.InstanceURL(), "/")
	if t := provider.Client.GetToken(); t != nil {
		a.accessToken = t.AccessToken
	}
	return nil
}

// versionedURL returns the full URL for a versioned REST path.
func (a *salesforceAdapter) versionedURL(path string) string {
	return fmt.Sprintf("%s/services/data/%s%s", a.instanceURL, a.apiVersion, path)
}

// doJSON makes an authenticated REST call and decodes the response.
func (a *salesforceAdapter) doJSON(ctx context.Context, method, fullURL string, body any) (map[string]any, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("crm/salesforce: marshal body: %w", err)
		}
		bodyReader = strings.NewReader(string(data))
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+a.accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("crm/salesforce: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNoContent {
		return map[string]any{"success": true}, nil
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("crm/salesforce: HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	if len(respBody) == 0 {
		return map[string]any{"success": true}, nil
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		// Salesforce sometimes returns arrays; wrap if needed.
		var arr []any
		if json.Unmarshal(respBody, &arr) == nil {
			return map[string]any{"records": arr}, nil
		}
		return nil, fmt.Errorf("crm/salesforce: decode: %w", err)
	}
	return result, nil
}

func (a *salesforceAdapter) CreateRecord(ctx context.Context, objectType string, fields map[string]any) (*RecordResult, error) {
	u := a.versionedURL(fmt.Sprintf("/sobjects/%s", url.PathEscape(objectType)))
	result, err := a.doJSON(ctx, http.MethodPost, u, fields)
	if err != nil {
		return nil, err
	}
	rr := &RecordResult{
		ID:      fmt.Sprintf("%v", result["id"]),
		Success: result["success"] == true,
	}
	if errs, ok := result["errors"].([]any); ok {
		for _, e := range errs {
			rr.Errors = append(rr.Errors, fmt.Sprintf("%v", e))
		}
	}
	return rr, nil
}

func (a *salesforceAdapter) GetRecord(ctx context.Context, objectType, id string) (map[string]any, error) {
	u := a.versionedURL(fmt.Sprintf("/sobjects/%s/%s", url.PathEscape(objectType), url.PathEscape(id)))
	return a.doJSON(ctx, http.MethodGet, u, nil)
}

func (a *salesforceAdapter) UpdateRecord(ctx context.Context, objectType, id string, fields map[string]any) error {
	u := a.versionedURL(fmt.Sprintf("/sobjects/%s/%s", url.PathEscape(objectType), url.PathEscape(id)))
	_, err := a.doJSON(ctx, http.MethodPatch, u, fields)
	return err
}

func (a *salesforceAdapter) UpsertRecord(ctx context.Context, objectType, extField, extValue string, fields map[string]any) (*RecordResult, error) {
	u := a.versionedURL(fmt.Sprintf("/sobjects/%s/%s/%s", url.PathEscape(objectType), url.PathEscape(extField), url.PathEscape(extValue)))
	result, err := a.doJSON(ctx, http.MethodPatch, u, fields)
	if err != nil {
		return nil, err
	}
	rr := &RecordResult{Success: true}
	if id, ok := result["id"]; ok {
		rr.ID = fmt.Sprintf("%v", id)
	}
	return rr, nil
}

func (a *salesforceAdapter) DeleteRecord(ctx context.Context, objectType, id string) error {
	u := a.versionedURL(fmt.Sprintf("/sobjects/%s/%s", url.PathEscape(objectType), url.PathEscape(id)))
	_, err := a.doJSON(ctx, http.MethodDelete, u, nil)
	return err
}

func (a *salesforceAdapter) Query(ctx context.Context, query string) (*QueryResult, error) {
	u := a.versionedURL("/query?q=" + url.QueryEscape(query))
	result, err := a.doJSON(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	qr := &QueryResult{}
	if recs, ok := result["records"].([]any); ok {
		for _, r := range recs {
			if m, ok := r.(map[string]any); ok {
				qr.Records = append(qr.Records, m)
			}
		}
	}
	if ts, ok := result["totalSize"].(float64); ok {
		qr.TotalSize = int(ts)
	}
	if d, ok := result["done"].(bool); ok {
		qr.Done = d
	}
	return qr, nil
}

func (a *salesforceAdapter) Search(ctx context.Context, query string) (*SearchResult, error) {
	u := a.versionedURL("/search?q=" + url.QueryEscape(query))
	result, err := a.doJSON(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	sr := &SearchResult{}
	if recs, ok := result["searchRecords"].([]any); ok {
		for _, r := range recs {
			if m, ok := r.(map[string]any); ok {
				sr.Results = append(sr.Results, m)
			}
		}
	} else if recs, ok := result["records"].([]any); ok {
		for _, r := range recs {
			if m, ok := r.(map[string]any); ok {
				sr.Results = append(sr.Results, m)
			}
		}
	}
	return sr, nil
}

func (a *salesforceAdapter) BulkOperation(ctx context.Context, op BulkOp) (*BulkResult, error) {
	if len(op.Records) == 0 {
		return nil, fmt.Errorf("crm/salesforce: bulk operation requires at least one record")
	}

	// 1. Create the Bulk API v2 ingest job.
	body := map[string]any{
		"object":              op.ObjectType,
		"operation":           op.Operation,
		"contentType":         "CSV",
		"lineEnding":          "LF",
		"columnDelimiter":     "COMMA",
	}
	u := a.versionedURL("/jobs/ingest")
	result, err := a.doJSON(ctx, http.MethodPost, u, body)
	if err != nil {
		return nil, err
	}
	jobID := fmt.Sprintf("%v", result["id"])

	// 2. Serialize records to CSV and upload.
	csvData, err := recordsToCSV(op.Records)
	if err != nil {
		return nil, fmt.Errorf("crm/salesforce: csv encode: %w", err)
	}

	batchURL := a.versionedURL(fmt.Sprintf("/jobs/ingest/%s/batches", url.PathEscape(jobID)))
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, batchURL, bytes.NewReader(csvData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+a.accessToken)
	req.Header.Set("Content-Type", "text/csv")
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("crm/salesforce: upload csv: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("crm/salesforce: upload csv: HTTP %d", resp.StatusCode)
	}

	// 3. Close the job by setting state to UploadComplete.
	closeURL := a.versionedURL(fmt.Sprintf("/jobs/ingest/%s", url.PathEscape(jobID)))
	result, err = a.doJSON(ctx, http.MethodPatch, closeURL, map[string]any{"state": "UploadComplete"})
	if err != nil {
		return nil, err
	}

	br := &BulkResult{
		JobID: jobID,
		State: fmt.Sprintf("%v", result["state"]),
	}
	if rp, ok := result["numberRecordsProcessed"].(float64); ok {
		br.RecordsProcessed = int(rp)
	}
	if rf, ok := result["numberRecordsFailed"].(float64); ok {
		br.RecordsFailed = int(rf)
	}
	return br, nil
}

// recordsToCSV converts a slice of record maps to CSV bytes. Column order
// is derived from the first record; missing values become empty strings.
func recordsToCSV(records []map[string]any) ([]byte, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("no records")
	}
	// Stable column list from first record.
	var cols []string
	for k := range records[0] {
		cols = append(cols, k)
	}

	var buf bytes.Buffer
	// Header row.
	for i, c := range cols {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(csvEscape(c))
	}
	buf.WriteByte('\n')
	// Data rows.
	for _, rec := range records {
		for i, c := range cols {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(csvEscape(fmt.Sprintf("%v", rec[c])))
		}
		buf.WriteByte('\n')
	}
	return buf.Bytes(), nil
}

// csvEscape wraps a value in double-quotes if it contains comma, quote, or newline.
func csvEscape(s string) string {
	if strings.ContainsAny(s, ",\"\n\r") {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}

func (a *salesforceAdapter) DescribeObject(ctx context.Context, objectType string) (map[string]any, error) {
	u := a.versionedURL(fmt.Sprintf("/sobjects/%s/describe", url.PathEscape(objectType)))
	return a.doJSON(ctx, http.MethodGet, u, nil)
}

func (a *salesforceAdapter) GetLimits(ctx context.Context) (map[string]any, error) {
	u := a.versionedURL("/limits")
	return a.doJSON(ctx, http.MethodGet, u, nil)
}

func (a *salesforceAdapter) Close() error {
	return nil
}
