#!/bin/bash

# activate virtual environment
source venv/bin/activate

pip install -r requirements.txt

# note: record python installation as follows
# pip freeze > requirements.txt

# ------------------------- 
# key package installations (for manual installation)
# pip install python-dotenv
# pip install Flask bleach # Only needed for web application

# (Optional) for Gemini API testing
# pip install google-generativeai

# Gemma specific software
# pip install keras-nlp
# pip install "jax[cuda12]" # install jax for CUDA 12 drivers

# pip install datasets # only needed for tuning
# pip install matplotlib # only needed for tuning evaluation
