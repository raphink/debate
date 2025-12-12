import React from 'react';
import PropTypes from 'prop-types';
import PanelistCard from './PanelistCard';
import styles from './PanelistGrid.module.css';

/**
 * PanelistGrid component displays a grid of panelist cards with selection functionality.
 * Supports keyboard navigation and enforces 2-5 panelist selection limit.
 * 
 * @param {Object} props - Component props
 * @param {Array} props.panelists - Array of panelist objects
 * @param {Array} props.selectedPanelists - Array of currently selected panelist objects
 * @param {Function} props.onToggleSelection - Callback when a panelist is selected/deselected
 * @param {number} props.maxSelection - Maximum number of panelists that can be selected (default: 5)
 * @param {boolean} props.isLocked - Whether the selection is locked (from cache)
 */
const PanelistGrid = ({ panelists, selectedPanelists, onToggleSelection, maxSelection = 5, isLocked = false }) => {
  const selectedIds = selectedPanelists.map(p => p.id);
  const isMaxReached = selectedPanelists.length >= maxSelection;

  if (!panelists || panelists.length === 0) {
    return (
      <div className={styles.emptyState}>
        <p>No panelists available. Please validate a topic first.</p>
      </div>
    );
  }

  return (
    <div className={styles.gridContainer}>
      {isLocked && (
        <div className={styles.lockedNotice}>
          Panelists pre-selected from historical debate. Click "Modify Panelists" to make changes.
        </div>
      )}
      <div className={styles.grid}>
        {panelists.map((panelist) => {
          const isSelected = selectedIds.includes(panelist.id);
          const isDisabled = (isMaxReached && !isSelected) || isLocked;

          return (
            <PanelistCard
              key={panelist.id}
              panelist={panelist}
              isSelected={isSelected}
              onToggle={onToggleSelection}
              disabled={isDisabled}
            />
          );
        })}
      </div>
    </div>
  );
};

PanelistGrid.propTypes = {
  panelists: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      name: PropTypes.string.isRequired,
      tagline: PropTypes.string.isRequired,
      bio: PropTypes.string.isRequired,
      avatarUrl: PropTypes.string.isRequired,
      position: PropTypes.string,
    })
  ).isRequired,
  selectedPanelists: PropTypes.array.isRequired,
  onToggleSelection: PropTypes.func.isRequired,
  maxSelection: PropTypes.number,
  isLocked: PropTypes.bool,
};

export default PanelistGrid;
