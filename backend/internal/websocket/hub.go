package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/minsix/backend/internal/models"
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected. Total clients: %d", len(h.clients))
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastAlert sends an alert to all connected clients
func (h *Hub) BroadcastAlert(flagged *models.FlaggedTransaction) {
	msg := models.WebSocketMessage{
		Type: "fraud_alert",
		Payload: models.AlertPayload{
			TxHash:    flagged.TxHash,
			RiskScore: flagged.RiskScore,
			Reasons:   flagged.Reasons,
			Timestamp: flagged.FlaggedAt,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal alert: %v", err)
		return
	}

	h.broadcast <- data
}

// BroadcastTransaction sends a new transaction to all connected clients
func (h *Hub) BroadcastTransaction(tx *models.Transaction) {
	msg := models.WebSocketMessage{
		Type:    "new_transaction",
		Payload: tx,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal transaction: %v", err)
		return
	}

	h.broadcast <- data
}

// BroadcastStats sends updated statistics to all connected clients
func (h *Hub) BroadcastStats(stats map[string]float64) {
	msg := models.WebSocketMessage{
		Type:    "stats_update",
		Payload: stats,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal stats: %v", err)
		return
	}

	h.broadcast <- data
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
