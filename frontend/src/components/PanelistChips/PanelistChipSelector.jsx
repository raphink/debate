import React, { useState } from 'react';
import PropTypes from 'prop-types';
import PanelistChip from './PanelistChip';
import styles from './PanelistChipSelector.module.css';

const FAMOUS_PANELISTS = [
  { 
    id: 'albert-einstein', 
    name: 'Albert Einstein', 
    slug: 'albert-einstein',
    avatarUrl: 'einstein.jpg',
    tagline: 'Theoretical Physicist',
    biography: 'Developer of the theory of relativity and one of the most influential physicists of the 20th century.'
  },
  { 
    id: 'marie-curie', 
    name: 'Marie Curie', 
    slug: 'marie-curie',
    avatarUrl: 'curie.jpg',
    tagline: 'Pioneering Physicist and Chemist',
    biography: 'First woman to win a Nobel Prize and the only person to win Nobel Prizes in two scientific fields.'
  },
  { 
    id: 'stephen-hawking', 
    name: 'Stephen Hawking', 
    slug: 'stephen-hawking',
    avatarUrl: 'hawking.jpg',
    tagline: 'Theoretical Physicist and Cosmologist',
    biography: 'Renowned for his work on black holes and cosmology, author of "A Brief History of Time".'
  },
  { 
    id: 'ada-lovelace', 
    name: 'Ada Lovelace', 
    slug: 'ada-lovelace',
    avatarUrl: 'lovelace.jpg',
    tagline: 'Mathematics Pioneer',
    biography: 'First computer programmer, known for her work on Charles Babbage\'s Analytical Engine.'
  },
  { 
    id: 'nikola-tesla', 
    name: 'Nikola Tesla', 
    slug: 'nikola-tesla',
    avatarUrl: 'tesla.jpg',
    tagline: 'Inventor and Electrical Engineer',
    biography: 'Pioneer of alternating current electricity and wireless communication technologies.'
  },
  { 
    id: 'isaac-newton', 
    name: 'Isaac Newton', 
    slug: 'isaac-newton',
    avatarUrl: 'newton.jpg',
    tagline: 'Mathematician and Physicist',
    biography: 'Formulated the laws of motion and universal gravitation, invented calculus.'
  },
  { 
    id: 'charles-darwin', 
    name: 'Charles Darwin', 
    slug: 'charles-darwin',
    avatarUrl: 'darwin.jpg',
    tagline: 'Naturalist and Biologist',
    biography: 'Proposed the theory of evolution by natural selection in "On the Origin of Species".'
  },
  { 
    id: 'galileo-galilei', 
    name: 'Galileo Galilei', 
    slug: 'galileo-galilei',
    avatarUrl: 'galileo.jpg',
    tagline: 'Astronomer and Physicist',
    biography: 'Father of observational astronomy and modern physics, championed heliocentrism.'
  },
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
