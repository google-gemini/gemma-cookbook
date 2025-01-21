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
learning_rate_value = 1e-4
train_epochs = 20
lora_name = "gemma2-2b_k-tuned"

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


def generate_from_base_model(prompt_text):
    """Prints 'Starting generation run...' to the console."""
    print("Starting generation run with base model...")
    set_environment()

    # create instance
    gemma = keras_nlp.models.GemmaCausalLM.from_preset(model_id)
    gemma.summary()

    input = f"<start_of_turn>user\n{prompt_text}<end_of_turn>\n<start_of_turn>model\n"
    output = gemma.generate(input, max_length=token_limit)
    print("\nGemma output:")
    print(output)


def prepare_tuning_dataset():
    tokenizer = keras_nlp.models.GemmaTokenizer.from_preset(model_id)

    # prompt structure:
    # <start_of_turn>user
    # 다음에 대한 이메일 답장을 작성해줘.
    # "{EMAIL CONTENT FROM THE CUSTOMER}"
    # <end_of_turn>
    # <start_of_turn>model
    # {MODEL ANSWER}<end_of_turn>

    # load data from repository (or local directory)
    from datasets import load_dataset
    ds = load_dataset(
        # Dataset : https://huggingface.co/datasets/bebechien/korean_cake_boss
        "bebechien/korean_cake_boss",
        split="train",
    )
    print(ds)
    data = ds.with_format("np", columns=["input", "output"], output_all_columns=False)
    tuning_dataset = []

    for x in data:
        item = f"<start_of_turn>user\n다음에 대한 이메일 답장을 작성해줘.\n\"{x['input']}\"<end_of_turn>\n<start_of_turn>model\n{x['output']}<end_of_turn>"
        length = len(tokenizer(item))
        # skip data if the token length is longer than our limit
        if length < token_limit:
            tuning_dataset.append(item)
            if(len(tuning_dataset)>=num_data_limit):
                break

    # FOR TESTING ONLY:
    print(len(tuning_dataset))
    print(tuning_dataset[0])
    print(tuning_dataset[1])
    print(tuning_dataset[2])

    return tuning_dataset


def tune_model_with_lora():
    set_environment()

    # Prepate the dataset
    tuning_dataset = prepare_tuning_dataset()

    # initialize model
    gemma = keras_nlp.models.GemmaCausalLM.from_preset(model_id)

    # Enable LoRA for the model and set the LoRA rank to 4.
    gemma.backbone.enable_lora(rank=lora_rank)
    gemma.summary()

    # Limit the input sequence length (to control memory usage).
    gemma.preprocessor.sequence_length = token_limit
    
    # Use AdamW (a common optimizer for transformer models).
    optimizer = keras.optimizers.AdamW(
        learning_rate=learning_rate_value,
        weight_decay=0.01,
    )

    # Exclude layernorm and bias terms from decay.
    optimizer.exclude_from_weight_decay(var_names=["bias", "scale"])

    gemma.compile(
        loss=keras.losses.SparseCategoricalCrossentropy(from_logits=True),
        optimizer=optimizer,
        weighted_metrics=[keras.metrics.SparseCategoricalAccuracy()],
    )

    class CustomCallback(keras.callbacks.Callback):
        def on_epoch_end(self, epoch, logs=None):
            model_name = f"weights/{lora_name}_{lora_rank}_epoch{epoch+1}.lora.h5"
            gemma.backbone.save_lora_weights(model_name)

    print("Starting tuning run...")
    history = gemma.fit(tuning_dataset, epochs=train_epochs, batch_size=batch_size_value, callbacks=[CustomCallback()])


# default method -----------------------------
if __name__ == "__main__":
    print("Starting the default method")
    # test generation with base model:
    #generate_from_base_model("roses are red")

    # conduct a model tuning run
    tune_model_with_lora()