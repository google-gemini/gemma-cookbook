#!/bin/bash

# update apt repository
sudo apt update

# install software
sudo apt install git python3-venv

#ll create a virtual environment for the project
python3 -m venv venv

# to activate virtual environment
source venv/bin/activate
# to deactivate:
# deactivate

# check version of CUDA drivers
nvcc --version