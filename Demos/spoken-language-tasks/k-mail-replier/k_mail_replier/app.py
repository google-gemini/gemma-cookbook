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
from flask import Flask, render_template, request
#from k_mail_replier.models.gemini import create_message_processor
from k_mail_replier.models.gemma import create_message_processor

app = Flask(__name__, static_url_path='/static', static_folder='static')
customer_request = None
model_processor = create_message_processor() # initialize model

@app.route('/', methods=['GET', 'POST'])
def index():
    global customer_request
    global model_processor
    """Set up web interface and handle POST input."""
    # First run behavior: load a test email
    if customer_request is None:
        customer_request = get_test_email()
        return render_template('index.html', request=customer_request)

    # Process email data
    if request.method == 'POST':
        prompt = get_prompt()
        customer_request = request.form['request']
        prompt += customer_request
        result = model_processor(prompt)
        # re-render page with data:
        return render_template('index.html', request=customer_request, result=result)

    return render_template('index.html')

if __name__ == '__main__':
    app.run(debug=True)

def get_prompt():
    """Write a polite reply to this email thanks the sender for the request and saying that we will reply with more detail soon:"""    
    return "발신자에게 요청에 대한 감사를 전하고, 곧 자세한 내용을 알려드리겠다고 정중하게 답장해 주세요:\n"

def get_test_email():
    try:
        with open('data/email-001-ko.txt', 'r') as file:
            email_content = file.read()
    except FileNotFoundError:
        email_content = "Error: File not found!"
    return email_content