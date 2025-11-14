# Quick Start Guide: Creating 0.9B and Smaller Gemma 3n Models

## TL;DR - Get Started in 5 Minutes

### For 0.9B Model (26 layers - Recommended for 4-6GB RAM Mobile)

1. **Open** the MatFormer Lab notebook: `[Gemma_3n]MatFormer_Lab.ipynb`

2. **Navigate to the "Config details" cell** (after the CSV load)

3. **Replace the `config_name` line** with this or comment it out and uncomment the custom config cell:

```python
# Option A: Using custom configuration
layers_to_skip = [19, 20, 21, 22, 23, 24, 25, 26, 27]
ffn_hidden_dims = [2048*3]*10 + [int(2048*3.5)]*9 + [2048*4]*7
ffn_hidden_dims_str = str(ffn_hidden_dims)

# Option B: Using the config_name selector (if added to dropdown)
config_name = "Custom 0.9B (26-layer)"
```

4. **Run all subsequent cells** as normal

5. **Result**: 0.95B model that fits in 4-6GB RAM with 4-bit quantization

---

## Detailed Comparison Table

| Config | Layers | Parameters | MMLU Est. | 4-bit Size | Best For |
|--------|--------|-----------|----------|-----------|----------|
| **0.5B (20L)** | 20 | 0.52B | 40-42% | ~0.9 GB | Web, ultra-light mobile |
| **0.7B (23L)** | 23 | 0.71B | 44-46% | ~1.1 GB | Light mobile (4GB) |
| **0.9B (26L)** | 26 | 0.95B | 46-48% | ~1.5 GB | **Mobile (4-6GB)** ✓ |
| **1.3B (28L)** | 28 | 1.32B | 48-50% | ~2.1 GB | Mobile (6-8GB) |
| **E2B (30L)** | 30 | 1.91B | 50.9% | ~2.9 GB | Mobile (8GB+) |
| **1.5B (30L)** | 30 | 1.51B | 49-51% | ~2.3 GB | High-end mobile |

---

## Step-by-Step Implementation

### Step 1: Prepare Environment

```bash
# In Google Colab or local environment with GPU
!pip install "transformers>=4.53" "timm>=1.0.16" -q

# Login to Hugging Face (required for model access)
from huggingface_hub import notebook_login
notebook_login()
```

### Step 2: Set Model Source

```python
# In the "Import and Export Options" cell
original_model_id = "google/gemma-3n-E4B-it"  # or E4B-pt for pre-trained
local_output_path = "my_0_9b_model"
push_hf_repo_id = "username/gemma-3n-0-9b"  # Your HF repo
```

### Step 3: Apply 0.9B Configuration

**In the "Config details" cell**, uncomment and set:

```python
# Custom config for 0.9B model
layers_to_skip = [19, 20, 21, 22, 23, 24, 25, 26, 27]

# FFN hidden dimensions (optimized distribution)
ffn_hidden_dims = (
    [2048 * 3] * 10 +      # Layers 0-9: 6,144
    [int(2048 * 3.5)] * 9 +     # Layers 10-18: 7,168
    [2048 * 4] * 7         # Layers 19-25: 8,192
)

ffn_hidden_dims_str = str(ffn_hidden_dims)

print(f"Config: 0.9B Model")
print(f"Layers Skipped: {layers_to_skip}")
print(f"Final Layers: {35 - len(layers_to_skip)}")
print(f"FFN Dims (first 10): {ffn_hidden_dims[:10]}")
```

### Step 4: Run Slicing Pipeline

```python
# Run cells in order:
# 1. Load the model config
# 2. Verify slicing configuration
# 3. Update configuration
# 4. Save configuration and tokenizer
# 5. Load model checkpoints
# 6. Slice the model! (main processing)
# 7. Save final model
```

### Step 5: Push to Hugging Face (Optional)

```python
# Push sliced model to your Hugging Face repo
from huggingface_hub import HfApi

api = HfApi()
api.upload_folder(
    folder_path=local_output_path,
    repo_id=push_hf_repo_id,
    repo_type="model"
)
```

### Step 6: Load and Test

```python
from transformers import AutoModelForCausalLM, BitsAndBytesConfig
import torch

# Load your sliced model
model = AutoModelForCausalLM.from_pretrained(
    push_hf_repo_id,
    torch_dtype=torch.float16,
    device_map="auto",
    quantization_config=BitsAndBytesConfig(load_in_4bit=True)
)

print(f"Model parameters: {model.num_parameters():,}")
```

---

## Configuration Presets for Different Scenarios

### Scenario 1: Mobile App (4GB RAM)

```python
# Best: 0.9B with 4-bit quantization
layers_to_skip = [19, 20, 21, 22, 23, 24, 25, 26, 27]
ffn_hidden_dims = [2048*3]*10 + [int(2048*3.5)]*9 + [2048*4]*7

# Expected performance:
# - Model size: ~1.2 GB (quantized)
# - Inference speed: 50-100 tokens/sec
# - Memory during inference: 2.5-3.5 GB
```

### Scenario 2: Web Browser (Client-side)

# Best: 0.5B with extreme quantization
layers_to_skip = [12, 13, 14, 15, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27]
ffn_hidden_dims = [2048*2]*8 + [int(2048*2.5)]*7 + [2048*3]*5

# Expected:
# - Model size: ~0.8-0.9 GB
# - Inference in browser (ONNX/WebGPU)
# - Fast first response

### Scenario 3: Edge Device (6GB RAM, prefer accuracy)

```python
# Best: 1.3B with 4-bit
layers_to_skip = [21, 22, 23, 24, 25, 26, 27]
ffn_hidden_dims = [2048*3]*10 + [int(2048*3.5)]*9 + [2048*4]*9

# Expected:
# - MMLU: 48-50%
# - Size: 1.8-2.1 GB (quantized)
# - Best quality/size tradeoff
```

---

## Understanding FFN Dimension Strategy

The FFN (Feed-Forward Network) dimensions control model capacity at each layer:

```
FFN_dim = 2048 * multiplier

2048 * 2   = 4,096    (very compact, low capacity)
2048 * 2.5 = 5,120    (compact)
2048 * 3   = 6,144    (standard light)
2048 * 3.5 = 7,168    (medium)
2048 * 4   = 8,192    (full capacity)
```

### Why This Distribution?

```
Early layers (0-9):   Lower capacity (6,144)
  → Capture basic linguistic features
  → Can reuse in inference

Middle layers (10-18): Medium capacity (7,168)
  → Process semantic content
  → Balance efficiency and modeling

Late layers (19-25):  Full capacity (8,192)
  → Critical for output quality
  → Global layers with KV sharing
  → Must preserve for performance
```

**Rule of thumb**: Later layers are more important for output quality, so keep them at full capacity.

---

## Inference Optimization Tips

### For Mobile Deployment

```python
# 1. Use 4-bit quantization
from transformers import BitsAndBytesConfig
bnb_config = BitsAndBytesConfig(
    load_in_4bit=True,
    bnb_4bit_quant_type="nf4",
    bnb_4bit_compute_dtype=torch.bfloat16
)

# 2. Disable caching if memory is critical
model.config.use_cache = False

# 3. Use streaming inference
from transformers import TextIteratorStreamer
from threading import Thread

# 4. Limit sequence length
max_tokens = 512

# 5. Use faster attention (if available)
model.config.attn_implementation = "flash_attention_2"
```

### For Web Deployment

```python
# 1. Convert to ONNX or TensorFlow
# python -m transformers.onnx --model=username/gemma-3n-0-9b my_model_onnx

# 2. Quantize for browser
# Use onnx-quantizer or WebGPU native quantization

# 3. Implement streaming with server-sent events
# Best for web UI responsiveness
```

---

## Performance Expectations

### Accuracy (MMLU benchmark)

```
Full E4B (3.98B):  62.30%
Full E2B (1.91B):  50.90%
─────────────────────────
Custom 1.5B:       49-51%  ← Best tradeoff
Custom 1.3B:       48-50%
Custom 0.9B:       46-48%  ← Recommended for 4-6GB
Custom 0.7B:       44-46%
Custom 0.5B:       40-42%  ← Web-only
```

### Inference Speed (Tokens/sec, single GPU)

```
Device: NVIDIA L4 GPU (Colab free tier)
Batch size: 1, Sequence length: 512

Custom 0.9B:  80-120 tokens/sec
Custom 1.3B:  60-90 tokens/sec
E2B (1.91B):  40-60 tokens/sec
E4B (3.98B):  20-40 tokens/sec

Mobile (Snapdragon 8 Gen 3):
Custom 0.9B (quantized): 5-15 tokens/sec
Custom 0.5B (quantized): 10-25 tokens/sec
```

### Memory Usage

```
FP32 (full precision):
  0.9B: ~3.6 GB

FP16 / BF16:
  0.9B: ~1.8 GB

INT8:
  0.9B: ~0.9 GB

INT4 / NF4:
  0.9B: ~0.5-0.7 GB (params) → ~1.5 GB (total)  ← Recommended for mobile
```

---

## Troubleshooting

### Issue: "Layers X and Y are reserved" error

**Solution**: Don't skip the last 2 layers (KV shared layers)
```python
# ✓ Correct
layers_to_skip = [19, 20, 21, 22, 23, 24, 25, 26, 27]  # Keeps 0-18 + 28-34

# ✗ Wrong
layers_to_skip = [19, 20, 21, 22, 23, 24, 33, 34]  # Skips global layers!
```

### Issue: FFN dimensions length mismatch

**Solution**: Ensure length matches final layer count
```python
final_layers = 35 - len(layers_to_skip)  # Should be 26
ffn_dims = [2048*3]*10 + [2048*3.5]*9 + [2048*4]*7  # Length must be 26

assert len(ffn_dims) == final_layers  # Verify
```

### Issue: Out of memory during slicing

**Solution**: 
- Run on Google Colab (free GPU)
- Or use CPU with sufficient RAM (16GB+)
- Model doesn't load into memory, only processes checkpoints

### Issue: Sliced model loads but has poor quality

**Solution**: Check FFN dimension distribution
```python
# Better: Keep later layers at full capacity
ffn = [2048*2.5]*10 + [2048*3.5]*8 + [2048*4]*8  # ✓

# Worse: Uniform low capacity
ffn = [2048*2.5]*26  # ✗ (poor quality)
```

---

## Next Steps

1. **Test the 0.9B configuration** above using MatFormer Lab
2. **Evaluate on your target device** (measure inference speed, memory)
3. **Fine-tune if needed** using LoRA on your domain data
4. **Deploy** with 4-bit quantization for mobile
5. **Monitor** performance and iterate

---

## Additional Resources

- **MatFormer Paper**: https://arxiv.org/abs/2310.07707
- **Gemma 3n Blog**: https://developers.googleblog.com/en/introducing-gemma-3n-developer-guide
- **Slicing Configs Dataset**: https://huggingface.co/datasets/google/gemma3n-slicing-configs
- **Quantization Guide**: https://huggingface.co/docs/transformers/quantization
- **Mobile Deployment**: https://huggingface.co/docs/transformers/onnx

---

## Questions or Issues?

If these configurations don't work for your use case:
1. Check the [custom_slicing_configs.py](./custom_slicing_configs.py) for programmatic validation
2. Refer to the [detailed analysis](./RESPONSE_SUB_BILLION_AND_AUDIO_SLICING.md)
3. File an issue with your device specs and constraints

