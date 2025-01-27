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
from dotenv import load_dotenv
import keras
import datasets
import json

# Set the backbend before importing Keras
os.environ["KERAS_BACKEND"] = "jax"
# Avoid memory fragmentation on JAX backend.
os.environ["XLA_PYTHON_CLIENT_MEM_FRACTION"] = "1.00"
import keras_nlp

# Model and tuning configuration
model_id = "gemma2_instruct_2b_en"
token_limit = 256
batch_size_value = 1
num_data_limit = 100
lora_rank = 4
learning_rate_value = 9e-4
weight_decay_value = 0.004
train_epochs = 3
lora_name = "gemma2-2b_inquiry_tuned"

def set_environment():
    """Loads environment variables needed for execution."""
    load_dotenv()
    # load Kaggle account info for downloading Gemma
    kaggle_username = os.getenv('KAGGLE_USERNAME')
    if not kaggle_username:
        raise ValueError("KAGGLE_USERNAME environment variable not found. Did you set it in your .env file?")
    kaggle_key = os.getenv('KAGGLE_KEY')
    if not kaggle_key:
        raise ValueError("KAGGLE_KEY environment variable not found. Did you set it in your .env file?")

def read_json_files_to_dicts(directory_path):
    """Finds all JSON files in a directory and reads them into a 
        list of dictionary objects.

    Args:
        directory: The directory to search for JSON files.

    Returns:
        An array dictionary objects each containing a "prompt" and "response".
    """
    json_data = []
    for filename in os.listdir(directory_path):
        if filename.endswith(".json"):
            file_path = os.path.join(directory_path, filename)
            with open(file_path, 'r') as f:
                try:
                    data = json.load(f)                    
                    json_data.append(
                        dict(
                            prompt = data["prompt"],
                            response = data["response"]
                        )
                    )
                except json.JSONDecodeError as e:
                    print(f"Error decoding JSON in file {filename}: {e}")
    return json_data

def prepare_tuning_dataset():
    # collect data
    prompt_data = read_json_files_to_dicts("./data")

    # prepare data for tuning
    tuning_dataset = []
    template = "{instruction}\n{response}"

    for prompt in prompt_data:
        tuning_dataset.append(template.format(instruction=prompt["prompt"],response=prompt["response"]))
    
    # print(tuning_dataset) # FOR TESTING ONLY

    return tuning_dataset

def tune_model_with_lora():
    set_environment()

    # Prepate the dataset
    tuning_dataset = prepare_tuning_dataset()

    # initialize model
    gemma = keras_nlp.models.GemmaCausalLM.from_preset(model_id)

    # Limit the input sequence length (to control memory usage).
    gemma.preprocessor.sequence_length = token_limit

    # Enable LoRA for the model and set the LoRA rank to 4.
    gemma.backbone.enable_lora(rank=lora_rank)
    
    # Use AdamW (a common optimizer for transformer models).
    optimizer = keras.optimizers.AdamW(
        learning_rate=learning_rate_value,
        weight_decay=weight_decay_value,
    )

    # Exclude layernorm and bias terms from decay.
    optimizer.exclude_from_weight_decay(var_names=["bias", "scale"])

    gemma.compile(
        loss=keras.losses.SparseCategoricalCrossentropy(from_logits=True),
        optimizer=optimizer,
        weighted_metrics=[keras.metrics.SparseCategoricalAccuracy()],
    )

    # callback to write out weights for each epoch
    class CustomCallback(keras.callbacks.Callback):
        def on_epoch_end(self, epoch, logs=None):
            model_name = f"weights/{lora_name}_{lora_rank}_epoch{epoch+1}.lora.h5"
            gemma.backbone.save_lora_weights(model_name)

    print("Starting tuning run...")
    history = gemma.fit(
        tuning_dataset, 
        epochs=train_epochs, 
        batch_size=batch_size_value, 
        callbacks=[CustomCallback()]
    )

def generate_from_model(prompt_text, use_tuned_weights):
    """Generates a response using the Gemma base (untuned)  model."""
    print("Starting generation run with Gemma base model...")
    set_environment()

    # create instance
    gemma = keras_nlp.models.GemmaCausalLM.from_preset(model_id)

    if use_tuned_weights:
        # load and compile tuned model weights
        gemma.backbone.enable_lora(rank=4)
        gemma.backbone.load_lora_weights(f"./weights/gemma2-2b_inquiry_tuned_4_epoch3.lora.h5")

    # For this use case, the greedy sampler is best
    gemma.compile(sampler="greedy")

    gemma.summary()

    input = f"<start_of_turn>user\n{prompt_text}<end_of_turn>\n<start_of_turn>model\n"
    output = gemma.generate(input, max_length=token_limit)
    print("\nGemma output:")
    print(output)

# default method -----------------------------
if __name__ == "__main__":
    print("Starting the default method")
    # prepare dataset
    #prepare_tuning_dataset()

    # conduct a model tuning run
    tune_model_with_lora()
    