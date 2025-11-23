.PHONY: build-api deploy gen

build-api:
	GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go

gen:
	# Assumes tygo is installed: go install github.com/gzuidhof/tygo@latest
	tygo generate

deploy: build-api
	cd web && npm run build
	cd infra && cdk deploy
