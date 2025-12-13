import React from 'react';
import PropTypes from 'prop-types';
import styles from './SearchInput.module.css';

const SearchInput = ({ value, onChange, placeholder = 'Search...' }) => {
  const handleClear = () => {
    onChange('');
  };

  return (
    <div className={styles.searchContainer}>
      <span className={styles.searchIcon} aria-hidden="true">🔍</span>
      <input
        type="text"
        className={styles.searchInput}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        aria-label="Search debates"
      />
      {value && (
        <button
          type="button"
          className={styles.clearButton}
          onClick={handleClear}
          aria-label="Clear search"
        >
          ×
        </button>
      )}
    </div>
  );
};

SearchInput.propTypes = {
  value: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
  placeholder: PropTypes.string,
};

export default SearchInput;
