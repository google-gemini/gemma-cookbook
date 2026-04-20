import express from 'express';
import { AutoProcessor, PaliGemmaForConditionalGeneration, load_image } from "@huggingface/transformers";
import { createCanvas, loadImage } from "canvas"; // Import canvas for drawing
import fs from "fs";
import { createRequire } from 'module';
const require = createRequire(import.meta.url);
const path = require('path');

import { Buffer } from 'node:buffer';
global.Buffer = Buffer;

// Function to generate a random RGB color
function getRandomColor() {
    const r = Math.floor(Math.random() * 256);
    const g = Math.floor(Math.random() * 256);
    const b = Math.floor(Math.random() * 256);
    return `rgb(${r},${g},${b})`;
}

// Load processor and model
const model_id = "NSTiwari/paligemma2-3b-mix-224-onnx"; // Change this to use a different PaliGemma model
const processor = await AutoProcessor.from_pretrained(model_id);
const model = await PaliGemmaForConditionalGeneration.from_pretrained(
    model_id,
    {
        dtype: {
            embed_tokens: "fp16", // or 'fp16'
            vision_encoder: "fp16", // or 'q4', 'fp16'
            decoder_model_merged: "q4", // or 'q4f16'
        },
    }
);
console.log("Model and processor loaded successfully.");

const app = express();
const port = 3000;

// Serve static files from the 'public' directory
app.use(express.static('public'));
app.use(express.json({ limit: '10mb' })); // Increase limit if needed for large images

// API endpoint for image processing
app.post('/process-image', async (req, res) => {
    try {
        const base64Image = req.body.image;
        const prompt = req.body.prompt || "<image>caption en"; // Default prompt if not sent
        const targetWidth = req.body.targetWidth;   //Added targetWidth from body to get the targetWidth
        const targetHeight = req.body.targetHeight; //Added targetHeight from body to get the targetHeight
        const originalWidth = req.body.originalWidth; //Add the originalWidth
        const originalHeight = req.body.originalHeight; //Add the originalHeight

        // Create a buffer from the base64 image data
        const buffer = Buffer.from(base64Image, 'base64');

        // Save the buffer to disk as a JPEG file
        const tempImagePath = path.join(process.cwd(), 'temp_image.jpg'); // Use a unique filename for each image
        fs.writeFileSync(tempImagePath, buffer);

        // Load the image
        const raw_image = await loadImage(tempImagePath);

        // Resize the image using canvas
        const canvas = createCanvas(targetWidth, targetHeight);
        const ctx = canvas.getContext('2d');
        ctx.drawImage(raw_image, 0, 0, targetWidth, targetHeight);

        // Convert canvas to a data URL
        const resizedImageDataURL = canvas.toDataURL('image/jpeg');

        // Extract base64 from resized image
        const resizedBase64Image = resizedImageDataURL.split(',')[1];
        const resizedBuffer = Buffer.from(resizedBase64Image, 'base64');
        const resizedTempImagePath = path.join(process.cwd(), 'resized_temp_image.jpg');
        fs.writeFileSync(resizedTempImagePath, resizedBuffer);

        // Load the resized image with huggingface
        const resizedImage = await load_image(resizedTempImagePath);

        const inputs = await processor(resizedImage, prompt);

        // Generate a response from the model
        const response = await model.generate({
            ...inputs,
            max_new_tokens: 100,
        });

        // Extract the generated IDs from the response
        const generatedIds = response.slice(null, [inputs.input_ids.dims[1], null]);

        // Decode the generated IDs to get the answer
        const decodedAnswer = processor.batch_decode(generatedIds, {
            skip_special_tokens: true,
        });

        console.log("Prompt: ", prompt);
        console.log("Response: ", decodedAnswer[0]);

        if (prompt.includes("<image>detect")) {
            // Parse the response to extract bounding box coordinates
            const boundingBoxes = decodedAnswer[0].match(/<loc(\d+)>/g);

            if (boundingBoxes && boundingBoxes.length === 4) {
                // Extract numbers from <locXXXX> tags and cast to integers
                const coordinates = boundingBoxes.map(tag => parseInt(tag.replace("<loc", "").replace(">", "")));
                const [y1, x1, y2, x2] = coordinates.map(coord => Math.floor(coord));

                // Normalize the bounding box coordinates to the RESIZED image dimensions
                const normX1 = Math.round((x1 / 1024) * targetWidth);
                const normY1 = Math.round((y1 / 1024) * targetHeight);
                const normX2 = Math.round((x2 / 1024) * targetWidth);
                const normY2 = Math.round((y2 / 1024) * targetHeight);

                // Extract the label from the prompt
                const labelMatch = prompt.match(/detect (.*)/);
                const label = labelMatch ? labelMatch[1] : "Unknown";

                // Remove the temp image file
                fs.unlinkSync(tempImagePath);
                fs.unlinkSync(resizedTempImagePath);
                res.json({ success: true, boundingBox: { x1: normX1, y1: normY1, x2: normX2, y2: normY2, label: label}, message: decodedAnswer[0] });
            }
            else {
                // Remove the temp image file if the bounding box coordinates were not found
                fs.unlinkSync(tempImagePath);
                fs.unlinkSync(resizedTempImagePath);
                res.json({ success: false, message: decodedAnswer[0] });
            }
        }
        else {
            // Remove the temp image file if the bounding box coordinates were not found
            fs.unlinkSync(tempImagePath);
            fs.unlinkSync(resizedTempImagePath);
            res.json({ success: false, message: decodedAnswer[0] });
        }

    } catch (error) {
        console.error("Error processing image:", error);
        res.status(500).json({ success: false, message: 'Error processing image' });
    }
});

app.listen(port, () => {
    console.log(`Server listening at http://localhost:${port}`);
});
