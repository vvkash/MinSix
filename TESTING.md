# Testing Guide for Minsix

This guide covers different testing approaches to verify the platform works correctly.

## Testing Levels

### 1. Unit Testing (Automated)

**Backend Tests:**
```bash
cd backend
go test ./... -v -cover
```

Tests fraud detection heuristics, database operations, and utility functions.

**Frontend Tests:**
```bash
cd frontend
npm test
```

### 2. Integration Testing

Test the full stack locally before deploying.

#### Option A: Test with Ethereum Mainnet (Live Data)

This is the most realistic test but uses your Alchemy credits:

```bash
# 1. Start services
make dev

# 2. Monitor backend logs
# You should see:
# - "Connected to Ethereum network"
# - "Subscribed to new blocks"
# - "New block: [number]"

# 3. Open dashboard
# http://localhost:3000
# - Should show "Live" status (green)
# - Will display real transactions as they occur
```

**Pros:** Real-world data, tests actual fraud detection
**Cons:** Slow (waits for real blocks), uses Alchemy credits

#### Option B: Test with Ethereum Testnet (Sepolia/Goerli)

Faster block times, free to use:

```bash
# 1. Change network in backend/.env
ALCHEMY_NETWORK=eth-sepolia

# 2. Start services
make dev
```

**Pros:** Faster blocks (~12s), free
**Cons:** Less realistic transaction patterns

#### Option C: Mock Transaction Testing (Recommended for Development)

Test fraud detection without waiting for real blocks.

### 3. Manual API Testing

Test backend endpoints directly:

```bash
# Health check
curl http://localhost:8080/api/health

# Get flagged transactions
curl http://localhost:8080/api/transactions?limit=10

# Get statistics
curl http://localhost:8080/api/stats

# Wallet analysis (replace with real address)
curl http://localhost:8080/api/wallets/0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
```

### 4. WebSocket Testing

Test real-time connections:

```bash
# Install wscat
npm install -g wscat

# Connect to WebSocket
wscat -c ws://localhost:8080/ws

# Should receive messages like:
# {"type":"fraud_alert","payload":{...}}
```

### 5. Database Testing

Verify data persistence:

```bash
# Connect to PostgreSQL
docker exec -it minsix-postgres psql -U postgres -d minsix

# Check tables
\dt

# View flagged transactions
SELECT tx_hash, risk_score, reasons FROM flagged_transactions LIMIT 5;

# View statistics
SELECT * FROM statistics;

# Exit
\q
```

## Creating Test Data

Since waiting for real Ethereum blocks is slow, here's how to create test data:

### Method 1: Manual Database Insertion

```sql
-- Connect to database
docker exec -it minsix-postgres psql -U postgres -d minsix

-- Insert test transaction
INSERT INTO transactions (tx_hash, block_number, from_address, to_address, value, gas_price, gas_used, timestamp)
VALUES (
  '0xtest123456789abcdef',
  12345678,
  '0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb',
  '0x0000000000000000000000000000000000000000',
  '15000000000000000000', -- 15 ETH (should trigger large transfer)
  '30000000000',
  21000,
  NOW()
);

-- Flag it manually
INSERT INTO flagged_transactions (transaction_id, tx_hash, risk_score, reasons, status)
VALUES (
  currval('transactions_id_seq'),
  '0xtest123456789abcdef',
  70,
  ARRAY['Large transfer: 15 ETH', 'Transaction to null/burn address'],
  'pending'
);
```

### Method 2: Test Data Generator (Recommended)

Use the built-in test data generator:

```bash
# Make sure backend is NOT running (to avoid conflicts)
cd backend
go run cmd/test-data/main.go
```

This creates:
- Large transfer (50 ETH) - **FLAGGED**
- Null address transaction - **FLAGGED**
- Burn address transaction - **FLAGGED**  
- Rapid fire transactions - **FLAGGED**
- High gas price - **FLAGGED**
- Suspicious contract interaction - **FLAGGED**
- Normal transaction - **NOT FLAGGED**

After running, start the frontend to see results:
```bash
cd frontend
npm run dev
# Open http://localhost:3000
```

### Method 3: Load Testing

For performance testing:

```bash
# Install Apache Bench
brew install ab  # macOS
# or apt-get install apache2-utils  # Linux

# Test API performance
ab -n 1000 -c 10 http://localhost:8080/api/transactions

# Test WebSocket connections
# Use a tool like Artillery or k6
```

## Test Scenarios

### Scenario 1: End-to-End Test (Full System)

**Goal:** Verify complete flow from Ethereum to dashboard

```bash
# Terminal 1: Start PostgreSQL
docker-compose up

# Terminal 2: Start backend
cd backend
go run cmd/server/main.go

# Terminal 3: Start frontend
cd frontend
npm run dev

# Terminal 4: Generate test data
cd backend
go run cmd/test-data/main.go
```

**Verify:**
1. Dashboard shows "Live" status
2. Flagged transactions appear in the list
3. Statistics update correctly
4. Each transaction card shows correct risk score
5. Etherscan links work (use mainnet addresses)
6. Responsive design works on mobile

### Scenario 2: Fraud Detection Accuracy

Test each heuristic individually:

```go
// Run specific tests
cd backend
go test -v ./internal/detector -run TestCheckLargeTransfer
go test -v ./internal/detector -run TestCheckNullAddress
go test -v ./internal/detector -run TestCheckRapidTransactions
```

### Scenario 3: WebSocket Reliability

Test connection stability:

1. Open dashboard
2. Check WebSocket status (should be green "Live")
3. Stop backend: `Ctrl+C`
4. Dashboard should show red "Disconnected"
5. Restart backend
6. Dashboard should reconnect automatically

### Scenario 4: Database Performance

Test query performance:

```sql
-- Connect to database
docker exec -it minsix-postgres psql -U postgres -d minsix

-- Check query performance
EXPLAIN ANALYZE SELECT * FROM flagged_transactions 
ORDER BY flagged_at DESC LIMIT 50;

-- Should use indexes and be fast (<10ms)
```

## Testing Checklist

Before deploying to production:

- [ ] Backend unit tests pass (`go test ./...`)
- [ ] Frontend builds successfully (`npm run build`)
- [ ] Docker images build (`docker-compose build`)
- [ ] Database migrations run successfully
- [ ] API endpoints return correct data
- [ ] WebSocket connects and receives messages
- [ ] Dashboard displays data correctly
- [ ] Fraud detection flags suspicious transactions
- [ ] Statistics calculate correctly
- [ ] Error handling works (test with invalid inputs)
- [ ] CORS is configured properly
- [ ] Environment variables are set
- [ ] Alchemy API key works
- [ ] PostgreSQL connection is stable

## Common Test Issues

### Issue: No transactions appearing

**Solutions:**
1. Check Alchemy API key is valid
2. Verify network setting (mainnet vs testnet)
3. Check backend logs for errors
4. Mainnet blocks are ~12s apart, be patient
5. Use test data generator instead

### Issue: WebSocket won't connect

**Solutions:**
1. Check backend is running on port 8080
2. Verify CORS settings allow localhost:3000
3. Check browser console for errors
4. Try different browser (disable extensions)

### Issue: Database connection fails

**Solutions:**
1. Ensure PostgreSQL container is running
2. Check connection string format
3. Verify port 5432 is available
4. Try restarting: `docker-compose restart`

### Issue: Frontend shows stale data

**Solutions:**
1. Hard refresh browser (Cmd+Shift+R)
2. Clear browser cache
3. Check WebSocket is connected
4. Verify API_URL in .env.local

## Automated Testing with GitHub Actions

The project includes CI/CD:

```bash
# Runs automatically on push/PR
# - Backend tests
# - Frontend build
# - Docker builds

# View in .github/workflows/ci.yml
```

## Performance Benchmarks

Expected performance metrics:

| Metric | Target | Notes |
|--------|--------|-------|
| API Response Time | <100ms | /api/transactions |
| WebSocket Latency | <50ms | Alert delivery |
| Database Query | <10ms | With indexes |
| Fraud Detection | <5ms | Per transaction |
| Frontend Load | <2s | Initial page load |
| Memory Usage (Backend) | <200MB | Steady state |

## Testing in Production

### Health Monitoring

```bash
# Set up monitoring endpoint
curl https://your-domain.com/api/health

# Expected response:
# {"status":"healthy","clients":N}
```

### Log Analysis

```bash
# Backend logs to check:
# - "Connected to Ethereum network"
# - "Subscribed to new blocks"
# - "New block: X"
# - "Flagged transaction: 0x..."

# Frontend logs to check:
# - "WebSocket connected"
# - No CORS errors
# - No 404s
```

### A/B Testing Heuristics

To test new fraud detection rules:

1. Deploy to staging environment
2. Run against historical data
3. Measure false positive rate
4. Compare with production heuristics
5. Gradually roll out if better

## Quick Test Commands

```bash
# Run all tests
make test-backend
make test-frontend

# Generate test data
cd backend && go run cmd/test-data/main.go

# Check database
docker exec -it minsix-postgres psql -U postgres -d minsix -c "SELECT COUNT(*) FROM flagged_transactions;"

# Test WebSocket
wscat -c ws://localhost:8080/ws

# Load test
ab -n 100 -c 10 http://localhost:8080/api/health
```

## Next Steps

After testing locally:

1. **Deploy to Staging** - Test in production-like environment
2. **Monitor Metrics** - Set up Datadog/New Relic
3. **Load Test** - Verify performance at scale
4. **Security Audit** - Check for vulnerabilities
5. **Deploy to Production** - See DEPLOYMENT.md

Happy testing!

