# Gemma Research

This directory contains research experiments and examples using Google's Gemma models.

## Projects

*   **[VaultGemma](#vaultgemma)**: Privacy-focused fine-tuning with differential privacy.
*   **[T5Gemma](#t5gemma)**: Encoder-decoder variant of Gemma.

---

## VaultGemma

VaultGemma is a privacy-focused variant of Google's Gemma model family, designed for secure fine-tuning and deployment with differential privacy guarantees. This implementation demonstrates how to fine-tune VaultGemma 1B on medical data using LoRA (Low-Rank Adaptation) and Opacus differential privacy.

### Features

- **4-bit Quantization**: Memory-efficient training using BitsAndBytes
- **LoRA Fine-tuning**: Parameter-efficient adaptation with <2% trainable parameters
- **Differential Privacy**: Privacy-preserving training with configurable ε and δ budgets
- **Medical Q&A**: Fine-tuned on medical flashcard dataset for healthcare applications

### Repository Structure

This repository contains code for both training and inference:

* [Fine-tuning](#fine-tuning): How to fine-tune VaultGemma with differential privacy
* [Inference](#inference): How to load and run fine-tuned VaultGemma models

### Fine-tuning

| Notebook Name | Description |
:-------------- | ----------- |
| [[VaultGemma]FineTuning_Inference_Huggingface.ipynb]([VaultGemma]FineTuning_Inference_Huggingface.ipynb) | Complete pipeline for fine-tuning VaultGemma 1B on medical data using LoRA adapters and differential privacy, with inference example |

#### Training Features
- Medical Meadow Medical Flashcards dataset
- 4-bit NF4 quantization for reduced memory footprint
- LoRA adapters targeting all projection layers
- Opacus differential privacy (ε=3.0, δ=1e-5)
- Cosine learning rate schedule with warmup
- Automatic checkpointing based on loss thresholds

### Inference

The same notebook includes inference code to:
- Load fine-tuned LoRA adapters
- Generate responses to medical questions
- Process single or batch queries
- Adjust generation parameters (temperature, top_p)

#### Quick Start

```python
from transformers import AutoModelForCausalLM, GemmaTokenizer
from peft import PeftModel

# Load model and adapters
model = AutoModelForCausalLM.from_pretrained("google/vaultgemma-1b")
tokenizer = GemmaTokenizer.from_pretrained("google/vaultgemma-1b")
peft_model = PeftModel.from_pretrained(model, "path/to/adapters")

# Generate response
question = "What are the symptoms of diabetes?"
response = generate_response(question)
```

### Requirements

```
torch
transformers
peft
opacus
datasets
bitsandbytes
kagglehub
```

### Privacy Guarantees

This implementation provides (ε, δ)-differential privacy guarantees:
- **Target ε**: 3.0 (configurable)
- **Target δ**: 1e-5 (inverse of dataset size)
- **Gradient clipping**: Max norm of 1.0
- **Privacy accounting**: Automatic epsilon tracking via Opacus

---

## T5Gemma

T5Gemma (aka encoder-decoder Gemma) is a family of encoder-decoder large language models, developed by adapting pretrained decoder-only models into an encoder-decoder architecture.

### Notebooks

| Notebook Name | Description |
| :--- | :--- |
| [[T5Gemma]Example.ipynb]([T5Gemma]Example.ipynb) | Guide to sampling and fine-tuning T5Gemma using Flax and Hugging Face |

### Features

- **Encoder-Decoder Architecture**: Adapts decoder-only Gemma models to T5-style architecture.
- **Scales**:
    - **Gemma 2 scale**: 2B-2B, 9B-2B, and 9B-9B.
    - **T5 scale**: Small, Base, Large, XL, and ML.
- **Frameworks**: Examples provided for both **Hugging Face** (PyTorch) and **Flax** (Kauldron).
- **Tasks**:
    - **Sampling**: Basic text generation examples.
    - **Fine-tuning**: Example of fine-tuning for machine translation (English to French) using the MTNT dataset.

### Requirements

```
gemma
kauldron
etils
optax
treescope
kagglehub
transformers
datasets
```
