# genie - Your CLI assistant 🧞‍♂️

[![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/harshalranjhani/genie?logo=github&style=for-the-badge)](https://github.com/harshalranjhani)
[![GitHub last commit](https://img.shields.io/github/last-commit/harshalranjhani/genie?style=for-the-badge&logo=git)](https://github.com/harshalranjhani)
[![GitHub stars](https://img.shields.io/github/stars/harshalranjhani/genie?style=for-the-badge)](https://github.com/harshalranjhani/genie)
[![My stars](https://img.shields.io/github/stars/harshalranjhani?affiliations=OWNER%2CCOLLABORATOR&style=for-the-badge&label=My%20stars)](https://github.com/harshalranjhani/genie)
[![GitHub forks](https://img.shields.io/github/forks/harshalranjhani/genie?style=for-the-badge&logo=git)](https://github.com/harshalranjhani/network)
[![Code size](https://img.shields.io/github/languages/code-size/harshalranjhani/genie?style=for-the-badge)](https://github.com/harshalranjhani)
[![Languages](https://img.shields.io/github/languages/count/harshalranjhani/genie?style=for-the-badge)](https://github.com/harshalranjhani)
[![Top](https://img.shields.io/github/languages/top/harshalranjhani/genie?style=for-the-badge&label=Top%20Languages)](https://github.com/harshalranjhani)
[![Issues](https://img.shields.io/github/issues/harshalranjhani/genie?style=for-the-badge&label=Issues)](https://github.com/harshalranjhani)
[![Watchers](https://img.shields.io/github/watchers/harshalranjhani/genie?label=Watch&style=for-the-badge)](https://github.com/harshalranjhani/)

Your personal assistant for the CLI that helps you:

- run commands
- generate images
- generate music
- summarize comments (supports multiple languages)
- document code
- get information about anything related to tech directly from the CLI
- maintain a chat session with the genie for advanced context understanding

<p align="center">
<a href="https://genie.harshalranjhani.in">
<img src="https://cdn.hashnode.com/res/hashnode/image/upload/v1716281685360/_2uaTNTl5.webp?auto=format" alt="genie-logo"/>
</a>
</p>

<!-- [![Generic badge](https://img.shields.io/badge/view-demo-blue?style=for-the-badge&label=View%20Demo%20Video)](https://youtu.be/OKKK1GOnlIU)  -->

## Features

### Welcome

![Welcome](https://cdn.hashnode.com/res/hashnode/image/upload/v1718473561830/j4aVeAVll.png?auto=format)

## Overview

The Genie CLI provides a set of commands to help developers use AI to automate tasks, generate documentation, and improve their workflow. This documentation covers the available commands, their usage, and examples to help you get started with the Genie CLI.

[!<img src="https://files.edgestore.dev/gre1kolpt9w3vnwd/publicImages/_public/7d74e1a7-32fc-4266-b19a-badb997da2ba.png" alt="Support Genie!" width="300" height="70">](https://rzp.io/rzp/JUFHvoxw)


You can support the development of Genie by contributing to the project or making a donation. You can also make donations using the `support` command in the Genie CLI.

Once you have made a donation, you can use the `verify` command to verify your donation and get access to additional features and benefits later on.

If you've donated, a big thank you from genie!

## Configurations

`genie init`

The `init` command is used to store your API keys, session IDs, and ignore list file paths. By running this command, you can configure the Genie CLI to access external services and customize its behavior.

### Ignore List File Paths

When configuring the Genie CLI, you can specify ignore patterns to exclude certain files or directories from processing. This can be useful for skipping test files, build directories, or other irrelevant content.

The ignore list is a file that contains a list of files and directories that you want to exclude from your project. This is highly important to reduce the token count and improve the performance of the analysis. Create an ignore list file anywhere on your system and provide the path to it during the `genie init` command.

This is a text file with the name `ignorelist.txt` that contains a list of files and directories to ignore. For example:

```text
node_modules
dist
build
```

### API Keys

The Genie CLI requires API keys to access external services for text-to-image generation, text-to-music generation, and other features. You can obtain API keys from the respective service providers and store them securely using the `genie init` command.

## Commands

### 1. `do`

The `do` command allows you to execute commands that you might not remember. Leveraging the power of AI, Genie can help you run commands without having to remember them. Just type in what you want to do, and Genie will take care of the rest.

**Usage:**

```bash
genie do "prompt"
```

**Flags:**

- `--safe`: Run the command in safe mode, which ensures that the command is executed only if it is safe to do so.

**Description:**

- **AI-Powered Command Execution**: Execute commands using natural language prompts.
- **Helpful for Beginners**: Ideal for beginners who may not be familiar with command-line syntax.

### 2. `image`

The `image` command is used to generate images from text. This can be useful for generating any kind of image from text, such as diagrams, charts, or illustrations.

**Usage:**

```bash
genie image "prompt"
```

**Description:**

- **Text-to-Image Generation**: Converts text prompts into images.
- **Versatile Use Cases**: Can be used for creating diagrams, charts, illustrations, and more.

Note: This is only supported for the GPT engine.

### 3. `music`

The `music` command allows you to generate music based on a text prompt. This can be useful for creating background music, soundtracks, or other audio content.

**Usage:**

```bash
genie music "prompt"
```

**Flags:**

- `--d`: Specify the duration of the generated music. (Default: 8 seconds, max: 15 seconds)
- `--logs`: Display logs during music generation.

**Description:**

- **Text-to-Music Generation**: Converts text prompts into music.
- **Customizable Duration**: Set the duration of the generated music.
- **Real-Time Logs**: Option to display logs during music generation.

### 4. `tell`

The `tell` command is used to generate text responses to questions or prompts. This can be useful for generating responses to queries, providing information, or creating conversational content about the CLI.

**Usage:**

```bash
genie tell "prompt"
```

**Flags:**

- `--include-dir`: Include the current directory snapshot in the request for better context.

**Description:**

- **Text Response Generation**: Generates text responses to prompts.
- **Conversational AI**: Provides information and answers questions in a conversational format.

### 5. `summarize`

The `summarize` command generates a structured markdown summary of comments within your project files. This is useful for creating documentation and reviewing code.

For this command to work, you need to have comments in your code that are marked as headings and subheadings. Genie will automatically detect these comments and generate a markdown summary based on them.

Example of comments in code:

In python:

```python
# genie:heading: This is a heading
# genie:subheading: This is a subheading
```

or in javascript:

```javascript
// genie:heading: This is a heading
// genie:subheading: This is a subheading
```

Make sure to match the exact format for the comments to be detected correctly. The format is `genie:heading:` for headings and `genie:subheading:` for subheadings. Remember to add a space after the colon and before the text. Also add a space after the comment marker (`#` or `//`) and before the `genie` keyword.

This command can be used in relation to the `document` command to generate summaries of the codebase.

**Usage:**

```bash
genie summarize
```

**Flags:**

- `--email`: Send the generated markdown summary as a PDF via email.
- `--support`: Lists the supported languages for comment detection.
- `--filename`: Specify the filename for the generated markdown summary.

**Description:**

- **Automatic Detection**: Scans project files for comments marked as headings and subheadings.
- **Multi-Language Support**: Supports multiple programming languages by recognizing various comment markers.
- **Email Integration**: Option to send the generated markdown summary as a PDF via email.
- **Ignore Patterns**: Customizable ignore patterns to exclude specific files or directories.

### 6. `document`

The `document` command generates documentation for your file and integrates it with genie comments. This can then be useful to get easier summaries of your code and to understand the codebase better.

**Usage:**

```bash
genie document
```

**Flags:**

- `--file`: Specify the file to generate documentation for. (Required flag)

**Description:**

- **Automatic Documentation Generation**: Generates documentation for your codebase.
- **Integrates with Genie Comments**: Utilizes genie comments to create structured documentation.
- **Customizable Output**: Specify the file to generate documentation for.

### 7. `chat`

Note: The `chat` command is only available for supporters of the Genie project.

The `chat` command opens a chat interface where you can interact with Genie in a conversational manner. This can be useful for asking questions, getting help, or exploring the capabilities of the Genie CLI.

**Usage:**

```bash
genie chat
```

**Flags:**

- `--safe`: Run the command in safe mode, which ensures that the conversation is safe and appropriate.

**Description:**

- **Conversational Interface**: Interact with Genie in a chat-like environment.
- **AI-Powered Responses**: Get answers to questions and prompts in a conversational format.

## Conclusion

The Genie CLI is a powerful tool that helps streamline your development workflow by automating tasks, generating documentation, and more. By using the available commands, you can improve your productivity and maintain a consistent project structure.

For more information, visit the [Genie About Page](https://genie.harshalranjhani.in/about).


## Points to remember while testing the app

1. Do not forget to run the `init` command before running any other command

2. The `generate` command with Gemini Engine currently only works on local and not on the build, i am working on a fix. It works on the build directly with the GPT engine.

3. Make sure you provide a valid **ignorelist.txt** file path when you `init` the app. This file is like `.gitignore` and contains the files that you want to ignore when passing prompts to the model. This is done to make sure to stay within the model token limits.

4. The `music` command uses the music-gen model from the replicate API. Make sure you have the correct API key stored in the keyring.

5. Run the `docs` command to open the documentation in the browser.

## Instructions

1. Clone the repository

```bash
git clone https://github.com/harshalranjhani/genie.git
```

2. Enter the directory

```bash
cd genie
```

3. Install the dependencies

```bash
go mod tidy
```

4. Get going

```bash
go run main.go
```

## Useful Links

- [Genie website](https://genie.harshalranjhani.in)

## Need help?

Feel free to contact me on [LinkedIn](https://www.linkedin.com/in/harshal-ranjhani/)

## [![Twitter](https://img.shields.io/badge/Twitter-blue.svg?logo=twitter&logoColor=white)](https://twitter.com/ranjhaniharshal) [![Dev.to](https://img.shields.io/badge/Dev.to-black.svg?logo=dev.to&logoColor=white)](https://dev.to/harshalranjhani) [![LinkedIn](https://img.shields.io/badge/LinkedIn-blue.svg?logo=linkedin&logoColor=white)](https://www.linkedin.com/in/harshal-ranjhani/) [![Hashnode](https://img.shields.io/badge/Hashnode-black.svg?logo=hashnode&logoColor=white)](https://hashnode.com/@harshalranjhani)

```javascript
if (youEnjoyed) {
  starThisRepository();
}
```

---
