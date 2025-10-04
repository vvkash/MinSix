package ethereum

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/minsix/backend/internal/models"
)

type Client struct {
	client    *ethclient.Client
	apiURL    string
	networkID *big.Int
}

func NewClient(alchemyAPIKey, network string) (*Client, error) {
	apiURL := fmt.Sprintf("wss://eth-mainnet.g.alchemy.com/v2/%s", alchemyAPIKey)
	if network != "eth-mainnet" {
		apiURL = fmt.Sprintf("wss://%s.g.alchemy.com/v2/%s", network, alchemyAPIKey)
	}

	client, err := ethclient.Dial(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get network ID: %w", err)
	}

	log.Printf("Connected to Ethereum network (ID: %s)", networkID.String())

	return &Client{
		client:    client,
		apiURL:    apiURL,
		networkID: networkID,
	}, nil
}

// SubscribeToBlocks subscribes to new blocks and processes transactions
func (c *Client) SubscribeToBlocks(ctx context.Context, txHandler func(*models.Transaction)) error {
	headers := make(chan *types.Header)
	sub, err := c.client.SubscribeNewHead(ctx, headers)
	if err != nil {
		return fmt.Errorf("failed to subscribe to new blocks: %w", err)
	}

	log.Println("Subscribed to new blocks")

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Printf("Subscription error: %v", err)
				return
			case header := <-headers:
				log.Printf("New block: %d", header.Number.Uint64())
				
				// Process block transactions
				block, err := c.client.BlockByNumber(ctx, header.Number)
				if err != nil {
					log.Printf("Failed to get block: %v", err)
					continue
				}

				// Process each transaction in the block
				for _, tx := range block.Transactions() {
					modelTx, err := c.convertTransaction(tx, block)
					if err != nil {
						log.Printf("Failed to convert transaction: %v", err)
						continue
					}
					txHandler(modelTx)
				}
			case <-ctx.Done():
				log.Println("Block subscription stopped")
				return
			}
		}
	}()

	return nil
}

// MonitorAddress monitors transactions for a specific address
func (c *Client) MonitorAddress(ctx context.Context, address string, txHandler func(*models.Transaction)) error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(address)},
	}

	logs := make(chan types.Log)
	sub, err := c.client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to address logs: %w", err)
	}

	log.Printf("Monitoring address: %s", address)

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Printf("Address monitoring error: %v", err)
				return
			case vLog := <-logs:
				tx, _, err := c.client.TransactionByHash(ctx, vLog.TxHash)
				if err != nil {
					log.Printf("Failed to get transaction: %v", err)
					continue
				}

				block, err := c.client.BlockByHash(ctx, vLog.BlockHash)
				if err != nil {
					log.Printf("Failed to get block: %v", err)
					continue
				}

				modelTx, err := c.convertTransaction(tx, block)
				if err != nil {
					log.Printf("Failed to convert transaction: %v", err)
					continue
				}
				txHandler(modelTx)
			case <-ctx.Done():
				log.Println("Address monitoring stopped")
				return
			}
		}
	}()

	return nil
}

// GetTransaction retrieves a transaction by hash
func (c *Client) GetTransaction(ctx context.Context, txHash string) (*models.Transaction, error) {
	hash := common.HexToHash(txHash)
	tx, pending, err := c.client.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	if pending {
		return nil, fmt.Errorf("transaction is pending")
	}

	receipt, err := c.client.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get receipt: %w", err)
	}

	block, err := c.client.BlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get block: %w", err)
	}

	return c.convertTransaction(tx, block)
}

// GetLatestBlock gets the latest block number
func (c *Client) GetLatestBlock(ctx context.Context) (uint64, error) {
	header, err := c.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

// convertTransaction converts eth transaction to internal model
func (c *Client) convertTransaction(tx *types.Transaction, block *types.Block) (*models.Transaction, error) {
	from, err := types.Sender(types.LatestSignerForChainID(c.networkID), tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender: %w", err)
	}

	modelTx := &models.Transaction{
		TxHash:      tx.Hash().Hex(),
		BlockNumber: block.Number().Int64(),
		FromAddress: from.Hex(),
		Value:       tx.Value().String(),
		GasPrice:    tx.GasPrice().String(),
		Timestamp:   time.Unix(int64(block.Time()), 0),
	}

	if tx.To() != nil {
		to := tx.To().Hex()
		modelTx.ToAddress = &to
	}

	// Get gas used from receipt (in real impl)
	modelTx.GasUsed = int64(tx.Gas())

	if len(tx.Data()) > 0 {
		data := fmt.Sprintf("0x%x", tx.Data())
		modelTx.InputData = &data
	}

	return modelTx, nil
}

// Close closes the Ethereum client connection
func (c *Client) Close() {
	c.client.Close()
}
