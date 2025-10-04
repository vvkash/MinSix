package database

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/minsix/backend/internal/models"
)

// SaveTransaction inserts a new transaction
func (db *DB) SaveTransaction(tx *models.Transaction) error {
	query := `
		INSERT INTO transactions (tx_hash, block_number, from_address, to_address, value, gas_price, gas_used, input_data, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (tx_hash) DO NOTHING
		RETURNING id
	`
	err := db.QueryRow(query, tx.TxHash, tx.BlockNumber, tx.FromAddress, tx.ToAddress, tx.Value, tx.GasPrice, tx.GasUsed, tx.InputData, tx.Timestamp).Scan(&tx.ID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to save transaction: %w", err)
	}
	return nil
}

// FlagTransaction creates a flagged transaction record
func (db *DB) FlagTransaction(flag *models.FlaggedTransaction) error {
	query := `
		INSERT INTO flagged_transactions (transaction_id, tx_hash, risk_score, reasons, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, flagged_at
	`
	err := db.QueryRow(query, flag.TransactionID, flag.TxHash, flag.RiskScore, pq.Array(flag.Reasons), flag.Status).Scan(&flag.ID, &flag.FlaggedAt)
	if err != nil {
		return fmt.Errorf("failed to flag transaction: %w", err)
	}
	return nil
}

// GetFlaggedTransactions retrieves flagged transactions with optional limit
func (db *DB) GetFlaggedTransactions(limit int) ([]*models.FlaggedTransaction, error) {
	query := `
		SELECT ft.id, ft.transaction_id, ft.tx_hash, ft.risk_score, ft.reasons, ft.flagged_at, ft.status,
		       t.block_number, t.from_address, t.to_address, t.value, t.gas_price, t.timestamp
		FROM flagged_transactions ft
		LEFT JOIN transactions t ON ft.transaction_id = t.id
		ORDER BY ft.flagged_at DESC
		LIMIT $1
	`
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get flagged transactions: %w", err)
	}
	defer rows.Close()

	var results []*models.FlaggedTransaction
	for rows.Next() {
		ft := &models.FlaggedTransaction{Transaction: &models.Transaction{}}
		var reasons pq.StringArray
		err := rows.Scan(
			&ft.ID, &ft.TransactionID, &ft.TxHash, &ft.RiskScore, &reasons, &ft.FlaggedAt, &ft.Status,
			&ft.Transaction.BlockNumber, &ft.Transaction.FromAddress, &ft.Transaction.ToAddress,
			&ft.Transaction.Value, &ft.Transaction.GasPrice, &ft.Transaction.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan flagged transaction: %w", err)
		}
		ft.Reasons = reasons
		ft.Transaction.TxHash = ft.TxHash
		results = append(results, ft)
	}
	return results, nil
}

// IsBlacklisted checks if an address is blacklisted
func (db *DB) IsBlacklisted(address string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM blacklisted_addresses WHERE address = $1)`
	err := db.QueryRow(query, address).Scan(&exists)
	return exists, err
}

// GetWalletTransactions gets transactions for a specific wallet
func (db *DB) GetWalletTransactions(address string, limit int) ([]*models.Transaction, error) {
	query := `
		SELECT id, tx_hash, block_number, from_address, to_address, value, gas_price, gas_used, timestamp
		FROM transactions
		WHERE from_address = $1 OR to_address = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`
	rows, err := db.Query(query, address, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet transactions: %w", err)
	}
	defer rows.Close()

	var results []*models.Transaction
	for rows.Next() {
		tx := &models.Transaction{}
		err := rows.Scan(&tx.ID, &tx.TxHash, &tx.BlockNumber, &tx.FromAddress, &tx.ToAddress, &tx.Value, &tx.GasPrice, &tx.GasUsed, &tx.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		results = append(results, tx)
	}
	return results, nil
}

// UpdateStatistic updates a platform statistic
func (db *DB) UpdateStatistic(name string, value float64) error {
	query := `
		INSERT INTO statistics (metric_name, metric_value, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (metric_name) DO UPDATE SET metric_value = $2, updated_at = NOW()
	`
	_, err := db.Exec(query, name, value)
	return err
}

// GetStatistics retrieves all platform statistics
func (db *DB) GetStatistics() (map[string]float64, error) {
	query := `SELECT metric_name, metric_value FROM statistics`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]float64)
	for rows.Next() {
		var name string
		var value float64
		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}
		stats[name] = value
	}
	return stats, nil
}

// IncrementStatistic increments a statistic by a delta value
func (db *DB) IncrementStatistic(name string, delta float64) error {
	query := `
		INSERT INTO statistics (metric_name, metric_value, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (metric_name) DO UPDATE SET 
			metric_value = statistics.metric_value + $2,
			updated_at = NOW()
	`
	_, err := db.Exec(query, name, delta)
	return err
}
