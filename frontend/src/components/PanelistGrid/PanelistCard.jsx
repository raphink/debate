import React from 'react';
import PropTypes from 'prop-types';
import styles from './PanelistCard.module.css';

/**
 * PanelistCard component displays a single panelist with avatar, name, tagline, bio, and position.
 * Supports selection state and keyboard accessibility.
 * 
 * @param {Object} props - Component props
 * @param {Object} props.panelist - Panelist data object
 * @param {string} props.panelist.id - Unique panelist identifier
 * @param {string} props.panelist.name - Full name of the panelist
 * @param {string} props.panelist.tagline - Brief descriptor (max 150 chars)
 * @param {string} props.panelist.bio - Detailed background (max 500 chars)
 * @param {string} props.panelist.avatarUrl - Path to avatar image
 * @param {string} props.panelist.position - Stance on the topic (max 200 chars)
 * @param {boolean} props.isSelected - Whether this panelist is currently selected
 * @param {Function} props.onToggle - Callback when panelist is clicked/selected
 * @param {boolean} props.disabled - Whether selection is disabled (e.g., max limit reached)
 */
const PanelistCard = ({ panelist, isSelected, onToggle, disabled }) => {
  const { id, name, tagline, bio, avatarUrl, position } = panelist;

  const handleClick = () => {
    if (!disabled || isSelected) {
      onToggle(panelist);
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      handleClick();
    }
  };

  return (
    <div
      className={`${styles.card} ${isSelected ? styles.selected : ''} ${disabled && !isSelected ? styles.disabled : ''}`}
      onClick={handleClick}
      onKeyDown={handleKeyDown}
      role="button"
      tabIndex={0}
      aria-pressed={isSelected}
      aria-label={`${isSelected ? 'Deselect' : 'Select'} ${name} for debate`}
    >
      <div className={styles.header}>
        <img
          src={avatarUrl.startsWith('http') || avatarUrl.startsWith('/') ? avatarUrl : `${process.env.PUBLIC_URL}/avatars/${avatarUrl}`}
          alt={`${name} avatar`}
          className={styles.avatar}
          onError={(e) => {
            e.target.src = `${process.env.PUBLIC_URL}/avatars/placeholder-avatar.svg`;
          }}
        />
        <div className={styles.headerText}>
          <h3 className={styles.name}>{name}</h3>
          <p className={styles.tagline}>{tagline}</p>
        </div>
        {isSelected && (
          <div className={styles.checkmark} aria-hidden="true">
            ✓
          </div>
        )}
      </div>
      
      <p className={styles.bio}>{bio}</p>
      
      {position && (
        <div className={styles.position}>
          <strong>Position:</strong> {position}
        </div>
      )}
    </div>
  );
};

PanelistCard.propTypes = {
  panelist: PropTypes.shape({
    id: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    tagline: PropTypes.string.isRequired,
    bio: PropTypes.string.isRequired,
    avatarUrl: PropTypes.string.isRequired,
    position: PropTypes.string,
  }).isRequired,
  isSelected: PropTypes.bool.isRequired,
  onToggle: PropTypes.func.isRequired,
  disabled: PropTypes.bool,
};

PanelistCard.defaultProps = {
  disabled: false,
};

export default PanelistCard;
