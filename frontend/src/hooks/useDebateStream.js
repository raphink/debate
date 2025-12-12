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
  const [debateId, setDebateId] = useState(null);
  
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
    setDebateId(null);

    // Handle debate ID
    const handleDebateId = (id) => {
      setDebateId(id);
      // Update URL without page reload
      window.history.pushState({}, '', `/d/${id}`);
    };

    // Handle incoming message chunks
    const handleMessage = (panelistId, text) => {
      setCurrentPanelistId(panelistId);
      setMessages((prev) => {
        const lastMessage = prev[prev.length - 1];
        
        // If same speaker, append to existing bubble
        if (lastMessage && lastMessage.panelistId === panelistId) {
          return [
            ...prev.slice(0, -1),
            {
              panelistId,
              text: lastMessage.text + text,
            },
          ];
        }
        
        // Different speaker - create new bubble
        return [
          ...prev,
          {
            panelistId,
            text,
          },
        ];
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
      handleComplete,
      handleDebateId
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
    debateId,
    startDebate,
    stopDebate,
    reset,
  };
};

export default useDebateStream;
