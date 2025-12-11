import { MIN_TOPIC_LENGTH, MAX_TOPIC_LENGTH } from './constants';

/**
 * Validates topic length (client-side validation)
 * @param {string} topic - The topic to validate
 * @returns {{isValid: boolean, error: string|null}}
 */
export const validateTopicLength = (topic) => {
  const trimmed = topic.trim();
  
  if (trimmed.length < MIN_TOPIC_LENGTH) {
    return {
      isValid: false,
      error: `Topic must be at least ${MIN_TOPIC_LENGTH} characters long`,
    };
  }
  
  if (trimmed.length > MAX_TOPIC_LENGTH) {
    return {
      isValid: false,
      error: `Topic must not exceed ${MAX_TOPIC_LENGTH} characters`,
    };
  }
  
  return {
    isValid: true,
    error: null,
  };
};

/**
 * Checks if topic contains potentially dangerous characters
 * @param {string} topic - The topic to check
 * @returns {boolean} - True if topic appears safe
 */
export const isTopicSafe = (topic) => {
  // Check for HTML tags
  if (topic.includes('<') || topic.includes('>')) {
    return false;
  }
  
  // Check for common HTML entities
  if (topic.includes('&lt;') || topic.includes('&gt;') || 
      topic.includes('&#') || topic.includes('&amp;')) {
    return false;
  }
  
  return true;
};

/**
 * Validates panelist selection
 * @param {Array} selectedPanelists - Array of selected panelist IDs
 * @param {number} min - Minimum number of panelists
 * @param {number} max - Maximum number of panelists
 * @returns {{isValid: boolean, error: string|null}}
 */
export const validatePanelistSelection = (selectedPanelists, min, max) => {
  if (selectedPanelists.length < min) {
    return {
      isValid: false,
      error: `Please select at least ${min} panelists`,
    };
  }
  
  if (selectedPanelists.length > max) {
    return {
      isValid: false,
      error: `You can select up to ${max} panelists`,
    };
  }
  
  return {
    isValid: true,
    error: null,
  };
};
