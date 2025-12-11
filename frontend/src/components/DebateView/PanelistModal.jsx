import React, { useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import styles from './PanelistModal.module.css';

/**
 * PanelistModal component displays detailed information about a panelist.
 * Shows avatar, name, tagline, and full biography in a modal overlay.
 * Accessible with keyboard controls (Escape to close, focus trap).
 * 
 * @param {Object} props - Component props
 * @param {Object} props.panelist - Panelist object (id, name, avatarUrl, tagline, bio)
 * @param {Function} props.onClose - Callback to close the modal
 */
const PanelistModal = ({ panelist, onClose }) => {
  const modalRef = useRef(null);
  const closeButtonRef = useRef(null);

  // Focus close button when modal opens
  useEffect(() => {
    if (closeButtonRef.current) {
      closeButtonRef.current.focus();
    }
  }, []);

  // Handle Escape key to close
  useEffect(() => {
    const handleEscape = (e) => {
      if (e.key === 'Escape') {
        onClose();
      }
    };

    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, [onClose]);

  // Handle click outside modal
  const handleOverlayClick = (e) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  // Handle keyboard events on overlay
  const handleOverlayKeyDown = (e) => {
    if (e.key === 'Enter' || e.key === ' ') {
      if (e.target === e.currentTarget) {
        onClose();
      }
    }
  };

  if (!panelist) return null;

  return (
    <div 
      className={styles.overlay} 
      onClick={handleOverlayClick}
      onKeyDown={handleOverlayKeyDown}
      role="button"
      tabIndex={0}
      aria-label="Close modal"
    >
      <div className={styles.modal} ref={modalRef} role="dialog" aria-modal="true" aria-labelledby="modal-title">
        <button
          ref={closeButtonRef}
          className={styles.closeButton}
          onClick={onClose}
          aria-label="Close panelist details"
        >
          ✕
        </button>

        <div className={styles.content}>
          <img
            src={panelist.avatarUrl.startsWith('http') || panelist.avatarUrl.startsWith('/') ? panelist.avatarUrl : `${process.env.PUBLIC_URL}/avatars/${panelist.avatarUrl}`}
            alt={`${panelist.name} avatar`}
            className={styles.avatar}
            onError={(e) => {
              e.target.src = '/avatars/placeholder-avatar.svg';
            }}
          />
          
          <h2 id="modal-title" className={styles.name}>{panelist.name}</h2>
          
          {panelist.tagline && (
            <p className={styles.tagline}>{panelist.tagline}</p>
          )}
          
          {panelist.bio && (
            <p className={styles.bio}>{panelist.bio}</p>
          )}
        </div>
      </div>
    </div>
  );
};

PanelistModal.propTypes = {
  panelist: PropTypes.shape({
    id: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    avatarUrl: PropTypes.string.isRequired,
    tagline: PropTypes.string,
    bio: PropTypes.string,
  }),
  onClose: PropTypes.func.isRequired,
};

export default PanelistModal;
