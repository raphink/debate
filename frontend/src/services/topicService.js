import apiClient from './api';

const VALIDATE_TOPIC_URL = process.env.REACT_APP_VALIDATE_TOPIC_URL || 'http://localhost:8080';

/**
 * Validates a debate topic with the backend API
 * @param {string} topic - The topic to validate
 * @param {string[]} suggestedNames - Optional array of suggested panelist names (max 5)
 * @returns {Promise<{isRelevant: boolean, message: string, topic: string, suggestedPanelists: Array}>}
 */
export const validateTopic = async (topic, suggestedNames = []) => {
  const payload = { topic };
  if (suggestedNames && suggestedNames.length > 0) {
    payload.suggestedNames = suggestedNames.slice(0, 5);
  }
  const response = await apiClient.post(VALIDATE_TOPIC_URL, payload);
  return response.data;
};

export default {
  validateTopic,
};
