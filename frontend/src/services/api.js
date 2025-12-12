import axios from 'axios';
import { TIMEOUTS } from '../utils/constants';

/**
 * Configured Axios instance for API calls
 */
const apiClient = axios.create({
  timeout: TIMEOUTS.validation,
  headers: {
    'Content-Type': 'application/json',
  },
});

/**
 * Request interceptor - add any auth headers or logging
 */
apiClient.interceptors.request.use(
  (config) => {
    // Could add request logging here
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

/**
 * Response interceptor - handle errors globally
 */
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // Transform error to user-friendly format
    if (error.code === 'ECONNABORTED') {
      return Promise.reject({
        error: 'Request timed out. Please try again.',
        code: 'TIMEOUT',
        retryable: true,
      });
    }
    
    if (error.response) {
      // Server responded with error
      return Promise.reject(error.response.data);
    }
    
    if (error.request) {
      // Request made but no response
      return Promise.reject({
        error: 'Network error. Please check your connection.',
        code: 'NETWORK_ERROR',
        retryable: true,
      });
    }
    
    // Something else went wrong
    return Promise.reject({
      error: 'An unexpected error occurred.',
      code: 'UNKNOWN_ERROR',
      retryable: true,
    });
  }
);

/**
 * Retrieve a saved debate by UUID
 * @param {string} uuid - The debate UUID
 * @returns {Promise<Object>} The debate data
 */
export const getDebateById = async (uuid) => {
  const baseURL = process.env.REACT_APP_GET_DEBATE_URL || 'http://localhost:8084';
  const response = await fetch(`${baseURL}?id=${uuid}`);
  
  if (!response.ok) {
    if (response.status === 404) {
      throw new Error('Debate not found');
    }
    if (response.status === 400) {
      throw new Error('Invalid debate ID');
    }
    throw new Error('Failed to load debate');
  }
  
  return response.json();
};

export default apiClient;
