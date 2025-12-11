import DOMPurify from 'dompurify';

/**
 * Sanitizes HTML content to prevent XSS attacks
 * @param {string} dirtyHTML - Potentially unsafe HTML string
 * @returns {string} - Sanitized HTML safe for rendering
 */
export const sanitizeHTML = (dirtyHTML) => {
  return DOMPurify.sanitize(dirtyHTML, {
    ALLOWED_TAGS: [], // No HTML tags allowed - text only
    ALLOWED_ATTR: [],
    KEEP_CONTENT: true,
  });
};

/**
 * Sanitizes text content from Claude API responses
 * @param {string} text - Raw text from API
 * @returns {string} - Sanitized text
 */
export const sanitizeAPIResponse = (text) => {
  if (!text || typeof text !== 'string') {
    return '';
  }
  
  // Remove any HTML tags
  const sanitized = sanitizeHTML(text);
  
  // Trim whitespace
  return sanitized.trim();
};

/**
 * Sanitizes user input before sending to API
 * @param {string} input - User input text
 * @returns {string} - Sanitized input
 */
export const sanitizeUserInput = (input) => {
  if (!input || typeof input !== 'string') {
    return '';
  }
  
  // Remove HTML tags and trim
  return sanitizeHTML(input).trim();
};
