/* eslint-disable @typescript-eslint/naming-convention */
/**
 * Copyright 2024 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import * as vscode from 'vscode';

const http = require('http');

export async function gemmaServiceRequest(requestText: String) {

  // Get Gemma Service Host and port number from local user configuration
  const gemmaHost = vscode.workspace.getConfiguration().get<string>('gemma.service.host');
  if (!gemmaHost) {
    vscode.window.showErrorMessage('Gemma Service Host not configured. Check Pipet Code Agent settings.');
    return;
  }

  // Construct the http POST request
  const jsonRequest = {
    text: requestText,
  };

  const options = {
    hostname: gemmaHost, // web service hostname
    port: 8000, // web service port number
    path: '/gemma_request/', // web service endpoint path
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Content-Length': Buffer.byteLength(JSON.stringify(jsonRequest)),
    },
  };

  return new Promise((resolve, reject) => {
    const req = http.request(options, (res: any) => {
      let data = '';
      res.on('data', (chunk: any) => data += chunk);
      res.on('end', () => {
        try {
          const jsonResponse = JSON.parse(data);
          resolve(jsonResponse.text); // Resolve the Promise
        } catch (error) {
          reject(error);
        }
      });
    });

    req.on('error', reject); // Reject the Promise on request error
    req.write(JSON.stringify(jsonRequest));
    req.end();
  });
}