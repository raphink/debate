import React from 'react';
import styles from './LoadingSpinner.module.css';

/**
 * Loading spinner component for async operations with engaging animation
 */
const LoadingSpinner = ({ size = 'medium', message = 'Loading...' }) => {
  const sizeClass = styles[`spinner${size.charAt(0).toUpperCase() + size.slice(1)}`];
  
  return (
    <div className={styles.spinnerContainer} role="status" aria-live="polite">
      <div className={`${styles.spinner} ${sizeClass}`} aria-hidden="true">
        <div className={styles.spinnerRing}></div>
        <div className={styles.spinnerRing}></div>
        <div className={styles.spinnerRing}></div>
        <div className={styles.spinnerDot}></div>
      </div>
      <span className={styles.spinnerMessage}>{message}</span>
    </div>
  );
};

export default LoadingSpinner;
