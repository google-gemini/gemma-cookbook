# VaultGemma

VaultGemma is a privacy-focused variant of Google's Gemma model family, designed for secure fine-tuning and deployment with differential privacy guarantees. This implementation demonstrates how to fine-tune VaultGemma 1B on medical data using LoRA (Low-Rank Adaptation) and Opacus differential privacy.

## Features

- **4-bit Quantization**: Memory-efficient training using BitsAndBytes
- **LoRA Fine-tuning**: Parameter-efficient adaptation with <2% trainable parameters
- **Differential Privacy**: Privacy-preserving training with configurable ε and δ budgets
- **Medical Q&A**: Fine-tuned on medical flashcard dataset for healthcare applications

## Repository Structure

This repository contains code for both training and inference:

* [Fine-tuning](#fine-tuning): How to fine-tune VaultGemma with differential privacy
* [Inference](#inference): How to load and run fine-tuned VaultGemma models

## Fine-tuning

| Notebook Name | Description |
:-------------- | ----------- |
| [[VaultGemma]FineTuning_Inference_Huggingface.ipynb]([VaultGemma]FineTuning_Inference_Huggingface.ipynb) | Complete pipeline for fine-tuning VaultGemma 1B on medical data using LoRA adapters and differential privacy, with inference example |

### Training Features
- Medical Meadow Medical Flashcards dataset
- 4-bit NF4 quantization for reduced memory footprint
- LoRA adapters targeting all projection layers
- Opacus differential privacy (ε=3.0, δ=1e-5)
- Cosine learning rate schedule with warmup
- Automatic checkpointing based on loss thresholds

## Inference

The same notebook includes inference code to:
- Load fine-tuned LoRA adapters
- Generate responses to medical questions
- Process single or batch queries
- Adjust generation parameters (temperature, top_p)

### Quick Start

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

## Requirements

```
torch
transformers
peft
opacus
datasets
bitsandbytes
kagglehub
```

## Privacy Guarantees

This implementation provides (ε, δ)-differential privacy guarantees:
- **Target ε**: 3.0 (configurable)
- **Target δ**: 1e-5 (inverse of dataset size)
- **Gradient clipping**: Max norm of 1.0
- **Privacy accounting**: Automatic epsilon tracking via Opacus

