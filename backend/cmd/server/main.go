package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/minsix/backend/internal/database"
	"github.com/minsix/backend/internal/detector"
	"github.com/minsix/backend/internal/ethereum"
	"github.com/minsix/backend/internal/handlers"
	"github.com/minsix/backend/internal/models"
	"github.com/minsix/backend/internal/websocket"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get configuration
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/minsix?sslmode=disable")
	alchemyKey := getEnv("ALCHEMY_API_KEY", "")
	alchemyNetwork := getEnv("ALCHEMY_NETWORK", "eth-mainnet")
	port := getEnv("PORT", "8080")

	if alchemyKey == "" {
		log.Fatal("ALCHEMY_API_KEY is required")
	}

	// Initialize database
	db, err := database.NewDatabase(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	migrationSQL := readMigrationFile("migrations/001_initial_schema.sql")
	if err := db.RunMigrations(migrationSQL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize fraud detector
	fraudDetector := detector.NewFraudDetector(db)

	// Initialize Ethereum client
	ethClient, err := ethereum.NewClient(alchemyKey, alchemyNetwork)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum: %v", err)
	}
	defer ethClient.Close()

	// Set up transaction handler
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transactionHandler := func(tx *models.Transaction) {
		// Save transaction to database
		if err := db.SaveTransaction(tx); err != nil {
			log.Printf("Failed to save transaction: %v", err)
			return
		}

		// Increment total transactions
		db.IncrementStatistic("total_transactions", 1)

		// Analyze for fraud
		flagged, err := fraudDetector.AnalyzeTransaction(tx)
		if err != nil {
			log.Printf("Failed to analyze transaction: %v", err)
			return
		}

		// If flagged, save and broadcast
		if flagged != nil {
			if err := db.FlagTransaction(flagged); err != nil {
				log.Printf("Failed to flag transaction: %v", err)
				return
			}

			// Increment flagged count
			db.IncrementStatistic("total_flagged", 1)

			// Broadcast alert
			hub.BroadcastAlert(flagged)
			log.Printf("WARNING: Flagged transaction %s (Risk: %d)", flagged.TxHash, flagged.RiskScore)
		}

		// Broadcast transaction update
		hub.BroadcastTransaction(tx)
	}

	// Subscribe to new blocks
	if err := ethClient.SubscribeToBlocks(ctx, transactionHandler); err != nil {
		log.Fatalf("Failed to subscribe to blocks: %v", err)
	}

	// Set up HTTP server
	router := mux.NewRouter()
	handler := handlers.NewHandler(db, hub)

	// API routes
	router.HandleFunc("/api/health", handler.HealthCheck).Methods("GET")
	router.HandleFunc("/api/transactions", handler.GetFlaggedTransactions).Methods("GET")
	router.HandleFunc("/api/wallets/{address}", handler.GetWalletAnalysis).Methods("GET")
	router.HandleFunc("/api/stats", handler.GetStatistics).Methods("GET")
	router.HandleFunc("/ws", handler.HandleWebSocket)

	// CORS configuration
	corsOrigins := getEnv("CORS_ORIGINS", "http://localhost:3000")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{corsOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Start server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      c.Handler(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func readMigrationFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}
	return string(data)
}
