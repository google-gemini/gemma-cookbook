#!/bin/bash

# activate virtual environment
source venv/bin/activate

# package installations
pip install -r requirements.txt

# note: update requirements list with:
# pip freeze > requirements.txt