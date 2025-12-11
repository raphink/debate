import React, { useEffect, useRef, useState } from 'react';
import PropTypes from 'prop-types';
import DebateBubble from './DebateBubble';
import TypingIndicator from './TypingIndicator';
import PanelistModal from './PanelistModal';
import styles from './DebateView.module.css';

/**
 * DebateView component displays the debate conversation as a scrollable chat interface.
 * Auto-scroll can be toggled on/off.
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
  const [autoScroll, setAutoScroll] = useState(false);
  const [selectedPanelist, setSelectedPanelist] = useState(null);

  // Auto-scroll to latest message (only if enabled)
  useEffect(() => {
    if (autoScroll && messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth', block: 'end' });
    }
  }, [messages, isStreaming, autoScroll]);

  // Create a map of panelists by ID for quick lookup
  const panelistMap = panelists.reduce((acc, panelist) => {
    acc[panelist.id] = panelist;
    return acc;
  }, {
    // Add moderator with default avatar
    moderator: {
      id: 'moderator',
      name: 'Moderator',
      avatarUrl: `${process.env.PUBLIC_URL}/avatars/moderator-avatar.svg`,
      tagline: 'Neutral Facilitator',
      bio: 'Guiding the conversation'
    }
  });

  const currentPanelist = currentPanelistId ? panelistMap[currentPanelistId] : null;

  const handleAvatarClick = (panelist) => {
    setSelectedPanelist(panelist);
  };

  const handleCloseModal = () => {
    setSelectedPanelist(null);
  };

  return (
    <div className={styles.container} ref={containerRef}>
      <div className={styles.controls}>
        <label className={styles.autoScrollToggle}>
          <input
            type="checkbox"
            checked={autoScroll}
            onChange={(e) => setAutoScroll(e.target.checked)}
          />
          <span>Auto-scroll</span>
        </label>
      </div>
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
              onAvatarClick={handleAvatarClick}
            />
          );
        })}

        {isStreaming && <TypingIndicator />}

        <div ref={messagesEndRef} />
      </div>

      {selectedPanelist && (
        <PanelistModal 
          panelist={selectedPanelist} 
          onClose={handleCloseModal}
        />
      )}
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
