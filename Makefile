GEN_TARGET=internal/api
SPEC=openapi.yaml
PACKAGE=api

.PHONY: generate run-server run-client verify tidy

generate:
	@echo "Generating ogen code..."
	@ogen --target $(GEN_TARGET) --package $(PACKAGE) --clean $(SPEC)

run-server:
	@go run -tags=ogen ./cmd/server

run-client:
	@go run -tags=ogen ./cmd/client

# End-to-end verification: generate, start server, run client, stop server
verify:
	@set -e; \
	echo "[verify] Generating code..."; \
	ogen --target $(GEN_TARGET) --package $(PACKAGE) --clean $(SPEC); \
	echo "[verify] Starting server..."; \
	go run -tags=ogen ./cmd/server & SERVER_PID=$$!; \
	sleep 1; \
	echo "[verify] Running client (CRUD)..."; \
	go run -tags=ogen ./cmd/client; \
	echo "[verify] Stopping server (PID $$SERVER_PID)..."; \
	kill $$SERVER_PID || true

tidy:
	go mod tidy
