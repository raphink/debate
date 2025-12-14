import { useState, useEffect } from 'react';
import { useDebounce } from './useDebounce';
import { fetchDebateHistory } from '../services/debateService';

/**
 * Topic autocomplete hook - manages autocomplete state and API calls
 * 
 * @param {string} query - The search query from user input
 * @param {boolean} enabled - Whether autocomplete is enabled (default true)
 * @returns {Object} { suggestions, isLoading, error }
 * 
 * @example
 * const { suggestions, isLoading, error } = useTopicAutocomplete(topicInput);
 */
export const useTopicAutocomplete = (query, enabled = true) => {
  const [suggestions, setSuggestions] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  // Debounce query to prevent excessive API calls
  const debouncedQuery = useDebounce(query, 300);

  useEffect(() => {
    // Reset state if autocomplete disabled or query too short
    if (!enabled || !debouncedQuery || debouncedQuery.length < 3) {
      setSuggestions([]);
      setIsLoading(false);
      setError(null);
      return;
    }

    // Set up abort controller for cleanup
    const abortController = new AbortController();
    let isMounted = true;

    const fetchSuggestions = async () => {
      try {
        setIsLoading(true);
        setError(null);

        const data = await fetchDebateHistory({ query: debouncedQuery });
        
        if (isMounted && !abortController.signal.aborted) {
          setSuggestions(data.debates || []);
          setIsLoading(false);
        }
      } catch (err) {
        if (isMounted && !abortController.signal.aborted) {
          console.error('Autocomplete error:', err);
          setError(err.message || 'Failed to load suggestions');
          setSuggestions([]);
          setIsLoading(false);
        }
      }
    };

    fetchSuggestions();

    // Cleanup function
    return () => {
      isMounted = false;
      abortController.abort();
    };
  }, [debouncedQuery, enabled]);

  return { suggestions, isLoading, error };
};

export default useTopicAutocomplete;
