import React from 'prop-types';
import styles from './ValidationResult.module.css';

const ValidationResult = ({ isRelevant, message, topic, onContinue, onTryAgain }) => {
  return (
    <div className={`${styles.container} ${isRelevant ? styles.success : styles.failure}`}>
      <div className={styles.icon} aria-hidden="true">
        {isRelevant ? '✓' : '×'}
      </div>
      
      <div className={styles.content}>
        <h2 className={styles.title}>
          {isRelevant ? 'Topic Validated!' : 'Topic Not Suitable'}
        </h2>
        
        <p className={styles.message}>{message}</p>
        
        {topic && (
          <blockquote className={styles.topic}>
            "{topic}"
          </blockquote>
        )}
      </div>
      
      <div className={styles.actions}>
        {isRelevant ? (
          <button
            className={styles.primaryButton}
            onClick={onContinue}
            aria-label="Continue to select panelists"
          >
            Select Panelists →
          </button>
        ) : (
          <button
            className={styles.secondaryButton}
            onClick={onTryAgain}
            aria-label="Try a different topic"
          >
            Try Another Topic
          </button>
        )}
      </div>
    </div>
  );
};

ValidationResult.propTypes = {
  isRelevant: PropTypes.bool.isRequired,
  message: PropTypes.string.isRequired,
  topic: PropTypes.string,
  onContinue: PropTypes.func,
  onTryAgain: PropTypes.func,
};

ValidationResult.defaultProps = {
  topic: '',
  onContinue: () => {},
  onTryAgain: () => {},
};

export default ValidationResult;
