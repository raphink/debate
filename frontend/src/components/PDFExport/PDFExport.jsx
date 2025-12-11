import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { exportDebatePDF } from './pdfGenerator';
import styles from './PDFExport.module.css';

/**
 * PDFExport component provides a button to export the debate as a PDF
 * @param {Object} props - Component props
 * @param {string} props.topic - The debate topic
 * @param {Array} props.panelists - Array of panelist objects
 * @param {Array} props.messages - Array of message objects
 */
const PDFExport = ({ topic, panelists, messages }) => {
  const [isExporting, setIsExporting] = useState(false);
  const [error, setError] = useState(null);

  const handleExport = () => {
    setIsExporting(true);
    setError(null);

    try {
      const debateData = { topic, panelists, messages };
      const result = exportDebatePDF(debateData);
      
      if (!result.success) {
        setError('Failed to generate PDF. Please try again.');
      }
    } catch (err) {
      console.error('PDF export error:', err);
      setError('An unexpected error occurred during PDF generation.');
    } finally {
      setIsExporting(false);
    }
  };

  // Don't show export button if debate is empty
  if (!messages || messages.length === 0) {
    return null;
  }

  return (
    <div className={styles.container}>
      <button
        className={styles.exportButton}
        onClick={handleExport}
        disabled={isExporting}
        aria-label="Export debate as PDF"
      >
        {isExporting ? (
          <>
            <span className={styles.spinner} />
            Generating PDF...
          </>
        ) : (
          <>
            <svg 
              className={styles.icon} 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path 
                strokeLinecap="round" 
                strokeLinejoin="round" 
                strokeWidth={2} 
                d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" 
              />
            </svg>
            Export as PDF
          </>
        )}
      </button>

      {error && (
        <div className={styles.error} role="alert">
          {error}
        </div>
      )}
    </div>
  );
};

PDFExport.propTypes = {
  topic: PropTypes.string.isRequired,
  panelists: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      name: PropTypes.string.isRequired,
      tagline: PropTypes.string,
      bio: PropTypes.string,
    })
  ).isRequired,
  messages: PropTypes.arrayOf(
    PropTypes.shape({
      panelistId: PropTypes.string.isRequired,
      text: PropTypes.string.isRequired,
    })
  ).isRequired,
};

export default PDFExport;
