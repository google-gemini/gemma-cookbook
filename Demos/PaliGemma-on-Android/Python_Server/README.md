# Python Django Ninja based API server

**Note:** The Python server is productionized and hosted on [Render](https://paligemma.onrender.com). If you wish to run this locally, follow the steps below:

## Steps to setup
1. Create your virtual environment using `virtualenv` or `venv`.
2. Install requirements using `pip install -r requirements.txt`.
3. Go to directory /apiserver where manage.py file resides.
4. Execute the command `python manage.py migrate` for setting up the database.
5. Execute the command `mkdir -p apiserver/media/images`.
6. Go to [settings.py](apiserver/apiserver/settings.py) and set the following:
  - SECRET_KEY = 'YOUR_DJANGO_SECRET_KEY'
  - ALLOWED HOSTS = ['localhost']
    
7. Execute the command `python manage.py runserver` to start the server.
8. Go to the URL [http://127.0.0.1:8000/api/docs#/default/api_views_detect](http://127.0.0.1:8000/api/docs#/default/api_views_detect) to test the API.
