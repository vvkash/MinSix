# Quick Start Guide

Get Minsix up and running in 5 minutes.

## Prerequisites

- Go 1.21+ ([Download](https://golang.org/dl/))
- Node.js 18+ ([Download](https://nodejs.org/))
- Docker ([Download](https://www.docker.com/get-started))
- Alchemy API Key ([Sign up](https://www.alchemy.com/))

## 1. Clone & Setup

```bash
# Clone the repository
git clone <your-repo-url>
cd MinsixCrypto

# Run automated setup
make setup
```

This will:
- Install Go dependencies
- Install Node.js dependencies  
- Start PostgreSQL in Docker
- Run database migrations
- Create environment files

## 2. Configure API Key

Edit `backend/.env` and add your Alchemy API key:

```bash
ALCHEMY_API_KEY=your_actual_api_key_here
ALCHEMY_NETWORK=eth-mainnet
DATABASE_URL=postgres://postgres:postgres@localhost:5432/minsix?sslmode=disable
PORT=8080
CORS_ORIGINS=http://localhost:3000
```

## 3. Start Development Servers

```bash
make dev
```

This starts:
- PostgreSQL (port 5432)
- Backend API (port 8080)
- Frontend Dashboard (port 3000)

## 4. Access Dashboard

Open [http://localhost:3000](http://localhost:3000) in your browser.

You should see:
- Real-time transaction monitoring
- Fraud detection alerts
- Platform statistics
- WebSocket connection status

## 5. Verify It's Working

### Check Backend Health
```bash
curl http://localhost:8080/api/health
```

Expected response:
```json
{
  "status": "healthy",
  "clients": 1
}
```

### Check Database Connection
```bash
docker ps | grep postgres
```

Should show `minsix-postgres` running.

### Check Frontend
Visit http://localhost:3000 - you should see the Minsix dashboard with a "Live" indicator (green) in the top right.

## 6. Understanding the Flow

1. **Ethereum Monitoring**: Backend connects to Ethereum via Alchemy and subscribes to new blocks
2. **Transaction Analysis**: Each transaction is analyzed using fraud detection heuristics
3. **Flagging**: Suspicious transactions are flagged and saved to PostgreSQL
4. **Real-time Updates**: WebSocket broadcasts alerts to the dashboard
5. **Visualization**: Dashboard displays flagged transactions with risk scores and reasons

## Troubleshooting

### Backend won't start
- Check if Alchemy API key is set correctly
- Verify PostgreSQL is running: `docker ps`
- Check logs for specific errors

### Frontend shows "Disconnected"
- Ensure backend is running on port 8080
- Check browser console for WebSocket errors
- Verify CORS settings in backend/.env

### Database errors
- Ensure PostgreSQL container is running
- Try restarting: `docker-compose restart`
- Check connection string in .env

### No transactions appearing
- Ethereum mainnet may have periods of low activity
- Check backend logs to ensure it's connected to Alchemy
- Verify your Alchemy API key has available credits

## Next Steps

- **Add Wallets**: Modify `MonitorAddress()` in backend to track specific wallets
- **Customize Heuristics**: Edit `backend/internal/detector/detector.go` to adjust detection logic
- **Explore API**: Try other endpoints like `/api/transactions` and `/api/stats`
- **Deploy**: See [DEPLOYMENT.md](DEPLOYMENT.md) for production deployment options

## Common Commands

```bash
# Start everything
make dev

# Start only backend
make dev-backend

# Start only frontend  
make dev-frontend

# Run migrations
make migrate

# Stop all services
docker-compose down

# View logs
docker-compose logs -f

# Rebuild everything
make clean && make setup
```

## Getting Help

- Check [README.md](README.md) for detailed documentation
- See [DEPLOYMENT.md](DEPLOYMENT.md) for hosting options
- Read [CONTRIBUTING.md](CONTRIBUTING.md) to contribute
- Open an issue on GitHub for bugs

## What's Next?

The system is now monitoring Ethereum mainnet for fraudulent transactions. As blocks are mined and transactions occur, the fraud detection engine will:

1. Analyze each transaction
2. Flag suspicious activity based on heuristics
3. Display alerts on your dashboard in real-time
4. Store historical data in PostgreSQL

You can customize the detection rules, add more heuristics, or integrate with external threat intelligence sources.

Happy fraud hunting!
