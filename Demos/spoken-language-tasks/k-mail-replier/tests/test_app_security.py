
import unittest
from unittest.mock import MagicMock, patch
import os
import sys

# Mock necessary modules
sys.modules['keras_nlp'] = MagicMock()
sys.modules['dotenv'] = MagicMock()
gemma_mock = MagicMock()
gemma_mock.create_message_processor.return_value = lambda x: "Mocked Response"
sys.modules['k_mail_replier.models.gemma'] = gemma_mock

# Add the project root to sys.path
current_dir = os.path.dirname(os.path.abspath(__file__))
project_root = os.path.abspath(os.path.join(current_dir, '..'))
if project_root not in sys.path:
    sys.path.insert(0, project_root)

# Import app after mocking
try:
    from k_mail_replier.app import app
    import k_mail_replier.app as app_module
except ImportError:
    pass

class TestAppSecurity(unittest.TestCase):

    def setUp(self):
        app.config['TESTING'] = True
        self.client = app.test_client()

    def test_no_global_customer_request(self):
        """Test that customer_request is not a global variable in the module."""
        # This test expects the variable to be removed from module scope.
        # We check hasattr on the module.
        has_global = hasattr(app_module, 'customer_request')
        self.assertFalse(has_global, "Global variable 'customer_request' should be removed to prevent race conditions.")

    @patch('k_mail_replier.app.get_test_email')
    def test_data_isolation(self, mock_get_test_email):
        """Test that one user's data does not leak to another user."""
        mock_get_test_email.return_value = "Default Test Email"

        # 1. User A submits data (POST)
        user_a_secret = "User A Secret Data"
        self.client.post('/', data={'request': user_a_secret})

        # 2. User B visits the page (GET)
        response = self.client.get('/')
        content = response.data.decode('utf-8')

        # 3. Verify User B does not see User A's data
        self.assertNotIn(user_a_secret, content, "User A's data leaked into User B's response.")

        # 4. Verify User B sees default content
        # With the fix, User B should see "Default Test Email" because GET invokes get_test_email()
        self.assertIn("Default Test Email", content, "User B should see the default test email content.")

    def test_get_request_has_content(self):
        """Test that GET request returns a page with content populated."""
        response = self.client.get('/')
        self.assertEqual(response.status_code, 200)
        content = response.data.decode('utf-8')
        # Check for textarea
        self.assertIn('<textarea', content)

if __name__ == '__main__':
    unittest.main()
