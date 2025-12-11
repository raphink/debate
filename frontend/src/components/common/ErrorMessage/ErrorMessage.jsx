import React from 'react';
import { Button } from '../Button/Button';
import './ErrorMessage.module.css';

/**
 * Error message component with optional retry button
 */
export const ErrorMessage = ({ 
  message, 
  onRetry, 
  retryable = true,
  code,
}) => {
  return (
    <div className="error-container" role="alert" aria-live="assertive">
      <div className="error-icon" aria-hidden="true">⚠️</div>
      <div className="error-content">
        <h3 className="error-title">Error</h3>
        <p className="error-message">{message}</p>
        {code && <p className="error-code">Error code: {code}</p>}
      </div>
      {retryable && onRetry && (
        <Button onClick={onRetry} variant="primary">
          Try Again
        </Button>
      )}
    </div>
  );
};

export default ErrorMessage;
