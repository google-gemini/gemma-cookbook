# Run a fine-tuned Gemma 3 270M model in the browser with Transformers.js

This app demonstrates how to generate emojis from a text input using a fine-tuned Gemma 3 270M model running directly in the browser. For this demo you'll just change one line of code to point to your ONNX model.

## Run the demo
1. Download the app files in this directory.
2. In the worker.js file, update the model string in the `pipeline()` function call to the model on Hugging Face Hub.
    1. Alternatively, download and place the model files in a new subdirectory i.e. `app-transformersjs/myemoji-gemma-3-270m-it-onnx/` for full offline use.
3. Open terminal on your computer and navigate (`cd`) to the app folder.
4. Run `npx serve` to start the local server.
5. Open the provided `localhost` address in your browser to run the app.

**Requirements:** Browser with [WebGPU support](https://caniuse.com/webgpu)

## How it works
This demo sets up a simple web server to host a frontend where users can enter a text prompt. This starts a generation process in a web worker to avoid blocking the main UI thread. The worker uses [Transformers.js](https://huggingface.co/docs/transformers.js/index) to generate a response from the model and send it back to the user.
 
## Resources
* [Notebook: Fine-tune Gemma 3 270M](https://github.com/google-gemini/gemma-cookbook/blob/main/Demos/Emoji-Gemma-on-Web/resources/Fine_tune_Gemma_3_270M_for_emoji_generation.ipynb)
* [Notebook: Convert Gemma 3 270M to ONNX](https://github.com/google-gemini/gemma-cookbook/blob/main/Demos/Emoji-Gemma-on-Web/resources/Convert_Gemma_3_270M_to_ONNX.ipynb)
* [Hugging Face Transformers.js documentation](https://huggingface.co/docs/transformers.js/index)
