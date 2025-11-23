#!/bin/bash

# Stop on error
set -e

echo "🚀 Initializing Prayer Journal Monorepo..."

# 1. Create Directory Structure
mkdir -p cmd/api
mkdir -p internal/api/handlers
mkdir -p internal/api/middleware
mkdir -p internal/models
mkdir -p infra

# 2. Initialize Go Modules
# We use a workspace to keep the backend and infra logic friendly but separate if needed,
# though for this size, a single module at root is often easier. 
# Let's stick to a single root module for simplicity in a monorepo.

if [ ! -f go.mod ]; then
    go mod init prayer-journal
    echo "✅ Go module initialized"
fi

# 3. Create Backend Entrypoint
cat <<EOF > cmd/api/main.go
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
)

var chiLambda *chiadapter.ChiLambda

func init() {
	r := chi.NewRouter()
	
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// TODO: Mount middleware and routes here

	chiLambda = chiadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	lambda.Start(Handler)
}
EOF

# 4. Create Infra Entrypoint (CDK)
cat <<EOF > infra/main.go
package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type PrayerJournalStackProps struct {
	awscdk.StackProps
}

func NewPrayerJournalStack(scope constructs.Construct, id string, props *PrayerJournalStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// TODO: Add S3, CloudFront, and Lambda here

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewPrayerJournalStack(app, "PrayerJournalStack", &PrayerJournalStackProps{
		StackProps: awscdk.StackProps{
			Env: &awscdk.Environment{
				Region: jsii.String("us-east-1"),
			},
		},
	})

	app.Synth(nil)
}
EOF

# 5. Initialize Vue Frontend (using npm create vue@latest non-interactively if possible, 
# but standard scaffold is safer to just stub directory)
if [ ! -d "web" ]; then
    echo "🎨 Creating Vue 3 + TypeScript scaffolding..."
    npm create vue@latest web -- --typescript --router --pinia --eslint --default
    echo "✅ Vue app created in /web"
fi

# 6. Create Makefile
cat <<EOF > Makefile
.PHONY: build-api deploy gen

build-api:
	GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go

gen:
	# Assumes tygo is installed: go install github.com/gzuidhof/tygo@latest
	tygo generate

deploy: build-api
	cd web && npm run build
	cd infra && cdk deploy
EOF

echo "🎉 Scaffolding complete. Run 'go mod tidy' to fetch dependencies."
