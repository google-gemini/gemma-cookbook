# Deploying Gemma Models to Google Cloud Run

This guide shows you how to deploy Google Gemma LLM to Google Cloud Run.  These pre-built containers leverage Ollama for serving, with additional added support for the Google GenAI SDK . Use off the shelf, or fine-tune for your own use-cases.  These containers have both Google GenAI SDK and OpenAI SDK compatible.

## Supported Models and Pre-Built Docker Images
Our service supports the following Gemma models:
* gemma-3-1b-it
* gemma-3-4b-it
* gemma-3-12b-it
* gemma-3-27b-it

You can provide your own fine-tuned models following [section below](#deploying-and-using-fine-tuned-gemma3-models)

These models are exposed through two Gemma3 API endpoints:
* [`/v1beta/{model=models/*}:generateContent`](https://ai.google.dev/api/generate-content#method:-models.generatecontent) - Generates a model response given an input GenerateContentRequest.
* [`/v1beta/{model=models/*}:streamGenerateContent`](https://ai.google.dev/api/generate-content#method:-models.streamgeneratecontent) - Generates a streamed response from the model given an input GenerateContentRequest.
* Additionally, an OpenAI-compatible endpoint is available at `/v1/chat/completions`.

### Pre-built Docker Images
We provide pre-built Docker images for convenience. These images have the respective Gemma models bundled:
* `us-docker.pkg.dev/cloudrun/container/gemma/gemma3-1b`
* `us-docker.pkg.dev/cloudrun/container/gemma/gemma3-4b`
* `us-docker.pkg.dev/cloudrun/container/gemma/gemma3-12b`
* `us-docker.pkg.dev/cloudrun/container/gemma/gemma3-27b`

## Deploying to Cloud Run
This section guides you through deploying a Cloud Run service using our provided Docker images or your own custom-built image.

Use the following `gcloud run deploy` command to deploy your Cloud Run service:
```
gcloud run deploy {SERVICE_NAME} \
 --image {IMAGE} \
 --concurrency 4 \
 --cpu 8 \
 --set-env-vars OLLAMA_NUM_PARALLEL=4 \
 --set-env-vars=API_KEY={YOUR_API_KEY} \
 --gpu 1 \
 --gpu-type nvidia-l4 \
 --max-instances 1 \
 --memory 32Gi \
 --allow-unauthenticated \
 --no-cpu-throttling \
 --timeout=600 \
 --region {REGION}
```

Explanation of Variables:
* `SERVICE_NAME`: The unique name for your Cloud Run service.
* `IMAGE`: The Docker image to deploy. This can be one of our [pre-built images](#pre-built-docker-images) or an image you built yourself.
* `YOUR_API_KEY`: **Crucial for authentication**. Set this to a strong, unique API key string of your choice. This key will be required to access your service. See the [Authentication](#authentication) section below for more details.
* `REGION`: The Google Cloud region where your Cloud Run service will be deployed (e.g., us-central1). Ensure this region supports the specified GPU type. See [GPU support for Cloud Run services](https://cloud.google.com/run/docs/configuring/services/gpu) for more details.
* For other flags and optimizing setting, see [Run LLM inference on Cloud Run GPUs with Gemma 3 and Ollama](https://cloud.google.com/run/docs/tutorials/gpu-gemma-with-ollama#build-and-deploy) for more details.

After successful deployment, the gcloud command will output the Cloud Run service URL. Save this URL as `<cloud_run_url>` for interacting with your service.

## Authentication
To get started quickly, you can deploy the Cloud Run service with public (unauthenticated) access using `--allow-unauthenticated`.  The service will validate the `API_KEY` environment variable you set during deployment.  Longer term, we recommend enabling IAM authentication and using the google-auth SDK.

#### Setting the API Key
* Environment Variable: As shown in the deployment command: `--set-env-vars=API_KEY={YOUR_API_KEY}`.
* **Secret Manager (recommended for production)**:
For enhanced security, store your API key in Google Cloud Secret Manager and expose it as an environment variable `--update-secrets=API_KEY={yourSecrete}:latest`. For more details, refer to the [Cloud Run Secrets documentation](https://cloud.google.com/run/docs/configuring/services/secrets#access-secrets).

#### Using the API Key in Requests
You will need to include this `YOUR_API_KEY` in every request to your Cloud Run service, as shown in the [Interacting with the Service](#interacting-with-the-cloud-run-service) section.

## Interacting with the Cloud Run Service
Once your Cloud Run service is deployed, you can interact with it using curl, Google's GenAI SDK, or OpenAI's SDK.

Placeholders:
* `<cloud_run_url>`: The URL of your deployed Cloud Run service.
* `<YOUR_API_KEY>`: The API key you configured during deployment.
* `<model>`: The model name you deployed (e.g., gemma-3-1b-it, gemma-3-4b-it, or your fine-tuned model name).

### 1. Using curl:

Generate Content
```
curl "<cloud_run_url>/v1beta/models/<model>:generateContent?key={YOUR_API_KEY}" \
   -H 'Content-Type: application/json' \
   -X POST \
   -d '{
     "contents": [{
       "parts":[{"text": "Write a story about a magic backpack. You are the narrator of an interactive text adventure game."}]
       }]
      }'
```

Stream Generate Content
```
curl "<cloud_run_url>/v1beta/models/<model>:streamGenerateContent?key={YOUR_API_KEY}" \
   -H 'Content-Type: application/json' \
   -X POST \
   -d '{
     "contents": [{
       "parts":[{"text": "Write a story about a magic backpack. You are the narrator of an interactive text adventure game."}]
       }]
      }'
```

### 2. Using Google GenAI SDK (Python)

Refer to the [official GenAI SDK documentation](https://ai.google.dev/gemini-api/docs/libraries) for more details.

#### 2.1 Install GenAI SDK
```
pip install --upgrade google-genai
```

#### 2.2 Python example:
```
from google import genai
from google.genai.types import HttpOptions

# Configure the client to use your Cloud Run endpoint and API key
client = genai.Client(api_key="<YOUR_API_KEY>", http_options=HttpOptions(base_url="<cloud_run_url>"))


# Example: Generate content (non-streaming)
response = client.models.generate_content(
   model="<model>", # Example: "gemma-3-4b-it" or your custom model name
   contents=["How does AI work?"]
)
print(response.text)


# Example: Stream generate content
response = client.models.generate_content_stream(
   model="<model>", # Example: "gemma-3-4b-it" or your custom model name
   contents=["Write a story about a magic backpack. You are the narrator of an interactive text adventure game."]
)
for chunk in response:
   print(chunk.text, end="")
```

### 3. Using OpenAI SDK (Python)

Refer to the [official OpenAI SDK documentation](https://platform.openai.com/docs/libraries#install-an-official-sdk) for more details.

#### 3.1 Install OpenAI SDK
```
pip install openai
```

#### 3.2 Python Code Example
```
from openai import OpenAI

# Configure the OpenAI client to point to your Cloud Run endpoint
openAIclient = OpenAI(
   api_key="<YOUR_API_KEY>",
   base_url="<cloud_run_url>/v1" # Note: Add /v1 to the base_url
)

completion = openAIclient.chat.completions.create(
   model="<model>", # Example: "gemma3:4b" or your custom model name
   messages=[
     {            
       "role": "developer",                                                                                                  
       "content": "You are a helpful assistant."
     },
     {
       "role": "user",
       "content": "Hello!"
     }
   ]
)

print(completion.choices[0].message.content)
```

#### 3.3 `curl` Example (OpenAI Compatible)
```
curl <cloud_run_url>/v1/chat/completions \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer <YOUR_API_KEY>" \
 -d '{
   "model": "<model>",
   "messages": [
     {
       "role": "developer",
       "content": "You are a helpful assistant."
     },
     {
       "role": "user",
       "content": "Hello!"
     }
   ]
 }'
```

## Deploying and Using Fine-Tuned Gemma3 Models

This section details how to deploy and use your own custom fine-tuned Gemma models with the Cloud Run service. This involves creating a custom Ollama model, uploading its components to GCS, and mounting that GCS bucket to your Cloud Run service.

Steps:

#### 1. Customize the Model using Ollama Locally:

Follow https://github.com/ollama/ollama?tab=readme-ov-file#customize-a-model to import GGUF model in the Modelfile, and create the model in Ollama
```
ollama create <your-custom-model-name> -f Modelfile
```
This command will process your GGUF file and create the necessary blobs and manifests for Ollama in your [local Ollama models directory](https://github.com/ollama/ollama/blob/main/docs/faq.md#where-are-models-stored).

#### 2. Locate and Upload Ollama Model Files to GCS:

Navigate to your [local Ollama models directory](https://github.com/ollama/ollama/blob/main/docs/faq.md#where-are-models-stored). You will find `blobs/` and `manifests/` subdirectories. These contain the components of your newly created custom model.

Upload the contents generated for your custom model to your GCS bucket in corresponding `blobs/` and `manifests/` subdirectories. This will ensure the correct structure for Ollama to find your model.

For example,
```
cd <your-local-ollama-model-dir>
gsutil -m cp -r . gs://YOUR_MODEL_BUCKET_NAME
```

#### 3. Deploy Cloud Run Service with GCS Volume Mount

```
gcloud run deploy {SERVICE_NAME} \
 --image {IMAGE} \
 --concurrency 4 \
 --cpu 8 \
 --set-env-vars OLLAMA_NUM_PARALLEL=4 \
 --set-env-vars=API_KEY={YOUR_API_KEY} \
 --gpu 1 \
 --gpu-type nvidia-l4 \
 --max-instances 1 \
 --memory 32Gi \
 --allow-unauthenticated \
 --no-cpu-throttling \
 --timeout=600 \
 --region {REGION} \
 --add-volume name={VOLUME_NAME},type=cloud-storage,bucket={YOUR_MODEL_BUCKET_NAME} \
 --add-volume-mount volume={VOLUME_NAME},mount-path=/models
```

Explanation of Variables:
* `VOLUME_NAME`: A name for your volume (e.g., my-gemma-volume).
* `YOUR_MODEL_BUCKET_NAME`: The name of your GCS bucket.

Note that you have two options for the `{IMAGE}`:

Option A: Using a Pre-built Image with GCS Mount (Recommended for convenience):

You can use one of our [pre-built images](#pre-built-docker-images) (which contain baked-in models). When you mount a GCS bucket to `/models`, the content of the GCS bucket will override and replace any models that were originally baked into the Docker image at that path. This means your custom models from GCS will be served.

Option B: Building Your Own Image (without baked-in models) with GCS Mount:

If you prefer a smaller Docker image or want full control over the image contents, you can build your own image based on our [Dockerfile](./Dockerfile). In this scenario, you will remove the `ollama pull` command from your Dockerfile, ensuring no models are baked into the image. Then, you deploy this custom-built image with the GCS volume mount.

```
# Example Dockerfile snippet (within your build stage)
# ... (other instructions) ...
# Store model weight files in /models (this is where the GCS volume will be mounted)
ENV OLLAMA_MODELS /models
# REMOVE THIS LINE:
# RUN /bin/ollama serve & sleep 5 && ollama pull $MODEL
# ... (rest of your Dockerfile) ...
```
Then [build your custom image](https://cloud.google.com/run/docs/building/containers#use-dockerfile).

#### 4. Interact with Your Custom Model:
Now you can use your custom fine-tuned models by specifying the `<your-custom-model-name>` (the name you used in `ollama create`) in your API requests, just like with the pre-built models.

```
from google import genai
from google.genai.types import HttpOptions

client = genai.Client(api_key="<YOUR_API_KEY>", http_options=HttpOptions(base_url="<cloud_run_url>"))

response = client.models.generate_content_stream(
   model="<your-custom-model-name>", # Use the name you defined in ollama create
   contents=["Write a story about a magic backpack. You are the narrator of an interactive text adventure game."]
)
for chunk in response:
   print(chunk.text, end="")
```
Similarly for OpenAI SDK and curl examples, replace `<model>` with `<your-custom-model-name>`.

