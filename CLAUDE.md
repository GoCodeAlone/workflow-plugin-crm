# CLAUDE.md — Workflow Plugin CRM

Vendor-neutral CRM plugin for the GoCodeAlone/workflow engine. Salesforce adapter wraps `workflow-plugin-salesforce`'s exported provider.

## Build & Test

```sh
go build ./...
go test ./... -v -race -count=1
```

## Cross-compile for deployment

```sh
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o workflow-plugin-crm ./cmd/workflow-plugin-crm/
```

## Structure

- `cmd/workflow-plugin-crm/main.go` — Plugin entry point (calls `sdk.Serve`)
- `internal/plugin.go` — Plugin manifest, module/step providers
- `internal/crm.go` — Vendor-neutral CRMProvider interface
- `internal/salesforce_adapter.go` — Salesforce implementation of CRMProvider
- `internal/module_provider.go` — crm.provider module (lifecycle + registry)
- `internal/registry.go` — Global provider registry
- `internal/step_record.go` — CRUD step types (create/get/update/upsert/delete)
- `internal/step_query.go` — Query and search step types
- `internal/step_ops.go` — Bulk import, describe object, get limits steps
- `internal/helpers.go` — Parameter resolution helpers
- `internal/step_registry.go` — Step type → constructor dispatch
- `plugin.json` — Capability manifest for the workflow registry
- `.goreleaser.yaml` — GoReleaser v2 config for cross-platform releases

## Releasing

```sh
git tag v0.1.0
git push origin v0.1.0
```
GoReleaser builds cross-platform binaries and creates a GitHub Release automatically.
