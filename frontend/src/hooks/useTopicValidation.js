import { useState, useCallback } from 'react';
import { validateTopic } from '../services/topicService';

/**
 * Custom hook for topic validation with streaming panelist results
 * Manages state for topic validation including loading, error, result, and progressive panelist loading
 */
const useTopicValidation = () => {
  const [isValidating, setIsValidating] = useState(false);
  const [validationResult, setValidationResult] = useState(null);
  const [panelists, setPanelists] = useState([]);
  const [error, setError] = useState(null);

  /**
   * Validates a topic with the backend API using streaming
   * @param {string} topic - The topic to validate
   * @param {string[]} suggestedNames - Optional array of suggested panelist names
   */
  const validate = useCallback(async (topic, suggestedNames = []) => {
    setIsValidating(true);
    setError(null);
    setValidationResult(null);
    setPanelists([]);

    try {
      await validateTopic(
        topic,
        suggestedNames,
        // onValidation callback
        (data) => {
          console.log('Setting validation result:', data);
          setValidationResult({
            isRelevant: data.isRelevant,
            message: data.message,
            topic: topic,
          });
          // If not relevant, stop loading immediately
          if (!data.isRelevant) {
            setIsValidating(false);
          }
        },
        // onPanelist callback
        (panelist) => {
          console.log('Adding panelist to state:', panelist);
          // First panelist means validation passed
          setValidationResult(prev => prev || {
            isRelevant: true,
            message: '',
            topic: topic,
          });
          setPanelists(prev => {
            const updated = [...prev, panelist];
            console.log('Panelists state updated, count:', updated.length);
            return updated;
          });
        },
        // onError callback
        (err) => {
          setError(err);
          setIsValidating(false);
        },
        // onComplete callback
        () => {
          setIsValidating(false);
        }
      );
    } catch (err) {
      setError(err);
      setIsValidating(false);
      throw err;
    }
  }, []);

  /**
   * Resets the validation state
   */
  const reset = useCallback(() => {
    setIsValidating(false);
    setValidationResult(null);
    setPanelists([]);
    setError(null);
  }, []);

  return {
    isValidating,
    validationResult,
    panelists,
    error,
    validate,
    reset,
  };
};

export default useTopicValidation;
