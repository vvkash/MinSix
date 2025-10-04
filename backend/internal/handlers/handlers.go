package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/minsix/backend/internal/database"
	ws "github.com/minsix/backend/internal/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type Handler struct {
	db  *database.DB
	hub *ws.Hub
}

func NewHandler(db *database.DB, hub *ws.Hub) *Handler {
	return &Handler{
		db:  db,
		hub: hub,
	}
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "healthy",
		"clients": h.hub.GetClientCount(),
	}
	respondJSON(w, http.StatusOK, response)
}

// GetFlaggedTransactions returns recent flagged transactions
func (h *Handler) GetFlaggedTransactions(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	transactions, err := h.db.GetFlaggedTransactions(limit)
	if err != nil {
		log.Printf("Error getting flagged transactions: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	respondJSON(w, http.StatusOK, transactions)
}

// GetWalletAnalysis returns analysis for a specific wallet
func (h *Handler) GetWalletAnalysis(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	if address == "" {
		respondError(w, http.StatusBadRequest, "Address is required")
		return
	}

	transactions, err := h.db.GetWalletTransactions(address, 100)
	if err != nil {
		log.Printf("Error getting wallet transactions: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to fetch wallet data")
		return
	}

	isBlacklisted, err := h.db.IsBlacklisted(address)
	if err != nil {
		log.Printf("Error checking blacklist: %v", err)
	}

	response := map[string]interface{}{
		"address":       address,
		"transactions":  transactions,
		"blacklisted":   isBlacklisted,
		"tx_count":      len(transactions),
	}

	respondJSON(w, http.StatusOK, response)
}

// GetStatistics returns platform statistics
func (h *Handler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.db.GetStatistics()
	if err != nil {
		log.Printf("Error getting statistics: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to fetch statistics")
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

// HandleWebSocket handles WebSocket connections
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := ws.NewClient(h.hub, conn)
	h.hub.Register(client)

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
