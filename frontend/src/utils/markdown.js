/**
 * Parse inline Markdown formatting (bold, italic, bold-italic)
 * and convert to HTML or plain text with markers
 */

/**
 * Convert Markdown inline formatting to HTML
 * Supports: ***bold italic***, **bold**, *italic*
 * @param {string} text - Text with Markdown formatting
 * @returns {string} HTML string with formatting
 */
export const markdownToHtml = (text) => {
  if (!text) return '';
  
  // Escape HTML first to prevent XSS
  let escaped = text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
  
  // Convert Markdown formatting (order matters!)
  // 1. Bold italic: ***text*** or ___text___
  escaped = escaped.replace(/\*\*\*(.+?)\*\*\*/g, '<strong><em>$1</em></strong>');
  escaped = escaped.replace(/___(.+?)___/g, '<strong><em>$1</em></strong>');
  
  // 2. Bold: **text** or __text__
  escaped = escaped.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
  escaped = escaped.replace(/__(.+?)__/g, '<strong>$1</strong>');
  
  // 3. Italic: *text* or _text_
  escaped = escaped.replace(/\*(.+?)\*/g, '<em>$1</em>');
  escaped = escaped.replace(/_(.+?)_/g, '<em>$1</em>');
  
  return escaped;
};

/**
 * Parse Markdown formatting into segments for PDF rendering
 * Returns array of {text, bold, italic} objects
 * @param {string} text - Text with Markdown formatting
 * @returns {Array} Array of text segments with formatting flags
 */
export const parseMarkdownSegments = (text) => {
  if (!text) return [{ text: '', bold: false, italic: false }];
  
  const segments = [];
  let currentPos = 0;
  
  // Regex to match all formatting patterns
  const formatRegex = /(\*\*\*(.+?)\*\*\*|___(.+?)___|__(.+?)__|_(.+?)_|\*\*(.+?)\*\*|\*(.+?)\*)/g;
  let match;
  
  while ((match = formatRegex.exec(text)) !== null) {
    // Add plain text before this match
    if (match.index > currentPos) {
      segments.push({
        text: text.substring(currentPos, match.index),
        bold: false,
        italic: false
      });
    }
    
    // Determine formatting type and content
    let content, bold = false, italic = false;
    
    if (match[2] || match[3]) {
      // ***bold italic*** or ___bold italic___
      content = match[2] || match[3];
      bold = true;
      italic = true;
    } else if (match[4] || match[6]) {
      // **bold** or __bold__
      content = match[4] || match[6];
      bold = true;
    } else if (match[5] || match[7]) {
      // *italic* or _italic_
      content = match[5] || match[7];
      italic = true;
    }
    
    segments.push({ text: content, bold, italic });
    currentPos = match.index + match[0].length;
  }
  
  // Add remaining plain text
  if (currentPos < text.length) {
    segments.push({
      text: text.substring(currentPos),
      bold: false,
      italic: false
    });
  }
  
  return segments.length > 0 ? segments : [{ text, bold: false, italic: false }];
};

/**
 * Strip all Markdown formatting for plain text output
 * @param {string} text - Text with Markdown formatting
 * @returns {string} Plain text without formatting markers
 */
export const stripMarkdown = (text) => {
  if (!text) return '';
  
  return text
    .replace(/\*\*\*(.+?)\*\*\*/g, '$1')
    .replace(/___(.+?)___/g, '$1')
    .replace(/\*\*(.+?)\*\*/g, '$1')
    .replace(/__(.+?)__/g, '$1')
    .replace(/\*(.+?)\*/g, '$1')
    .replace(/_(.+?)_/g, '$1');
};
