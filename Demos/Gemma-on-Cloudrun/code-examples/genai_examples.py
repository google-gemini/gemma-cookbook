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
