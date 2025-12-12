import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import UnifiedTopicInput from '../components/TopicInput/UnifiedTopicInput';
import PanelistChipSelector from '../components/PanelistChips/PanelistChipSelector';
import { validateTopic } from '../services/topicService';
import styles from './Home.module.css';

const Home = () => {
  const navigate = useNavigate();
  const [topic, setTopic] = useState('');
  const [selectedPanelists, setSelectedPanelists] = useState([]);
  const [validationError, setValidationError] = useState(null);
  const [isValidating, setIsValidating] = useState(false);

  const handleAutocompleteSelect = (debate) => {
    // Navigate to panelist selection with pre-filled data from autocomplete
    navigate('/select-panelists', {
      state: {
        debateId: debate.id,
        topic: debate.topic,
        panelists: debate.panelists,
        skipValidation: true,
      },
    });
  };

  const handleFindPanelists = async () => {
    setIsValidating(true);
    setValidationError(null);

    try {
      // Validate topic with optional panelist context
      const result = await validateTopic(topic, selectedPanelists);

      if (result.isRelevant) {
        // Navigate to panelist selection page
        navigate('/select-panelists', {
          state: {
            topic: result.topic,
            panelists: selectedPanelists,
            suggestedNames: result.suggestedNames,
            skipValidation: false,
          },
        });
      } else {
        setValidationError(result.validationMessage || 'This topic is not suitable for a debate. Please try another.');
      }
    } catch (error) {
      console.error('Validation error:', error);
      setValidationError('Failed to validate topic. Please try again.');
    } finally {
      setIsValidating(false);
    }
  };

  const isSubmitDisabled = topic.trim().length < 3 || isValidating;

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1 className={styles.title}>AI-Powered Debate Generator</h1>
        <p className={styles.subtitle}>
          Search previous debates or create a new topic to generate engaging discussions
        </p>
      </header>

      <main className={styles.main}>
        <div className={styles.inputSection}>
          <div className={styles.formGroup}>
            <label htmlFor="topic-input" className={styles.label}>
              Debate Topic
            </label>
            <UnifiedTopicInput
              value={topic}
              onChange={setTopic}
              onSelectSuggestion={handleAutocompleteSelect}
              error={validationError}
              disabled={isValidating}
              placeholder="What should they debate?"
            />
          </div>

          <div className={styles.formGroup}>
            <div className={styles.label}>
              Select Panelists <span className={styles.optional}>(Optional)</span>
            </div>
            <PanelistChipSelector
              value={selectedPanelists}
              onChange={setSelectedPanelists}
              disabled={isValidating || topic.trim().length < 3}
              max={10}
            />
          </div>

          <button
            onClick={handleFindPanelists}
            disabled={isSubmitDisabled}
            className={styles.submitButton}
          >
            {isValidating ? 'Validating...' : 'Find Panelists →'}
          </button>
        </div>
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
