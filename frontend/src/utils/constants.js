// App constants
export const MAX_PANELISTS = 5;
export const MIN_PANELISTS = 2;
export const MIN_TOPIC_LENGTH = 10;
export const MAX_TOPIC_LENGTH = 500;
export const MAX_PANELIST_SUGGESTIONS = 20;

// API endpoints (configured via environment variables)
export const API_ENDPOINTS = {
  validateTopic: process.env.REACT_APP_VALIDATE_TOPIC_URL || 'http://localhost:8080',
  suggestPanelists: process.env.REACT_APP_SUGGEST_PANELISTS_URL || 'http://localhost:8081',
  generateDebate: process.env.REACT_APP_GENERATE_DEBATE_URL || 'http://localhost:8082',
};

// Timeouts (milliseconds)
export const TIMEOUTS = {
  validation: 30000,  // 30 seconds
  panelists: 30000,   // 30 seconds
  debate: 120000,     // 2 minutes for complete debate
};

// Debate session states
export const DEBATE_STATUS = {
  GENERATING: 'generating',
  COMPLETE: 'complete',
  ERROR: 'error',
};

// Error codes
export const ERROR_CODES = {
  NETWORK_ERROR: 'NETWORK_ERROR',
  TIMEOUT: 'TIMEOUT',
  RATE_LIMIT: 'RATE_LIMIT_EXCEEDED',
  INVALID_TOPIC: 'INVALID_TOPIC_LENGTH',
  API_ERROR: 'API_ERROR',
};
