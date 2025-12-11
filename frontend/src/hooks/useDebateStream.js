import { useState, useCallback, useRef, useEffect } from 'react';
import { generateDebateStream } from '../services/debateService';

/**
 * Custom hook for managing debate streaming state.
 * Handles SSE connection, message accumulation, and error states.
 * 
 * @returns {Object} Debate streaming state and control functions
 */
const useDebateStream = () => {
  const [messages, setMessages] = useState([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [error, setError] = useState(null);
  const [isComplete, setIsComplete] = useState(false);
  const [currentPanelistId, setCurrentPanelistId] = useState(null);
  
  const cleanupRef = useRef(null);

  /**
   * Start debate generation stream
   * 
   * @param {string} topic - The debate topic
   * @param {Array} selectedPanelists - Array of selected panelist objects
   */
  const startDebate = useCallback((topic, selectedPanelists) => {
    // Reset state
    setMessages([]);
    setError(null);
    setIsComplete(false);
    setIsStreaming(true);
    setCurrentPanelistId(null);

    // Handle incoming message chunks
    const handleMessage = (panelistId, text) => {
      console.log('[DEBUG] Frontend received chunk:', { panelistId, text: text.substring(0, 100) });
      setCurrentPanelistId(panelistId);
      setMessages((prev) => {
        // Check if this is a NEW message from same panelist (backend sends complete messages)
        // Backend now sends complete messages, not incremental chunks
        // So each call to handleMessage should create a new message entry
        const lastMessage = prev[prev.length - 1];
        
        // If it's the exact same panelist AND the text looks like a continuation (no [ID]:),
        // it might be a streaming chunk (though backend should send complete messages now)
        if (lastMessage && lastMessage.panelistId === panelistId && !text.includes('[') && !text.includes(']:')) {
          // Treat as potential continuation for robustness
          console.log('[DEBUG] Treating as continuation for', panelistId);
          return [
            ...prev.slice(0, -1),
            {
              panelistId,
              text: lastMessage.text + text, // No space, direct append
            },
          ];
        } else {
          // New message from this panelist
          console.log('[DEBUG] New message from', panelistId);
          return [
            ...prev,
            {
              panelistId,
              text,
            },
          ];
        }
      });
    };

    // Handle errors
    const handleError = (err) => {
      setError(err);
      setIsStreaming(false);
      setCurrentPanelistId(null);
    };

    // Handle completion
    const handleComplete = () => {
      setIsStreaming(false);
      setIsComplete(true);
      setCurrentPanelistId(null);
    };

    // Start the stream
    cleanupRef.current = generateDebateStream(
      topic,
      selectedPanelists,
      handleMessage,
      handleError,
      handleComplete
    );
  }, []);

  /**
   * Stop the debate stream
   */
  const stopDebate = useCallback(() => {
    if (cleanupRef.current) {
      cleanupRef.current();
      cleanupRef.current = null;
    }
    setIsStreaming(false);
    setCurrentPanelistId(null);
  }, []);

  /**
   * Reset debate state
   */
  const reset = useCallback(() => {
    stopDebate();
    setMessages([]);
    setError(null);
    setIsComplete(false);
  }, [stopDebate]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (cleanupRef.current) {
        cleanupRef.current();
      }
    };
  }, []);

  return {
    messages,
    isStreaming,
    error,
    isComplete,
    currentPanelistId,
    startDebate,
    stopDebate,
    reset,
  };
};

export default useDebateStream;
