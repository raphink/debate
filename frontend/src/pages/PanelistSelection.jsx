import React, { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import PanelistGrid from '../components/PanelistGrid/PanelistGrid';
import PanelistSelector from '../components/PanelistSelector/PanelistSelector';
import ErrorMessage from '../components/common/ErrorMessage/ErrorMessage';
import usePanelistSelection from '../hooks/usePanelistSelection';
import styles from './PanelistSelection.module.css';

/**
 * PanelistSelection page allows users to browse and select panelists for the debate.
 * Receives panelist data from topic validation and manages selection state.
 * Can also receive pre-filled panelists from autocomplete selection.
 */
const PanelistSelection = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { topic, panelists, source, preFilled } = location.state || {};

  const {
    selectedPanelists,
    toggleSelection,
    clearSelection,
    isValidSelection,
    setSelection,
  } = usePanelistSelection();

  // Pre-fill panelists if coming from autocomplete
  useEffect(() => {
    if (source === 'autocomplete' && preFilled && preFilled.length > 0) {
      // Pre-select the panelists from the autocomplete selection
      setSelection(preFilled);
    }
  }, [source, preFilled, setSelection]);

  // Redirect if no panelists data available
  useEffect(() => {
    if (!panelists || panelists.length === 0) {
      navigate('/', { replace: true });
    }
  }, [panelists, navigate]);

  const handleProceed = () => {
    if (isValidSelection()) {
      navigate('/debate', {
        state: {
          topic,
          selectedPanelists,
        },
      });
    }
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
            onToggleSelection={toggleSelection}
          />
        </div>

        <aside className={styles.selectorSection}>
          <PanelistSelector
            selectedPanelists={selectedPanelists}
            onClear={clearSelection}
            onProceed={handleProceed}
          />
        </aside>
      </div>
    </div>
  );
};

export default PanelistSelection;
