import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { MIN_TOPIC_LENGTH, MAX_TOPIC_LENGTH } from '../../utils/constants';
import { validateTopicLength, isTopicSafe } from '../../utils/validation';
import styles from './TopicInput.module.css';

const TopicInput = ({ onSubmit, isLoading }) => {
  const [topic, setTopic] = useState('');
  const [suggestedNames, setSuggestedNames] = useState('');
  const [clientError, setClientError] = useState(null);

  const handleInputChange = (e) => {
    const value = e.target.value;
    setTopic(value);
    setClientError(null);
  };

  const handleNamesChange = (e) => {
    const value = e.target.value;
    setSuggestedNames(value);
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

    // Parse suggested names (comma-separated, max 5)
    const names = suggestedNames
      .split(',')
      .map(n => n.trim())
      .filter(n => n.length > 0)
      .slice(0, 5);

    // Clear error and submit
    setClientError(null);
    onSubmit(topic, names);
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
          value={suggestedNames}
          onChange={handleNamesChange}
          placeholder="E.g., Martin Luther King Jr., Gandhi, Mother Teresa"
          className={styles.input}
          disabled={isLoading}
          aria-describedby="names-help"
        />
        <span id="names-help" className={styles.help}>
          Comma-separated names (up to 5). AI may include them if they have relevant views.
        </span>
      </div>

      <button
        type="submit"
        className={styles.submitButton}
        disabled={isLoading || topic.trim().length < MIN_TOPIC_LENGTH || isOverLimit}
        aria-busy={isLoading}
      >
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
