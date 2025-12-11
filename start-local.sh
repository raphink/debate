#!/bin/bash

# Quick start script for local development with Docker Compose and Summon

set -e

# Load .gcloudrc if it exists (for multi-account gcloud configurations)
if [ -f .gcloudrc ]; then
    echo "📋 Loading gcloud configuration from .gcloudrc..."
    set -a
    source .gcloudrc
    set +a
fi

echo "🚀 Starting Debate Generator with Docker Compose + Summon..."

# Check if summon is installed
if ! command -v summon &> /dev/null; then
    echo "❌ summon is not installed"
    echo "   Install from: https://cyberark.github.io/summon/"
    exit 1
fi

# Check if summon-gcloud plugin is installed
if [ ! -f /usr/local/lib/summon/gcloud ]; then
    echo "❌ summon-gcloud plugin not found at /usr/local/lib/summon/gcloud"
    echo "   Please install the gcloud plugin for summon"
    exit 1
fi

# Check if secrets.yml exists
if [ ! -f secrets.yml ]; then
    echo "❌ secrets.yml not found"
    echo "   This file should define secret paths in GCP Secret Manager"
    exit 1
fi

# Build and start services with summon
echo "🔨 Building and starting services with summon..."
summon -p gcloud docker-compose up --build -d

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
