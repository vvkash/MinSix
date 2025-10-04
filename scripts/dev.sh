#!/bin/bash

echo "Starting Minsix Development Environment"
echo "=========================================="
echo ""

# Function to cleanup background processes
cleanup() {
    echo ""
    echo "Shutting down services..."
    kill $(jobs -p) 2>/dev/null
    exit
}

trap cleanup SIGINT SIGTERM

# Check if PostgreSQL is running
if ! docker ps | grep -q minsix-postgres; then
    echo "Starting PostgreSQL..."
    docker-compose up -d
    sleep 3
fi

# Start backend
echo "Starting backend..."
cd backend
go run cmd/server/main.go &
BACKEND_PID=$!

# Wait for backend to start
sleep 3

# Start frontend
echo "Starting frontend..."
cd ../frontend
npm run dev &
FRONTEND_PID=$!

echo ""
echo "All services started"
echo "=========================================="
echo "Dashboard:  http://localhost:3000"
echo "Backend:    http://localhost:8080"
echo "Database:   localhost:5432"
echo "=========================================="
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

# Wait for processes
wait
