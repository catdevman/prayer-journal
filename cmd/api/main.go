package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/catdevman/prayer-journal/internal/api/handlers"
	"github.com/catdevman/prayer-journal/internal/api/middleware"
	"github.com/catdevman/prayer-journal/internal/repository"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/go-chi/chi/v5" // Assuming Chi, replace with your router of choice
)

func main() {
	ctx := context.Background()

	// 1. Initialize AWS Config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// 2. Initialize DynamoDB Client
	// If running locally, you might want to point to LocalStack or DynamoDB Local here
	ddbClient := dynamodb.NewFromConfig(cfg)

	// 3. Initialize Dependencies
	repo := repository.NewDynamoRepository(ddbClient)
	prayerHandler := handlers.NewPrayerHandler(repo)

	// 4. Setup Router
	r := chi.NewRouter()

	// Add your Auth Middleware here
	// r.Use(middleware.AuthMiddleware)
	r.Use(middleware.CorsMiddleware) // Assuming you have this from your file list

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/prayers", func(r chi.Router) {
		r.Post("/", prayerHandler.CreatePrayer)
		r.Get("/", prayerHandler.ListPrayers)
	})

	// 5. Start Server
	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		slog.Info("🚀 Local server starting on http://localhost:8080")
		if err := http.ListenAndServe(":8080", r); err != nil {
			slog.Error("Local server failed", "error", err)
			os.Exit(1)
		}
		return
	} else {
		lambda.Start(r)
	}
}
