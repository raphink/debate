import React from 'react';
import { useNavigate } from 'react-router-dom';
import TopicInput from '../components/TopicInput/TopicInput';
import ValidationResult from '../components/ValidationResult/ValidationResult';
import LoadingSpinner from '../components/common/LoadingSpinner/LoadingSpinner';
import ErrorMessage from '../components/common/ErrorMessage/ErrorMessage';
import PanelistGrid from '../components/PanelistGrid/PanelistGrid';
import PanelistSelector from '../components/PanelistSelector/PanelistSelector';
import useTopicValidation from '../hooks/useTopicValidation';
import usePanelistSelection from '../hooks/usePanelistSelection';
import styles from './Home.module.css';

const Home = () => {
  const navigate = useNavigate();
  const { isValidating, validationResult, panelists, error, validate, reset } = useTopicValidation();
  const {
    selectedPanelists,
    toggleSelection,
    clearSelection,
    isValidSelection,
  } = usePanelistSelection();

  // Debug logging
  console.log('Home render - validationResult:', validationResult, 'panelists count:', panelists.length, 'isValidating:', isValidating);

  const handleSubmit = async (topic, suggestedNames = []) => {
    console.log('Home handleSubmit - topic:', topic, 'suggestedNames:', suggestedNames);
    try {
      await validate(topic, suggestedNames);
    } catch (err) {
      // Error is already set in the hook
      console.error('Validation error:', err);
    }
  };

  const handleTryAgain = () => {
    reset();
    clearSelection();
  };

  const handleProceedToDebate = () => {
    if (isValidSelection() && validationResult) {
      navigate('/debate', {
        state: {
          topic: validationResult.topic,
          selectedPanelists,
        },
      });
    }
  };

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1 className={styles.title}>AI-Powered Debate Generator</h1>
        <p className={styles.subtitle}>
          Generate engaging debates between historical theological and philosophical figures
        </p>
      </header>

      <main className={styles.main}>
        {/* Show input section only if not validating and no results yet */}
        {!isValidating && !validationResult && (
          <div className={styles.inputSection}>
            <TopicInput onSubmit={handleSubmit} isLoading={isValidating} />
          </div>
        )}

        {error && (
          <div className={styles.errorContainer}>
            <ErrorMessage
              message={error.error || 'Failed to validate topic. Please try again.'}
              retryable={error.retryable}
              onRetry={reset}
            />
          </div>
        )}

        {validationResult && !validationResult.isRelevant && (
          <ValidationResult
            isRelevant={validationResult.isRelevant}
            message={validationResult.message}
            topic={validationResult.topic}
            onTryAgain={handleTryAgain}
          />
        )}

        {(isValidating || (validationResult && validationResult.isRelevant)) && (
          <div className={styles.panelistSection}>
            <div className={styles.topicDisplay}>
              <h2 className={styles.sectionTitle}>Select Debate Panelists</h2>
              {validationResult && (
                <p className={styles.validatedTopic}>
                  Topic: <em>&ldquo;{validationResult.topic}&rdquo;</em>
                </p>
              )}
            </div>

            {!isValidating && panelists.length < 5 && panelists.length > 0 && (
              <ErrorMessage
                message={`Only ${panelists.length} panelists were suggested for this topic. You may want to refine your topic to get more diverse perspectives.`}
                type="warning"
              />
            )}

            {panelists.length > 0 && (
              <div className={styles.panelistContent}>
                <div className={styles.gridSection}>
                  <PanelistGrid
                    panelists={panelists}
                    selectedPanelists={selectedPanelists}
                    onToggleSelection={toggleSelection}
                  />
                  
                  {/* Show loading animation at the end of the list while streaming */}
                  {isValidating && (
                    <div className={styles.loadingAtEnd}>
                      <LoadingSpinner />
                      <p className={styles.loadingText}>Looking for more panelists...</p>
                    </div>
                  )}
                </div>

                {!isValidating && (
                  <aside className={styles.selectorSection}>
                    <PanelistSelector
                      selectedPanelists={selectedPanelists}
                      onClear={clearSelection}
                      onProceed={handleProceedToDebate}
                    />
                  </aside>
                )}
              </div>
            )}

            {/* Show loading when no panelists yet */}
            {isValidating && panelists.length === 0 && (
              <div className={styles.loadingContainer}>
                <LoadingSpinner />
                <p className={styles.loadingText}>Looking for panelists...</p>
              </div>
            )}
          </div>
        )}
      </main>

      <footer className={styles.footer}>
        <p className={styles.footerText}>
          Powered by Claude AI | For educational and entertainment purposes
        </p>
      </footer>
    </div>
  );
};

export default Home;
