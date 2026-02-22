import os
import sys
import unittest
from unittest.mock import MagicMock, patch
import django
from django.conf import settings

# We need to make sure we can import 'api'
sys.path.append(os.getcwd())

# Configure minimal Django settings
if not settings.configured:
    settings.configure(
        INSTALLED_APPS=[
            'api',
            'django.contrib.contenttypes',
            'django.contrib.auth',
            'ninja',
        ],
        DATABASES={'default': {'ENGINE': 'django.db.backends.sqlite3', 'NAME': ':memory:'}},
        SECRET_KEY='dummy',
        MEDIA_ROOT='/tmp/media',
    )
    django.setup()

# Import the view
from api import views

class TestCleanup(unittest.TestCase):
    def setUp(self):
        # Patch Client in views
        self.patcher_client = patch('api.views.Client')
        self.mock_client_class = self.patcher_client.start()
        self.mock_client_instance = self.mock_client_class.return_value

        # Patch ImageDetection
        self.patcher_models = patch('api.views.ImageDetection')
        self.mock_image_detection = self.patcher_models.start()

        # Patch os.listdir and os.remove
        self.patcher_listdir = patch('os.listdir')
        self.mock_listdir = self.patcher_listdir.start()

        self.patcher_remove = patch('os.remove')
        self.mock_remove = self.patcher_remove.start()

        self.patcher_isfile = patch('os.path.isfile')
        self.mock_isfile = self.patcher_isfile.start()

        self.patcher_exists = patch('os.path.exists')
        self.mock_exists = self.patcher_exists.start()

        # Patch Image.open
        self.patcher_image = patch('api.views.Image')
        self.mock_image = self.patcher_image.start()

        # Patch handle_file
        self.patcher_handle_file = patch('api.views.handle_file')
        self.mock_handle_file = self.patcher_handle_file.start()
        # Make handle_file return its input or something harmless
        self.mock_handle_file.side_effect = lambda x: x

        # Mock PIL Image object
        self.mock_img_obj = MagicMock()
        self.mock_image.open.return_value = self.mock_img_obj
        self.mock_img_obj.convert.return_value = self.mock_img_obj
        self.mock_img_obj.resize.return_value = self.mock_img_obj

        # Set up default mocks
        self.mock_listdir.return_value = ['file1.jpg', 'file2.jpg', 'target.jpg']
        self.mock_isfile.return_value = True

        # Mock ImageDetection.objects.create return value
        mock_obj = MagicMock()
        # Mock URL behavior.
        # views.py: image_path = pathlib.Path(prompt_obj.image.url[1:])
        # If url is "/media/images/target.jpg", then [1:] is "media/images/target.jpg"
        mock_obj.image.url = "/media/images/target.jpg"
        self.mock_image_detection.objects.create.return_value = mock_obj

        # Mock Client predict result
        mock_result = [
            {"value": []}, # index 0
            {}, # index 1
            {"width": 100, "height": 100} # index 2
        ]
        self.mock_client_instance.predict.return_value = mock_result

    def tearDown(self):
        patch.stopall()

    def test_cleanup_optimized(self):
        # This test checks for the optimization.
        # It should verify that os.listdir is NOT called
        # And os.remove IS called for specific files.

        request = MagicMock()
        prompt = "detect cat"
        image = MagicMock()
        image.__str__.return_value = "target.jpg"
        width = 100
        height = 100

        views.detect(request, prompt, image, width, height)

        # Verify os.listdir was NOT called
        self.mock_listdir.assert_not_called()

        # Verify os.remove was called for specific files
        media_path = os.getcwd() + '/media/images/'
        # img_path is a Path object, resized_img_path is a string
        import pathlib
        expected_calls = [
            unittest.mock.call(pathlib.Path(os.getcwd()) / 'media/images/target.jpg'), # img_path
            unittest.mock.call(os.path.join(media_path, 'resized_target.jpg')) # resized_img_path
        ]

        # We check that these calls were made.
        # Note: img_path construction:
        # image_path = pathlib.Path(prompt_obj.image.url[1:]) -> media/images/target.jpg
        # img_path = pathlib.Path(cwd , image_path) -> cwd/media/images/target.jpg
        # So it matches expected_calls[0]

        # We can relax exact order if needed, but we expect both calls.
        self.assertEqual(self.mock_remove.call_count, 2)
        self.mock_remove.assert_has_calls(expected_calls, any_order=True)

if __name__ == '__main__':
    unittest.main()
