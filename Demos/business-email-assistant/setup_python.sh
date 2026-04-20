#!/bin/bash
# install Python components

# activate virtual environment
source venv/bin/activate

pip install -r requirements.txt

# note: record python installation as follows
# pip freeze > requirements.txt
