import { jsPDF } from 'jspdf';

/**
 * Load an image from URL and convert to base64 data URL
 * @param {string} url - Image URL (absolute or relative)
 * @returns {Promise<string>} Base64 data URL
 */
const loadImageAsDataURL = (url) => {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.crossOrigin = 'anonymous'; // Enable CORS for Wikimedia images
    
    img.onload = () => {
      const canvas = document.createElement('canvas');
      canvas.width = img.width;
      canvas.height = img.height;
      const ctx = canvas.getContext('2d');
      ctx.drawImage(img, 0, 0);
      try {
        const dataURL = canvas.toDataURL('image/jpeg', 0.8);
        resolve(dataURL);
      } catch (error) {
        console.error('Failed to convert image to data URL:', error);
        reject(error);
      }
    };
    
    img.onerror = (error) => {
      console.error('Failed to load image:', url, error);
      reject(error);
    };
    
    // Handle relative paths
    if (url.startsWith('http') || url.startsWith('//')) {
      img.src = url;
    } else {
      img.src = `${window.location.origin}${url.startsWith('/') ? '' : '/'}${url}`;
    }
  });
};

/**
 * Generates a PDF document from debate data
 * @param {Object} debateData - The debate data to export
 * @param {string} debateData.topic - The debate topic
 * @param {Array} debateData.panelists - Array of panelist objects with avatarUrl
 * @param {Array} debateData.messages - Array of message objects {panelistId, text}
 * @returns {Promise<jsPDF>} The generated PDF document
 */
export const generateDebatePDF = async (debateData) => {
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
  const avatarSize = 10; // Circular avatar diameter in mm
  let yPosition = margin;

  // Load all portrait images upfront
  const portraitCache = {};
  const defaultAvatar = `${window.location.origin}/avatars/placeholder-avatar.svg`;
  
  // Collect all unique avatar URLs
  const avatarUrls = new Set();
  panelists.forEach(p => {
    if (p.avatarUrl) avatarUrls.add(p.avatarUrl);
  });
  avatarUrls.add(defaultAvatar); // Always load placeholder
  
  // Load all images
  for (const url of avatarUrls) {
    try {
      portraitCache[url] = await loadImageAsDataURL(url);
    } catch (error) {
      console.warn(`Failed to load portrait ${url}, using placeholder`);
      try {
        portraitCache[url] = await loadImageAsDataURL(defaultAvatar);
      } catch (fallbackError) {
        console.error('Failed to load placeholder avatar');
      }
    }
  }

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

  /**
   * Draw circular avatar image
   * @param {string} imageData - Base64 image data URL
   * @param {number} x - X position (center of circle)
   * @param {number} y - Y position (center of circle)
   * @param {number} radius - Radius of circle
   */
  const drawCircularAvatar = (imageData, x, y, radius) => {
    if (!imageData) return;
    
    // Save graphics state
    pdf.saveGraphicsState();
    
    // Create circular clipping path
    pdf.circle(x, y, radius);
    pdf.clip();
    
    // Draw image (square, will be clipped to circle)
    const size = radius * 2;
    pdf.addImage(imageData, 'JPEG', x - radius, y - radius, size, size);
    
    // Restore graphics state (removes clipping path)
    pdf.restoreGraphicsState();
    
    // Draw circle border AFTER restoring (outside clipping context)
    pdf.setDrawColor(200, 200, 200);
    pdf.setLineWidth(0.2);
    pdf.circle(x, y, radius);
    pdf.stroke();
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
    const bioLines = panelist.bio ? pdf.splitTextToSize(panelist.bio, contentWidth - 20) : [];
    const neededSpace = 15 + (panelist.tagline ? 5 : 0) + (bioLines.length * 5);
    
    checkPageBreak(neededSpace);

    // Draw avatar
    const avatarUrl = panelist.avatarUrl || defaultAvatar;
    const avatarData = portraitCache[avatarUrl] || portraitCache[defaultAvatar];
    if (avatarData) {
      drawCircularAvatar(avatarData, margin + 5, yPosition + 3, avatarSize / 2);
    }

    // Panelist name (aligned with top of avatar)
    const textStartX = margin + avatarSize + 8;
    pdf.setFont('helvetica', 'bold');
    pdf.setTextColor(31, 41, 55);
    pdf.text(panelist.name, textStartX, yPosition + 4);
    yPosition += 6;

    // Tagline
    if (panelist.tagline) {
      pdf.setFont('helvetica', 'italic');
      pdf.setTextColor(107, 114, 128);
      pdf.text(panelist.tagline, textStartX, yPosition);
      yPosition += 5;
    }

    // Bio
    if (panelist.bio) {
      pdf.setFont('helvetica', 'normal');
      pdf.setTextColor(75, 85, 99);
      pdf.text(bioLines, textStartX, yPosition);
      yPosition += (bioLines.length * 5);
    }
    
    yPosition += 8; // Spacing between panelists

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

    const isModerator = panelist.id === 'moderator';
    
    // Estimate space needed for this message
    const bubbleWidth = contentWidth * 0.7; // Use 70% of page width
    const messageLines = pdf.splitTextToSize(message.text, bubbleWidth - 12);
    const bubbleHeight = 12 + (messageLines.length * 5.5); // Increased line height
    const neededSpace = bubbleHeight + 6; // More spacing between bubbles
    
    // Check if we need a new page BEFORE starting to draw
    if (yPosition + neededSpace > pageHeight - margin) {
      pdf.addPage();
      yPosition = margin;
    }

    // Alternate bubble position for visual variety (except moderator - always centered)
    const bubbleX = isModerator 
      ? margin + (contentWidth - bubbleWidth) / 2 // Center moderator
      : (index % 2 === 0) 
        ? margin + 15 // Left-aligned for even messages
        : margin + contentWidth - bubbleWidth - 15; // Right-aligned for odd
    
    // Avatar position
    const avatarX = bubbleX - 8;
    const avatarY = yPosition + 3;

    // Draw chat bubble background with subtle shadow effect
    // Shadow
    pdf.setFillColor(220, 220, 220);
    pdf.roundedRect(bubbleX + 0.5, yPosition + 0.5, bubbleWidth, bubbleHeight, 3, 3, 'F');
    
    // Main bubble
    if (isModerator) {
      pdf.setFillColor(237, 233, 254); // purple-100 for moderator
    } else {
      pdf.setFillColor(249, 250, 251); // gray-50 for panelists
    }
    pdf.setDrawColor(209, 213, 219); // gray-300
    pdf.setLineWidth(0.3);
    pdf.roundedRect(bubbleX, yPosition, bubbleWidth, bubbleHeight, 3, 3, 'FD');

    // Draw circular avatar
    const avatarUrl = panelist.avatarUrl || defaultAvatar;
    const avatarData = portraitCache[avatarUrl] || portraitCache[defaultAvatar];
    if (avatarData) {
      drawCircularAvatar(avatarData, avatarX, avatarY, avatarSize / 2);
    }

    // Speaker name
    pdf.setFont('helvetica', 'bold');
    if (isModerator) {
      pdf.setTextColor(109, 40, 217); // purple-700 for moderator
    } else {
      pdf.setTextColor(5, 150, 105); // green-600 for panelists
    }
    pdf.text(panelist.name, bubbleX + 5, yPosition + 5);
    
    // Message text
    pdf.setFont('helvetica', 'normal');
    pdf.setTextColor(31, 41, 55); // gray-800
    pdf.text(messageLines, bubbleX + 5, yPosition + 10);
    
    yPosition += bubbleHeight + 6; // More spacing between messages

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
export const exportDebatePDF = async (debateData, filename) => {
  try {
    const pdf = await generateDebatePDF(debateData);
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
