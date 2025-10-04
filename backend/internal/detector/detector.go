package detector

import (
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/minsix/backend/internal/database"
	"github.com/minsix/backend/internal/models"
)

const (
	// Thresholds for fraud detection
	LargeTransferThresholdETH = 10.0  // ETH
	HighGasPriceMultiplier    = 3.0   // 3x average
	RapidTransactionWindow    = 60    // seconds
	MaxRapidTransactions      = 5     // transactions in window
)

type FraudDetector struct {
	db                *database.DB
	recentTxs         map[string][]time.Time // address -> timestamps
	averageGasPrice   *big.Int
}

func NewFraudDetector(db *database.DB) *FraudDetector {
	return &FraudDetector{
		db:              db,
		recentTxs:       make(map[string][]time.Time),
		averageGasPrice: big.NewInt(30000000000), // 30 Gwei default
	}
}

// AnalyzeTransaction runs all fraud detection heuristics
func (fd *FraudDetector) AnalyzeTransaction(tx *models.Transaction) (*models.FlaggedTransaction, error) {
	reasons := []string{}
	riskScore := 0

	// Heuristic 1: Check for blacklisted addresses
	if blacklisted, err := fd.checkBlacklist(tx); err != nil {
		log.Printf("Error checking blacklist: %v", err)
	} else if blacklisted {
		reasons = append(reasons, "Blacklisted address detected")
		riskScore += 40
	}

	// Heuristic 2: Large transfer detection
	if isLarge, amount := fd.checkLargeTransfer(tx); isLarge {
		reasons = append(reasons, fmt.Sprintf("Large transfer: %s ETH", amount))
		riskScore += 25
	}

	// Heuristic 3: Unusual timing patterns
	if rapidFire := fd.checkRapidTransactions(tx); rapidFire {
		reasons = append(reasons, "Rapid succession of transactions detected")
		riskScore += 20
	}

	// Heuristic 4: Gas price anomaly
	if unusual := fd.checkUnusualGasPrice(tx); unusual {
		reasons = append(reasons, "Unusual gas price detected")
		riskScore += 15
	}

	// Heuristic 5: Contract interaction patterns
	if suspicious := fd.checkContractInteraction(tx); suspicious {
		reasons = append(reasons, "Suspicious contract interaction")
		riskScore += 20
	}

	// Heuristic 6: Null or burn address
	if nullAddress := fd.checkNullAddress(tx); nullAddress {
		reasons = append(reasons, "Transaction to null/burn address")
		riskScore += 30
	}

	// Only flag if risk score is above threshold
	if riskScore >= 20 {
		flagged := &models.FlaggedTransaction{
			TxHash:    tx.TxHash,
			RiskScore: min(riskScore, 100),
			Reasons:   reasons,
			Status:    "pending",
		}

		if tx.ID > 0 {
			flagged.TransactionID = &tx.ID
		}

		return flagged, nil
	}

	return nil, nil
}

// checkBlacklist verifies if addresses are blacklisted
func (fd *FraudDetector) checkBlacklist(tx *models.Transaction) (bool, error) {
	fromBlacklisted, err := fd.db.IsBlacklisted(tx.FromAddress)
	if err != nil {
		return false, err
	}
	if fromBlacklisted {
		return true, nil
	}

	if tx.ToAddress != nil {
		toBlacklisted, err := fd.db.IsBlacklisted(*tx.ToAddress)
		if err != nil {
			return false, err
		}
		if toBlacklisted {
			return true, nil
		}
	}

	return false, nil
}

// checkLargeTransfer detects unusually large transfers
func (fd *FraudDetector) checkLargeTransfer(tx *models.Transaction) (bool, string) {
	value := new(big.Int)
	value.SetString(tx.Value, 10)

	// Convert to ETH
	ethValue := new(big.Float).Quo(
		new(big.Float).SetInt(value),
		big.NewFloat(1e18),
	)

	threshold := big.NewFloat(LargeTransferThresholdETH)
	if ethValue.Cmp(threshold) > 0 {
		return true, ethValue.Text('f', 4)
	}

	return false, ""
}

// checkRapidTransactions detects rapid succession of transactions
func (fd *FraudDetector) checkRapidTransactions(tx *models.Transaction) bool {
	now := tx.Timestamp
	address := tx.FromAddress

	// Clean old timestamps
	if timestamps, exists := fd.recentTxs[address]; exists {
		filtered := []time.Time{}
		for _, ts := range timestamps {
			if now.Sub(ts).Seconds() <= RapidTransactionWindow {
				filtered = append(filtered, ts)
			}
		}
		fd.recentTxs[address] = filtered
	}

	// Add current transaction
	fd.recentTxs[address] = append(fd.recentTxs[address], now)

	// Check if exceeds threshold
	return len(fd.recentTxs[address]) > MaxRapidTransactions
}

// checkUnusualGasPrice detects abnormal gas prices
func (fd *FraudDetector) checkUnusualGasPrice(tx *models.Transaction) bool {
	gasPrice := new(big.Int)
	gasPrice.SetString(tx.GasPrice, 10)

	// Check if gas price is unusually high
	threshold := new(big.Int).Mul(fd.averageGasPrice, big.NewInt(int64(HighGasPriceMultiplier)))
	if gasPrice.Cmp(threshold) > 0 {
		return true
	}

	// Check if gas price is suspiciously low (potential front-running)
	minPrice := new(big.Int).Div(fd.averageGasPrice, big.NewInt(10))
	if gasPrice.Cmp(minPrice) < 0 && gasPrice.Cmp(big.NewInt(0)) > 0 {
		return true
	}

	return false
}

// checkContractInteraction detects suspicious contract interactions
func (fd *FraudDetector) checkContractInteraction(tx *models.Transaction) bool {
	// If there's input data and it's not a simple transfer
	if tx.InputData != nil && len(*tx.InputData) > 10 {
		data := *tx.InputData
		
		// Check for common malicious patterns
		// Note: In production, this would be more sophisticated
		suspiciousPatterns := []string{
			"0xa9059cbb", // transfer (could be token draining)
			"0x095ea7b3", // approve (unlimited approvals are risky)
		}

		for _, pattern := range suspiciousPatterns {
			if strings.HasPrefix(strings.ToLower(data), pattern) {
				// Additional checks could be added here
				// For now, we'll flag high-value approvals
				if tx.ToAddress != nil && len(data) > 138 {
					return true
				}
			}
		}
	}

	return false
}

// checkNullAddress detects transactions to null or burn addresses
func (fd *FraudDetector) checkNullAddress(tx *models.Transaction) bool {
	if tx.ToAddress == nil {
		return false
	}

	nullAddresses := []string{
		"0x0000000000000000000000000000000000000000",
		"0x000000000000000000000000000000000000dead",
	}

	toAddr := strings.ToLower(*tx.ToAddress)
	for _, null := range nullAddresses {
		if toAddr == null {
			return true
		}
	}

	return false
}

// UpdateAverageGasPrice updates the rolling average gas price
func (fd *FraudDetector) UpdateAverageGasPrice(gasPrice *big.Int) {
	// Simple exponential moving average
	alpha := big.NewFloat(0.1) // Weight for new value
	oldWeight := big.NewFloat(0.9)

	oldAvg := new(big.Float).SetInt(fd.averageGasPrice)
	newVal := new(big.Float).SetInt(gasPrice)

	weighted := new(big.Float).Mul(newVal, alpha)
	oldWeighted := new(big.Float).Mul(oldAvg, oldWeight)
	
	newAvg := new(big.Float).Add(weighted, oldWeighted)
	fd.averageGasPrice, _ = newAvg.Int(nil)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
