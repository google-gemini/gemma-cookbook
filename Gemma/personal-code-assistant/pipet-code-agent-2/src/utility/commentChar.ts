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

export function getLineCommentCharacter(editor: vscode.TextEditor) {
  const document = editor.document;
  const languageId = document.languageId;

  switch (languageId) {
    case 'javascript':
    case 'typescript':
    case 'java':
    case 'cpp':
      return '//';
    case 'python':
      return '#';
    case 'html':
    case 'xml':
      return '<!-- '; // Note: This starts full XML comment, adjust if needed
    default:
      return '//';; // Return a default comment character
  }
}