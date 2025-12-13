import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { formatRelativeTime, truncate } from '../../utils/formatters';
import PanelistModal from '../DebateView/PanelistModal';
import styles from './DebateCard.module.css';

const DebateCard = ({ debate, onClick }) => {
  const [selectedPanelist, setSelectedPanelist] = useState(null);

  const handleCardClick = (e) => {
    // Only trigger if clicking the card itself, not an avatar
    if (e.target.closest(`.${styles.avatar}`)) {
      return;
    }
    onClick(debate.id);
  };

  const handleCardKeyDown = (e) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      onClick(debate.id);
    }
  };

  const handleAvatarClick = (e, panelist) => {
    e.stopPropagation();
    setSelectedPanelist(panelist);
  };

  const handleAvatarKeyDown = (e, panelist) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      e.stopPropagation();
      setSelectedPanelist(panelist);
    }
  };

  const closeModal = () => {
    setSelectedPanelist(null);
  };

  return (
    <>
      <div
        className={styles.card}
        onClick={handleCardClick}
        onKeyDown={handleCardKeyDown}
        role="button"
        tabIndex={0}
        aria-label={`View debate: ${debate.topic}`}
      >
        <div className={styles.topicLabel}>Topic:</div>
        <h3 className={styles.topic}>{truncate(debate.topic, 120)}</h3>
        
        <div className={styles.panelists}>
          {debate.panelists.map((panelist) => (
            <div
              key={panelist.id}
              className={styles.avatar}
              onClick={(e) => handleAvatarClick(e, panelist)}
              onKeyDown={(e) => handleAvatarKeyDown(e, panelist)}
              role="button"
              tabIndex={0}
              aria-label={`View ${panelist.name} details`}
              title={panelist.name}
            >
              <img
                src={
                  panelist.avatarUrl?.startsWith('http') || panelist.avatarUrl?.startsWith('/')
                    ? panelist.avatarUrl
                    : `${process.env.PUBLIC_URL}/avatars/${panelist.avatarUrl || 'placeholder-avatar.svg'}`
                }
                alt={panelist.name}
                className={styles.avatarImage}
                onError={(e) => {
                  e.target.src = `${process.env.PUBLIC_URL}/avatars/placeholder-avatar.svg`;
                }}
              />
            </div>
          ))}
        </div>

        <div className={styles.timestamp}>
          {formatRelativeTime(debate.startedAt)}
        </div>
      </div>

      {selectedPanelist && (
        <PanelistModal panelist={selectedPanelist} onClose={closeModal} />
      )}
    </>
  );
};

DebateCard.propTypes = {
  debate: PropTypes.shape({
    id: PropTypes.string.isRequired,
    topic: PropTypes.string.isRequired,
    panelists: PropTypes.arrayOf(
      PropTypes.shape({
        id: PropTypes.string.isRequired,
        name: PropTypes.string.isRequired,
      })
    ).isRequired,
    startedAt: PropTypes.string.isRequired,
  }).isRequired,
  onClick: PropTypes.func.isRequired,
};

export default DebateCard;
