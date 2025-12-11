#!/bin/bash
set -e

# Deployment script for AI-Powered Debate Generator
# Deploys backend to GCP Cloud Functions and frontend to GitHub Pages

# Configuration
PROJECT_ID="debate-480911"
REGION="europe-west1"
RUNTIME="go123"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Load gcloud config from .gcloudrc if it exists
    if [ -f .gcloudrc ]; then
        log_info "Loading gcloud configuration from .gcloudrc"
        set -a
        source .gcloudrc
        set +a
    fi
    
    # Check gcloud
    if ! command -v gcloud &> /dev/null; then
        log_error "gcloud CLI not found. Install from: https://cloud.google.com/sdk/docs/install"
        exit 1
    fi
    
    # Check if authenticated
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
        log_error "Not authenticated with gcloud. Run: gcloud auth login"
        exit 1
    fi
    
    # Check project
    CURRENT_PROJECT=$(gcloud config get-value project 2>/dev/null)
    if [ "$CURRENT_PROJECT" != "$PROJECT_ID" ]; then
        log_warn "Current project is $CURRENT_PROJECT, switching to $PROJECT_ID"
        gcloud config set project "$PROJECT_ID"
    fi
    
    # Check secrets exist
    if ! gcloud secrets describe anthropic-api-key &>/dev/null; then
        log_error "Secret 'anthropic-api-key' not found. Create it with:"
        log_error "  echo -n 'YOUR_KEY' | gcloud secrets create anthropic-api-key --data-file=-"
        exit 1
    fi
    
    # Enable required APIs
    log_info "Checking and enabling required APIs..."
    REQUIRED_APIS=(
        "cloudfunctions.googleapis.com"
        "cloudbuild.googleapis.com"
        "run.googleapis.com"
        "secretmanager.googleapis.com"
        "artifactregistry.googleapis.com"
    )
    
    for api in "${REQUIRED_APIS[@]}"; do
        if ! gcloud services list --enabled --filter="name:$api" --format="value(name)" | grep -q "$api"; then
            log_info "Enabling $api..."
            gcloud services enable "$api" --quiet
        fi
    done
    
    log_info "Prerequisites check passed ✓"
}

# Deploy backend Cloud Functions
deploy_backend() {
    log_info "Deploying backend Cloud Functions to $REGION..."
    
    # Deploy validate-topic function
    log_info "Deploying validate-topic function..."
    gcloud functions deploy validate-topic \
        --gen2 \
        --runtime="$RUNTIME" \
        --region="$REGION" \
        --source=./backend/functions/validate-topic \
        --entry-point=HandleValidateTopic \
        --trigger-http \
        --allow-unauthenticated \
        --set-secrets=ANTHROPIC_API_KEY=anthropic-api-key:latest \
        --set-env-vars=ALLOWED_ORIGIN=https://raphink.github.io \
        --memory=256MB \
        --timeout=60s \
        --max-instances=100 \
        --min-instances=0 \
        --quiet
    
    # Get the URL
    VALIDATE_URL=$(gcloud functions describe validate-topic --region="$REGION" --gen2 --format="value(serviceConfig.uri)")
    log_info "validate-topic deployed: $VALIDATE_URL"
    
    # Note: suggest-panelists function removed - panelists now returned by validate-topic
    
    # Deploy generate-debate function
    log_info "Deploying generate-debate function..."
    gcloud functions deploy generate-debate \
        --gen2 \
        --runtime="$RUNTIME" \
        --region="$REGION" \
        --source=./backend/functions/generate-debate \
        --entry-point=HandleGenerateDebate \
        --trigger-http \
        --allow-unauthenticated \
        --set-secrets=ANTHROPIC_API_KEY=anthropic-api-key:latest \
        --set-env-vars=ALLOWED_ORIGIN=https://raphink.github.io \
        --memory=512MB \
        --timeout=300s \
        --max-instances=100 \
        --min-instances=0 \
        --quiet
    
    DEBATE_URL=$(gcloud functions describe generate-debate --region="$REGION" --gen2 --format="value(serviceConfig.uri)")
    log_info "generate-debate deployed: $DEBATE_URL"
    
    log_info "Backend deployment complete ✓"
    
    # Export URLs for frontend build (no suggest-panelists URL needed)
    export REACT_APP_VALIDATE_TOPIC_URL="$VALIDATE_URL"
    export REACT_APP_GENERATE_DEBATE_URL="$DEBATE_URL"
    
    # Save URLs to file for frontend deployment
    cat > frontend/.env.production << EOF
REACT_APP_VALIDATE_TOPIC_URL=$VALIDATE_URL
REACT_APP_GENERATE_DEBATE_URL=$DEBATE_URL
EOF
    
    log_info "Saved production URLs to frontend/.env.production"
}

# Deploy frontend to GitHub Pages
deploy_frontend() {
    log_info "Building and deploying frontend to GitHub Pages..."
    
    cd frontend
    
    # Check if gh-pages is installed
    if ! npm list gh-pages &>/dev/null; then
        log_info "Installing gh-pages..."
        npm install --save-dev gh-pages
    fi
    
    # Build the React app
    log_info "Building React app..."
    npm run build
    
    # Deploy to GitHub Pages
    log_info "Deploying to GitHub Pages..."
    npm run deploy
    
    cd ..
    
    log_info "Frontend deployment complete ✓"
    log_info "App will be available at: https://raphink.github.io/debate"
}

# Main deployment flow
main() {
    log_info "Starting deployment process..."
    
    # Parse command line arguments
    DEPLOY_BACKEND=true
    DEPLOY_FRONTEND=true
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --backend-only)
                DEPLOY_FRONTEND=false
                shift
                ;;
            --frontend-only)
                DEPLOY_BACKEND=false
                shift
                ;;
            --help)
                echo "Usage: $0 [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  --backend-only    Deploy only backend Cloud Functions"
                echo "  --frontend-only   Deploy only frontend to GitHub Pages"
                echo "  --help            Show this help message"
                echo ""
                echo "Environment variables:"
                echo "  PROJECT_ID        GCP project ID (default: debate-480911)"
                echo "  REGION            GCP region (default: europe-west1)"
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done
    
    check_prerequisites
    
    if [ "$DEPLOY_BACKEND" = true ]; then
        deploy_backend
    fi
    
    if [ "$DEPLOY_FRONTEND" = true ]; then
        if [ "$DEPLOY_BACKEND" = false ]; then
            log_warn "Deploying frontend only - make sure backend URLs are set in frontend/.env.production"
        fi
        deploy_frontend
    fi
    
    log_info "=========================================="
    log_info "Deployment complete! 🎉"
    log_info "=========================================="
    
    if [ "$DEPLOY_BACKEND" = true ]; then
        log_info "Backend Functions:"
        log_info "  - validate-topic: $VALIDATE_URL"
        log_info "  - suggest-panelists: $PANELISTS_URL"
        log_info "  - generate-debate: $DEBATE_URL"
    fi
    
    if [ "$DEPLOY_FRONTEND" = true ]; then
        log_info "Frontend: https://raphink.github.io/debate"
        log_info "Note: GitHub Pages deployment may take 1-2 minutes to propagate"
    fi
}

# Run main
main "$@"
