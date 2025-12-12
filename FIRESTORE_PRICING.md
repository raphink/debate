# Firestore Storage Pricing Estimate

**Feature**: Debate Storage with UUID URLs  
**Date**: 2025-12-12  
**Purpose**: Cost analysis for storing cached debate conversations in Firestore

## Overview

Store completed debates in Firestore with UUID-based URLs for:
- **Shareable links**: Users can share debate URLs with others
- **Caching**: Avoid regenerating identical debates
- **History**: Users can revisit past debates

## Data Structure

### Firestore Collection: `debates`

Each document uses UUID as document ID (e.g., `550e8400-e29b-41d4-a716-446655440000`)

**Document Structure**:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "topic": {
    "text": "Should Christians defy authorities when the law is unfair?",
    "suggestedNames": ["Martin Luther King Jr."],
    "isRelevant": true,
    "validationMessage": "This topic is well-suited for theological debate"
  },
  "panelists": [
    {
      "id": "Augustine354",
      "name": "Augustine of Hippo",
      "tagline": "4th-5th century theologian and philosopher",
      "biography": "Early Christian theologian known for 'Confessions' and 'City of God'...",
      "avatarUrl": "/avatars/Augustine354-avatar.png",
      "position": "Would argue that Christians have a duty to resist unjust laws"
    },
    {
      "id": "MLKJr",
      "name": "Martin Luther King Jr.",
      "tagline": "20th century civil rights leader",
      "biography": "American Baptist minister and activist...",
      "avatarUrl": "/avatars/MLKJr-avatar.png",
      "position": "Believed in nonviolent civil disobedience against unjust laws"
    }
  ],
  "messages": [
    {
      "id": "moderator-0",
      "panelistId": "moderator",
      "panelistName": "Moderator",
      "avatarUrl": "/avatars/moderator-avatar.png",
      "text": "Welcome to today's debate on whether Christians should defy unjust laws...",
      "timestamp": "2025-12-12T10:30:00Z",
      "sequence": 0,
      "isComplete": true
    }
    // ... 10-20 more messages (typical debate)
  ],
  "status": "complete",
  "startedAt": "2025-12-12T10:30:00Z",
  "completedAt": "2025-12-12T10:32:30Z",
  "metadata": {
    "createdBy": "anonymous",
    "userAgent": "Mozilla/5.0...",
    "shareCount": 0,
    "viewCount": 1
  }
}
```

## Size Calculations

### Per Debate Document

**Component Breakdown**:
- UUID: 36 bytes
- Topic object: ~250 bytes
  - text: ~100 bytes (avg)
  - suggestedNames: ~50 bytes
  - isRelevant/validationMessage: ~100 bytes
- Panelists array (3 panelists avg): ~1,500 bytes
  - Per panelist: ~500 bytes (name, tagline, bio, position, avatar URL)
- Messages array (15 messages avg): ~18,000 bytes
  - Per message: ~1,200 bytes
    - Text content: ~800 bytes (avg, with Markdown formatting)
    - Metadata: ~400 bytes (id, panelistId, name, timestamp, etc.)
- Status/timestamps: ~150 bytes
- Metadata: ~200 bytes

**Total per document**: ~20,000 bytes = **~20 KB**

### Conservative Estimates

- **Small debate** (2 panelists, 10 messages): ~15 KB
- **Average debate** (3 panelists, 15 messages): ~20 KB
- **Large debate** (5 panelists, 25 messages): ~35 KB

**Using 25 KB average** for pricing (conservative estimate with formatting overhead)

## Firestore Pricing (as of Dec 2025)

### Storage Costs
- **$0.18 per GB/month** (stored data)
- **$0.036 per GB/month** (metadata and index entries, estimated at 20% of document size)

**Total storage**: $0.216 per GB/month

### Operations Costs
- **Document writes**: $0.18 per 100,000 writes
  - 1 write per debate generation
- **Document reads**: $0.06 per 100,000 reads
  - 1 read per debate view (cached in browser after initial load)
- **Document deletes**: $0.02 per 100,000 deletes (if implementing cleanup)

### Network Costs
- **Network egress**: $0.12 per GB (downloads)
- First 10 GB/month free

## Usage Scenarios

### Scenario 1: Small App (100 debates/month, 500 views)

**Storage**:
- Data: 100 debates × 25 KB = 2.5 MB = 0.0025 GB
- Monthly cost: 0.0025 GB × $0.216 = **$0.0005/month** (~$0.01/year)

**Operations**:
- Writes: 100 debates = $0.00018
- Reads: 500 views = $0.00030
- Monthly cost: **$0.0005/month**

**Network**:
- Egress: 500 views × 25 KB = 12.5 MB = 0.0125 GB
- Monthly cost: **$0.00** (within free tier)

**Total**: **~$0.001/month** = **$0.01/year** 📊

---

### Scenario 2: Medium App (1,000 debates/month, 10,000 views)

**Storage**:
- Data: 1,000 debates × 25 KB = 25 MB = 0.025 GB
- Cumulative after 6 months: 0.15 GB
- Monthly cost: 0.15 GB × $0.216 = **$0.03/month**

**Operations**:
- Writes: 1,000 debates = $0.0018
- Reads: 10,000 views = $0.006
- Monthly cost: **$0.008/month**

**Network**:
- Egress: 10,000 views × 25 KB = 250 MB = 0.25 GB
- Monthly cost: **$0.00** (within free tier)

**Total**: **~$0.04/month** = **$0.48/year** 📊

---

### Scenario 3: Popular App (10,000 debates/month, 100,000 views)

**Storage**:
- Data: 10,000 debates × 25 KB = 250 MB = 0.25 GB
- Cumulative after 1 year: 3 GB
- Monthly cost: 3 GB × $0.216 = **$0.65/month**

**Operations**:
- Writes: 10,000 debates = $0.018
- Reads: 100,000 views = $0.06
- Monthly cost: **$0.08/month**

**Network**:
- Egress: 100,000 views × 25 KB = 2.5 GB = 2.5 GB
- Monthly cost: **$0.00** (within free tier)

**Total**: **~$0.73/month** = **$8.76/year** 📊

---

### Scenario 4: High-Volume App (50,000 debates/month, 500,000 views)

**Storage**:
- Data: 50,000 debates × 25 KB = 1.25 GB
- Cumulative after 1 year: 15 GB
- Monthly cost: 15 GB × $0.216 = **$3.24/month**

**Operations**:
- Writes: 50,000 debates = $0.09
- Reads: 500,000 views = $0.30
- Monthly cost: **$0.39/month**

**Network**:
- Egress: 500,000 views × 25 KB = 12.5 GB
- Monthly cost: (12.5 GB - 10 GB free) × $0.12 = **$0.30/month**

**Total**: **~$3.93/month** = **$47.16/year** 📊

---

## Free Tier Analysis

**Firestore Spark Plan (Free Tier)**:
- **Storage**: 1 GB free
- **Reads**: 50,000/day free (1.5M/month)
- **Writes**: 20,000/day free (600K/month)
- **Deletes**: 20,000/day free
- **Network egress**: 10 GB/month free

### Free Tier Capacity

With 25 KB average debate size:
- **Storage capacity**: 1 GB ÷ 25 KB = **~40,000 debates**
- **Daily writes**: 20,000/day = up to 20,000 debates/day
- **Daily reads**: 50,000/day = up to 50,000 views/day
- **Network egress**: 10 GB ÷ 25 KB = **~400,000 views/month**

**Free tier supports**:
- Up to 40,000 stored debates (or indefinitely with TTL cleanup)
- Up to 600,000 new debates/month
- Up to 1.5M views/month
- Up to 400,000 downloads/month

**Conclusion**: Free tier is sufficient for **99% of expected usage** 🎉

---

## Optimization Strategies

### 1. Time-To-Live (TTL) Policy
- Auto-delete debates older than 90 days
- Keeps storage under free tier
- Preserves recent/popular debates

```javascript
// Firestore rule example
{
  "metadata": {
    "expiresAt": "2026-03-12T10:30:00Z"  // 90 days from creation
  }
}
```

### 2. Compression
- Gzip compress message text before storing
- Reduce size by ~60-70%
- **New average**: ~10 KB per debate
- **Free tier capacity**: ~100,000 debates

### 3. Deduplication
- Hash topic + panelist IDs to detect duplicates
- Redirect to existing debate if exact match found
- Reduces redundant storage

```javascript
{
  "debateHash": "sha256(topic.text + panelist_ids_sorted)"
}
```

### 4. Lazy Loading Messages
- Store topic/panelists in main document
- Store messages in subcollection: `debates/{uuid}/messages/{sequence}`
- Load messages on-demand
- Reduces initial read size to ~2 KB

**Optimized structure**:
```
debates/{uuid}  // 2 KB: topic + panelists + metadata
  └─ messages/{sequence}  // 1 KB each: individual messages
```

### 5. CDN Caching
- Cache debate documents at edge (Cloudflare, Cloud CDN)
- Reduce Firestore reads by 80-90%
- Free tier on most CDNs

---

## Implementation Recommendations

### Phase 1: MVP (Use Free Tier)
- Store debates as-is (~25 KB each)
- No TTL (let free tier absorb first ~40K debates)
- Simple UUID generation on debate start
- URL pattern: `https://debate.app/d/{uuid}`

**Cost**: $0/month for first year (assuming <40K debates)

### Phase 2: Optimization (If needed)
- Implement 90-day TTL
- Add deduplication via hash
- Compress message text
- Add CDN caching

**Cost**: Still $0/month for most scenarios

### Phase 3: Scale (When exceeding free tier)
- Lazy load messages via subcollections
- Implement tiered storage (hot/cold)
- Archive old debates to Cloud Storage ($0.02/GB/month)

**Cost**: ~$5-10/month at 100K debates/month

---

## Security Considerations

### Firestore Security Rules

```javascript
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    match /debates/{debateId} {
      // Allow anyone to read debates (public sharing)
      allow read: if true;
      
      // Only allow writes from authenticated Cloud Functions
      allow create: if request.auth != null 
                    && request.auth.token.firebase.sign_in_provider == 'custom';
      
      // No updates or deletes from clients
      allow update, delete: if false;
    }
  }
}
```

### Access Control
- **Public reads**: Anyone with UUID can view debate
- **Private writes**: Only backend Cloud Functions can create debates
- **No client writes**: Prevents spam/abuse
- **UUID obscurity**: 128-bit UUID provides ~10^36 combinations (unguessable)

---

## Monitoring & Alerts

### Cloud Monitoring Metrics
- **Storage usage**: Alert at 80% of free tier (800 MB)
- **Read operations**: Alert at 80% of daily limit (40K reads/day)
- **Write operations**: Alert at 80% of daily limit (16K writes/day)
- **Document count**: Track total debates stored
- **Average document size**: Monitor for data bloat

### BigQuery Export (Optional)
- Export Firestore data to BigQuery for analytics
- Track popular topics, panelists, debate lengths
- Cost: Negligible for small datasets

---

## Summary

| Metric | Small | Medium | Popular | High-Volume |
|--------|-------|--------|---------|-------------|
| Debates/month | 100 | 1,000 | 10,000 | 50,000 |
| Views/month | 500 | 10,000 | 100,000 | 500,000 |
| **Monthly Cost** | **$0.00** | **$0.04** | **$0.73** | **$3.93** |
| **Annual Cost** | **$0.00** | **$0.48** | **$8.76** | **$47.16** |
| Free Tier? | ✅ Yes | ✅ Yes | ✅ Yes | ⚠️ Partial |

**Recommendation**: Implement Firestore storage with UUID URLs. Cost is **negligible** for expected usage, and free tier covers most scenarios. Add optimizations (TTL, compression, deduplication) only if needed.

**ROI**: Shareable URLs and caching provide significant UX value for essentially **zero cost** 🎯
