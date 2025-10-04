# Minsix - Ethereum Fraud Detection Platform

A real-time fraud detection platform for Ethereum wallets using Go, PostgreSQL, Next.js, and Alchemy.

## Architecture

- **Backend**: Go with WebSocket support for real-time monitoring
- **Database**: PostgreSQL for transaction storage and analytics
- **Frontend**: Next.js dashboard with real-time updates
- **Blockchain**: Alchemy API for Ethereum data indexing

## Features

- Real-time transaction monitoring via WebSockets
- Multi-heuristic fraud detection:
  - Large/unusual transfers
  - Blacklisted address detection
  - Unusual timing patterns
  - Contract interaction analysis
- Interactive dashboard with Etherscan integration
- Historical transaction analysis

Just create a Alchemy account and sign up for the free API
