import React from 'react';
import './LoadingSpinner.module.css';

/**
 * Loading spinner component for async operations
 */
const LoadingSpinner = ({ size = 'medium', message = 'Loading...' }) => {
  const sizeClass = `spinner-${size}`;
  
  return (
    <div className="spinner-container" role="status" aria-live="polite">
      <div className={`spinner ${sizeClass}`} aria-hidden="true"></div>
      <span className="spinner-message">{message}</span>
    </div>
  );
};

export default LoadingSpinner;
