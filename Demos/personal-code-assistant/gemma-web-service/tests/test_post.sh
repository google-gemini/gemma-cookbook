#!/bin/bash

curl -X POST -H "Content-Type: application/json" \
     -d '{"text": "roses are red"}' \
     http://localhost:8000/gemma_request/

echo