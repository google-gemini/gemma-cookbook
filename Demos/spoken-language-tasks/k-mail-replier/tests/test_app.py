import sys
import unittest
import os
from unittest.mock import MagicMock, patch, mock_open

class TestApp(unittest.TestCase):
    def setUp(self):
        # Prepare mocks
        self.mock_flask = MagicMock()
        self.mock_gemma = MagicMock()
        self.mock_keras_nlp = MagicMock()
        self.mock_dotenv = MagicMock()

        # Patch sys.modules
        self.patcher = patch.dict(sys.modules, {
            'flask': self.mock_flask,
            'k_mail_replier.models.gemma': self.mock_gemma,
            'keras_nlp': self.mock_keras_nlp,
            'dotenv': self.mock_dotenv
        })
        self.patcher.start()

        # Ensure path
        current_dir = os.path.dirname(os.path.abspath(__file__))
        parent_dir = os.path.dirname(current_dir)
        if parent_dir not in sys.path:
            sys.path.insert(0, parent_dir)

        # Remove app from sys.modules if present to ensure re-import with mocks
        if 'k_mail_replier.app' in sys.modules:
            del sys.modules['k_mail_replier.app']

        from k_mail_replier import app
        self.app = app

    def tearDown(self):
        self.patcher.stop()
        if 'k_mail_replier.app' in sys.modules:
            del sys.modules['k_mail_replier.app']

    def test_get_test_email_success(self):
        """Test happy path where file exists."""
        expected_content = "This is a test email."
        with patch('builtins.open', mock_open(read_data=expected_content)) as mock_file:
            result = self.app.get_test_email()

            # Assert result
            self.assertEqual(result, expected_content)

            # Assert open call
            mock_file.assert_called_once_with('data/email-001-ko.txt', 'r', encoding='utf-8')

    def test_get_test_email_file_not_found(self):
        """Test error path where file does not exist."""
        with patch('builtins.open', side_effect=FileNotFoundError):
            result = self.app.get_test_email()

            # Assert result
            self.assertEqual(result, "Error: File not found!")

if __name__ == '__main__':
    unittest.main()
