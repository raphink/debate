import React, { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { useTopicAutocomplete } from '../../hooks/useTopicAutocomplete';
import { TopicAutocompleteDropdown } from './TopicAutocompleteDropdown';
import styles from './UnifiedTopicInput.module.css';

const UnifiedTopicInput = ({
  value,
  onChange,
  onSelectSuggestion,
  error,
  disabled,
  placeholder,
}) => {
  const { suggestions, loading } = useTopicAutocomplete(value);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [showDropdown, setShowDropdown] = useState(false);
  const containerRef = useRef(null);

  useEffect(() => {
    setShowDropdown(value.length >= 3 && suggestions.length > 0);
    setSelectedIndex(-1);
  }, [value, suggestions]);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (containerRef.current && !containerRef.current.contains(event.target)) {
        setShowDropdown(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleKeyDown = (e) => {
    if (!showDropdown || suggestions.length === 0) return;

    switch (e.key) {
      case 'ArrowDown':
        setSelectedIndex((i) => Math.min(i + 1, suggestions.length - 1));
        e.preventDefault();
        break;
      case 'ArrowUp':
        setSelectedIndex((i) => Math.max(i - 1, -1));
        e.preventDefault();
        break;
      case 'Enter':
        if (selectedIndex >= 0) {
          onSelectSuggestion(suggestions[selectedIndex]);
          setShowDropdown(false);
          e.preventDefault();
        }
        break;
      case 'Escape':
        setShowDropdown(false);
        setSelectedIndex(-1);
        break;
      default:
        break;
    }
  };

  const handleSelect = (debate) => {
    onSelectSuggestion(debate);
    setShowDropdown(false);
  };

  return (
    <div className={styles.container} ref={containerRef}>
      <input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder={placeholder}
        disabled={disabled}
        className={`${styles.input} ${error ? styles.error : ''}`}
        role="combobox"
        aria-expanded={showDropdown}
        aria-controls="autocomplete-listbox"
        aria-activedescendant={
          selectedIndex >= 0 ? `option-${selectedIndex}` : undefined
        }
        aria-describedby={error ? 'input-error' : undefined}
      />
      {showDropdown && (
        <TopicAutocompleteDropdown
          suggestions={suggestions}
          loading={loading}
          selectedIndex={selectedIndex}
          onSelect={handleSelect}
          visible={showDropdown}
        />
      )}
      {error && (
        <div className={styles.errorMessage} id="input-error" role="alert">
          {error}
        </div>
      )}
    </div>
  );
};

UnifiedTopicInput.propTypes = {
  value: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
  onSelectSuggestion: PropTypes.func.isRequired,
  error: PropTypes.string,
  disabled: PropTypes.bool,
  placeholder: PropTypes.string,
};

UnifiedTopicInput.defaultProps = {
  error: null,
  disabled: false,
  placeholder: 'What should they debate?',
};

export default UnifiedTopicInput;
