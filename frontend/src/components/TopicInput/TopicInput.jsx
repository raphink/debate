import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { MIN_TOPIC_LENGTH, MAX_TOPIC_LENGTH } from '../../utils/constants';
import { validateTopicLength, isTopicSafe } from '../../utils/validation';
import styles from './TopicInput.module.css';

const TopicInput = ({ onSubmit, isLoading }) => {
  const [topic, setTopic] = useState('');
  const [clientError, setClientError] = useState(null);

  const handleInputChange = (e) => {
    const value = e.target.value;
    setTopic(value);
    setClientError(null);
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
    onSubmit(topic);
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

      <button
        type="submit"
        className={styles.submitButton}
        disabled={isLoading || topic.trim().length < MIN_TOPIC_LENGTH || isOverLimit}
        aria-busy={isLoading}
      >
        {isLoading ? 'Validating...' : 'Validate Topic'}
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
