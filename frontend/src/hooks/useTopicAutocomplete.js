import { useState, useEffect } from 'react';
import { autocompleteTopics } from '../services/api';

/**
 * Hook for topic autocomplete with debouncing
 * @param {string} query - Search query
 * @param {boolean} enabled - Whether autocomplete is enabled
 * @returns {Object} { suggestions, loading, error }
 */
export const useTopicAutocomplete = (query, enabled = true) => {
  const [suggestions, setSuggestions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Reset if query too short or disabled
    if (!enabled || !query || query.length < 3) {
      setSuggestions([]);
      setLoading(false);
      return;
    }

    // Debounce 300ms
    const timer = setTimeout(async () => {
      setLoading(true);
      setError(null);

      try {
        const data = await autocompleteTopics(query, 10);
        setSuggestions(data.debates || []);
      } catch (err) {
        console.error('Autocomplete error:', err);
        setError(err);
        setSuggestions([]);
      } finally {
        setLoading(false);
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [query, enabled]);

  return { suggestions, loading, error };
};
