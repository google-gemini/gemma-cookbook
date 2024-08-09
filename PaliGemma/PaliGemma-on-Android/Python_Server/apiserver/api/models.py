from django.db import models

# Create your models here.

class ImageDetection(models.Model):
    prompt = models.CharField(max_length=500)
    image = models.ImageField(upload_to="images/")