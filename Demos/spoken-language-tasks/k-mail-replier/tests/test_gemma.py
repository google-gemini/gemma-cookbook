
import unittest
from unittest.mock import MagicMock, patch
import os
import sys

# Mock external dependencies that are not installed in the environment
sys.modules['keras_nlp'] = MagicMock()
sys.modules['dotenv'] = MagicMock()

# Ensure the parent directory is in sys.path to import the module under test
current_dir = os.path.dirname(os.path.abspath(__file__))
parent_dir = os.path.dirname(current_dir)
if parent_dir not in sys.path:
    sys.path.insert(0, parent_dir)

from k_mail_replier.models import gemma

class TestGemma(unittest.TestCase):

    def setUp(self):
        # Create a mock for keras_nlp.models.GemmaCausalLM
        self.mock_gemma_causal_lm = MagicMock()
        # Mock the specific class inside the module
        gemma.keras_nlp.models = MagicMock()
        gemma.keras_nlp.models.GemmaCausalLM = self.mock_gemma_causal_lm

        # Setup the mock instance returned by from_preset
        self.mock_model_instance = MagicMock()
        self.mock_gemma_causal_lm.from_preset.return_value = self.mock_model_instance

        # Mock samplers
        gemma.keras_nlp.samplers = MagicMock()
        gemma.keras_nlp.samplers.TopKSampler = MagicMock()

    @patch('k_mail_replier.models.gemma.load_dotenv')
    @patch.dict(os.environ, {}, clear=True)
    def test_initialize_model_missing_kaggle_username(self, mock_load_dotenv):
        """Test that ValueError is raised when KAGGLE_USERNAME is missing."""
        with self.assertRaises(ValueError) as cm:
            gemma.initialize_model()
        self.assertEqual(str(cm.exception), "KAGGLE_USERNAME environment variable not found. Did you set it in your .env file?")

    @patch('k_mail_replier.models.gemma.load_dotenv')
    @patch.dict(os.environ, {'KAGGLE_USERNAME': 'test_user'}, clear=True)
    def test_initialize_model_missing_kaggle_key(self, mock_load_dotenv):
        """Test that ValueError is raised when KAGGLE_KEY is missing."""
        with self.assertRaises(ValueError) as cm:
            gemma.initialize_model()
        self.assertEqual(str(cm.exception), "KAGGLE_KEY environment variable not found. Did you set it in your .env file?")

    @patch('k_mail_replier.models.gemma.load_dotenv')
    @patch.dict(os.environ, {'KAGGLE_USERNAME': 'test_user', 'KAGGLE_KEY': 'test_key'}, clear=True)
    def test_initialize_model_success(self, mock_load_dotenv):
        """Test successful initialization of the model."""
        model = gemma.initialize_model()

        # Verify load_dotenv was called
        mock_load_dotenv.assert_called_once()

        # Verify model creation
        self.mock_gemma_causal_lm.from_preset.assert_called_once_with("gemma2_instruct_2b_en")

        # Verify LoRA configuration
        self.mock_model_instance.backbone.enable_lora.assert_called_once_with(rank=4)
        self.mock_model_instance.backbone.load_lora_weights.assert_called_once()
        args, _ = self.mock_model_instance.backbone.load_lora_weights.call_args
        self.assertTrue(args[0].endswith("gemma2-2b_k-tuned.lora.h5"))

        # Verify compilation
        self.mock_model_instance.compile.assert_called_once()

        # Verify return value
        self.assertEqual(model, self.mock_model_instance)

if __name__ == '__main__':
    unittest.main()
