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
import os
import re
from dotenv import load_dotenv

# Set the backbend before importing Keras
os.environ["KERAS_BACKEND"] = "jax"
# Avoid memory fragmentation on JAX backend.
os.environ["XLA_PYTHON_CLIENT_MEM_FRACTION"] = "1.00"

import keras_nlp

def initialize_model():
    """Loads environment variables and configures the Gemma model."""
    load_dotenv()
    # load Kaggle account info for downloading Gemma
    kaggle_username = os.getenv('KAGGLE_USERNAME')
    if not kaggle_username:
        raise ValueError("KAGGLE_USERNAME environment variable not found. Did you set it in your .env file?")
    kaggle_key = os.getenv('KAGGLE_KEY')
    if not kaggle_key:
        raise ValueError("KAGGLE_KEY environment variable not found. Did you set it in your .env file?")

    # create instance using Gemma 2 2B instruction tuned model
    gemma = keras_nlp.models.GemmaCausalLM.from_preset("gemma2_instruct_2b_en")
    #gemma.summary() # REMOVE: FOR TESTING ONLY

    # load and compile tuned model weights
    gemma.backbone.enable_lora(rank=4)
    gemma.backbone.load_lora_weights(f"./weights/gemma2-2b_k-tuned.lora.h5")
    #gemma.compile(sampler="top_k")
    gemma.compile(sampler=keras_nlp.samplers.TopKSampler(k=3, temperature=0.1))

    return gemma  # Return the initialized model

def create_message_processor():
    """Creates a message processor function with a persistent model."""
    model = initialize_model()

    def process_message(prompt_text):
        """Processes a message using a local Gemma model."""
        input = f"<start_of_turn>user\n{prompt_text}<end_of_turn>\n<start_of_turn>model\n"
        response = model.generate(input, max_length=512)
        # remove response tags
        response = extract_substring(response)

        print(response) # REMOVE: FOR TESTING ONLY
        return response

    return process_message

def extract_substring(text):
  """
  Extracts the substring between "<end_of_turn>\n<start_of_turn>model" and the next "<end_of_turn>".

  Args:
      text: The input text.

  Returns:
      The extracted substring, or None if no match is found.
  """
  match = re.search(r".*<end_of_turn>\n<start_of_turn>model\n(.*)<end_of_turn>", text, re.DOTALL)
  if match:
    return match.group(1).strip()
  else:
    return text

# default method
if __name__ == "__main__":
    process_message = create_message_processor()
    process_message("roses are red")
    #print(extract_substring("<start_of_turn>user\nTHE PROMPT<end_of_turn>\n<start_of_turn>model\nTHE TEXT RESPONSE"))