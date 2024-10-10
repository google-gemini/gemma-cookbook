#!/bin/bash

# activate virtual environment
source venv/bin/activate

cd gemma_service/
# to allow more than localhost access, add "--host 0.0.0.0":
#uvicorn gemma_service_main:app --host 0.0.0.0 --reload

python3 gemma_service_main.py
