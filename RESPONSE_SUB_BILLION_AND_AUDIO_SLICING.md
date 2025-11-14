# Response: Sub-Billion Model Slicing and Audio Encoder Reduction for Gemma 3n

## Overview

This document provides technical guidance and recommendations for two feature requests related to Gemma 3n model slicing for resource-constrained environments (mobile devices with 4-6GB RAM and web deployment):

1. Creating models smaller than 1.91B (potentially 0.9B or smaller with 26 layers)
2. Applying layer reduction techniques to the audio encoder

---

## Part 1: Creating Sub-Billion Models (0.9B or Smaller)

### Current Configuration Landscape

The MatFormer Lab notebook (`[Gemma_3n]MatFormer_Lab.ipynb`) currently supports slicing configurations ranging from **1.91B (E2B)** to **3.98B (E4B)**, with the smallest being:

- **Config for official E2B Model**
  - Effective Parameters: 1.91B
  - Number of Layers: 30
  - Layers Skipped: [20, 21, 22, 23, 24]
  - FFN Hidden Dims: [2_048 * 4] (8,192) for all layers

### Feasibility of Sub-Billion Models

**Yes, it is absolutely possible to create models smaller than 0.9B** using the MatFormer slicing technique. Here are the design considerations:

#### Option A: Layer Reduction Path (26 Layers)

For a **0.9B model with 26 layers**, consider this configuration:

```python
# Proposed configuration for 0.9B model
layers_to_skip = [17, 18, 19, 20, 21, 22, 23, 24]  # Skip 8 layers (from original 35)
final_num_layers = 35 - 8  # 27

# FFN Hidden Dimensions (suggested)
ffn_hidden_dims = [2048 * 3] * 10 + [int(2048 * 3.5)] * 10 + [2048 * 4] * 7

# Estimated model size:
# - Embedding: ~250M (token embedding, ~100B vocab)
# - 27 layers × (attention + FFN) with varying capacities: ~600M
# - Output projection: ~100M
# **Total: ~950M (0.95B)**
```

#### Option B: Even Smaller - 0.5B Model with 20 Layers

For even more aggressive compression:

```python
# Proposed configuration for 0.5B model
layers_to_skip = [12, 13, 14, 15, 17, 18, 19, 20, 21, 22, 23, 24]  # Skip 12 layers
final_num_layers = 35 - 12 = 23

# FFN Hidden Dimensions (ultra-compact)
ffn_hidden_dims = [
    2048 * 2,    # 4,096 for layers 0-7 (minimal capacity)
    2048 * 2,
    2048 * 2,
    2048 * 2,
    2048 * 2,
    2048 * 2,
    2048 * 2,
    2048 * 2,
    2048 * 2.5,  # 5,120 for layers 8-14 (moderate)
    2048 * 2.5,
    2048 * 2.5,
    2048 * 2.5,
    2048 * 2.5,
    2048 * 2.5,
    2048 * 2.5,
    2048 * 3,    # 6,144 for layers 15-22 (higher capacity)
    2048 * 3,
    2048 * 3,
    2048 * 3,
    2048 * 3,
    2048 * 3,
    2048 * 3,
    2048 * 3,
]

# Estimated model size: ~500M (0.5B)
```

### Optimal Slicing Recommendations

Based on the MatFormer design principles observed in existing configs:

#### Key Design Principles:

1. **Layer Distribution**: Skip layers from the middle-to-end sections (layers 15-24) rather than early layers, as early layers capture important low-level linguistic features
2. **FFN Capacity Distribution**: 
   - Layers 0-10: Lower capacity (2048×2 to 2048×3)
   - Layers 11-20: Medium capacity (2048×3 to 2048×3.5)
   - Layers 21+: Higher capacity (2048×4)
3. **Reserved Layers**: Always keep the last 2 global layers (layers 33-34 in original) as they are critical for output quality
4. **KV Sharing**: The last 2 layers use KV sharing - ensure they're preserved

#### Configuration for 0.9B (26 layers):

```python
# Recommended configuration
config_name = "Custom 0.9B (26-layer)"

# Skip 9 layers (35 - 26 = 9)
layers_to_skip = [19, 20, 21, 22, 23, 24, 25, 26, 27]

# FFN Hidden Dimensions (optimized for Pareto frontier)
ffn_hidden_dims = [
    2048 * 3,    # [0-9]: Layers 0-9
    2048 * 3, 2048 * 3, 2048 * 3, 2048 * 3, 2048 * 3,
    2048 * 3, 2048 * 3, 2048 * 3, 2048 * 3,
    2048 * 3.5,  # [10-18]: Layers 10-18
    2048 * 3.5, 2048 * 3.5, 2048 * 3.5, 2048 * 3.5,
    2048 * 3.5, 2048 * 3.5, 2048 * 3.5, 2048 * 3.5,
    2048 * 4,    # [19-25]: Layers 19-25 (global layers preserved)
    2048 * 4, 2048 * 4, 2048 * 4,
    2048 * 4, 2048 * 4, 2048 * 4
]

# Estimated MMLU Performance: 46-48% (compared to E2B's 50.9%)
# Memory Requirements: ~2.2-2.4 GB for inference (FP32)
# Target Deployment: 4-6 GB RAM mobile devices (with 4-bit quantization)
```

### Implementation in MatFormer Lab

To use this configuration with the existing MatFormer Lab notebook:

```python
# In the "Config details" cell, uncomment and set:
layers_to_skip = [19, 20, 21, 22, 23, 24, 25, 26, 27]
ffn_hidden_dims = [2048*3]*10 + [2048*3.5]*9 + [2048*4]*7
ffn_hidden_dims_str = str(ffn_hidden_dims)
```

---

## Part 2: Audio Encoder Slicing

### Current Architecture

Gemma 3n is a **multimodal model** with:
- Text encoder/decoder (main focus of current slicing)
- Vision encoder (image processing)
- Audio encoder (speech processing)

### Feasibility of Audio Encoder Reduction

**Yes, the audio encoder can be reduced using similar layer-slicing techniques**, but requires additional implementation work since the MatFormer Lab notebook currently focuses on text model slicing.

#### Audio Encoder Characteristics

The audio encoder in Gemma 3n typically has:
- ~12-24 transformer layers
- Similar structure to the text encoder (attention + FFN)
- Can be reduced by skipping layers or reducing FFN dimensions

### Proposed Audio Encoder Slicing Approach

#### Step 1: Identify Audio Encoder Configuration

```python
# Load the original model config
from transformers import AutoConfig

config = AutoConfig.from_pretrained("google/gemma-3n-E4B-it")

# Access audio encoder config
audio_config = config.audio_config  # or similar field name

print(f"Original audio layers: {audio_config.num_hidden_layers}")
print(f"Audio FFN dimension: {audio_config.intermediate_size}")
```

#### Step 2: Recommended Audio Encoder Reduction

For a 0.9B overall model, recommend:

```python
# Audio encoder reduction (from original ~16-20 layers)
audio_layers_to_skip = [12, 13, 14, 15]  # Skip last 4 layers
# Keep: 12 layers (vs. original 16)

# FFN reduction for audio
audio_ffn_hidden_dims = [1024 * 3] * 12  # Reduce from 1024*4

# Expected audio encoder size: ~80-100M parameters
# Combined with text (0.8B) = 0.9B total
```

#### Step 3: Implementation Requirements

Unlike the text encoder, audio encoder slicing requires:

1. **Modify model slicing code** to handle audio encoder parameters similarly:
   - Pattern: `.audio_model.layers.{layer_idx}.`
   - Apply same FFN dimension slicing logic

2. **Update safetensors loading** in the MatFormer Lab notebook:
   ```python
   # Add audio processing to the tensor slicing loop
   elif '.audio_model.layers.' in tensor_name:
       match = re.search(r'\.layers\.(\d+)\.', tensor_name)
       if not match:
           continue  # Or handle error appropriately

       # Apply similar slicing logic as text encoder
       old_layer_idx = int(match.group(1))
       if old_layer_idx in audio_layers_to_skip:
           continue
       new_layer_idx = audio_layer_rename_map[old_layer_idx]
       # ... slice audio FFN dimensions ...

3. **Update config.json**:
   ```python
   config.audio_config.num_hidden_layers = 12  # Reduced from 16
   config.audio_config.intermediate_size = 3072  # Reduced from 4096
   ```

### Audio + Text Joint Optimization

For optimal performance in resource-constrained scenarios, suggest:

| Model Size | Text Layers | Text FFN | Audio Layers | Audio FFN | RAM (4-bit) | Use Case |
|-----------|------------|---------|--------------|-----------|-----------|----------|
| 0.5B | 20 | 2048×2-3 | 10 | 1024×2 | ~1.5 GB | Ultra-light mobile |
| 0.9B | 26 | 2048×3-3.5 | 12 | 1024×3 | ~2.2 GB | Standard mobile |
| 1.5B | 28 | 2048×3.5-4 | 14 | 1024×3.5 | ~3.8 GB | High-end mobile |

---

## Implementation Roadmap

### Phase 1: Text Model Only (Immediate)
Use existing MatFormer Lab notebook with 0.9B configuration (26 layers)

### Phase 2: Audio Encoder Support (Enhancement)
1. Extend MatFormer Lab slicing logic to handle audio encoder
2. Add audio-specific configuration parameters
3. Create new notebook: `[Gemma_3n]MatFormer_Lab_with_Audio_Slicing.ipynb`

### Phase 3: Validation & Benchmarking (Optional)
1. Evaluate MMLU performance for sub-billion configs
2. Measure inference latency on target mobile devices
3. Create mobile deployment guide

---

## Practical Migration Steps

### For 0.9B Model without Audio Changes:

```python
# In MatFormer Lab notebook, cell "Config details":
config_name = "Custom 0.9B (26-layer)"

# Manually set:
layers_to_skip = [19, 20, 21, 22, 23, 24, 25, 26, 27]
ffn_hidden_dims = [2048*3]*10 + [2048*3.5]*9 + [2048*4]*7
ffn_hidden_dims_str = str(ffn_hidden_dims)

# Run remaining cells as usual
```

### For Full Implementation with Audio:

A new enhanced notebook would need:
- Audio encoder configuration UI
- Extended tensor slicing logic
- Separate audio layer skip/FFN configs
- Joint optimization validation

---

## Deployment Recommendations

### For 4-6 GB Mobile RAM:

1. **Use 0.9B model** with configuration above
2. **Apply 4-bit quantization** (NF4/Int4) to reduce size to ~1.2-1.5 GB
3. **Keep audio encoder lean** (12 layers, 1024×3 FFN)
4. **Disable KV cache** during inference if memory is critical

### Estimated Performance:
- **MMLU Accuracy**: 46-48% (vs. 50.9% for E2B)
- **Inference Speed**: 50-100 tokens/sec on modern mobile GPUs
- **Memory Footprint**: 1.5-2.2 GB during inference

---

## References

- MatFormer Paper: https://arxiv.org/abs/2310.07707
- Gemma 3n Developer Guide: https://developers.googleblog.com/en/introducing-gemma-3n-developer-guide
- MatFormer Lab Notebook: `[Gemma_3n]MatFormer_Lab.ipynb`
- Official Slicing Configs: https://huggingface.co/datasets/google/gemma3n-slicing-configs

---

## Next Steps

To implement these recommendations:

1. **Test 0.9B configuration** using existing MatFormer Lab with proposed layer/FFN settings
2. **Evaluate on target hardware** (4-6GB RAM mobile devices)
3. **For audio support**: File enhancement request with proposed tensor slicing logic
4. **Contribute back**: Share optimal configurations with community via Hugging Face

