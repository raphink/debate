import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import useDebateLoader from '../hooks/useDebateLoader';
import DebateView from '../components/DebateView/DebateView';
import PanelistModal from '../components/DebateView/PanelistModal';
import LoadingSpinner from '../components/common/LoadingSpinner/LoadingSpinner';
import ErrorMessage from '../components/common/ErrorMessage/ErrorMessage';
import Button from '../components/common/Button/Button';
import styles from './DebateViewer.module.css';

/**
 * DebateViewer page - displays cached debates from shareable URLs
 */
const DebateViewer = () => {
  const { uuid } = useParams();
  const navigate = useNavigate();
  const { debate, loading, error, retry } = useDebateLoader(uuid);
  const [selectedPanelist, setSelectedPanelist] = useState(null);

  if (loading) {
    return (
      <div className={styles.container}>
        <LoadingSpinner />
        <p className={styles.loadingText}>Loading debate...</p>
      </div>
    );
  }

  if (error) {
    const isNotFound = error.includes('not found');
    return (
      <div className={styles.container}>
        <ErrorMessage 
          message={isNotFound ? 'Debate Not Found' : 'Failed to Load Debate'}
        />
        <p className={styles.errorDetail}>
          {isNotFound 
            ? 'This debate may have been deleted or the link is invalid.'
            : error}
        </p>
        <div className={styles.actions}>
          {!isNotFound && (
            <Button onClick={retry} variant="secondary">
              Retry
            </Button>
          )}
          <Button onClick={() => navigate('/')}>
            Create New Debate
          </Button>
        </div>
      </div>
    );
  }

  if (!debate) {
    return (
      <div className={styles.container}>
        <p>No debate data</p>
      </div>
    );
  }

  // Transform debate data to match DebateView expected format
  const transformedData = {
    topic: debate.topic.text,
    panelists: debate.panelists || [],
    messages: debate.messages || [],
  };

  const handlePanelistClick = (panelist) => {
    setSelectedPanelist(panelist);
  };

  const handleCloseModal = () => {
    setSelectedPanelist(null);
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <button
          onClick={() => navigate('/')}
          className={styles.backButton}
          aria-label="Go back to home"
        >
          ← Back
        </button>
        <div className={styles.topicDisplay}>
          <h1 className={styles.title}>{transformedData.topic}</h1>
          <div className={styles.panelists}>
            <span className={styles.panelistsLabel}>Panelists:</span>
            {transformedData.panelists.map((p, idx) => (
              <React.Fragment key={p.id}>
                <button
                  className={styles.panelistName}
                  onClick={() => handlePanelistClick(p)}
                >
                  {p.name}
                </button>
                {idx < transformedData.panelists.length - 1 && <span className={styles.separator}>, </span>}
              </React.Fragment>
            ))}
          </div>
        </div>
      </div>

      <div className={styles.viewerContainer}>
        <DebateView 
          topic={transformedData.topic}
          panelists={transformedData.panelists}
          messages={transformedData.messages}
          isComplete={debate.status === 'complete'}
          debateId={debate.id}
          mode="viewer"
        />
      </div>

      {selectedPanelist && (
        <PanelistModal 
          panelist={selectedPanelist} 
          onClose={handleCloseModal}
        />
      )}
    </div>
  );
};

export default DebateViewer;
