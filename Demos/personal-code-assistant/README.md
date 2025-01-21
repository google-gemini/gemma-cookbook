# Personal AI Code Assistant with Gemma | Build with Google AI

This code project let's you create your own, personal AI coding assistant with Gemma by wrapping the model in a web service and creating a Visual Studio Code extention to communicate with it.

This project contains 2 sub-projects:

- **gemma-web-service** - A Gemma 2 2B model wrapped in a simple web service written with Python and the FastAPI library.
- **pipet-code-agent-2** - Visual Studio Code extension written in Node.js that connects to the Gemma service to handle code generation and other requests.
