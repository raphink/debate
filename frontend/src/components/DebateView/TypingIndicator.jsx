import React from 'react';
import styles from './TypingIndicator.module.css';

/**
 * TypingIndicator component shows an animated indicator
 * when waiting for the next debate response.
 */
const TypingIndicator = () => {
  return (
    <div className={styles.container}>
      <div className={styles.indicator}>
        <span className={styles.dot}></span>
        <span className={styles.dot}></span>
        <span className={styles.dot}></span>
      </div>
      <span className={styles.text}>Waiting for response...</span>
    </div>
  );
};

export default TypingIndicator;
