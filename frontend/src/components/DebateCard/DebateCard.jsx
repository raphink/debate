import React from 'react';
import PropTypes from 'prop-types';
import { formatPanelists, formatRelativeTime, truncate } from '../../utils/formatters';
import styles from './DebateCard.module.css';

const DebateCard = ({ debate, onClick }) => {
  const handleClick = () => {
    onClick(debate.id);
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      onClick(debate.id);
    }
  };

  return (
    <div
      className={styles.card}
      onClick={handleClick}
      onKeyDown={handleKeyDown}
      role="button"
      tabIndex={0}
      aria-label={`View debate: ${debate.topic}`}
    >
      <h3 className={styles.topic}>{truncate(debate.topic, 100)}</h3>
      <div className={styles.panelists}>
        {formatPanelists(debate.panelists)}
      </div>
      <div className={styles.timestamp}>
        {formatRelativeTime(debate.startedAt)}
      </div>
    </div>
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
