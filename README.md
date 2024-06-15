# genie - Your CLI assistant 🧞‍♂️

## Submission for MLH Fellowship 2024 🌟

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
- get information about anything related to tech directly from the CLI

<p align="center">
<a href="https://genie.harshalranjhani.in">
<img src="https://cdn.hashnode.com/res/hashnode/image/upload/v1716281685360/_2uaTNTl5.webp?auto=format" alt="genie-logo"/>
</a>
</p>

<!-- [![Generic badge](https://img.shields.io/badge/view-demo-blue?style=for-the-badge&label=View%20Demo%20Video)](https://youtu.be/OKKK1GOnlIU)  -->

## Features

### Welcome

![Welcome](https://cdn.hashnode.com/res/hashnode/image/upload/v1716285507443/fmlq1ls4j.png?auto=format)

1. The `init` command

   - start the genie CLI with the `init` command
   - store all your API keys when prompted
   - ![image](https://cdn.hashnode.com/res/hashnode/image/upload/v1716282272366/Lj-xUwtH1.png?auto=format)

2. The `do` command

   - used to run commands
   - use: `genie do <prompt>`
   - ![image](https://cdn.hashnode.com/res/hashnode/image/upload/v1716282661700/X74dam0Jl.png?auto=format)

   - flags
     - `--safe`: run the command in a safe mode, prevents direct execution of dangerous commands
     - ![image](https://cdn.hashnode.com/res/hashnode/image/upload/v1716282782898/ox37E_4wD.png?auto=format)

3. The `generate` command

   - used to generate images from text
   - this currently only works on local and not on the build, i am working on a fix
   - use: `genie generate <prompt>`
   - ![image](https://cdn.hashnode.com/res/hashnode/image/upload/v1716283662684/fzy4kQcyO.png?auto=format)
   - this was the generated image:
   - ![image](https://cdn.hashnode.com/res/hashnode/image/upload/v1716283748566/kC1lCPYMC.jpg?auto=format)

4. The `tell` command

   - used to get information about anything related to tech
   - use: `genie tell <prompt>`
   - ![image](https://cdn.hashnode.com/res/hashnode/image/upload/v1716283848962/irLWaz5mA.png?auto=format)

5. The `music` command

   - used to generate music, uses the music-gen model from the replicate API
   - use: `genie music <prompt>`
   - ![image](https://cdn.hashnode.com/res/hashnode/image/upload/v1716284145216/fx7b1R1nC.png?auto=format)

   - flags
     - `--d`: specify the duration of the audio
     - `--logs`: show logs while generating the audio

6. The `switch` command

   - used to switch between engines
   - currently supports gemini and gpt engines
   - supports the `do`, `tell` and `greet` commands
   - ![image](https://cdn.hashnode.com/res/hashnode/image/upload/v1717667356512/pN3kKPCio.png?auto=format)

7. The `chat` command

   - maintain a conversation with the genie
   - currently supports only the gemini engine
   - ![image](https://cdn.hashnode.com/res/hashnode/image/upload/v1718473330436/wR_UMtMCe.png?auto=format)

   - flags
     - `--safe`: have a conversation in safe mode


## Points to remember while testing the app

1. Do not forget to run the `init` command before running any other command

2. The `generate` command currently only works on local and not on the build, i am working on a fix.

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
