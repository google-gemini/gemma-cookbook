import unittest
from unittest.mock import MagicMock
import sys
import copy

# Mock pandas since it's not available in the environment
mock_pd = MagicMock()
sys.modules["pandas"] = mock_pd

from custom_slicing_configs import validate_config, SubBillionConfigs

class TestCustomSlicingConfigs(unittest.TestCase):
    def test_valid_configs(self):
        """Test that all predefined configs are valid."""
        configs = [
            SubBillionConfigs.CONFIG_0_5B_20LAYERS,
            SubBillionConfigs.CONFIG_0_7B_23LAYERS,
            SubBillionConfigs.CONFIG_0_9B_26LAYERS,
            SubBillionConfigs.CONFIG_1_3B_28LAYERS,
            SubBillionConfigs.CONFIG_1_5B_30LAYERS,
        ]
        for cfg in configs:
            is_valid, errors = validate_config(cfg)
            self.assertTrue(is_valid, f"Config {cfg['name']} should be valid, but errors found: {errors}")
            self.assertEqual(len(errors), 0)

    def test_layer_mismatch(self):
        """Test detection of mismatch between num_layers and layers_to_skip."""
        cfg = copy.deepcopy(SubBillionConfigs.CONFIG_0_9B_26LAYERS)
        cfg["num_layers"] = 25  # Should be 26 (35 - 9)
        is_valid, errors = validate_config(cfg)
        self.assertFalse(is_valid)
        self.assertIn("Layer mismatch: 26 expected but 25 specified", errors)

    def test_ffn_dims_length_mismatch(self):
        """Test detection of mismatch between ffn_hidden_dims length and num_layers."""
        cfg = copy.deepcopy(SubBillionConfigs.CONFIG_0_9B_26LAYERS)
        cfg["ffn_hidden_dims"] = cfg["ffn_hidden_dims"][:-1]
        is_valid, errors = validate_config(cfg)
        self.assertFalse(is_valid)
        self.assertIn("FFN dims length (25) != num_layers (26)", errors)

    def test_non_integer_ffn_dims(self):
        """Test detection of non-integer FFN dimensions."""
        cfg = copy.deepcopy(SubBillionConfigs.CONFIG_0_9B_26LAYERS)
        cfg["ffn_hidden_dims"][0] = 2048.5
        is_valid, errors = validate_config(cfg)
        self.assertFalse(is_valid)
        self.assertIn("FFN dimensions must be integers.", errors)

    def test_ffn_dim_too_small(self):
        """Test detection of FFN dimension below base_hidden_size."""
        cfg = copy.deepcopy(SubBillionConfigs.CONFIG_0_9B_26LAYERS)
        cfg["ffn_hidden_dims"][0] = 1024  # Less than 2048
        is_valid, errors = validate_config(cfg)
        self.assertFalse(is_valid)
        self.assertTrue(any("outside reasonable range" in e for e in errors))

    def test_ffn_dim_too_large(self):
        """Test detection of FFN dimension above max_ffn_dim."""
        cfg = copy.deepcopy(SubBillionConfigs.CONFIG_0_9B_26LAYERS)
        cfg["ffn_hidden_dims"][0] = 20000  # More than 16384
        is_valid, errors = validate_config(cfg)
        self.assertFalse(is_valid)
        self.assertTrue(any("outside reasonable range" in e for e in errors))

    def test_custom_base_parameters(self):
        """Test validate_config with custom base model parameters."""
        # Create a config for a base model with 40 layers, hidden size 1024, max FFN 4096
        cfg = {
            "num_layers": 38,
            "layers_to_skip": [38, 39],
            "ffn_hidden_dims": [2048] * 38
        }

        # Should be valid with these custom parameters
        is_valid, errors = validate_config(
            cfg,
            base_model_num_layers=40,
            base_hidden_size=1024,
            max_ffn_dim=4096
        )
        self.assertTrue(is_valid, f"Errors: {errors}")

        # Should be invalid if we use defaults (35 layers)
        is_valid, errors = validate_config(cfg)
        self.assertFalse(is_valid)
        self.assertIn("Layer mismatch: 33 expected but 38 specified", errors)

if __name__ == "__main__":
    unittest.main()
