#
# Copyright 2024 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
import google.generativeai as genai
from dotenv import load_dotenv
import os

def initialize_model():
    """Loads environment variables and configures the GenAI client."""
    load_dotenv()
    api_key = os.getenv('API_KEY')
    if not api_key:
        raise ValueError("API_KEY environment variable not found. Did you set it in your .env file?")
    genai.configure(api_key=api_key)
    return genai.GenerativeModel('gemini-1.5-flash')  # Return the initialized model

def create_message_processor():
    """Creates a message processor function with a persistent model."""
    model = initialize_model()

    def process_message(message):
        """Processes a message using the GenAI model."""
        response = model.generate_content(message)
        print(response.text) # REMOVE: FOR TESTING ONLY
        return response.text
    return process_message