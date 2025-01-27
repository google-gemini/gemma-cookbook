/**
 * @fileoverview Functions used in events attached in templates/index.html.
 */

/**
 * Opens a file select dialog to import the file text into a text area.
 *
 * @param {Element} inputTextArea - Text area to import the file text into.
 */
function openFileSelectDialog(inputTextArea) {
  const fileInput = document.createElement('input');
  fileInput.type = 'file';
  fileInput.accept = '.txt';

  fileInput.addEventListener('change', function() {
    importTextFile(fileInput, inputTextArea);
  });

  fileInput.click();
}

/**
 * Reads an input text file and imports the content into a text area.
 *
 * @param input - File to be imported
 * @param {Element} inputTextArea - Text area to import the file text into
 */
function importTextFile(input, inputTextArea) {
  const file = input.files[0];

  if (file) {
    const reader = new FileReader();

    reader.onload = function(e) {
      inputTextArea.value = e.target.result;
    }

    reader.onerror = function(e) {
      console.error("Error reading file:", e);
    }

    reader.readAsText(file);
  } else {
    console.warn("Please select a text file.");
  }
}

/**
 * Copy the text of the output text area to the user's clipboard.
 *
 * @param {Element} copyTextArea - Text area that has the content to be copied
 * @param {Element} copyTooltip - Tooltip associated with the copy button
 */
function copyToClipboard(copyTextArea, copyTooltip) {
  copyTextArea.select();
  document.execCommand('copy');
  copyTooltip.textContent = 'Copied!';
}

/**
 * Resets the copy to clipboard tooltip text when the user moves the cursor
 * away from the clipboard button.
 *
 * @param {Element} copyTooltip - Tooltip associated with the copy button
 */
function copyTooltipReset(copyTooltip) {
  copyTooltip.textContent = 'Copy to clipboard';
}
