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

import unittest
import sys
import os
from unittest.mock import MagicMock, patch

# Mock dependencies to allow running tests without them installed
# and to avoid side effects like loading heavy models.
sys.modules['dotenv'] = MagicMock()
mock_keras_nlp = MagicMock()
sys.modules['keras_nlp'] = mock_keras_nlp
# Ensure consistency between import keras_nlp; keras_nlp.models and import keras_nlp.models
sys.modules['keras_nlp.models'] = mock_keras_nlp.models

# Add parent directory to path to allow importing gemma_service
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

# Import after mocks
from gemma_service.gemma_model import initialize_model, gemma_model_id

class TestGemmaModelInitialization(unittest.TestCase):

    @patch('gemma_service.gemma_model.os.getenv')
    @patch('gemma_service.gemma_model.load_dotenv')
    @patch('keras_nlp.models.GemmaCausalLM.from_preset')
    def test_initialize_model_success(self, mock_from_preset, mock_load_dotenv, mock_getenv):
        """Test successful initialization when environment variables are set."""
        # Setup mocks
        def getenv_side_effect(key):
            if key == 'KAGGLE_USERNAME':
                return 'test_user'
            if key == 'KAGGLE_KEY':
                return 'test_key'
            return None
        mock_getenv.side_effect = getenv_side_effect

        # Call function
        model = initialize_model()

        # Verify
        mock_load_dotenv.assert_called_once()
        mock_getenv.assert_any_call('KAGGLE_USERNAME')
        mock_getenv.assert_any_call('KAGGLE_KEY')
        mock_from_preset.assert_called_once_with(gemma_model_id)
        self.assertEqual(model, mock_from_preset.return_value)

    @patch('gemma_service.gemma_model.os.getenv')
    @patch('gemma_service.gemma_model.load_dotenv')
    def test_initialize_model_missing_username(self, mock_load_dotenv, mock_getenv):
        """Test initialization fails when KAGGLE_USERNAME is missing."""
        # Setup mocks
        def getenv_side_effect(key):
            if key == 'KAGGLE_USERNAME':
                return None
            if key == 'KAGGLE_KEY':
                return 'test_key'
            return None
        mock_getenv.side_effect = getenv_side_effect

        # Call function and expect error
        with self.assertRaises(ValueError) as context:
            initialize_model()

        self.assertIn("KAGGLE_USERNAME environment variable not found", str(context.exception))
        mock_load_dotenv.assert_called_once()

    @patch('gemma_service.gemma_model.os.getenv')
    @patch('gemma_service.gemma_model.load_dotenv')
    def test_initialize_model_missing_key(self, mock_load_dotenv, mock_getenv):
        """Test initialization fails when KAGGLE_KEY is missing."""
        # Setup mocks
        def getenv_side_effect(key):
            if key == 'KAGGLE_USERNAME':
                return 'test_user'
            if key == 'KAGGLE_KEY':
                return None
            return None
        mock_getenv.side_effect = getenv_side_effect

        # Call function and expect error
        with self.assertRaises(ValueError) as context:
            initialize_model()

        self.assertIn("KAGGLE_KEY environment variable not found", str(context.exception))
        mock_load_dotenv.assert_called_once()

if __name__ == '__main__':
    unittest.main()
