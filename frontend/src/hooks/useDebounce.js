import { useState, useEffect } from 'react';

/**
 * Debounce hook - delays updating value until after delay period
 * Useful for autocomplete queries to prevent excessive API calls
 * 
 * @param {any} value - The value to debounce
 * @param {number} delay - Delay in milliseconds (default 300ms)
 * @returns {any} The debounced value
 * 
 * @example
 * const searchQuery = useDebounce(inputValue, 300);
 * 
 * useEffect(() => {
 *   if (searchQuery) {
 *     fetchResults(searchQuery);
 *   }
 * }, [searchQuery]);
 */
export const useDebounce = (value, delay = 300) => {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    // Set up timeout to update debounced value after delay
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    // Cleanup: cancel timeout if value changes before delay expires
    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
};

export default useDebounce;
