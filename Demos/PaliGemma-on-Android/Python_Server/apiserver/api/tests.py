from unittest.mock import patch, MagicMock
from django.test import TestCase, Client
from django.core.files.uploadedfile import SimpleUploadedFile
from api.views import api

class PerformanceTests(TestCase):
    @patch('api.views.Client')
    @patch('api.views.handle_file')
    @patch('api.views.ImageDetection.objects.create')
    @patch('api.views.Image.open')
    @patch('api.views.os.getcwd')
    @patch('api.views.os.remove')
    @patch('api.views.os.listdir')
    @patch('api.views.os.path.isfile')
    def test_client_initialization_count(self, mock_isfile, mock_listdir, mock_remove, mock_getcwd, mock_open, mock_create, mock_handle_file, mock_client):
        # Setup mocks
        mock_client_instance = MagicMock()
        mock_client.return_value = mock_client_instance
        mock_client_instance.predict.return_value = (
             {"value": [{"token": "<loc0><loc0><loc0><loc0>", "class_or_confidence": "cat"}]},
             None,
             {"width": 100, "height": 100}
        )

        mock_create.return_value.image.url = "/media/test.jpg"
        mock_getcwd.return_value = "/tmp"
        mock_open.return_value.convert.return_value.resize.return_value.save = MagicMock()
        mock_listdir.return_value = []

        # mock handle_file to return something
        mock_handle_file.return_value = "mock_file_handle"

        client = Client()

        # We need to simulate the file upload
        image = SimpleUploadedFile("test.jpg", b"file_content", content_type="image/jpeg")

        # First request
        response = client.post(
            "/api/detect",
            data={"prompt": "detect cat", "width": 100, "height": 100, "image": image},
            format='multipart' # Standard Django Client format
        )

        # Second request
        image.seek(0)
        response = client.post(
            "/api/detect",
            data={"prompt": "detect cat", "width": 100, "height": 100, "image": image},
            format='multipart'
        )

        # Verify Client was called only once (cached)
        self.assertEqual(mock_client.call_count, 1)
