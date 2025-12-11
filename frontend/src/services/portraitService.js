const GET_PORTRAIT_URL = process.env.REACT_APP_GET_PORTRAIT_URL || 'http://localhost:8082';

/**
 * Fetches portrait URL for a panelist from Wikimedia Commons
 * @param {string} panelistId - Unique panelist identifier
 * @param {string} panelistName - Full name for Wikimedia search
 * @returns {Promise<string>} Portrait URL or placeholder
 */
export const getPortrait = async (panelistId, panelistName) => {
  try {
    const response = await fetch(GET_PORTRAIT_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        panelistId,
        panelistName,
      }),
    });

    if (!response.ok) {
      console.error(`Portrait fetch failed for ${panelistName}:`, response.status);
      return 'placeholder-avatar.svg';
    }

    const data = await response.json();
    return data.portraitUrl || 'placeholder-avatar.svg';
  } catch (error) {
    console.error(`Error fetching portrait for ${panelistName}:`, error);
    return 'placeholder-avatar.svg';
  }
};

export default {
  getPortrait,
};
