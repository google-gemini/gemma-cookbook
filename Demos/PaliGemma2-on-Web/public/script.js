document.addEventListener('DOMContentLoaded', () => {
    const imageUpload = document.getElementById('imageUpload');
    const processButton = document.getElementById('processButton');
    const originalImage = document.getElementById('originalImage');
    const processedCanvas = document.getElementById('processedCanvas');
    const promptInput = document.getElementById('promptInput');
    const responseTextDiv = document.getElementById('responseText');
    const ctx = processedCanvas.getContext('2d');

    let imageBase64 = '';
    let originalImageURL = '';
    let originalImageObj;

    // Initially hide the original image
    originalImage.style.display = 'none';

    // Handle image upload
    imageUpload.addEventListener('change', (event) => {
        const file = event.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e) => {
                imageBase64 = e.target.result.split(',')[1];
                originalImageURL = e.target.result;
                originalImage.src = originalImageURL;

                originalImageObj = new Image();
                originalImageObj.onload = () => {
                    // Keep aspect ratio while scaling
                    const originalWidth = originalImageObj.width;
                    const originalHeight = originalImageObj.height;

                    // Maximum dimensions for the display area
                    const maxWidth = 600;
                    const maxHeight = 400;

                    // Calculate aspect ratio
                    const aspectRatio = originalWidth / originalHeight;

                    let displayWidth = maxWidth;
                    let displayHeight = maxWidth / aspectRatio;

                    if (displayHeight > maxHeight) {
                        displayHeight = maxHeight;
                        displayWidth = maxHeight * aspectRatio;
                    }

                    processedCanvas.width = displayWidth;
                    processedCanvas.height = displayHeight;
                    ctx.clearRect(0, 0, processedCanvas.width, processedCanvas.height);
                    ctx.drawImage(originalImageObj, 0, 0, displayWidth, displayHeight); // Draw the original image on the canvas

                    // Show the original image after it's loaded for preview
                    originalImage.style.display = 'block';
                };
                originalImageObj.src = originalImageURL;

                // Enable the process button
                processButton.disabled = false;
                responseTextDiv.style.display = 'none';
                responseTextDiv.innerHTML = '';
                promptInput.value = '';
            };
            reader.readAsDataURL(file);
        }
    });

    // Handle process button click
    processButton.addEventListener('click', async () => {
        if (!imageBase64) {
            alert("Please upload an image first.");
            return;
        }

        // Clear off previous results.
        ctx.clearRect(0, 0, processedCanvas.width, processedCanvas.height);
        responseTextDiv.innerHTML = 'Analyzing...';
        responseTextDiv.style.display = 'block';


        let prompt = promptInput.value || "";
        if (prompt.toLowerCase().includes("detect")) {
            const labelMatch = prompt.match(/detect\s+(.*)/i);
            const label = labelMatch ? labelMatch[1] : 'Unknown';
            prompt = `<image>detect ${label}`;
        } else {
            prompt = `<image>${prompt}`;
        }

        try {
            const response = await fetch('/process-image', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ image: imageBase64, prompt: prompt, targetWidth: processedCanvas.width, targetHeight: processedCanvas.height, originalWidth: originalImageObj.width, originalHeight: originalImageObj.height }),
            });

            if (response.ok) {
                const data = await response.json();

                if (data.success) {
                    if (prompt.includes("<image>detect")) {
                        const { boundingBox } = data;
                        drawBoundingBox(boundingBox, ctx)

                        responseTextDiv.style.display = 'block';
                        responseTextDiv.innerHTML = "Response: " + escapeHtml(data.message);
                    }
                    else {
                        processedCanvas.width = 0;
                        processedCanvas.height = 0;
                        ctx.clearRect(0, 0, processedCanvas.width, processedCanvas.height);
                        responseTextDiv.style.display = 'block';
                        responseTextDiv.innerHTML = data.message;
                    }
                }
                else {
                    processedCanvas.width = 0;
                    processedCanvas.height = 0;
                    ctx.clearRect(0, 0, processedCanvas.width, processedCanvas.height);
                    responseTextDiv.style.display = 'block';
                    responseTextDiv.innerHTML = "Response: " + data.message;
                }
            }
            else {
                alert('Error processing image.');
            }

        } catch (error) {
            console.error('Error:', error);
            alert('Error processing image.');
        } finally {

        }
    });


    // Function to draw the bounding box on canvas
    function drawBoundingBox(boundingBox, ctx) {

        const { x1, y1, x2, y2, label } = boundingBox;

        // Generate random color for the bounding box and label background
        const randomColor = getRandomColor();

        // Set styles for the bounding box (random color stroke)
        ctx.strokeStyle = randomColor;
        ctx.lineWidth = 5;
        ctx.strokeRect(x1, y1, x2 - x1, y2 - y1);

        // Adjust label background height to fit the text properly
        const labelPadding = 10;
        const textWidth = ctx.measureText(label.charAt(0).toUpperCase() + label.slice(1)).width;
        const labelWidth = textWidth * 3;
        const labelHeight = 30;
        const labelY = y1 - labelHeight;

        // Draw background for the label (same random color as bounding box)
        ctx.fillStyle = randomColor;
        ctx.fillRect(x1, labelY, labelWidth, labelHeight);

        // Set the text color to white
        ctx.fillStyle = "white";
        ctx.font = "bold 20px Arial";
        ctx.fillText(label.charAt(0).toUpperCase() + label.slice(1), x1 + labelPadding, labelY + labelHeight - labelPadding);
    }

    // Function to generate a random RGB color
    function getRandomColor() {
        const r = Math.floor(Math.random() * 256);
        const g = Math.floor(Math.random() * 256);
        const b = Math.floor(Math.random() * 256);
        return `rgb(${r},${g},${b})`;
    }

    function escapeHtml(unsafe) {
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }
    
});
