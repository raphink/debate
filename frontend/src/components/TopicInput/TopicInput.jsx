import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { MIN_TOPIC_LENGTH, MAX_TOPIC_LENGTH } from '../../utils/constants';
import { validateTopicLength, isTopicSafe } from '../../utils/validation';
import styles from './TopicInput.module.css';

const TopicInput = ({ onSubmit, isLoading }) => {
  const [topic, setTopic] = useState('');
  const [suggestedNames, setSuggestedNames] = useState([]);
  const [nameInput, setNameInput] = useState('');
  const [clientError, setClientError] = useState(null);

  const handleInputChange = (e) => {
    const value = e.target.value;
    setTopic(value);
    setClientError(null);
  };

  const handleNameInputChange = (e) => {
    const value = e.target.value;
    
    // Check if user typed comma followed by space or just comma at the end
    if (value.endsWith(', ') || (value.endsWith(',') && value.length > 1)) {
      const newName = value.replace(/,\s*$/, '').trim();
      if (newName && suggestedNames.length < 5 && !suggestedNames.includes(newName)) {
        setSuggestedNames([...suggestedNames, newName]);
        setNameInput('');
      } else {
        setNameInput('');
      }
    } else {
      setNameInput(value);
    }
  };

  const handleNameInputKeyDown = (e) => {
    // Also handle Enter key to add a name
    if (e.key === 'Enter' && nameInput.trim()) {
      e.preventDefault();
      const newName = nameInput.trim();
      if (suggestedNames.length < 5 && !suggestedNames.includes(newName)) {
        setSuggestedNames([...suggestedNames, newName]);
        setNameInput('');
      }
    }
    // Handle backspace on empty input to remove last chip
    else if (e.key === 'Backspace' && !nameInput && suggestedNames.length > 0) {
      setSuggestedNames(suggestedNames.slice(0, -1));
    }
  };

  const removeNameChip = (nameToRemove) => {
    setSuggestedNames(suggestedNames.filter(name => name !== nameToRemove));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    
    // Client-side validation
    const validation = validateTopicLength(topic);
    if (!validation.isValid) {
      setClientError(validation.error);
      return;
    }

    if (!isTopicSafe(topic)) {
      setClientError('Topic contains invalid characters');
      return;
    }

    // Clear error and submit
    setClientError(null);
    onSubmit(topic, suggestedNames);
  };

  const charactersRemaining = MAX_TOPIC_LENGTH - topic.length;
  const isNearLimit = charactersRemaining < 50;
  const isOverLimit = charactersRemaining < 0;

  return (
    <form onSubmit={handleSubmit} className={styles.form}>
      <div className={styles.inputGroup}>
        <label htmlFor="topic-input" className={styles.label}>
          Enter your debate topic
        </label>
        <textarea
          id="topic-input"
          value={topic}
          onChange={handleInputChange}
          placeholder="E.g., Should Christians defy authorities when the law is unfair?"
          className={styles.textarea}
          disabled={isLoading}
          rows={4}
          aria-describedby="character-count topic-help"
          aria-invalid={clientError ? 'true' : 'false'}
        />
        <div className={styles.meta}>
          <span id="topic-help" className={styles.help}>
            Topics related to theology, philosophy, ethics, or morality
          </span>
          <span
            id="character-count"
            className={`${styles.count} ${isNearLimit ? styles.warning : ''} ${isOverLimit ? styles.error : ''}`}
            aria-live="polite"
          >
            {charactersRemaining} characters remaining
          </span>
        </div>
        {clientError && (
          <div className={styles.errorMessage} role="alert">
            {clientError}
          </div>
        )}
      </div>

      <div className={styles.inputGroup}>
        <label htmlFor="suggested-names" className={styles.label}>
          Suggest panelists (optional)
        </label>
        <input
          id="suggested-names"
          type="text"
      <div className={styles.inputGroup}>
        <label htmlFor="suggested-names" className={styles.label}>
          Suggest panelists (optional)
        </label>
        <div className={styles.chipsContainer}>
          {suggestedNames.map((name) => (
            <div key={name} className={styles.chip}>
              <span className={styles.chipText}>{name}</span>
              <button
                type="button"
                onClick={() => removeNameChip(name)}
                className={styles.chipRemove}
                aria-label={`Remove ${name}`}
                disabled={isLoading}
              >
                ×
              </button>
            </div>
          ))}
          <input
            id="suggested-names"
            type="text"
            value={nameInput}
            onChange={handleNameInputChange}
            onKeyDown={handleNameInputKeyDown}
            placeholder={suggestedNames.length === 0 ? "E.g., Martin Luther King Jr., Gandhi, Mother Teresa" : "Add another..."}
            className={styles.chipInput}
            disabled={isLoading || suggestedNames.length >= 5}
            aria-describedby="names-help"
          />
        </div>
        <span id="names-help" className={styles.help}>
          Type a name and press comma+space or Enter to add (up to 5). AI will prioritize them if relevant.
        </span>
      </div>
        {isLoading ? 'Looking for Panelists...' : 'Find Panelists'}
      </button>
    </form>
  );
};

TopicInput.propTypes = {
  onSubmit: PropTypes.func.isRequired,
  isLoading: PropTypes.bool,
};

TopicInput.defaultProps = {
  isLoading: false,
};

export default TopicInput;
