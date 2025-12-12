import React, { useEffect, useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import PanelistGrid from '../components/PanelistGrid/PanelistGrid';
import PanelistSelector from '../components/PanelistSelector/PanelistSelector';
import ErrorMessage from '../components/common/ErrorMessage/ErrorMessage';
import usePanelistSelection from '../hooks/usePanelistSelection';
import { isCacheHit } from '../utils/cacheDetection';
import styles from './PanelistSelection.module.css';

/**
 * PanelistSelection page allows users to browse and select panelists for the debate.
 * Receives panelist data from topic validation and manages selection state.
 * Supports pre-filling panelists from historical debates via autocomplete.
 */
const PanelistSelection = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const {
    topic,
    panelists,
    debateId: originalDebateId,
    skipValidation,
  } = location.state || {};

  const {
    selectedPanelists,
    toggleSelection,
    clearSelection,
    isValidSelection,
    setSelection,
  } = usePanelistSelection();

  // Cache detection state
  const [isLocked, setIsLocked] = useState(false);
  const [showCacheIndicator, setShowCacheIndicator] = useState(false);
  const [originalDebateData, setOriginalDebateData] = useState(null);

  // Pre-fill panelists if provided from autocomplete
  useEffect(() => {
    if (panelists && panelists.length > 0 && originalDebateId) {
      // This is from autocomplete - pre-fill and lock
      setSelection(panelists);
      setIsLocked(true);
      setShowCacheIndicator(true);
      setOriginalDebateData({
        id: originalDebateId,
        topic,
        panelists,
      });
    }
  }, [panelists, originalDebateId, topic, setSelection]);

  // Redirect if no panelists data available
  useEffect(() => {
    if (!panelists || panelists.length === 0) {
      navigate('/', { replace: true });
    }
  }, [panelists, navigate]);

  const handleModifyPanelists = () => {
    setIsLocked(false);
    setShowCacheIndicator(false);
  };

  const handleProceed = () => {
    if (!isValidSelection()) return;

    // Check for cache hit if we have original debate data
    if (originalDebateData) {
      const cacheHit = isCacheHit(originalDebateData, topic, selectedPanelists);

      if (cacheHit) {
        // Load cached debate directly
        console.log('Cache hit detected, loading debate:', originalDebateData.id);
        navigate(`/d/${originalDebateData.id}`);
        return;
      }
    }

    // No cache hit or not from autocomplete - proceed with generation
    navigate('/debate', {
      state: {
        topic,
        selectedPanelists,
      },
    });
  };

  const handleBack = () => {
    navigate('/', { state: { topic } });
  };

  if (!panelists || panelists.length === 0) {
    return null; // Will redirect in useEffect
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <button
          onClick={handleBack}
          className={styles.backButton}
          aria-label="Go back to topic input"
        >
          ← Back
        </button>
        <div className={styles.topicDisplay}>
          <h1 className={styles.title}>Select Debate Panelists</h1>
          <p className={styles.topic}>Topic: <em>&ldquo;{topic}&rdquo;</em></p>
        </div>
      </div>

      {showCacheIndicator && (
        <div className={styles.cacheIndicator}>
          <span>✓ Using cached debate - panelists pre-filled from history</span>
        </div>
      )}

      {panelists.length < 5 && (
        <ErrorMessage
          message={`Only ${panelists.length} panelists were suggested for this topic. You may want to refine your topic to get more diverse perspectives.`}
          type="warning"
        />
      )}

      <div className={styles.content}>
        <div className={styles.gridSection}>
          <PanelistGrid
            panelists={panelists}
            selectedPanelists={selectedPanelists}
            onToggleSelection={isLocked ? () => {} : toggleSelection}
            isLocked={isLocked}
          />
        </div>

        <aside className={styles.selectorSection}>
          {isLocked && (
            <button
              onClick={handleModifyPanelists}
              className={styles.modifyButton}
            >
              Modify Panelists
            </button>
          )}
          <PanelistSelector
            selectedPanelists={selectedPanelists}
            onClear={isLocked ? undefined : clearSelection}
            onProceed={handleProceed}
            proceedButtonText={showCacheIndicator ? 'Load Debate' : 'Generate Debate'}
          />
        </aside>
      </div>
    </div>
  );
};

export default PanelistSelection;
