import { useState, useCallback } from 'react';
import { MAX_PANELISTS, MIN_PANELISTS } from '../utils/constants';

/**
 * Custom hook for managing panelist selection state.
 * Enforces 2-5 panelist selection limit and provides selection management functions.
 * 
 * @param {Array} initialPanelists - Initial array of selected panelists (optional)
 * @returns {Object} Selection state and management functions
 */
const usePanelistSelection = (initialPanelists = []) => {
  const [selectedPanelists, setSelectedPanelists] = useState(initialPanelists);

  /**
   * Toggle panelist selection.
   * Adds if not selected and under max limit, removes if already selected.
   * 
   * @param {Object} panelist - Panelist object to toggle
   */
  const toggleSelection = useCallback((panelist) => {
    setSelectedPanelists((prev) => {
      const isSelected = prev.some(p => p.id === panelist.id);
      
      if (isSelected) {
        // Remove panelist
        return prev.filter(p => p.id !== panelist.id);
      } else {
        // Add panelist if under max limit
        if (prev.length < MAX_PANELISTS) {
          return [...prev, panelist];
        }
        // Max limit reached, don't add
        return prev;
      }
    });
  }, []);

  /**
   * Clear all selected panelists
   */
  const clearSelection = useCallback(() => {
    setSelectedPanelists([]);
  }, []);

  /**
   * Check if a specific panelist is selected
   * 
   * @param {string} panelistId - ID of the panelist to check
   * @returns {boolean} True if panelist is selected
   */
  const isSelected = useCallback((panelistId) => {
    return selectedPanelists.some(p => p.id === panelistId);
  }, [selectedPanelists]);

  /**
   * Check if selection is valid for debate generation (2-5 panelists)
   * 
   * @returns {boolean} True if selection meets requirements
   */
  const isValidSelection = useCallback(() => {
    const count = selectedPanelists.length;
    return count >= MIN_PANELISTS && count <= MAX_PANELISTS;
  }, [selectedPanelists]);

  /**
   * Check if maximum selection limit is reached
   * 
   * @returns {boolean} True if max limit reached
   */
  const isMaxReached = useCallback(() => {
    return selectedPanelists.length >= MAX_PANELISTS;
  }, [selectedPanelists]);

  /**
   * Set selected panelists directly (useful for loading saved state)
   * 
   * @param {Array} panelists - Array of panelist objects to set
   */
  const setSelection = useCallback((panelists) => {
    if (Array.isArray(panelists) && panelists.length <= MAX_PANELISTS) {
      setSelectedPanelists(panelists);
    }
  }, []);

  return {
    selectedPanelists,
    toggleSelection,
    clearSelection,
    isSelected,
    isValidSelection,
    isMaxReached,
    setSelection,
    count: selectedPanelists.length,
  };
};

export default usePanelistSelection;
