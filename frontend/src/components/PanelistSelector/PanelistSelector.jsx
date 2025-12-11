import React from 'react';
import PropTypes from 'prop-types';
import Button from '../common/Button/Button';
import styles from './PanelistSelector.module.css';
import { MAX_PANELISTS, MIN_PANELISTS } from '../../utils/constants';

/**
 * PanelistSelector component displays selection status and controls.
 * Shows selected count, min/max limits, and provides navigation and clear actions.
 * 
 * @param {Object} props - Component props
 * @param {Array} props.selectedPanelists - Array of currently selected panelist objects
 * @param {Function} props.onClear - Callback to clear all selections
 * @param {Function} props.onProceed - Callback to proceed to debate generation
 * @param {boolean} props.isLoading - Whether debate generation is in progress
 */
const PanelistSelector = ({ selectedPanelists, onClear, onProceed, isLoading = false }) => {
  const count = selectedPanelists.length;
  const canProceed = count >= MIN_PANELISTS && count <= MAX_PANELISTS;
  const isMaxReached = count >= MAX_PANELISTS;

  return (
    <div className={styles.container}>
      <div className={styles.status}>
        <h3 className={styles.title}>Selected Panelists</h3>
        <div className={styles.counter}>
          <span className={`${styles.count} ${canProceed ? styles.valid : styles.invalid}`}>
            {count}
          </span>
          <span className={styles.limit}>/ {MAX_PANELISTS}</span>
        </div>
      </div>

      <div className={styles.selectedList}>
        {selectedPanelists.length === 0 ? (
          <p className={styles.emptyMessage}>
            Select {MIN_PANELISTS}-{MAX_PANELISTS} panelists to generate a debate
          </p>
        ) : (
          <>
            {selectedPanelists.map((panelist) => (
              <div key={panelist.id} className={styles.selectedItem}>
                <img
                  src={panelist.avatarUrl}
                  alt={`${panelist.name} avatar`}
                  className={styles.selectedAvatar}
                  onError={(e) => {
                    e.target.src = '/avatars/placeholder-avatar.svg';
                  }}
                />
                <span className={styles.selectedName}>{panelist.name}</span>
              </div>
            ))}
          </>
        )}
      </div>

      {count > 0 && count < MIN_PANELISTS && (
        <p className={styles.warning}>
          Select at least {MIN_PANELISTS - count} more panelist{MIN_PANELISTS - count !== 1 ? 's' : ''}
        </p>
      )}

      {isMaxReached && (
        <p className={styles.info}>
          Maximum {MAX_PANELISTS} panelists selected
        </p>
      )}

      <div className={styles.actions}>
        <Button
          onClick={onClear}
          variant="secondary"
          disabled={count === 0 || isLoading}
          aria-label="Clear all selected panelists"
        >
          Clear Selection
        </Button>
        <Button
          onClick={onProceed}
          variant="primary"
          disabled={!canProceed || isLoading}
          aria-label={`Generate debate with ${count} selected panelists`}
        >
          {isLoading ? 'Generating...' : 'Generate Debate'}
        </Button>
      </div>
    </div>
  );
};

PanelistSelector.propTypes = {
  selectedPanelists: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      name: PropTypes.string.isRequired,
      avatarUrl: PropTypes.string.isRequired,
    })
  ).isRequired,
  onClear: PropTypes.func.isRequired,
  onProceed: PropTypes.func.isRequired,
  isLoading: PropTypes.bool,
};

export default PanelistSelector;
