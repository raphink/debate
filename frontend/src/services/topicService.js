const VALIDATE_TOPIC_URL = process.env.REACT_APP_VALIDATE_TOPIC_URL || 'http://localhost:8080';

/**
 * Validates a debate topic with the backend API (simplified for unified input)
 * @param {string} topic - The topic to validate
 * @param {Array} panelists - Optional array of pre-selected panelists
 * @returns {Promise<Object>} Validation result
 */
export const validateTopic = async (topic, panelists = []) => {
  const payload = { topic };
  
  if (panelists && panelists.length > 0) {
    payload.panelists = panelists.map(p => p.id || p.slug);
  }

  try {
    const response = await fetch(VALIDATE_TOPIC_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    console.error('Failed to validate topic:', error);
    throw error;
  }
};

/**
 * Validates a debate topic with the backend API using streaming
 * @param {string} topic - The topic to validate
 * @param {string[]} suggestedNames - Optional array of suggested panelist names (max 5)
 * @param {Function} onValidation - Callback for validation result {isRelevant, message}
 * @param {Function} onPanelist - Callback for each panelist received
 * @param {Function} onError - Callback for errors
 * @param {Function} onComplete - Callback when stream completes
 * @returns {Promise<void>}
 */
export const validateTopicStream = async (
  topic,
  suggestedNames = [],
  onValidation,
  onPanelist,
  onError,
  onComplete
) => {
  const payload = { topic };
  if (suggestedNames && suggestedNames.length > 0) {
    payload.suggestedNames = suggestedNames.slice(0, 5);
  }

  try {
    const response = await fetch(VALIDATE_TOPIC_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';

    // eslint-disable-next-line no-constant-condition
    while (true) {
      const { done, value } = await reader.read();

      if (done) {
        if (onComplete) onComplete();
        break;
      }

      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split('\n');
      buffer = lines.pop(); // Keep incomplete line in buffer

      for (const line of lines) {
        if (line.trim()) {
          try {
            const chunk = JSON.parse(line);
            
            if (chunk.type === 'validation') {
              const data = JSON.parse(chunk.data);
              if (onValidation) onValidation(data);
            } else if (chunk.type === 'panelist') {
              const panelist = JSON.parse(chunk.data);
              if (onPanelist) onPanelist(panelist);
            } else if (chunk.type === 'error') {
              if (onError) onError(new Error(chunk.error));
            } else if (chunk.type === 'done') {
              if (onComplete) onComplete();
            }
          } catch (err) {
            console.error('Failed to parse chunk:', err, line);
          }
        }
      }
    }
  } catch (err) {
    if (onError) onError(err);
    throw err;
  }
};

export default {
  validateTopic,
};
