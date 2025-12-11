import React from 'react';
import styles from './TypingIndicator.module.css';

/**
 * TypingIndicator component shows an animated indicator
 * when waiting for the next debate response.
 * 
 * @param {Object} props - Component props
 * @param {string} props.panelistName - Name of the panelist who is "typing"
 */
const TypingIndicator = ({ panelistName }) => {
  return (
    <div className={styles.container}>
      <div className={styles.indicator}>
        <span className={styles.dot}></span>
        <span className={styles.dot}></span>
        <span className={styles.dot}></span>
      </div>
      {panelistName && (
        <span className={styles.text}>{panelistName} is responding...</span>
      )}
    </div>
  );
};

export default TypingIndicator;
