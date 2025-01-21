#!/bin/bash

# activate virtual environment
source ../venv/bin/activate

# delete any previously generated weights
rm -f weights/*.h5

python3 main.py