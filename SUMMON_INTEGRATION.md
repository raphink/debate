# Summon Integration Summary

## Overview

The debate generator now uses **Summon** with a custom GCP Secret Manager plugin for secure secret management in both local development and production environments.

## Architecture

```
Local Development:
summon -p gcloud → /usr/local/lib/summon/gcloud → GCP Secret Manager → Environment Variables → Docker Compose

Production (Cloud Functions):
Cloud Functions → GCP Secret Manager (native integration)
```

## Benefits

1. **Unified Secret Source**: Same secrets used in local dev and production
2. **No .env Files**: Eliminates risk of committing secrets to version control
3. **Secure by Default**: Secrets never stored on disk in plaintext
4. **Easy Rotation**: Update secrets in GCP Secret Manager, applies everywhere
5. **Audit Trail**: GCP Secret Manager provides access logging
6. **Team Consistency**: All developers use identical secret workflow
7. **Production Parity**: Local development matches production secret access pattern

## Files Modified

### 1. `secrets.yml` (NEW)
Maps environment variables to GCP Secret Manager paths:
```yaml
ANTHROPIC_API_KEY: gcp/secrets/anthropic-api-key
GCP_PROJECT_ID: gcp/secrets/gcp-project-id
```

### 2. `start-local.sh` (UPDATED)
Added prerequisite checks:
- Verifies `summon` binary exists
- Verifies `/usr/local/lib/summon/gcloud` plugin exists
- Verifies `secrets.yml` configuration exists
- Runs: `summon -p gcloud docker-compose up --build -d`

### 3. `specs/001-debate-generator/research.md` (UPDATED)
Added Section 8: "Local Development with Docker"
- Documented summon workflow
- Explained plugin installation requirements
- Listed benefits of unified secret management

### 4. `specs/001-debate-generator/quickstart.md` (UPDATED)
**Docker Option Prerequisites:**
- Added summon installation requirement
- Added GCP plugin requirement
- Added GCP Secret Manager setup instructions

**Manual Option:**
- Updated to use `summon -p gcloud env` or `gcloud secrets` commands
- Removed `.env` file creation instructions
- Added `.env.local` temporary file pattern for manual development

**All Commands:**
- `docker-compose up` → `summon -p gcloud docker-compose up`
- `docker-compose logs` → `summon -p gcloud docker-compose logs`
- `docker-compose down` → `summon -p gcloud docker-compose down`

### 5. `README.md` (UPDATED)
**Prerequisites:**
- Added summon installation
- Added GCP plugin requirement
- Added GCP Secret Manager access requirement

**Setup:**
- Replaced `.env` creation with GCP secret creation commands
- Added `gcloud secrets create` examples

**Docker Commands:**
- All docker-compose commands now use `summon -p gcloud` prefix
- Added explanation of summon benefits

**Manual Setup:**
- Updated secret export commands to use summon or gcloud CLI
- Removed `.env` file references

### 6. `.env.example` (UPDATED)
Added comprehensive header:
- Notes file is for reference only
- Documents recommended GCP Secret Manager + Summon approach
- Shows alternative manual setup pattern
- Links to `secrets.yml` for secret mappings
- Reinforces never committing `.env` files

## Prerequisites for Developers

### Required Software
1. **Summon**: Install from https://github.com/cyberark/summon#install
2. **Summon GCP Plugin**: Install to `/usr/local/lib/summon/gcloud`
3. **gcloud CLI**: For GCP Secret Manager access
4. **Docker & Docker Compose**: For containerized development

### GCP Secret Manager Setup
Create required secrets (one-time per project):

```bash
# Store Anthropic API key
echo -n "sk-ant-api03-YOUR_KEY_HERE" | gcloud secrets create anthropic-api-key --data-file=-

# Store GCP project ID
echo -n "your-project-id" | gcloud secrets create gcp-project-id --data-file=-
```

### Verify Installation
```bash
# Check summon
which summon

# Check GCP plugin
ls -l /usr/local/lib/summon/gcloud

# Check GCP authentication
gcloud auth list

# Test secret access
summon -p gcloud -f secrets.yml env
```

## Usage

### Starting Local Development
```bash
# Recommended: Use convenience script (validates prerequisites)
./start-local.sh

# Or manually with summon
summon -p gcloud docker-compose up --build

# Or detached mode
summon -p gcloud docker-compose up -d --build
```

### Viewing Logs
```bash
summon -p gcloud docker-compose logs -f

# Specific service
summon -p gcloud docker-compose logs -f validate-topic
```

### Stopping Services
```bash
summon -p gcloud docker-compose down
```

### Manual Development (No Docker)
```bash
# Export secrets to temporary file
summon -p gcloud env | grep -E "ANTHROPIC_API_KEY|GCP_PROJECT_ID" > .env.local
source .env.local

# Or export directly
export ANTHROPIC_API_KEY=$(gcloud secrets versions access latest --secret="anthropic-api-key")
export GCP_PROJECT_ID=$(gcloud secrets versions access latest --secret="gcp-project-id")

# Then run services normally
cd backend/functions/validate-topic && go run main.go
```

## Security Notes

1. **Never commit secrets**: `.env` and `.env.local` are in `.gitignore`
2. **Temporary files only**: `.env.local` is only for manual development sessions
3. **GCP IAM**: Ensure proper IAM permissions for Secret Manager access
4. **Audit access**: Use GCP Secret Manager audit logs to track secret access
5. **Principle of least privilege**: Grant Secret Manager access only to required developers

## Troubleshooting

### "summon: command not found"
Install summon from https://github.com/cyberark/summon#install

### "plugin not found: gcloud"
Ensure custom plugin is installed at `/usr/local/lib/summon/gcloud`

### "permission denied" when accessing secrets
```bash
# Check GCP authentication
gcloud auth list

# Re-authenticate if needed
gcloud auth login

# Verify IAM permissions
gcloud projects get-iam-policy YOUR_PROJECT_ID \
  --flatten="bindings[].members" \
  --filter="bindings.members:user:YOUR_EMAIL"
```

### Secrets not found
```bash
# List secrets in project
gcloud secrets list

# Create missing secret
echo -n "value" | gcloud secrets create SECRET_NAME --data-file=-
```

## Migration from .env Files

If you previously used `.env` files:

1. **Migrate secrets to GCP:**
   ```bash
   # Read from existing .env
   export ANTHROPIC_API_KEY=$(grep ANTHROPIC_API_KEY .env | cut -d '=' -f2)
   
   # Store in GCP Secret Manager
   echo -n "$ANTHROPIC_API_KEY" | gcloud secrets create anthropic-api-key --data-file=-
   ```

2. **Remove .env file:**
   ```bash
   rm .env  # Already in .gitignore, won't be committed
   ```

3. **Use summon workflow:**
   ```bash
   summon -p gcloud docker-compose up --build
   ```

## Future Enhancements

- [ ] Add summon setup tasks to Phase 1 in `tasks.md`
- [ ] Create automated GCP Secret Manager setup script
- [ ] Add CI/CD pipeline secret injection using summon
- [ ] Document secret rotation procedure
- [ ] Add integration tests for summon workflow

## References

- Summon: https://github.com/cyberark/summon
- GCP Secret Manager: https://cloud.google.com/secret-manager/docs
- Debate Generator Research Doc: `specs/001-debate-generator/research.md#8-local-development-with-docker`
