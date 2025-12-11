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
  const { isValidating, validationResult, error, validate, reset } = useTopicValidation();
  const {
    selectedPanelists,
    toggleSelection,
    clearSelection,
    isValidSelection,
  } = usePanelistSelection();

  const handleSubmit = async (topic) => {
    try {
      await validate(topic);
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
        <div className={styles.inputSection}>
          <TopicInput onSubmit={handleSubmit} isLoading={isValidating} />
          {isValidating && (
            <div className={styles.loadingContainer}>
              <LoadingSpinner />
              <p className={styles.loadingText}>Validating your topic...</p>
            </div>
          )}
        </div>

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

        {validationResult && validationResult.isRelevant && validationResult.suggestedPanelists && (
          <div className={styles.panelistSection}>
            <div className={styles.topicDisplay}>
              <h2 className={styles.sectionTitle}>Select Debate Panelists</h2>
              <p className={styles.validatedTopic}>
                Topic: <em>&ldquo;{validationResult.topic}&rdquo;</em>
              </p>
            </div>

            {validationResult.suggestedPanelists.length < 5 && (
              <ErrorMessage
                message={`Only ${validationResult.suggestedPanelists.length} panelists were suggested for this topic. You may want to refine your topic to get more diverse perspectives.`}
                type="warning"
              />
            )}

            <div className={styles.panelistContent}>
              <div className={styles.gridSection}>
                <PanelistGrid
                  panelists={validationResult.suggestedPanelists}
                  selectedPanelists={selectedPanelists}
                  onToggleSelection={toggleSelection}
                />
              </div>

              <aside className={styles.selectorSection}>
                <PanelistSelector
                  selectedPanelists={selectedPanelists}
                  onClear={clearSelection}
                  onProceed={handleProceedToDebate}
                />
              </aside>
            </div>
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
