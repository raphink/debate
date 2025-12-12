import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import UnifiedTopicInput from '../components/TopicInput/UnifiedTopicInput';
import ChipInput from '../components/ChipInput/ChipInput';
import { validateTopic } from '../services/topicService';
import styles from './Home.module.css';

const Home = () => {
  const navigate = useNavigate();
  const [topic, setTopic] = useState('');
  const [panelistChips, setPanelistChips] = useState([]);
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

  // Parse chip names into panelist objects
  const parsePanelistChips = (chips) => {
    if (chips.length === 0) return [];
    
    return chips.map(name => ({
      id: name.toLowerCase().replace(/\s+/g, '-'),
      name: name,
      slug: name.toLowerCase().replace(/\s+/g, '-'),
    }));
  };

  const handleFindPanelists = async () => {
    setIsValidating(true);
    setValidationError(null);

    try {
      // Validate topic
      const result = await validateTopic(topic);

      // Parse panelist chips if provided
      const panelists = parsePanelistChips(panelistChips);

      // Navigate to panelist selection page
      navigate('/select-panelists', {
        state: {
          topic: result.topic || topic,
          panelists: panelists,
          suggestedNames: result.suggestedNames || [],
          skipValidation: false,
        },
      });
    } catch (error) {
      console.error('Validation error:', error);
      setValidationError('Failed to validate topic. Please try again.');
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
              Panelists <span className={styles.optional}>(Optional)</span>
            </div>
            <ChipInput
              value={panelistChips}
              onChange={setPanelistChips}
              disabled={isValidating}
              max={5}
              placeholder="e.g., Albert Einstein"
            />
            <p className={styles.hint}>
              Type names separated by commas, Tab, or Enter (up to 5 panelists)
            </p>
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
