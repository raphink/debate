#!/bin/bash

# Quick start script for local development with Docker

set -e

echo "🚀 Starting Debate Generator with Docker Compose..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  No .env file found. Creating from .env.example..."
    cp .env.example .env
    echo "📝 Please edit .env and add your ANTHROPIC_API_KEY"
    echo "   Then run this script again."
    exit 1
fi

# Check if ANTHROPIC_API_KEY is set
if grep -q "your-key-here" .env; then
    echo "⚠️  Please set your ANTHROPIC_API_KEY in .env file"
    exit 1
fi

# Build and start services
echo "🔨 Building and starting services..."
docker-compose up --build -d

echo ""
echo "✅ All services started successfully!"
echo ""
echo "📍 Application URLs:"
echo "   Frontend:              http://localhost:3000"
echo "   Topic Validation:      http://localhost:8080"
echo "   Panelist Suggestions:  http://localhost:8081"
echo "   Debate Generation:     http://localhost:8082"
echo ""
echo "📊 View logs with: docker-compose logs -f"
echo "🛑 Stop services with: docker-compose down"
echo ""
