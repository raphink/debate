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
 * @param {Function} onKeyboardNav - Callback for arrow key navigation (direction: 'up' | 'down')
 */
const TopicAutocompleteDropdown = ({ 
  suggestions, 
  isLoading, 
  onSelect, 
  selectedIndex,
  onClose,
  onKeyboardNav
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
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
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
    const date = new Date(dateString);
    const now = new Date();
    const diffDays = Math.floor((now - date) / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) return 'Today';
    if (diffDays === 1) return 'Yesterday';
    if (diffDays < 7) return `${diffDays} days ago`;
    
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
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
            tabIndex={0}
            onClick={() => onSelect(debate)}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                onSelect(debate);
              }
            }}
            onMouseEnter={() => {
              // Optional: update selectedIndex on hover for keyboard consistency
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
                {formatDate(debate.startedAt)}
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
    startedAt: PropTypes.string.isRequired,
  })),
  isLoading: PropTypes.bool,
  onKeyboardNav: PropTypes.func,
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
