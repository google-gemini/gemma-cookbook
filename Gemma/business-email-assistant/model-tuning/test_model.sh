#!/bin/bash

# activate virtual environment
source ../venv/bin/activate

# test the model: set last parameter to "True" to use generated weights
python3 -c 'import main; main.generate_from_model("I`d like a vanilla cake with blueberry filling and a unicorn.", True)'