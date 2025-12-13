import { formatDistanceToNow } from 'date-fns';

/**
 * Format a list of panelists for display
 * Shows first 3 names, then "+N more" if there are additional panelists
 * @param {Array} panelists - Array of panelist objects with name property
 * @returns {string} Formatted string like "Name1, Name2, Name3 +2 more"
 */
export function formatPanelists(panelists) {
  if (!panelists || panelists.length === 0) {
    return 'No panelists';
  }

  const displayCount = Math.min(3, panelists.length);
  const displayNames = panelists
    .slice(0, displayCount)
    .map(p => p.name)
    .join(', ');

  const remaining = panelists.length - displayCount;
  if (remaining > 0) {
    return `${displayNames} +${remaining} more`;
  }

  return displayNames;
}

/**
 * Format a timestamp as relative time (e.g., "2 hours ago", "3 days ago")
 * @param {string|Date} timestamp - ISO timestamp or Date object
 * @returns {string} Relative time string
 */
export function formatRelativeTime(timestamp) {
  if (!timestamp) {
    return 'Unknown';
  }

  try {
    const date = typeof timestamp === 'string' ? new Date(timestamp) : timestamp;
    return formatDistanceToNow(date, { addSuffix: true });
  } catch (error) {
    console.error('Failed to format timestamp:', error);
    return 'Unknown';
  }
}

/**
 * Truncate text with ellipsis if it exceeds maxLength
 * @param {string} text - Text to truncate
 * @param {number} maxLength - Maximum length before truncation
 * @returns {string} Truncated text with "..." appended if needed
 */
export function truncate(text, maxLength) {
  if (!text || text.length <= maxLength) {
    return text || '';
  }

  return text.substring(0, maxLength) + '...';
}
