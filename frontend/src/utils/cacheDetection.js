/**
 * Determine if current topic + panelist combination matches a historical debate
 * @param {Object} originalDebate - The debate from autocomplete/navigation state
 * @param {string} currentTopic - Current topic text
 * @param {Array} currentPanelists - Current panelist array
 * @returns {boolean} True if cache hit (exact match)
 */
export const isCacheHit = (originalDebate, currentTopic, currentPanelists) => {
  if (!originalDebate) return false;

  // Topic must match exactly (case-sensitive)
  if (originalDebate.topic !== currentTopic) {
    return false;
  }

  // Panelist count must match
  if (!originalDebate.panelists || !currentPanelists) {
    return false;
  }

  if (originalDebate.panelists.length !== currentPanelists.length) {
    return false;
  }

  // Deep equality check on panelist IDs (order-independent)
  const originalIds = originalDebate.panelists
    .map((p) => p.id)
    .sort()
    .join(',');
  const currentIds = currentPanelists
    .map((p) => p.id)
    .sort()
    .join(',');

  return originalIds === currentIds;
};
