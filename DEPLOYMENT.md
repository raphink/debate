# Deployment Guide

## Overview

This application deploys to a hybrid architecture:
- **Backend**: GCP Cloud Functions (europe-west1)
- **Frontend**: GitHub Pages

## Prerequisites

1. **GCP Setup**:
   ```bash
   # Install gcloud CLI
   # https://cloud.google.com/sdk/docs/install
   
   # Authenticate
   gcloud auth login
   
   # Set project
   gcloud config set project debate-480911
   
   # Create secret (if not exists)
   echo -n "YOUR_ANTHROPIC_API_KEY" | gcloud secrets create anthropic-api-key --data-file=-
   ```

2. **GitHub Setup**:
   ```bash
   # Enable GitHub Pages in repository settings
   # Settings → Pages → Source: gh-pages branch
   ```

3. **Local Dependencies**:
   ```bash
   # Node.js 18+ for frontend build
   # gcloud CLI for backend deployment
   ```

## Quick Deploy

Deploy everything:
```bash
./deploy.sh
```

Deploy backend only:
```bash
./deploy.sh --backend-only
```

Deploy frontend only:
```bash
./deploy.sh --frontend-only
```

## Manual Deployment

### Backend (GCP Cloud Functions)

```bash
# Deploy validate-topic
gcloud functions deploy validate-topic \
  --gen2 \
  --runtime=go123 \
  --region=europe-west1 \
  --source=./backend/functions/validate-topic \
  --entry-point=main \
  --trigger-http \
  --allow-unauthenticated \
  --set-secrets=ANTHROPIC_API_KEY=anthropic-api-key:latest \
  --memory=256MB \
  --timeout=60s

# Deploy suggest-panelists
gcloud functions deploy suggest-panelists \
  --gen2 \
  --runtime=go123 \
  --region=europe-west1 \
  --source=./backend/functions/suggest-panelists \
  --entry-point=main \
  --trigger-http \
  --allow-unauthenticated \
  --set-secrets=ANTHROPIC_API_KEY=anthropic-api-key:latest \
  --memory=256MB \
  --timeout=60s

# Deploy generate-debate
gcloud functions deploy generate-debate \
  --gen2 \
  --runtime=go123 \
  --region=europe-west1 \
  --source=./backend/functions/generate-debate \
  --entry-point=main \
  --trigger-http \
  --allow-unauthenticated \
  --set-secrets=ANTHROPIC_API_KEY=anthropic-api-key:latest \
  --memory=512MB \
  --timeout=300s

# Deploy list-debates (supports autocomplete via ?q= parameter)
gcloud functions deploy list-debates \
  --gen2 \
  --runtime=go124 \
  --region=europe-west1 \
  --source=./backend/functions/list-debates \
  --entry-point=main \
  --trigger-http \
  --allow-unauthenticated \
  --set-env-vars=GCP_PROJECT_ID=debate-480911 \
  --memory=256MB \
  --timeout=30s
```

### Frontend (GitHub Pages)

1. **Configure environment**:
   ```bash
   # Create frontend/.env.production with Cloud Functions URLs
   cat > frontend/.env.production << EOF
   REACT_APP_VALIDATE_TOPIC_URL=https://europe-west1-debate-480911.cloudfunctions.net/validate-topic
   REACT_APP_SUGGEST_PANELISTS_URL=https://europe-west1-debate-480911.cloudfunctions.net/suggest-panelists
   REACT_APP_GENERATE_DEBATE_URL=https://europe-west1-debate-480911.cloudfunctions.net/generate-debate
   REACT_APP_GET_PORTRAIT_URL=https://europe-west1-debate-480911.cloudfunctions.net/get-portrait
   EOF
   ```

2. **Build and deploy**:
   ```bash
   cd frontend
   npm install
   npm run build
   npm run deploy
   ```

## CORS Configuration

All backend services use the `ALLOWED_ORIGIN` environment variable for CORS:

- **Development** (localhost): `http://localhost:3000` (default if not set)
- **Production** (GitHub Pages): `https://raphink.github.io`

Set during deployment:
```bash
# Deploying Cloud Functions with CORS configuration
gcloud functions deploy validate-topic \
  --set-env-vars ALLOWED_ORIGIN=https://raphink.github.io \
  ...

# Or update existing function
gcloud functions deploy validate-topic \
  --update-env-vars ALLOWED_ORIGIN=https://raphink.github.io
```

See `.env.production.example` for production configuration reference.

## Verify Deployment

### Backend
```bash
# Test validate-topic
curl -X POST https://europe-west1-debate-480911.cloudfunctions.net/validate-topic \
  -H "Content-Type: application/json" \
  -d '{"topic":"Should Christians defy unjust authorities?"}'

# Expected: {"isRelevant":true,"message":"...","topic":"..."}
```

### Frontend
Visit: https://raphink.github.io/debate

## Architecture

```
GitHub Pages (Frontend)
    ↓ HTTPS API calls
GCP Cloud Functions (Backend)
    ↓ Secret access
GCP Secret Manager (ANTHROPIC_API_KEY)
    ↓ API calls
Anthropic Claude API
```

## Regions

- **Cloud Functions**: europe-west1 (Belgium)
- **GitHub Pages**: Global CDN

## Cost Estimates

**Free Tier Coverage:**
- Cloud Functions: 2M requests/month
- Secret Manager: 10K accesses/month
- GitHub Pages: Unlimited for public repos

**Estimated Monthly Cost** (low traffic):
- GCP: $0 (within free tier)
- Anthropic API: ~$0.01-0.02 per debate

## Monitoring

View logs:
```bash
# All functions
gcloud functions logs read --region=europe-west1

# Specific function
gcloud functions logs read validate-topic --region=europe-west1 --limit=50
```

View metrics:
```bash
# Open Cloud Console
gcloud console functions
```

## Troubleshooting

**Issue**: Function deployment fails
```bash
# Check quotas
gcloud compute project-info describe --project=debate-480911

# Enable APIs
gcloud services enable cloudfunctions.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable secretmanager.googleapis.com
```

**Issue**: Secret not found
```bash
# List secrets
gcloud secrets list

# Create if missing
echo -n "YOUR_KEY" | gcloud secrets create anthropic-api-key --data-file=-

# Grant access to Cloud Functions service account
gcloud secrets add-iam-policy-binding anthropic-api-key \
  --member="serviceAccount:debate-480911@appspot.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

**Issue**: CORS errors
- Cloud Functions automatically set `Access-Control-Allow-Origin: *`
- Check browser console for specific error

**Issue**: GitHub Pages not updating
- Wait 1-2 minutes for CDN propagation
- Check Actions tab for deployment status
- Verify gh-pages branch exists

## Rollback

### Backend
```bash
# List versions
gcloud functions list --region=europe-west1

# Deploy previous version (manual redeploy from git tag)
git checkout <previous-tag>
./deploy.sh --backend-only
```

### Frontend
```bash
# Checkout previous commit
git checkout <previous-commit>
cd frontend
npm run deploy
```

## CI/CD

See `.github/workflows/deploy.yml` for automated deployment on push to main.

## Security Notes

- ✅ Secrets in GCP Secret Manager (never in code)
- ✅ HTTPS only (enforced by Cloud Functions and GitHub Pages)
- ✅ CORS configured for cross-origin requests
- ✅ Rate limiting via Cloud Functions quotas
- ✅ Input sanitization on both frontend and backend
- ✅ No authentication required (public app)

## Support

For issues:
1. Check logs: `gcloud functions logs read --region=europe-west1`
2. Verify secrets: `gcloud secrets describe anthropic-api-key`
3. Test locally: `./start-local.sh`
4. Review deployment output for errors
