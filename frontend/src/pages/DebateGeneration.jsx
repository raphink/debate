import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import DebateView from '../components/DebateView/DebateView';
import PanelistModal from '../components/DebateView/PanelistModal';
import PDFExport from '../components/PDFExport/PDFExport';
import Button from '../components/common/Button/Button';
import ErrorMessage from '../components/common/ErrorMessage/ErrorMessage';
import useDebateStream from '../hooks/useDebateStream';
import styles from './DebateGeneration.module.css';

/**
 * DebateGeneration page displays the streaming debate conversation.
 * Receives topic and selected panelists from navigation state.
 */
const DebateGeneration = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { topic, selectedPanelists } = location.state || {};
  const [selectedPanelistForModal, setSelectedPanelistForModal] = useState(null);

  const {
    messages,
    isStreaming,
    error,
    isComplete,
    currentPanelistId,
    debateId,
    startDebate,
    retry,
    reset,
  } = useDebateStream();

  // Redirect if no topic or panelists
  useEffect(() => {
    if (!topic || !selectedPanelists || selectedPanelists.length < 2) {
      navigate('/', { replace: true });
    }
  }, [topic, selectedPanelists, navigate]);

  // Auto-start debate on mount
  useEffect(() => {
    if (topic && selectedPanelists && selectedPanelists.length >= 2 && !isStreaming && !isComplete && !error) {
      startDebate(topic, selectedPanelists);
    }
  }, [topic, selectedPanelists, startDebate, isStreaming, isComplete, error]);

  const handleBack = () => {
    navigate('/', { state: { topic } });
  };

  const handleRetry = () => {
    retry(topic, selectedPanelists);
  };

  const handlePanelistClick = (panelist) => {
    setSelectedPanelistForModal(panelist);
  };

  const handleCloseModal = () => {
    setSelectedPanelistForModal(null);
  };

  if (!topic || !selectedPanelists) {
    return null; // Will redirect in useEffect
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <button
          onClick={handleBack}
          className={styles.backButton}
          aria-label="Go back to home"
        >
          ← Back
        </button>
        <div className={styles.topicDisplay}>
          <h1 className={styles.title}>Debate in Progress</h1>
          <p className={styles.topic}>Topic: <em>&ldquo;{topic}&rdquo;</em></p>
          <div className={styles.panelists}>
            <span className={styles.panelistsLabel}>Panelists:</span>
            {selectedPanelists.map((p, idx) => (
              <React.Fragment key={p.id}>
                <button
                  className={styles.panelistName}
                  onClick={() => handlePanelistClick(p)}
                >
                  {p.name}
                </button>
                {idx < selectedPanelists.length - 1 && <span className={styles.separator}>, </span>}
              </React.Fragment>
            ))}
          </div>
        </div>
      </div>

      {error && (
        <div className={styles.errorContainer}>
          <ErrorMessage
            message={error}
            retryable={true}
            onRetry={handleRetry}
          />
        </div>
      )}

      <div className={styles.debateContainer}>
        <DebateView
          messages={messages}
          panelists={selectedPanelists}
          isStreaming={isStreaming}
          currentPanelistId={currentPanelistId}
          debateId={debateId}
          isComplete={isComplete}
          mode="generation"
        />

        {isComplete && (
          <div className={styles.completeActions}>
            <p className={styles.completeMessage}>Debate concluded</p>
            <div className={styles.actionButtons}>
              <Button onClick={handleBack} variant="secondary">
                New Debate
              </Button>
            </div>
            <PDFExport 
              topic={topic}
              panelists={selectedPanelists}
              messages={messages}
            />
          </div>
        )}
      </div>

      {selectedPanelistForModal && (
        <PanelistModal 
          panelist={selectedPanelistForModal} 
          onClose={handleCloseModal}
        />
      )}
    </div>
  );
};

export default DebateGeneration;
