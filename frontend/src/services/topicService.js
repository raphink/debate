import apiClient from './api';

const VALIDATE_TOPIC_URL = process.env.REACT_APP_VALIDATE_TOPIC_URL || 'http://localhost:8080';

/**
 * Validates a debate topic with the backend API
 * @param {string} topic - The topic to validate
 * @returns {Promise<{isRelevant: boolean, message: string, topic: string}>}
 */
export const validateTopic = async (topic) => {
  try {
    const response = await apiClient.post(VALIDATE_TOPIC_URL, { topic });
    return response.data;
  } catch (error) {
    // Error is already processed by interceptor
    throw error;
  }
};

export default {
  validateTopic,
};
