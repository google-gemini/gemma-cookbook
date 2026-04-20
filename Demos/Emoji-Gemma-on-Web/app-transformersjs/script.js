// Set up DOM references
const generateBtn = document.getElementById("generate-btn");
const promptInput = document.getElementById("prompt-input");
const responseOutput = document.getElementById("response-output");
const statusMessageContainer = document.getElementById("status-message");

const workerPath = './worker.js'; 
let worker;
let modelReady = false;
let loadingDisclaimer = ""; // Variable to hold the WebGPU warning

// Handle the click-to-copy event for an emoji translation result
function handleEmojiButtonClick(combo) {
  navigator.clipboard.writeText(combo);
  statusMessageContainer.textContent = `Copied to clipboard: ${combo}`;
  console.log(`[UI] Copied "${combo}" to clipboard.`);
}

// Create and append a button for an emoji translation result
function createEmojiButton(combo) {
  const button = document.createElement("button");
  button.textContent = combo;
  button.onclick = () => handleEmojiButtonClick(combo);
  responseOutput.appendChild(button);
}

// Main function to run worker
async function initializeModelInWorker() {
  console.log("[UI] Initializing application...");
  
  // Create a progress loader that clears once the model loads
  statusMessageContainer.textContent = "Loading model (0%)";

  // Create a new worker to run the model in a separate thread
  worker = new Worker(workerPath, { type: "module" });
  console.log("[UI] Worker created.");

  // Listen for messages from the worker
  worker.onmessage = (event) => {
    const { type, data } = event.data;

    switch (type) {
      case "status_update":
        loadingDisclaimer = `<span style="display:block">${data}</span>`;
        break;

      case "progress":
        statusMessageContainer.innerHTML = `Loading model (${data.progress}%)${loadingDisclaimer}`;
        break;

      case "loaded":
        statusMessageContainer.innerHTML = `Loading model (100%)${loadingDisclaimer}`; // Show 100% briefly
        modelReady = true;
        generateBtn.disabled = false;               // Enable the generation button upon model load
        setTimeout(() => {           
          statusMessageContainer.innerHTML = ``;    // Then empty status message
        }, 500);
        break;
      
      case "result":
        const line = data.trim();
        if (line) {
          createEmojiButton(line);
        }
        break;
        
      case "complete":
        responseComplete();
        break;
      case "error":
        console.error("Worker error:", data);
        if (modelReady) {
          // Error during generation
          generateBtn.classList.remove('generating');
          statusMessageContainer.textContent =
            "Failed to generate response. Check the console for errors.";
          generateBtn.disabled = false; // Re-enable button
        } else {
          // Error during model loading
          statusMessageContainer.textContent = "Failed to load model. Please refresh the page.";
        }
        break;
    }
  };

  // Start loading the model in the worker
  console.log('[UI] Sending "load" message to worker.');
  worker.postMessage({ type: "load" });
}

function responseComplete() {
  generateBtn.classList.remove('generating');
  generateBtn.disabled = false;
  
  if (responseOutput.childElementCount === 0) {
    statusMessageContainer.textContent = "No results";
  } else {
    // Clear the status message if we have results
    statusMessageContainer.textContent = "Click to copy emojis";
    // For accessibility, move focus to the first result so keyboard users can navigate them
    responseOutput.firstChild.focus();
  }
}

function resetUI() {
  console.log("[UI] Resetting UI state.");
  generateBtn.classList.add('generating');
  generateBtn.disabled = true;
  responseOutput.innerHTML = ""; // Clear previous buttons
  statusMessageContainer.textContent = "Generating..."; // Set generating status
}

async function generateResponse() {
  const prompt = promptInput.value.trim();
  if (!prompt) {
    statusMessageContainer.textContent = "Please enter a prompt.";
    return;
  }

  if (!worker || !modelReady) {
    statusMessageContainer.textContent = "Model is not ready yet. Please wait.";
    return;
  }

  resetUI();

  // Send prompt to worker to start generation
  console.log(
    `[UI] Sending "generate" message to worker with prompt: "${prompt}"`
  );
  worker.postMessage({ type: "generate", data: { prompt } });
}

// After the script is loaded, initialize the text generator
async function initializeAndAttachListeners() {
  console.log("[UI] Starting initialization and attaching listeners.");

  await initializeModelInWorker();

  generateBtn.addEventListener("click", generateResponse);

  // Add event listener for "Enter" key press in the prompt input field
  promptInput.addEventListener("keydown", function (event) {
    if (event.key === "Enter" && !event.shiftKey && !generateBtn.disabled) {
      console.log('[UI] "Enter" key pressed, triggering generation.');
      event.preventDefault(); // Prevents the default action (form submission/new line)
      generateResponse();
    }
  });
}

initializeAndAttachListeners();

