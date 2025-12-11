import apiClient from './api';

/**
 * Generate a debate with streaming responses via Server-Sent Events
 * 
 * @param {string} topic - The debate topic
 * @param {Array} selectedPanelists - Array of selected panelist objects
 * @param {Function} onMessage - Callback for each message chunk: (panelistId, text) => void
 * @param {Function} onError - Callback for errors: (error) => void
 * @param {Function} onComplete - Callback when streaming completes
 * @returns {Function} Cleanup function to abort the stream
 */
export const generateDebateStream = (topic, selectedPanelists, onMessage, onError, onComplete) => {
  const baseURL = process.env.REACT_APP_GENERATE_DEBATE_URL || 'http://localhost:8081';
  const url = `${baseURL}/GenerateDebate`;

  const requestBody = {
    topic,
    selectedPanelists: selectedPanelists.map(p => ({
      id: p.id,
      name: p.name,
      tagline: p.tagline,
      bio: p.bio,
      avatarUrl: p.avatarUrl,
      position: p.position,
    })),
  };

  // Use fetch with ReadableStream for SSE
  const abortController = new AbortController();

  fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(requestBody),
    signal: abortController.signal,
  })
    .then(async (response) => {
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to generate debate');
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

        // Decode the chunk and add to buffer
        buffer += decoder.decode(value, { stream: true });

        // Process complete JSON objects from buffer
        const lines = buffer.split('\n');
        buffer = lines.pop() || ''; // Keep incomplete line in buffer

        for (const line of lines) {
          if (!line.trim()) continue;

          try {
            const chunk = JSON.parse(line);

            if (chunk.type === 'message' && chunk.panelistId && chunk.text) {
              onMessage(chunk.panelistId, chunk.text);
            } else if (chunk.type === 'error') {
              onError(new Error(chunk.error || 'Unknown error occurred'));
              return;
            } else if (chunk.type === 'done') {
              if (onComplete) onComplete();
              return;
            }
          } catch (err) {
            console.warn('Failed to parse chunk:', line, err);
          }
        }
      }
    })
    .catch((error) => {
      if (error.name !== 'AbortError') {
        onError(error);
      }
    });

  // Return cleanup function
  return () => {
    abortController.abort();
  };
};

export default {
  generateDebateStream,
};
