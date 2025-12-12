import React, { useState, useRef } from 'react';
import PropTypes from 'prop-types';
import styles from './ChipInput.module.css';

const ChipInput = ({ value, onChange, placeholder, disabled }) => {
  const [inputValue, setInputValue] = useState('');
  const inputRef = useRef(null);

  const handleKeyDown = (e) => {
    if (e.key === ',' || e.key === 'Tab' || e.key === 'Enter') {
      e.preventDefault();
      addChip();
    } else if (e.key === 'Backspace' && inputValue === '' && value.length > 0) {
      // Remove last chip on backspace if input is empty
      onChange(value.slice(0, -1));
    }
  };

  const addChip = () => {
    const trimmed = inputValue.trim();
    if (trimmed && !value.includes(trimmed)) {
      onChange([...value, trimmed]);
      setInputValue('');
    }
  };

  const removeChip = (chipToRemove) => {
    onChange(value.filter(chip => chip !== chipToRemove));
  };

  const handleInputChange = (e) => {
    setInputValue(e.target.value);
  };

  const handleContainerClick = () => {
    inputRef.current?.focus();
  };

  const handleContainerKeyDown = (e) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      inputRef.current?.focus();
    }
  };

  return (
    <div 
      className={`${styles.container} ${disabled ? styles.disabled : ''}`}
      onClick={handleContainerClick}
      onKeyDown={handleContainerKeyDown}
      role="button"
      tabIndex={0}
    >
      {value.map((chip, index) => (
        <div key={index} className={styles.chip}>
          <span className={styles.chipText}>{chip}</span>
          {!disabled && (
            <button
              type="button"
              className={styles.chipRemove}
              onClick={(e) => {
                e.stopPropagation();
                removeChip(chip);
              }}
              aria-label={`Remove ${chip}`}
            >
              ×
            </button>
          )}
        </div>
      ))}
      <input
        ref={inputRef}
        type="text"
        value={inputValue}
        onChange={handleInputChange}
        onKeyDown={handleKeyDown}
        onBlur={addChip}
        placeholder={value.length === 0 ? placeholder : ''}
        disabled={disabled}
        className={styles.input}
      />
    </div>
  );
};

ChipInput.propTypes = {
  value: PropTypes.arrayOf(PropTypes.string).isRequired,
  onChange: PropTypes.func.isRequired,
  placeholder: PropTypes.string,
  disabled: PropTypes.bool,
};

ChipInput.defaultProps = {
  placeholder: 'Type and press comma or tab',
  disabled: false,
};

export default ChipInput;
