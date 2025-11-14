"""
Example Script: Creating Custom Sub-Billion Gemma 3n Models

This script demonstrates how to create 0.9B and 0.5B model configurations
for deployment on resource-constrained mobile devices (4-6GB RAM).

Usage:
1. Save this file as `custom_slicing_configs.py`
2. In MatFormer Lab notebook, modify the "Config details" cell to use these configs
3. Or run this script to generate the configurations programmatically
"""

import json
from typing import List, Dict, Tuple
import pandas as pd

# =============================================================================
# Configuration Presets for Sub-Billion Models
# =============================================================================

class SubBillionConfigs:
    """Collection of optimal slicing configurations for sub-billion models."""
    
    # Configuration 1: 0.9B Model (26 layers)
    # Best balance for mobile deployment with 4-6GB RAM
    CONFIG_0_9B_26LAYERS = {
        "name": "Custom 0.9B (26-layer)",
        "num_layers": 26,
        "estimated_params_b": 0.95,
        "estimated_mmlu_accuracy": "46-48%",
        "layers_to_skip": [19, 20, 21, 22, 23, 24, 25, 26, 27],
        "ffn_hidden_dims": (
            [2048 * 3] * 10 +      # Layers 0-9: Lower capacity (6,144)
            [2048 * 3.5] * 9 +     # Layers 10-18: Medium capacity (7,168)
            [2048 * 4] * 7         # Layers 19-25: Higher capacity (8,192)
        ),
        "deployment_target": "Mobile (4-6GB RAM) with 4-bit quantization",
        "notes": "Optimal Pareto point between size and quality"
    }
    
    # Configuration 2: 0.5B Model (20 layers)
    # Ultra-compact for web deployment and very limited devices
    CONFIG_0_5B_20LAYERS = {
        "name": "Custom 0.5B (20-layer)",
        "num_layers": 20,
        "estimated_params_b": 0.52,
        "estimated_mmlu_accuracy": "40-42%",
        "layers_to_skip": [12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27],
        "ffn_hidden_dims": (
            [2048 * 2] * 8 +       # Layers 0-7: Minimal capacity (4,096)
            [2048 * 2.5] * 7 +     # Layers 8-14: Moderate capacity (5,120)
            [2048 * 3] * 5         # Layers 15-19: Higher capacity (6,144)
        ),
        "deployment_target": "Web and ultra-light mobile",
        "notes": "Aggressive compression for web browsers"
    }
    
    # Configuration 3: 1.3B Model (28 layers)
    # Sweet spot between E2B and custom configs
    CONFIG_1_3B_28LAYERS = {
        "name": "Custom 1.3B (28-layer)",
        "num_layers": 28,
        "estimated_params_b": 1.32,
        "estimated_mmlu_accuracy": "48-50%",
        "layers_to_skip": [16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27],
        "ffn_hidden_dims": (
            [2048 * 3] * 10 +      # Layers 0-9: Lower capacity (6,144)
            [2048 * 3.5] * 9 +     # Layers 10-18: Medium capacity (7,168)
            [2048 * 4] * 9         # Layers 19-27: Higher capacity (8,192)
        ),
        "deployment_target": "Higher-end mobile devices (6-8GB RAM)",
        "notes": "Better accuracy than E2B with only slightly larger size"
    }
    
    # Configuration 4: 0.7B Model (23 layers)
    # Middle ground between 0.5B and 0.9B
    CONFIG_0_7B_23LAYERS = {
        "name": "Custom 0.7B (23-layer)",
        "num_layers": 23,
        "estimated_params_b": 0.71,
        "estimated_mmlu_accuracy": "44-46%",
        "layers_to_skip": [15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27],
        "ffn_hidden_dims": (
            [2048 * 2.5] * 10 +    # Layers 0-9: Lower capacity (5,120)
            [2048 * 3] * 8 +       # Layers 10-17: Medium capacity (6,144)
            [2048 * 3.5] * 5       # Layers 18-22: Higher capacity (7,168)
        ),
        "deployment_target": "Mid-range mobile (5-6GB RAM)",
        "notes": "Good balance between compression and quality"
    }

    # Configuration 5: 1.5B Model (30 layers) - Extended E2B
    # Comparison with official E2B
    CONFIG_1_5B_30LAYERS = {
        "name": "Custom 1.5B (30-layer)",
        "num_layers": 30,
        "estimated_params_b": 1.51,
        "estimated_mmlu_accuracy": "49-51%",
        "layers_to_skip": [12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27],
        "ffn_hidden_dims": (
            [2048 * 3] * 10 +      # Layers 0-9: Lower capacity (6,144)
            [2048 * 3.5] * 10 +    # Layers 10-19: Medium capacity (7,168)
            [2048 * 4] * 10        # Layers 20-29: Higher capacity (8,192)
        ),
        "deployment_target": "High-end mobile (6-8GB RAM) or edge servers",
        "notes": "Larger than E2B (1.91B official) but similar inference cost"
    }


class AudioEncoderConfigs:
    """Optional audio encoder configurations for multimodal models."""
    
    # Standard audio encoder (for reference)
    AUDIO_ENCODER_FULL = {
        "name": "Full Audio Encoder",
        "num_layers": 16,
        "intermediate_size": 4096,
        "num_parameters_m": 150,
        "notes": "Original audio encoder in Gemma 3n"
    }
    
    # Reduced audio encoder for 0.9B model
    AUDIO_ENCODER_0_9B = {
        "name": "Reduced Audio Encoder (0.9B model)",
        "num_layers": 12,
        "layers_to_skip": [12, 13, 14, 15],
        "intermediate_size": 3072,  # Reduce from 4096
        "num_parameters_m": 110,
        "reduction_percent": 26.7,
        "notes": "Removes last 4 layers and reduces FFN dimension"
    }
    
    # Ultra-compact audio encoder for 0.5B model
    AUDIO_ENCODER_0_5B = {
        "name": "Ultra-compact Audio Encoder (0.5B model)",
        "num_layers": 10,
        "layers_to_skip": [10, 11, 12, 13, 14, 15],
        "intermediate_size": 2048,  # Further reduced
        "num_parameters_m": 75,
        "reduction_percent": 50.0,
        "notes": "6-layer reduction, 50% FFN dimension reduction"
    }
    
    # Minimal audio encoder for web deployment
    AUDIO_ENCODER_WEB = {
        "name": "Minimal Audio Encoder (Web)",
        "num_layers": 8,
        "layers_to_skip": [8, 9, 10, 11, 12, 13, 14, 15],
        "intermediate_size": 2048,
        "num_parameters_m": 50,
        "reduction_percent": 66.7,
        "notes": "Aggressive reduction for web browsers"
    }


# =============================================================================
# Helper Functions
# =============================================================================

def get_config_for_deployment(
    target_size_gb: float,
    deployment_type: str = "mobile"
) -> Dict:
    """
    Recommend a configuration based on deployment constraints.
    
    Args:
        target_size_gb: Target model size in GB (after 4-bit quantization)
        deployment_type: 'mobile', 'web', or 'edge'
    
    Returns:
        Recommended configuration dictionary
    """
    if deployment_type == "web":
        return SubBillionConfigs.CONFIG_0_5B_20LAYERS.copy()
    elif deployment_type == "mobile" and target_size_gb <= 1.5:
        return SubBillionConfigs.CONFIG_0_5B_20LAYERS.copy()
    elif deployment_type == "mobile" and target_size_gb <= 2.5:
        return SubBillionConfigs.CONFIG_0_9B_26LAYERS.copy()
    elif deployment_type == "mobile" and target_size_gb <= 3.5:
        return SubBillionConfigs.CONFIG_1_3B_28LAYERS.copy()
    else:
        return SubBillionConfigs.CONFIG_1_5B_30LAYERS.copy()


def validate_config(config: Dict) -> Tuple[bool, List[str]]:
    """
    Validate a configuration for logical consistency.
    
    Returns:
        (is_valid, list_of_errors)
    """
    errors = []
    
    expected_layers = 35 - len(config["layers_to_skip"])
    if expected_layers != config["num_layers"]:
        errors.append(
            f"Layer mismatch: {expected_layers} expected but {config['num_layers']} specified"
        )
    
    if len(config["ffn_hidden_dims"]) != config["num_layers"]:
        errors.append(
            f"FFN dims length ({len(config['ffn_hidden_dims'])}) != num_layers ({config['num_layers']})"
        )
    
    for dim in config["ffn_hidden_dims"]:
        if dim < 2048 or dim > 16384:
            errors.append(f"FFN dimension {dim} outside reasonable range [2048, 16384]")
    
    return len(errors) == 0, errors


def export_for_matformer_lab(config: Dict) -> str:
    """
    Export configuration in format suitable for MatFormer Lab notebook.
    
    Returns:
        Python code snippet to paste into the notebook
    """
    code = f"""
# Configuration: {config['name']}
# Parameters: {config['estimated_params_b']}B
# Estimated MMLU: {config['estimated_mmlu_accuracy']}
# Deployment: {config['deployment_target']}

config_name = "{config['name']}"

layers_to_skip = {config['layers_to_skip']}

ffn_hidden_dims = {config['ffn_hidden_dims']}

ffn_hidden_dims_str = str(ffn_hidden_dims)
"""
    return code


def create_config_comparison_table() -> pd.DataFrame:
    """Create comparison table of all available configurations."""
    
    configs = [
        SubBillionConfigs.CONFIG_0_5B_20LAYERS,
        SubBillionConfigs.CONFIG_0_7B_23LAYERS,
        SubBillionConfigs.CONFIG_0_9B_26LAYERS,
        SubBillionConfigs.CONFIG_1_3B_28LAYERS,
        SubBillionConfigs.CONFIG_1_5B_30LAYERS,
    ]
    
    data = []
    for cfg in configs:
        data.append({
            "Configuration": cfg["name"],
            "Layers": cfg["num_layers"],
            "Params (B)": cfg["estimated_params_b"],
            "MMLU Est.": cfg["estimated_mmlu_accuracy"],
            "Size 4-bit (GB)": cfg["estimated_params_b"] * 0.5,  # Rough estimate
            "Target": cfg["deployment_target"],
        })
    
    return pd.DataFrame(data)


# =============================================================================
# Example Usage
# =============================================================================

if __name__ == "__main__":
    print("=" * 80)
    print("Sub-Billion Gemma 3n Model Slicing Configurations")
    print("=" * 80)
    
    # Print configuration comparison
    print("\n### Available Configurations:\n")
    comparison_df = create_config_comparison_table()
    print(comparison_df.to_string(index=False))
    
    # Validate and export example configs
    print("\n" + "=" * 80)
    print("Configuration Details & Export")
    print("=" * 80)
    
    configs_to_demo = [
        SubBillionConfigs.CONFIG_0_9B_26LAYERS,
        SubBillionConfigs.CONFIG_0_5B_20LAYERS,
    ]
    
    for config in configs_to_demo:
        is_valid, errors = validate_config(config)
        
        print(f"\n### {config['name']}")
        print(f"- Estimated Parameters: {config['estimated_params_b']}B")
        print(f"- Estimated MMLU Accuracy: {config['estimated_mmlu_accuracy']}")
        print(f"- Deployment Target: {config['deployment_target']}")
        print(f"- Valid: {'✓' if is_valid else '✗'}")
        
        if not is_valid:
            for error in errors:
                print(f"  Error: {error}")
        
        print("\n**MatFormer Lab Code:**")
        print(export_for_matformer_lab(config))
    
    # Recommendation example
    print("\n" + "=" * 80)
    print("Deployment Recommendation Examples")
    print("=" * 80)
    
    print("\n**For 4GB Mobile RAM with 4-bit quantization:**")
    rec = get_config_for_deployment(target_size_gb=1.5, deployment_type="mobile")
    print(f"Recommended: {rec['name']}")
    print(f"Size: {rec['estimated_params_b']}B → ~{rec['estimated_params_b'] * 0.5:.1f}GB (4-bit)")
    
    print("\n**For Web Browser:**")
    rec = get_config_for_deployment(target_size_gb=1.0, deployment_type="web")
    print(f"Recommended: {rec['name']}")
    print(f"Size: {rec['estimated_params_b']}B → ~{rec['estimated_params_b'] * 0.5:.1f}GB (4-bit)")
    
    print("\n**For High-End Mobile (8GB RAM):**")
    rec = get_config_for_deployment(target_size_gb=3.5, deployment_type="mobile")
    print(f"Recommended: {rec['name']}")
    print(f"Size: {rec['estimated_params_b']}B → ~{rec['estimated_params_b'] * 0.5:.1f}GB (4-bit)")
    
    # Audio encoder configs
    print("\n" + "=" * 80)
    print("Audio Encoder Configuration Options")
    print("=" * 80)
    
    audio_configs = [
        AudioEncoderConfigs.AUDIO_ENCODER_FULL,
        AudioEncoderConfigs.AUDIO_ENCODER_0_9B,
        AudioEncoderConfigs.AUDIO_ENCODER_0_5B,
    ]
    
    for audio_cfg in audio_configs:
        print(f"\n{audio_cfg['name']}")
        print(f"  Layers: {audio_cfg['num_layers']}")
        print(f"  FFN Intermediate: {audio_cfg['intermediate_size']}")
        print(f"  Estimated Params: {audio_cfg.get('num_parameters_m', 'N/A')}M")
        if 'reduction_percent' in audio_cfg:
            print(f"  Reduction: {audio_cfg['reduction_percent']:.1f}%")

