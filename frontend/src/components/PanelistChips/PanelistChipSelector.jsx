import React, { useState } from 'react';
import PropTypes from 'prop-types';
import PanelistChip from './PanelistChip';
import styles from './PanelistChipSelector.module.css';

const FAMOUS_PANELISTS = [
  { id: 'albert-einstein', name: 'Albert Einstein', avatarUrl: 'einstein.jpg' },
  { id: 'marie-curie', name: 'Marie Curie', avatarUrl: 'curie.jpg' },
  { id: 'stephen-hawking', name: 'Stephen Hawking', avatarUrl: 'hawking.jpg' },
  { id: 'ada-lovelace', name: 'Ada Lovelace', avatarUrl: 'lovelace.jpg' },
  { id: 'nikola-tesla', name: 'Nikola Tesla', avatarUrl: 'tesla.jpg' },
  { id: 'isaac-newton', name: 'Isaac Newton', avatarUrl: 'newton.jpg' },
  { id: 'charles-darwin', name: 'Charles Darwin', avatarUrl: 'darwin.jpg' },
  { id: 'galileo-galilei', name: 'Galileo Galilei', avatarUrl: 'galileo.jpg' },
];

const PanelistChipSelector = ({ value, onChange, disabled, max }) => {
  const [showPicker, setShowPicker] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');

  const handleRemove = (id) => {
    onChange(value.filter((p) => p.id !== id));
  };

  const handleAdd = (panelist) => {
    if (value.length >= max) {
      return;
    }
    if (!value.find((p) => p.id === panelist.id)) {
      onChange([...value, panelist]);
    }
    setShowPicker(false);
    setSearchQuery('');
  };

  const availablePanelists = FAMOUS_PANELISTS.filter(
    (p) =>
      !value.find((v) => v.id === p.id) &&
      p.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className={styles.container}>
      <div className={styles.chipList}>
        {value.map((panelist) => (
          <PanelistChip
            key={panelist.id}
            panelist={panelist}
            onRemove={handleRemove}
            disabled={disabled}
          />
        ))}
        {value.length < max && !disabled && (
          <button
            type="button"
            className={styles.addButton}
            onClick={() => setShowPicker(!showPicker)}
            aria-label="Add panelist"
          >
            +
          </button>
        )}
      </div>

      {showPicker && (
        <div className={styles.picker}>
          <div className={styles.pickerHeader}>
            <input
              type="text"
              placeholder="Search panelists..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className={styles.searchInput}
              autoFocus
            />
            <button
              type="button"
              onClick={() => {
                setShowPicker(false);
                setSearchQuery('');
              }}
              className={styles.closeButton}
              aria-label="Close"
            >
              ×
            </button>
          </div>
          <div className={styles.pickerGrid}>
            {availablePanelists.length > 0 ? (
              availablePanelists.map((panelist) => (
                <button
                  key={panelist.id}
                  type="button"
                  className={styles.pickerItem}
                  onClick={() => handleAdd(panelist)}
                >
                  <img
                    src={`${process.env.PUBLIC_URL}/avatars/${panelist.avatarUrl}`}
                    alt={panelist.name}
                    className={styles.pickerAvatar}
                    onError={(e) => {
                      e.target.src = `${process.env.PUBLIC_URL}/avatars/placeholder-avatar.svg`;
                    }}
                  />
                  <span className={styles.pickerName}>{panelist.name}</span>
                </button>
              ))
            ) : (
              <div className={styles.emptyState}>No panelists found</div>
            )}
          </div>
        </div>
      )}

      {value.length === 0 && !showPicker && (
        <p className={styles.hint}>
          Optional: Add panelists to get better topic suggestions
        </p>
      )}
    </div>
  );
};

PanelistChipSelector.propTypes = {
  value: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      name: PropTypes.string.isRequired,
      avatarUrl: PropTypes.string,
    })
  ),
  onChange: PropTypes.func.isRequired,
  disabled: PropTypes.bool,
  max: PropTypes.number,
};

PanelistChipSelector.defaultProps = {
  value: [],
  disabled: false,
  max: 10,
};

export default PanelistChipSelector;
