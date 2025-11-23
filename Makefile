.PHONY: build-api deploy gen dev dev-be dev-fe setup

# --- Setup & Helpers ---

setup:
	go mod tidy
	cd web && npm install
	@echo "✅ Setup complete"

gen:
	# Assumes tygo is installed: go install github.com/gzuidhof/tygo@latest
	tygo generate

# --- Development ---

dev-be:
	# Runs the Go backend locally on port 8080
	go run cmd/api/main.go

dev-fe:
	# Runs the Vue frontend locally (usually port 5173)
	cd web && npm run dev

dev:
	# Runs both in parallel using make's -j flag
	# Note: Logs will be interleaved. Use Ctrl+C to stop both.
	@make -j2 dev-be dev-fe

# --- Deployment ---

build-api:
	GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go

deploy: build-api
	cd web && npm run build
	cd infra && cdk deploy
