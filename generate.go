//go:build tools

package tools

//go:generate ogen --target internal/api --package api --clean openapi.yaml
