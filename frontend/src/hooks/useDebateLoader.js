import { useState, useEffect } from 'react';
import { getDebateById } from '../services/api';

/**
 * Hook to load a debate from backend by UUID
 * @param {string} uuid - The debate UUID
 * @returns {Object} Loading state and debate data
 */
const useDebateLoader = (uuid) => {
  const [debate, setDebate] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!uuid) {
      setError('No debate ID provided');
      setLoading(false);
      return;
    }

    const loadDebate = async () => {
      try {
        setLoading(true);
        setError(null);
        const data = await getDebateById(uuid);
        setDebate(data);
      } catch (err) {
        setError(err.message || 'Failed to load debate');
      } finally {
        setLoading(false);
      }
    };

    loadDebate();
  }, [uuid]);

  const retry = () => {
    if (uuid) {
      setLoading(true);
      setError(null);
      getDebateById(uuid)
        .then(setDebate)
        .catch(err => setError(err.message))
        .finally(() => setLoading(false));
    }
  };

  return { debate, loading, error, retry };
};

export default useDebateLoader;
