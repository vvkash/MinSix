package detector

import (
	"testing"
	"time"

	"github.com/minsix/backend/internal/models"
)

func TestCheckLargeTransfer(t *testing.T) {
	fd := &FraudDetector{}

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "Small transfer",
			value:    "1000000000000000000", // 1 ETH
			expected: false,
		},
		{
			name:     "Large transfer",
			value:    "15000000000000000000", // 15 ETH
			expected: true,
		},
		{
			name:     "Very large transfer",
			value:    "100000000000000000000", // 100 ETH
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &models.Transaction{
				Value: tt.value,
			}
			result, _ := fd.checkLargeTransfer(tx)
			if result != tt.expected {
				t.Errorf("checkLargeTransfer() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCheckNullAddress(t *testing.T) {
	fd := &FraudDetector{}

	tests := []struct {
		name      string
		toAddress *string
		expected  bool
	}{
		{
			name:      "Null address",
			toAddress: stringPtr("0x0000000000000000000000000000000000000000"),
			expected:  true,
		},
		{
			name:      "Burn address",
			toAddress: stringPtr("0x000000000000000000000000000000000000dead"),
			expected:  true,
		},
		{
			name:      "Normal address",
			toAddress: stringPtr("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"),
			expected:  false,
		},
		{
			name:      "No address (contract creation)",
			toAddress: nil,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &models.Transaction{
				ToAddress: tt.toAddress,
			}
			result := fd.checkNullAddress(tx)
			if result != tt.expected {
				t.Errorf("checkNullAddress() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCheckRapidTransactions(t *testing.T) {
	fd := &FraudDetector{
		recentTxs: make(map[string][]time.Time),
	}

	address := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
	now := time.Now()

	// Add transactions rapidly
	for i := 0; i < 6; i++ {
		tx := &models.Transaction{
			FromAddress: address,
			Timestamp:   now.Add(time.Duration(i) * time.Second),
		}
		result := fd.checkRapidTransactions(tx)
		
		// Should only flag after exceeding threshold
		if i <= MaxRapidTransactions && result {
			t.Errorf("Transaction %d: expected false, got true", i)
		}
		if i > MaxRapidTransactions && !result {
			t.Errorf("Transaction %d: expected true, got false", i)
		}
	}
}

func TestCheckContractInteraction(t *testing.T) {
	fd := &FraudDetector{}

	tests := []struct {
		name      string
		inputData *string
		expected  bool
	}{
		{
			name:      "No input data",
			inputData: nil,
			expected:  false,
		},
		{
			name:      "Simple transfer",
			inputData: stringPtr("0x"),
			expected:  false,
		},
		{
			name:      "Transfer function with sufficient data",
			inputData: stringPtr("0xa9059cbb000000000000000000000000742d35cc6634c0532925a3b844bc9e7595f0beb0000000000000000000000000000000000000000000000000de0b6b3a7640000"),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toAddress := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
			tx := &models.Transaction{
				ToAddress: &toAddress,
				InputData: tt.inputData,
			}
			result := fd.checkContractInteraction(tx)
			if result != tt.expected {
				t.Errorf("checkContractInteraction() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
