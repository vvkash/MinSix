package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/minsix/backend/internal/database"
	"github.com/minsix/backend/internal/detector"
	"github.com/minsix/backend/internal/models"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/minsix?sslmode=disable"
	}

	// Connect to database
	db, err := database.NewDatabase(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize fraud detector
	fraudDetector := detector.NewFraudDetector(db)

	log.Println("Creating test transactions...")

	// Test Case 1: Large transfer (should be flagged)
	createTestTransaction(db, fraudDetector, &models.Transaction{
		TxHash:      "0xtest_large_transfer_123",
		BlockNumber: 18000000,
		FromAddress: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		ToAddress:   stringPtr("0x1234567890123456789012345678901234567890"),
		Value:       "50000000000000000000", // 50 ETH
		GasPrice:    "30000000000",
		GasUsed:     21000,
		Timestamp:   time.Now(),
	}, "Large transfer")

	// Test Case 2: Null address (should be flagged)
	createTestTransaction(db, fraudDetector, &models.Transaction{
		TxHash:      "0xtest_null_address_456",
		BlockNumber: 18000001,
		FromAddress: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		ToAddress:   stringPtr("0x0000000000000000000000000000000000000000"),
		Value:       "5000000000000000000", // 5 ETH
		GasPrice:    "25000000000",
		GasUsed:     21000,
		Timestamp:   time.Now(),
	}, "Null address")

	// Test Case 3: Burn address (should be flagged)
	createTestTransaction(db, fraudDetector, &models.Transaction{
		TxHash:      "0xtest_burn_address_789",
		BlockNumber: 18000002,
		FromAddress: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		ToAddress:   stringPtr("0x000000000000000000000000000000000000dead"),
		Value:       "3000000000000000000", // 3 ETH
		GasPrice:    "28000000000",
		GasUsed:     21000,
		Timestamp:   time.Now(),
	}, "Burn address")

	// Test Case 4: Rapid transactions (should be flagged after threshold)
	baseTime := time.Now()
	for i := 0; i < 7; i++ {
		createTestTransaction(db, fraudDetector, &models.Transaction{
			TxHash:      fmt.Sprintf("0xtest_rapid_%d", i),
			BlockNumber: 18000003 + int64(i),
			FromAddress: "0xRapidSender1234567890123456789012345678",
			ToAddress:   stringPtr("0x1234567890123456789012345678901234567890"),
			Value:       "1000000000000000000", // 1 ETH
			GasPrice:    "30000000000",
			GasUsed:     21000,
			Timestamp:   baseTime.Add(time.Duration(i*10) * time.Second),
		}, fmt.Sprintf("Rapid transaction #%d", i+1))
	}

	// Test Case 5: High gas price (should be flagged)
	createTestTransaction(db, fraudDetector, &models.Transaction{
		TxHash:      "0xtest_high_gas_abc",
		BlockNumber: 18000010,
		FromAddress: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		ToAddress:   stringPtr("0x1234567890123456789012345678901234567890"),
		Value:       "2000000000000000000", // 2 ETH
		GasPrice:    "200000000000",        // Very high gas
		GasUsed:     21000,
		Timestamp:   time.Now(),
	}, "High gas price")

	// Test Case 6: Token transfer (contract interaction)
	createTestTransaction(db, fraudDetector, &models.Transaction{
		TxHash:      "0xtest_token_transfer_def",
		BlockNumber: 18000011,
		FromAddress: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		ToAddress:   stringPtr("0xdAC17F958D2ee523a2206206994597C13D831ec7"), // USDT contract
		Value:       "0",
		GasPrice:    "30000000000",
		GasUsed:     65000,
		InputData:   stringPtr("0xa9059cbb000000000000000000000000742d35cc6634c0532925a3b844bc9e7595f0beb0000000000000000000000000000000000000000000000000de0b6b3a7640000"),
		Timestamp:   time.Now(),
	}, "Token transfer")

	// Test Case 7: Normal transaction (should NOT be flagged)
	createTestTransaction(db, fraudDetector, &models.Transaction{
		TxHash:      "0xtest_normal_ghi",
		BlockNumber: 18000012,
		FromAddress: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		ToAddress:   stringPtr("0x1234567890123456789012345678901234567890"),
		Value:       "500000000000000000", // 0.5 ETH
		GasPrice:    "25000000000",
		GasUsed:     21000,
		Timestamp:   time.Now(),
	}, "Normal transaction")

	// Update statistics
	db.UpdateStatistic("total_transactions", 14)
	db.IncrementStatistic("total_flagged", 0) // Will be incremented by flagged txs

	log.Println("Test data created successfully")
	log.Println("")
	log.Println("View results at http://localhost:3000")
	log.Println("Or query the database directly:")
	log.Println("  docker exec -it minsix-postgres psql -U postgres -d minsix")
	log.Println("  SELECT * FROM flagged_transactions;")
}

func createTestTransaction(db *database.DB, detector *detector.FraudDetector, tx *models.Transaction, description string) {
	// Save transaction
	if err := db.SaveTransaction(tx); err != nil {
		log.Printf("ERROR: Failed to save %s: %v", description, err)
		return
	}

	// Analyze for fraud
	flagged, err := detector.AnalyzeTransaction(tx)
	if err != nil {
		log.Printf("ERROR: Failed to analyze %s: %v", description, err)
		return
	}

	if flagged != nil {
		if err := db.FlagTransaction(flagged); err != nil {
			log.Printf("ERROR: Failed to flag %s: %v", description, err)
			return
		}
		log.Printf("FLAGGED: %s (Risk: %d, Reasons: %v)", description, flagged.RiskScore, flagged.Reasons)
	} else {
		log.Printf("OK: %s", description)
	}
}

func stringPtr(s string) *string {
	return &s
}
