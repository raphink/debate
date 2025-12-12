import React from 'react';
import PropTypes from 'prop-types';
import styles from './PanelistChip.module.css';

const PanelistChip = ({ panelist, onRemove, disabled }) => {
  const { id, name, avatarUrl } = panelist;

  const handleRemove = (e) => {
    e.stopPropagation();
    if (!disabled) {
      onRemove(id);
    }
  };

  const getAvatarSrc = () => {
    if (!avatarUrl) {
      return `${process.env.PUBLIC_URL}/avatars/placeholder-avatar.svg`;
    }
    if (avatarUrl.startsWith('http') || avatarUrl.startsWith('/')) {
      return avatarUrl;
    }
    return `${process.env.PUBLIC_URL}/avatars/${avatarUrl}`;
  };

  return (
    <div className={`${styles.chip} ${disabled ? styles.disabled : ''}`}>
      <img
        src={getAvatarSrc()}
        alt={`${name} avatar`}
        className={styles.avatar}
        onError={(e) => {
          e.target.src = `${process.env.PUBLIC_URL}/avatars/placeholder-avatar.svg`;
        }}
      />
      <span className={styles.name}>{name}</span>
      {!disabled && (
        <button
          type="button"
          className={styles.removeButton}
          onClick={handleRemove}
          aria-label={`Remove ${name}`}
          tabIndex={0}
        >
          ×
        </button>
      )}
    </div>
  );
};

PanelistChip.propTypes = {
  panelist: PropTypes.shape({
    id: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    avatarUrl: PropTypes.string,
  }).isRequired,
  onRemove: PropTypes.func.isRequired,
  disabled: PropTypes.bool,
};

PanelistChip.defaultProps = {
  disabled: false,
};

export default React.memo(PanelistChip);
