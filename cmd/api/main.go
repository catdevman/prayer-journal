package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var chiLambda *chiadapter.ChiLambda
var r *chi.Mux

func init() {
	r = chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// If NOT running in Lambda, start local server
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		slog.Info("🚀 Local server starting on http://localhost:8080")
		if err := http.ListenAndServe(":8080", r); err != nil {
			slog.Error("Local server failed", "error", err)
			os.Exit(1)
		}
		return
	}

	lambda.Start(Handler)
}
