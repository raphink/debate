import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import TopicInput from '../components/TopicInput/TopicInput';
import ValidationResult from '../components/ValidationResult/ValidationResult';
import LoadingSpinner from '../components/common/LoadingSpinner/LoadingSpinner';
import ErrorMessage from '../components/common/ErrorMessage/ErrorMessage';
import useTopicValidation from '../hooks/useTopicValidation';
import styles from './Home.module.css';

const Home = () => {
  const navigate = useNavigate();
  const { isValidating, validationResult, error, validate, reset } = useTopicValidation();

  const handleSubmit = async (topic) => {
    try {
      await validate(topic);
    } catch (err) {
      // Error is already set in the hook
      console.error('Validation error:', err);
    }
  };

  // Auto-navigate to panelist selection when topic is validated as relevant
  useEffect(() => {
    if (validationResult && validationResult.isRelevant) {
      navigate('/panelists', {
        state: {
          topic: validationResult.topic,
          panelists: validationResult.suggestedPanelists || [],
        },
      });
    }
  }, [validationResult, navigate]);

  const handleTryAgain = () => {
    reset();
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
        {!validationResult && !error && (
          <div className={styles.inputSection}>
            <TopicInput onSubmit={handleSubmit} isLoading={isValidating} />
            {isValidating && (
              <div className={styles.loadingContainer}>
                <LoadingSpinner />
                <p className={styles.loadingText}>Validating your topic...</p>
              </div>
            )}
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
