# 5-Minute Testing Guide

The fastest way to test Minsix without waiting for real Ethereum blocks.

## Quick Test (Recommended)

### 1. Setup
```bash
make setup
# Edit backend/.env and add your ALCHEMY_API_KEY
```

### 2. Start Database
```bash
docker-compose up -d
```

### 3. Generate Test Data
```bash
cd backend
go run cmd/test-data/main.go
```

This creates realistic test transactions including:
- Large transfers (flagged)
- Null/burn addresses (flagged)
- Rapid transactions (flagged)
- Normal transactions (not flagged)

### 4. View Results

**Option A: View in Database**
```bash
docker exec -it minsix-postgres psql -U postgres -d minsix

SELECT tx_hash, risk_score, reasons 
FROM flagged_transactions 
ORDER BY flagged_at DESC;

\q
```

**Option B: View in Dashboard**
```bash
# Terminal 1: Start backend
cd backend
go run cmd/server/main.go

# Terminal 2: Start frontend
cd frontend
npm run dev

# Open http://localhost:3000
```

## What You Should See

### In the Dashboard:
- "Live" status indicator (green)
- 6+ flagged transactions in the list
- Each with risk scores 40-70
- Reasons like "Large transfer", "Null address", etc.
- Statistics showing total transactions and flagged count
- Responsive cards with transaction details

### Expected Test Results:
```
FLAGGED: Large transfer (Risk: 70)
FLAGGED: Null address (Risk: 70)
FLAGGED: Burn address (Risk: 70)
OK: Normal transaction
```

## API Testing

Test the backend directly:

```bash
# Health check
curl http://localhost:8080/api/health

# Get flagged transactions
curl http://localhost:8080/api/transactions?limit=10 | json_pp

# Get statistics  
curl http://localhost:8080/api/stats | json_pp
```

## WebSocket Testing

```bash
# Install wscat
npm install -g wscat

# Connect to WebSocket
wscat -c ws://localhost:8080/ws

# You should see connection confirmation
# Generate more test data to see real-time alerts
```

## Full System Test (with Real Ethereum)

If you want to test with live blockchain data:

```bash
# 1. Edit backend/.env
ALCHEMY_NETWORK=eth-sepolia  # Use testnet for faster blocks

# 2. Start everything
make dev

# 3. Wait for blocks (Sepolia: ~12s, Mainnet: ~12s)
# Watch the backend logs for "New block: X"
```

**Note:** Mainnet testing uses your Alchemy credits and is slower.

## Troubleshooting

**No data appearing?**
```bash
# Check if test data was created
docker exec -it minsix-postgres psql -U postgres -d minsix -c "SELECT COUNT(*) FROM transactions;"
```

**Frontend errors?**
```bash
# Install dependencies if you haven't
cd frontend
npm install
```

**Backend won't start?**
- Check that port 8080 is available
- Verify ALCHEMY_API_KEY is set in backend/.env
- Ensure PostgreSQL is running

## Next Steps

Once testing works:
1. See [TESTING.md](TESTING.md) for comprehensive testing
2. See [DEPLOYMENT.md](DEPLOYMENT.md) for hosting options
3. Customize fraud detection in `backend/internal/detector/detector.go`

## Quick Commands

```bash
make test-data       # Generate test transactions
make dev             # Start everything
make test-backend    # Run unit tests
docker-compose logs  # View database logs
```

That's it! You now have a working fraud detection system with test data.
