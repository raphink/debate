import React from 'react';
import './Button.module.css';

/**
 * Accessible button component with keyboard navigation support
 */
const Button = ({ 
  children, 
  onClick, 
  disabled = false, 
  variant = 'primary',
  type = 'button',
  ariaLabel,
  className = '',
}) => {
  const variantClass = variant === 'secondary' ? 'button-secondary' : 'button-primary';
  
  return (
    <button
      type={type}
      className={`button ${variantClass} ${className}`}
      onClick={onClick}
      disabled={disabled}
      aria-label={ariaLabel || undefined}
      aria-disabled={disabled}
    >
      {children}
    </button>
  );
};

export default Button;
