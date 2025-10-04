package models

import (
	"time"
)

type Transaction struct {
	ID          int       `json:"id"`
	TxHash      string    `json:"tx_hash"`
	BlockNumber int64     `json:"block_number"`
	FromAddress string    `json:"from_address"`
	ToAddress   *string   `json:"to_address"`
	Value       string    `json:"value"`
	GasPrice    string    `json:"gas_price"`
	GasUsed     int64     `json:"gas_used"`
	InputData   *string   `json:"input_data"`
	Timestamp   time.Time `json:"timestamp"`
	CreatedAt   time.Time `json:"created_at"`
}

type FlaggedTransaction struct {
	ID            int       `json:"id"`
	TransactionID *int      `json:"transaction_id"`
	TxHash        string    `json:"tx_hash"`
	RiskScore     int       `json:"risk_score"`
	Reasons       []string  `json:"reasons"`
	FlaggedAt     time.Time `json:"flagged_at"`
	Status        string    `json:"status"`
	Transaction   *Transaction `json:"transaction,omitempty"`
}

type BlacklistedAddress struct {
	ID      int       `json:"id"`
	Address string    `json:"address"`
	Reason  string    `json:"reason"`
	Source  string    `json:"source"`
	AddedAt time.Time `json:"added_at"`
}

type MonitoredWallet struct {
	ID          int        `json:"id"`
	Address     string     `json:"address"`
	Label       *string    `json:"label"`
	AddedAt     time.Time  `json:"added_at"`
	LastChecked *time.Time `json:"last_checked"`
}

type Statistics struct {
	ID          int       `json:"id"`
	MetricName  string    `json:"metric_name"`
	MetricValue float64   `json:"metric_value"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type AlertPayload struct {
	TxHash    string   `json:"tx_hash"`
	RiskScore int      `json:"risk_score"`
	Reasons   []string `json:"reasons"`
	Timestamp time.Time `json:"timestamp"`
}
