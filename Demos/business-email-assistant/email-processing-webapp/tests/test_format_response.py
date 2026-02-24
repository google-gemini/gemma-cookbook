import unittest
from unittest.mock import MagicMock
import sys
import os
import json

# Mock external dependencies to avoid environment issues and side effects
# 1. Mock flask
sys.modules['flask'] = MagicMock()

# 2. Mock models.gemma (to avoid model loading)
mock_gemma = MagicMock()
sys.modules['models.gemma'] = mock_gemma
mock_gemma.create_message_processor.return_value = MagicMock()

# Add the parent directory to sys.path to allow importing app
# app.py is located at ../app.py relative to this test file
current_dir = os.path.dirname(os.path.abspath(__file__))
parent_dir = os.path.dirname(current_dir)
if parent_dir not in sys.path:
    sys.path.insert(0, parent_dir)

from app import format_response

class TestFormatResponse(unittest.TestCase):
    def test_valid_json(self):
        """Test with valid JSON string."""
        input_text = '{"key": "value"}'
        # Expected output is pretty-printed JSON with 2-space indentation
        expected = '{\n  "key": "value"\n}'
        self.assertEqual(format_response(input_text), expected)

    def test_markdown_json(self):
        """Test with JSON wrapped in markdown code blocks."""
        input_text = '```json\n{"key": "value"}\n```'
        expected = '{\n  "key": "value"\n}'
        self.assertEqual(format_response(input_text), expected)

    def test_plain_markdown(self):
        """Test with JSON wrapped in plain markdown code blocks."""
        input_text = '```\n{"key": "value"}\n```'
        expected = '{\n  "key": "value"\n}'
        self.assertEqual(format_response(input_text), expected)

    def test_single_quotes_json(self):
        """Test with JSON-like string using single quotes."""
        # The function replaces single quotes with double quotes to fix JSON
        input_text = "{'key': 'value'}"
        expected = '{\n  "key": "value"\n}'
        self.assertEqual(format_response(input_text), expected)

    def test_invalid_json(self):
        """Test with invalid JSON string."""
        input_text = "This is not JSON"
        # Invalid JSON is returned as-is (after processing)
        expected = "This is not JSON"
        self.assertEqual(format_response(input_text), expected)

    def test_whitespace_handling(self):
        """Test with surrounding whitespace."""
        input_text = '   {"key": "value"}   '
        expected = '{\n  "key": "value"\n}'
        self.assertEqual(format_response(input_text), expected)

    def test_empty_string(self):
        """Test with empty string."""
        input_text = ""
        expected = ""
        self.assertEqual(format_response(input_text), expected)

    def test_single_quote_in_content(self):
        """Test behavior when single quotes are part of the content."""
        # This test documents the current behavior: single quotes inside strings are replaced
        input_text = '{"message": "It\'s a trap"}'

        # Current implementation transforms this to: {"message": "It"s a trap"}
        # This becomes invalid JSON, so it catches the exception and returns the transformed string
        expected = '{"message": "It"s a trap"}'
        self.assertEqual(format_response(input_text), expected)

if __name__ == '__main__':
    unittest.main()
