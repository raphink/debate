import React, { useState, useEffect, useMemo } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { fetchDebateHistory } from '../services/debateService';
import SearchInput from '../components/SearchInput/SearchInput';
import DebateCard from '../components/DebateCard/DebateCard';
import styles from './DebateHistory.module.css';

const DebateHistory = () => {
  const navigate = useNavigate();
  const [debates, setDebates] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  // Fetch debates on mount
  useEffect(() => {
    const loadDebates = async () => {
      try {
        setIsLoading(true);
        setError(null);
        const data = await fetchDebateHistory(100, 0); // Fetch first 100 debates
        setDebates(data.debates || []);
      } catch (err) {
        console.error('Failed to fetch debates:', err);
        setError(err.message);
      } finally {
        setIsLoading(false);
      }
    };

    loadDebates();
  }, []);

  // Filter debates based on search query
  const filteredDebates = useMemo(() => {
    if (!searchQuery.trim()) {
      return debates;
    }

    const query = searchQuery.toLowerCase();
    return debates.filter(
      (debate) =>
        debate.topic.toLowerCase().includes(query) ||
        debate.panelists.some((p) => p.name.toLowerCase().includes(query))
    );
  }, [debates, searchQuery]);

  const handleDebateClick = (debateId) => {
    navigate(`/d/${debateId}`);
  };

  if (isLoading) {
    return (
      <div className={styles.container}>
        <div className={styles.loading}>
          <div className={styles.spinner}></div>
          <p>Loading debates...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.container}>
        <div className={styles.error}>
          <h2>Failed to load debates</h2>
          <p>{error}</p>
          <button onClick={() => window.location.reload()}>Retry</button>
        </div>
      </div>
    );
  }

  const showEmptyState = debates.length === 0;
  const showNoResults = filteredDebates.length === 0 && searchQuery.trim();

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <Link to="/" className={styles.backLink}>
          ← Back to Home
        </Link>
        <h1 className={styles.title}>Debate History</h1>
      </header>

      {showEmptyState ? (
        <div className={styles.emptyState}>
          <div className={styles.emptyIcon}>💬</div>
          <h2>No debates yet</h2>
          <p>Create your first debate to see it here!</p>
          <Link to="/" className={styles.createButton}>
            Create Your First Debate
          </Link>
        </div>
      ) : (
        <>
          <div className={styles.searchSection}>
            <SearchInput
              value={searchQuery}
              onChange={setSearchQuery}
              placeholder="Search debates by topic or panelist..."
            />
            <div className={styles.resultsInfo}>
              Showing {filteredDebates.length} of {debates.length} debates
            </div>
          </div>

          {showNoResults ? (
            <div className={styles.noResults}>
              <p>No debates found matching &ldquo;{searchQuery}&rdquo;</p>
              <button onClick={() => setSearchQuery('')} className={styles.clearSearchButton}>
                Clear search
              </button>
            </div>
          ) : (
            <div className={styles.debatesGrid}>
              {filteredDebates.map((debate) => (
                <DebateCard key={debate.id} debate={debate} onClick={handleDebateClick} />
              ))}
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default DebateHistory;
