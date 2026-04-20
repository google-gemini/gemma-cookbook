// This demo requires importScripts, so we use a local, pre-built JavaScript bundle.

importScripts('/mediapipe_genai_bundle.js');              // from npm @mediapipe/tasks-genai@0.10.24

const { FilesetResolver, LlmInference } = self.BundledCode;

let llmInference;
let modelPath = './myemoji-gemma-3-270m-it.task';         // UPDATE TO MATCH YOUR MODEL FILE

// Listen for messages from the main thread
self.onmessage = async (event) => {
  const { type, data } = event.data;
  console.log("[Worker] Received message:", { type, data });

  switch (type) {
    case "load":
      try {
        console.log("[Worker] Loading model...");
        // Load the model in the worker thread
        const genai = await FilesetResolver.forGenAiTasks('https://cdn.jsdelivr.net/npm/@mediapipe/tasks-genai@latest/wasm');
        llmInference = await LlmInference.createFromOptions(genai, {
          baseOptions: {
            modelAssetPath: modelPath
          },
          maxTokens: 32, 
          temperature: .8,
          forceF32: true,
        });
        console.log("[Worker] Model loaded successfully.");
        self.postMessage({ type: "loaded" });
      } catch (error) {
        console.error("[Worker] Error loading model:", error);
        self.postMessage({ type: "error", data: error.message });
      }
      break;

    case "generate":
      if (!llmInference) {
        console.error("[Worker] Generation failed: model not loaded yet.");
        self.postMessage({ type: "error", data: "Model not loaded yet." });
        return;
      }
      try {
        const generatedResponses = new Set();
        const prompt = `<start_of_turn>user\nTranslate this text to emoji: ${data.prompt}<end_of_turn>\n<start_of_turn>model\n`;
        // Request 3 unique, clean responses from the model
        for (let i = 0; i < 3; i++) {
          const modifiedPrompt = prompt + ' '.repeat(i);
          const rawResponse = await llmInference.generateResponse(modifiedPrompt);
          const cleanResponse = rawResponse.replace(/[^\p{Emoji}\s\u200D]/gu, '').trim();
          if (cleanResponse) {
            generatedResponses.add(cleanResponse);
          }
        }
        generatedResponses.forEach(response => {
        self.postMessage({ type: "result", data: response + '\n' });
        });
        self.postMessage({ type: "complete" });
      } catch (error) {
        console.error("[Worker] Error during generation:", error);
        self.postMessage({ type: "error", data: error.message });
      }
      break;
  }
};
