from ollama import Client
from ollama import chat

client = Client(
  host='<cloud_run_url>',
  headers={'Authorization': 'Bearer <YOUR_API_KEY>'}
)

# Example: non-streaming
response = client.chat(
    model='<model>', # Example: "gemma3:4b" or your custom model name
    messages=[
    {
        'role': 'user',
        'content': 'Why is the sky blue?',
    },
])

print(response['message']['content'])

# Example: streaming
stream = client.chat(
    model='<model>', # Example: "gemma3:4b" or your custom model name
    messages=[{'role': 'user', 'content': 'Why is the sky blue?'}],
    stream=True,
)

for chunk in stream:
  print(chunk['message']['content'], end='', flush=True)
