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
import { GoogleGenerativeAI } from '@google/generative-ai';
import { GenerateContentResult } from '@google/generative-ai';

export async function geminiAPIRequest(requestText: any) {

  const modelName = vscode.workspace.getConfiguration().get<string>('google.gemini.textModel', 'models/gemini-1.5-pro-latest');

  // Get API Key from local user configuration
  const apiKey = vscode.workspace.getConfiguration().get<string>('google.gemini.apiKey');
  if (!apiKey) {
    vscode.window.showErrorMessage('API key not configured. Check your settings.');
    return;
  }

  const genai = new GoogleGenerativeAI(apiKey);
  const model = genai.getGenerativeModel({ model: modelName });

  const result: Promise<GenerateContentResult> = model.generateContent(requestText);
  console.debug('Gemini result: ' + result); // FOR DEBUG ONLY 
  const response = (await result).response;

  return response.text;
}