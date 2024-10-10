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
import { getLineCommentCharacter } from './utility/commentChar';
// import { geminiAPIRequest } from './models/gemini-api';
import { gemmaServiceRequest } from './models/gemma-service';

// Provide instructions for the AI model
const PROMPT_INSTRUCTIONS = `
A good code comment describes the intent behind the code without
repeating information that's obvious from the code itself. Good comments
describe "why", explain any "magic" values and non-obvious behaviour.

Write a comment that explains the following code. Do NOT include a separate
explanation of the code, just the comment to added:
`;

export async function generateComment() {
    vscode.window.showInformationMessage('Generating code comment...');

    // Get selected text
    const editor = vscode.window.activeTextEditor;
    if (!editor) {
        console.debug('Abandon: no open text editor.');
        return;
    }

    const selection = editor.selection;
    const selectedCode = editor.document.getText(selection);

    // Build the full prompt
    const promptText = `${PROMPT_INSTRUCTIONS}${selectedCode}`;

    // Send the request and insert the response
    try {
        // const response: any = await geminiAPIRequest(promptText); // alternative
        const response: any = await gemmaServiceRequest(promptText);
        const responseText = response.toString();
        console.log(responseText);

        // Insert before selection.
        editor.edit((editBuilder) => {
            // Find the right comment character using the editor object  
            const commentChar = getLineCommentCharacter(editor);
            const commentPrefix = commentChar + ' ';

            // Copy the indent from the first line of the selection.
            const trimmed = selectedCode.trimStart();
            const padding = selectedCode.substring(0, selectedCode.length - trimmed.length);

            let pyComment = responseText.split('\n').map((l: string) => `${padding}${commentPrefix}${l}`).join('\n');
            if (pyComment.search(/\n$/) === -1) {
                // Add a final newline if necessary.
                pyComment += "\n";
            }
            let commentIntro = padding + commentPrefix + "Code comment: (generated)\n";
            editBuilder.insert(selection.start, commentIntro);
            editBuilder.insert(selection.start, pyComment);
        });
    } catch (error) {
        console.error('Error:', error);
    }
}
