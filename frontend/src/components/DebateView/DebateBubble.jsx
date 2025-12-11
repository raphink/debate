import React from 'react';
import PropTypes from 'prop-types';
import styles from './DebateBubble.module.css';

/**
 * DebateBubble component displays a single debate message from a panelist.
 * Shows avatar, name, and message text in a chat-style bubble.
 * 
 * @param {Object} props - Component props
 * @param {Object} props.panelist - Panelist object (id, name, avatarUrl, tagline, bio)
 * @param {string} props.text - The message text
 * @param {number} props.index - Message index in the conversation
 * @param {Function} props.onAvatarClick - Callback when avatar is clicked
 */
const DebateBubble = ({ panelist, text, index, onAvatarClick }) => {
  const isEven = index % 2 === 0;

  const handleAvatarClick = () => {
    if (onAvatarClick) {
      onAvatarClick(panelist);
    }
  };

  return (
    <div className={`${styles.container} ${isEven ? styles.left : styles.right}`}>
      <button
        className={styles.avatarButton}
        onClick={handleAvatarClick}
        aria-label={`View ${panelist.name} details`}
      >
        <img
          src={panelist.avatarUrl}
          alt={`${panelist.name} avatar`}
          className={styles.avatar}
          onError={(e) => {
            e.target.src = '/avatars/placeholder-avatar.svg';
          }}
        />
      </button>
      <div className={styles.bubble}>
        <div className={styles.header}>
          <span className={styles.name}>{panelist.name}</span>
        </div>
        <p className={styles.text}>{text}</p>
      </div>
    </div>
  );
};

DebateBubble.propTypes = {
  panelist: PropTypes.shape({
    id: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    avatarUrl: PropTypes.string.isRequired,
    tagline: PropTypes.string,
    bio: PropTypes.string,
  }).isRequired,
  text: PropTypes.string.isRequired,
  index: PropTypes.number.isRequired,
  onAvatarClick: PropTypes.func,
};

export default DebateBubble;
