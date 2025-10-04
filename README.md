# Minsix - Ethereum Fraud Detection Platform

A real-time fraud detection platform for Ethereum wallets using Go, PostgreSQL, Next.js, and Alchemy.

## Architecture

- **Backend**: Go with WebSocket support for real time data and figures
- **Database**: PostgreSQL for transaction storage and analytics
- **Frontend**: Next.js dashboard with constant updates
- **Blockchain**: Alchemy API for Ethereum data indexing

## Features

- Real time transaction monitoring via WebSockets
- layered fraud detection:
  - Large/unusual transfers
  - Blacklisted address detection
  - Unusual timing patterns
  - Contract interaction analysis
- Interactive dashboard with Etherscan integration
- Historical transaction analysis

Just create a Alchemy account and sign up for the free API

<img width="1384" height="812" alt="Screenshot 2025-10-04 at 3 19 26 PM" src="https://github.com/user-attachments/assets/de6a7e6c-c0e8-40d8-8ba7-febfa2397531" />
