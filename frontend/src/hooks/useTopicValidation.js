import { useState, useCallback } from 'react';
import { validateTopic } from '../services/topicService';

/**
 * Custom hook for topic validation
 * Manages state for topic validation including loading, error, and result
 */
const useTopicValidation = () => {
  const [isValidating, setIsValidating] = useState(false);
  const [validationResult, setValidationResult] = useState(null);
  const [error, setError] = useState(null);

  /**
   * Validates a topic with the backend API
   * @param {string} topic - The topic to validate
   * @param {string[]} suggestedNames - Optional array of suggested panelist names
   */
  const validate = useCallback(async (topic, suggestedNames = []) => {
    setIsValidating(true);
    setError(null);
    setValidationResult(null);

    try {
      const result = await validateTopic(topic, suggestedNames);
      setValidationResult(result);
      return result;
    } catch (err) {
      setError(err);
      throw err;
    } finally {
      setIsValidating(false);
    }
  }, []);

  /**
   * Resets the validation state
   */
  const reset = useCallback(() => {
    setIsValidating(false);
    setValidationResult(null);
    setError(null);
  }, []);

  return {
    isValidating,
    validationResult,
    error,
    validate,
    reset,
  };
};

export default useTopicValidation;
