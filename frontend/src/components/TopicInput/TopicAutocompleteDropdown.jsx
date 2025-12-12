import React, { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import styles from './TopicAutocompleteDropdown.module.css';

export const TopicAutocompleteDropdown = ({
  suggestions,
  onSelect,
  loading,
  visible,
  onClose,
}) => {
  const [selectedIndex, setSelectedIndex] = useState(0);
  const dropdownRef = useRef(null);

  // Reset selection when suggestions change
  useEffect(() => {
    setSelectedIndex(0);
  }, [suggestions]);

  // Keyboard navigation
  useEffect(() => {
    if (!visible) return;

    const handleKeyDown = (e) => {
      if (!suggestions.length) return;

      switch (e.key) {
        case 'ArrowDown':
          e.preventDefault();
          setSelectedIndex((prev) =>
            prev < suggestions.length - 1 ? prev + 1 : 0
          );
          break;
        case 'ArrowUp':
          e.preventDefault();
          setSelectedIndex((prev) =>
            prev > 0 ? prev - 1 : suggestions.length - 1
          );
          break;
        case 'Enter':
          e.preventDefault();
          if (suggestions[selectedIndex]) {
            onSelect(suggestions[selectedIndex]);
          }
          break;
        case 'Escape':
          e.preventDefault();
          onClose();
          break;
        default:
          break;
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [visible, suggestions, selectedIndex, onSelect, onClose]);

  // Click outside to close
  useEffect(() => {
    if (!visible) return;

    const handleClickOutside = (e) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target)) {
        // Check if click is on the input field (allow it)
        const inputField = document.querySelector('input[type="text"]');
        if (inputField && inputField.contains(e.target)) {
          return;
        }
        onClose();
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [visible, onClose]);

  if (!visible || (!loading && !suggestions.length)) {
    return null;
  }

  return (
    <div className={styles.dropdown} ref={dropdownRef}>
      {loading && <div className={styles.loading}>Searching debates...</div>}
      {!loading && suggestions.length > 0 && (
        <ul className={styles.list} role="listbox">
          {suggestions.map((debate, index) => (
            <li
              key={debate.id}
              className={`${styles.item} ${
                index === selectedIndex ? styles.selected : ''
              }`}
              onClick={() => onSelect(debate)}
              onMouseEnter={() => setSelectedIndex(index)}
              onKeyPress={(e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  onSelect(debate);
                }
              }}
              role="option"
              aria-selected={index === selectedIndex}
              tabIndex={0}
            >
              <div className={styles.topic}>{debate.topic}</div>
              <div className={styles.badge}>
                {debate.panelistCount} panelist{debate.panelistCount !== 1 ? 's' : ''}
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

TopicAutocompleteDropdown.propTypes = {
  suggestions: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      topic: PropTypes.string.isRequired,
      panelists: PropTypes.array.isRequired,
      panelistCount: PropTypes.number.isRequired,
      createdAt: PropTypes.string.isRequired,
    })
  ).isRequired,
  onSelect: PropTypes.func.isRequired,
  loading: PropTypes.bool,
  visible: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
};
