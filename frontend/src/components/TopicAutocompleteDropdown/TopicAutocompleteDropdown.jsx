import React, { useRef, useEffect } from 'react';
import PropTypes from 'prop-types';
import styles from './TopicAutocompleteDropdown.module.css';

/**
 * TopicAutocompleteDropdown - displays autocomplete suggestions below topic input
 * Supports keyboard navigation (Arrow keys, Enter, Escape)
 * 
 * @param {Array} suggestions - Array of debate suggestions from API
 * @param {boolean} isLoading - Whether suggestions are loading
 * @param {Function} onSelect - Callback when user selects a suggestion
 * @param {number} selectedIndex - Currently focused suggestion index (for keyboard nav)
 * @param {Function} onClose - Callback to close dropdown
 */
const TopicAutocompleteDropdown = ({ 
  suggestions, 
  isLoading, 
  onSelect, 
  selectedIndex,
  onClose
}) => {
  const dropdownRef = useRef(null);
  const itemRefs = useRef([]);

  // Scroll selected item into view when keyboard navigation changes selection
  useEffect(() => {
    if (selectedIndex >= 0 && itemRefs.current[selectedIndex]) {
      itemRefs.current[selectedIndex].scrollIntoView({
        block: 'nearest',
        behavior: 'smooth'
      });
    }
  }, [selectedIndex]);

  // Handle click outside to close dropdown
  useEffect(() => {
    const handleClickOutside = (event) => {
      const dropdownEl = dropdownRef.current;
      const inputElement = dropdownEl?.parentElement?.querySelector('textarea');
      
      if (
        dropdownEl &&
        !dropdownEl.contains(event.target) &&
        (!inputElement || !inputElement.contains(event.target))
      ) {
        onClose();
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [onClose]);

  // Don't render if no suggestions and not loading
  if (!isLoading && (!suggestions || suggestions.length === 0)) {
    return null;
  }

  // Format date for display
  const formatDate = (dateString) => {
    // Zero out the time part for both dates to compare only the date
    const dateOnly = new Date(dateString);
    dateOnly.setHours(0, 0, 0, 0);
    
    const nowOnly = new Date();
    nowOnly.setHours(0, 0, 0, 0);
    
    const diffMs = nowOnly - dateOnly;
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) return 'Today';
    if (diffDays === 1) return 'Yesterday';
    if (diffDays > 1 && diffDays < 7) return `${diffDays} days ago`;
    
    // Handle future dates (invalid data)
    if (diffDays < 0) {
      console.warn('TopicAutocompleteDropdown: Detected future createdAt date:', dateString);
      return 'Invalid date';
    }
    
    return dateOnly.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
  };

  return (
    <div 
      ref={dropdownRef}
      className={styles.dropdown}
      role="listbox"
      aria-label="Topic suggestions"
    >
      {isLoading ? (
        <div className={styles.loadingContainer}>
          <div className={styles.spinner} />
          <span>Loading suggestions...</span>
        </div>
      ) : (
        suggestions.map((debate, index) => (
          <div
            key={debate.id}
            ref={(el) => (itemRefs.current[index] = el)}
            className={`${styles.suggestion} ${selectedIndex === index ? styles.selected : ''}`}
            role="option"
            aria-selected={selectedIndex === index}
            tabIndex={-1}
            onClick={() => onSelect(debate)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                e.preventDefault();
                onSelect(debate);
              }
            }}
          >
            <div className={styles.topicText}>{debate.topic}</div>
            
            <div className={styles.metadata}>
              <div className={styles.panelists}>
                {debate.panelists && debate.panelists.slice(0, 3).map((panelist, idx) => (
                  <img
                    key={panelist.id || idx}
                    src={panelist.avatarUrl || '/avatars/placeholder-avatar.svg'}
                    alt={panelist.name}
                    className={styles.avatar}
                    title={panelist.name}
                  />
                ))}
                {debate.panelists && debate.panelists.length > 3 && (
                  <span className={styles.moreCount}>+{debate.panelists.length - 3}</span>
                )}
              </div>
              
              <span className={styles.panelistCount}>
                {debate.panelists ? debate.panelists.length : 0} panelists
              </span>
              
              <span className={styles.date}>
                {formatDate(debate.createdAt)}
              </span>
            </div>
          </div>
        ))
      )}
    </div>
  );
};

TopicAutocompleteDropdown.propTypes = {
  suggestions: PropTypes.arrayOf(PropTypes.shape({
    id: PropTypes.string.isRequired,
    topic: PropTypes.string.isRequired,
    panelists: PropTypes.arrayOf(PropTypes.shape({
      id: PropTypes.string,
      name: PropTypes.string,
      avatarUrl: PropTypes.string,
    })),
    createdAt: PropTypes.string.isRequired,
  })),
  isLoading: PropTypes.bool,
  onSelect: PropTypes.func.isRequired,
  selectedIndex: PropTypes.number,
  onClose: PropTypes.func.isRequired,
};

TopicAutocompleteDropdown.defaultProps = {
  suggestions: [],
  isLoading: false,
  selectedIndex: -1,
};

export default TopicAutocompleteDropdown;
