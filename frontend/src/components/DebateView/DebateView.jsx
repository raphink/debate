import React, { useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import DebateBubble from './DebateBubble';
import TypingIndicator from './TypingIndicator';
import styles from './DebateView.module.css';

/**
 * DebateView component displays the debate conversation as a scrollable chat interface.
 * Auto-scrolls to show the latest messages.
 * 
 * @param {Object} props - Component props
 * @param {Array} props.messages - Array of message objects with { panelistId, text }
 * @param {Array} props.panelists - Array of panelist objects
 * @param {boolean} props.isStreaming - Whether debate is currently streaming
 * @param {string} props.currentPanelistId - ID of panelist currently responding
 */
const DebateView = ({ messages, panelists, isStreaming, currentPanelistId }) => {
  const messagesEndRef = useRef(null);
  const containerRef = useRef(null);

  // Auto-scroll to latest message
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth', block: 'end' });
    }
  }, [messages, isStreaming]);

  // Create a map of panelists by ID for quick lookup
  const panelistMap = panelists.reduce((acc, panelist) => {
    acc[panelist.id] = panelist;
    return acc;
  }, {});

  const currentPanelist = currentPanelistId ? panelistMap[currentPanelistId] : null;

  return (
    <div className={styles.container} ref={containerRef}>
      <div className={styles.messageList}>
        {messages.length === 0 && !isStreaming && (
          <div className={styles.emptyState}>
            <p>The debate will appear here...</p>
          </div>
        )}

        {messages.map((message, index) => {
          const panelist = panelistMap[message.panelistId];
          if (!panelist) {
            console.warn(`Panelist not found for ID: ${message.panelistId}`);
            return null;
          }

          return (
            <DebateBubble
              key={index}
              panelist={panelist}
              text={message.text}
              index={index}
            />
          );
        })}

        {isStreaming && <TypingIndicator panelistName={currentPanelist?.name} />}

        <div ref={messagesEndRef} />
      </div>
    </div>
  );
};

DebateView.propTypes = {
  messages: PropTypes.arrayOf(
    PropTypes.shape({
      panelistId: PropTypes.string.isRequired,
      text: PropTypes.string.isRequired,
    })
  ).isRequired,
  panelists: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      name: PropTypes.string.isRequired,
      avatarUrl: PropTypes.string.isRequired,
    })
  ).isRequired,
  isStreaming: PropTypes.bool,
  currentPanelistId: PropTypes.string,
};

DebateView.defaultProps = {
  isStreaming: false,
  currentPanelistId: null,
};

export default DebateView;
