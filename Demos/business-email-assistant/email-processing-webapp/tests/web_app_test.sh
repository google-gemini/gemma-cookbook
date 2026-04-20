#!/bin/bash
# test generation through web front end

curl -X POST -d "request=I'd like a vanilla cake with blueberries on top." \
     http://127.0.0.1:5000/

echo