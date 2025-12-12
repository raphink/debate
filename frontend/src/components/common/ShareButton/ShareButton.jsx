import React, { useState } from 'react';
import Button from '../Button/Button';
import styles from './ShareButton.module.css';

/**
 * ShareButton component - copies debate URL to clipboard
 */
const ShareButton = ({ debateId }) => {
  const [copied, setCopied] = useState(false);

  const handleShare = async () => {
    try {
      const url = `${window.location.origin}/d/${debateId}`;
      await navigator.clipboard.writeText(url);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy URL:', err);
      // Fallback for older browsers
      const textArea = document.createElement('textarea');
      textArea.value = `${window.location.origin}/d/${debateId}`;
      document.body.appendChild(textArea);
      textArea.select();
      try {
        document.execCommand('copy');
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
      } catch (fallbackErr) {
        console.error('Fallback copy failed:', fallbackErr);
      }
      document.body.removeChild(textArea);
    }
  };

  if (!debateId) {
    return null;
  }

  return (
    <div className={styles.container}>
      <Button 
        onClick={handleShare}
        variant="secondary"
        className={styles.shareButton}
      >
        {copied ? '✓ Link Copied!' : '🔗 Share Debate'}
      </Button>
      {copied && <span className={styles.toast}>Link copied to clipboard!</span>}
    </div>
  );
};

export default ShareButton;
