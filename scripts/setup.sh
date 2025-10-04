#!/bin/bash

echo "Minsix Setup Script"
echo "====================="
echo ""

# Check if required tools are installed
command -v go >/dev/null 2>&1 || { echo "ERROR: Go is not installed. Please install Go 1.21+"; exit 1; }
command -v node >/dev/null 2>&1 || { echo "ERROR: Node.js is not installed. Please install Node.js 18+"; exit 1; }
command -v docker >/dev/null 2>&1 || { echo "ERROR: Docker is not installed. Please install Docker"; exit 1; }

echo "All required tools are installed"
echo ""

# Setup backend
echo "Setting up backend..."
cd backend

if [ ! -f .env ]; then
    echo "Creating .env file from example..."
    cp .env.example .env
    echo "NOTE: Please edit backend/.env and add your ALCHEMY_API_KEY"
fi

echo "Installing Go dependencies..."
go mod download
echo "Backend setup complete"
echo ""

# Setup frontend
echo "Setting up frontend..."
cd ../frontend

if [ ! -f .env.local ]; then
    echo "Creating .env.local file from example..."
    cp .env.local.example .env.local
fi

echo "Installing npm dependencies..."
npm install
echo "Frontend setup complete"
echo ""

# Start PostgreSQL
echo "Starting PostgreSQL..."
cd ..
docker-compose up -d
echo "PostgreSQL started"
echo ""

# Wait for PostgreSQL
echo "Waiting for PostgreSQL to be ready..."
sleep 5

# Run migrations
echo "Running database migrations..."
cd backend
go run cmd/migrate/main.go
echo "Migrations complete"
echo ""

echo "Setup complete"
echo ""
echo "Next steps:"
echo "1. Edit backend/.env and add your ALCHEMY_API_KEY"
echo "2. Run 'npm run dev:all' to start all services"
echo ""
