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
