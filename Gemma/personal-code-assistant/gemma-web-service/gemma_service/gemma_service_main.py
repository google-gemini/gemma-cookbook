#
# Copyright 2024 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

from fastapi import FastAPI
import uvicorn
from pydantic import BaseModel
from gemma_model import create_model_instance
from gemma_model import get_model_id

app = FastAPI()
gemma_model = create_model_instance() # initialize model

class Request(BaseModel):
    text: str

class Response(BaseModel):
    text: str

@app.post("/gemma_request/")
async def process_text(request: Request):
    """
    Processes the input text and returns a modified version.
    """
    response_text = gemma_model(request.text)
    response = Response(text=response_text)
    return response

@app.get("/")
async def root():
    return "Gemma server: OK"

@app.get("/info")
async def info():
    return "Gemma service is using: " + get_model_id()

# Run the server 
if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
