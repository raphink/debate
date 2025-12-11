import { jsPDF } from 'jspdf';

/**
 * Generates a PDF document from debate data
 * @param {Object} debateData - The debate data to export
 * @param {string} debateData.topic - The debate topic
 * @param {Array} debateData.panelists - Array of panelist objects
 * @param {Array} debateData.messages - Array of message objects {panelistId, text}
 * @returns {jsPDF} The generated PDF document
 */
export const generateDebatePDF = (debateData) => {
  const { topic, panelists, messages } = debateData;
  
  // Create new PDF document (A4 size)
  const pdf = new jsPDF({
    orientation: 'portrait',
    unit: 'mm',
    format: 'a4'
  });

  const pageWidth = pdf.internal.pageSize.getWidth();
  const pageHeight = pdf.internal.pageSize.getHeight();
  const margin = 20;
  const contentWidth = pageWidth - (2 * margin);
  let yPosition = margin;

  // Create panelist lookup map
  const panelistMap = panelists.reduce((acc, panelist) => {
    acc[panelist.id] = panelist;
    return acc;
  }, {
    moderator: {
      id: 'moderator',
      name: 'Moderator',
      tagline: 'Neutral Facilitator',
      bio: 'Guiding the conversation'
    }
  });

  /**
   * Check if we need a new page
   */
  const checkPageBreak = (neededSpace) => {
    if (yPosition + neededSpace > pageHeight - margin) {
      pdf.addPage();
      yPosition = margin;
      addPageNumber();
      return true;
    }
    return false;
  };

  /**
   * Add page number to footer
   */
  const addPageNumber = () => {
    const pageNumber = pdf.internal.getCurrentPageInfo().pageNumber;
    pdf.setFontSize(10);
    pdf.setTextColor(128, 128, 128);
    pdf.text(
      `Page ${pageNumber}`,
      pageWidth / 2,
      pageHeight - 10,
      { align: 'center' }
    );
  };

  // ========== HEADER ==========
  
  // Title
  pdf.setFontSize(20);
  pdf.setTextColor(31, 41, 55); // gray-800
  pdf.setFont('helvetica', 'bold');
  pdf.text('Theological Debate', margin, yPosition);
  yPosition += 10;

  // Timestamp
  pdf.setFontSize(10);
  pdf.setTextColor(107, 114, 128); // gray-500
  pdf.setFont('helvetica', 'normal');
  const timestamp = new Date().toLocaleString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
  pdf.text(`Generated: ${timestamp}`, margin, yPosition);
  yPosition += 10;

  // Divider line
  pdf.setDrawColor(229, 231, 235); // gray-200
  pdf.line(margin, yPosition, pageWidth - margin, yPosition);
  yPosition += 10;

  // ========== TOPIC ==========
  
  pdf.setFontSize(14);
  pdf.setTextColor(31, 41, 55);
  pdf.setFont('helvetica', 'bold');
  pdf.text('Topic:', margin, yPosition);
  yPosition += 7;

  pdf.setFontSize(12);
  pdf.setFont('helvetica', 'normal');
  const topicLines = pdf.splitTextToSize(topic, contentWidth);
  pdf.text(topicLines, margin, yPosition);
  yPosition += (topicLines.length * 7) + 10;

  checkPageBreak(30);

  // ========== PANELISTS ==========
  
  pdf.setFontSize(14);
  pdf.setFont('helvetica', 'bold');
  pdf.text('Panelists:', margin, yPosition);
  yPosition += 10;

  pdf.setFontSize(10);
  pdf.setFont('helvetica', 'normal');

  panelists.forEach((panelist) => {
    checkPageBreak(20);

    // Panelist name
    pdf.setFont('helvetica', 'bold');
    pdf.text(panelist.name, margin + 5, yPosition);
    yPosition += 5;

    // Tagline
    if (panelist.tagline) {
      pdf.setFont('helvetica', 'italic');
      pdf.setTextColor(107, 114, 128);
      pdf.text(panelist.tagline, margin + 5, yPosition);
      yPosition += 5;
    }

    // Bio
    if (panelist.bio) {
      pdf.setFont('helvetica', 'normal');
      pdf.setTextColor(75, 85, 99);
      const bioLines = pdf.splitTextToSize(panelist.bio, contentWidth - 10);
      pdf.text(bioLines, margin + 5, yPosition);
      yPosition += (bioLines.length * 5) + 5;
    }

    pdf.setTextColor(31, 41, 55);
  });

  yPosition += 5;
  checkPageBreak(30);

  // Divider line
  pdf.setDrawColor(229, 231, 235);
  pdf.line(margin, yPosition, pageWidth - margin, yPosition);
  yPosition += 10;

  // ========== DEBATE CONVERSATION ==========
  
  pdf.setFontSize(14);
  pdf.setFont('helvetica', 'bold');
  pdf.text('Debate:', margin, yPosition);
  yPosition += 10;

  pdf.setFontSize(10);

  messages.forEach((message, index) => {
    const panelist = panelistMap[message.panelistId];
    if (!panelist) return;

    // Estimate space needed for this message
    const messageLines = pdf.splitTextToSize(message.text, contentWidth - 15);
    const neededSpace = 10 + (messageLines.length * 5);
    
    checkPageBreak(neededSpace);

    // Speaker name
    pdf.setFont('helvetica', 'bold');
    pdf.setTextColor(91, 138, 114); // green-700
    pdf.text(`${panelist.name}:`, margin, yPosition);
    yPosition += 6;

    // Message text
    pdf.setFont('helvetica', 'normal');
    pdf.setTextColor(44, 62, 80);
    pdf.text(messageLines, margin + 5, yPosition);
    yPosition += (messageLines.length * 5) + 8;

    pdf.setTextColor(31, 41, 55);
  });

  // Add page numbers to all pages
  const totalPages = pdf.internal.getNumberOfPages();
  for (let i = 1; i <= totalPages; i++) {
    pdf.setPage(i);
    addPageNumber();
  }

  return pdf;
};

/**
 * Download the PDF to the user's device
 * @param {jsPDF} pdf - The PDF document to download
 * @param {string} filename - The filename for the download
 */
export const downloadPDF = (pdf, filename = 'debate.pdf') => {
  pdf.save(filename);
};

/**
 * Generate and download a debate PDF in one step
 * @param {Object} debateData - The debate data to export
 * @param {string} filename - Optional filename for the PDF
 */
export const exportDebatePDF = (debateData, filename) => {
  try {
    const pdf = generateDebatePDF(debateData);
    const sanitizedTopic = debateData.topic
      .substring(0, 50)
      .replace(/[^a-z0-9]/gi, '-')
      .toLowerCase();
    const defaultFilename = `debate-${sanitizedTopic}-${Date.now()}.pdf`;
    downloadPDF(pdf, filename || defaultFilename);
    return { success: true };
  } catch (error) {
    console.error('PDF generation failed:', error);
    return { success: false, error: error.message };
  }
};
