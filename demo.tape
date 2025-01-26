# VHS documentation
#
# Output:
#   Output <path>.gif               Create a GIF output at the given <path>
#   Output <path>.mp4               Create an MP4 output at the given <path>
#   Output <path>.webm              Create a WebM output at the given <path>
#
# Require:
#   Require <string>                Ensure a program is on the $PATH to proceed
#
# Settings:
#   Set FontSize <number>           Set the font size of the terminal
#   Set FontFamily <string>         Set the font family of the terminal
#   Set Height <number>             Set the height of the terminal
#   Set Width <number>              Set the width of the terminal
#   Set LetterSpacing <float>       Set the font letter spacing (tracking)
#   Set LineHeight <float>          Set the font line height
#   Set LoopOffset <float>%         Set the starting frame offset for the GIF loop
#   Set Theme <json|string>         Set the theme of the terminal
#   Set Padding <number>            Set the padding of the terminal
#   Set Framerate <number>          Set the framerate of the recording
#   Set PlaybackSpeed <float>       Set the playback speed of the recording
#   Set MarginFill <file|#000000>   Set the file or color the margin will be filled with.
#   Set Margin <number>             Set the size of the margin. Has no effect if MarginFill isn't set.
#   Set BorderRadius <number>       Set terminal border radius, in pixels.
#   Set WindowBar <string>          Set window bar type. (one of: Rings, RingsRight, Colorful, ColorfulRight)
#   Set WindowBarSize <number>      Set window bar size, in pixels. Default is 40.
#   Set TypingSpeed <time>          Set the typing speed of the terminal. Default is 50ms.
#
# Sleep:
#   Sleep <time>                    Sleep for a set amount of <time> in seconds
#
# Type:
#   Type[@<time>] "<characters>"    Type <characters> into the terminal with a
#                                   <time> delay between each character
#
# Keys:
#   Escape[@<time>] [number]        Press the Escape key
#   Backspace[@<time>] [number]     Press the Backspace key
#   Delete[@<time>] [number]        Press the Delete key
#   Insert[@<time>] [number]        Press the Insert key
#   Down[@<time>] [number]          Press the Down key
#   Enter[@<time>] [number]         Press the Enter key
#   Space[@<time>] [number]         Press the Space key
#   Tab[@<time>] [number]           Press the Tab key
#   Left[@<time>] [number]          Press the Left Arrow key
#   Right[@<time>] [number]         Press the Right Arrow key
#   Up[@<time>] [number]            Press the Up Arrow key
#   Down[@<time>] [number]          Press the Down Arrow key
#   PageUp[@<time>] [number]        Press the Page Up key
#   PageDown[@<time>] [number]      Press the Page Down key
#   Ctrl+<key>                      Press the Control key + <key> (e.g. Ctrl+C)
#
# Display:
#   Hide                            Hide the subsequent commands from the output
#   Show                            Show the subsequent commands in the output

Output demo.gif

Require genie

Set Shell "bash"
Set FontSize 30
Set Width 1200
Set Height 800
Set Padding 20
Set Theme "Catppuccin Mocha"
# Set WindowBar Colorful
Set BorderRadius 10

# Title
Type "# Welcome to Genie - Your AI-Powered CLI Assistant! 🧞‍♂️" Sleep 500ms Enter
Sleep 3s
Sleep 2s
Ctrl+L

# Greet Command
Type "genie greet" Sleep 500ms Enter
Sleep 8s  # Increased wait time for API response
Sleep 2s  # Added pause after clearing
Ctrl+L

# Do Command
Type "genie do 'list all files in the current directory'" Sleep 500ms Enter
Sleep 12s  # Increased for command execution
Sleep 4s
Ctrl+L

# Tell Command
Type "genie tell 'what is the current codebase about and how do i use it?' --include-dir" Sleep 500ms Enter
Sleep 25s  # Longer wait for detailed response
Sleep 3s
Ctrl+L

# Document Command
Type "genie document --file main.go" Sleep 500ms Enter
Sleep 10s
Sleep 2s
Ctrl+L

# Switch Model
Type "genie switch --model deepseek-reasoner" Sleep 500ms Enter
Sleep 6s
Sleep 2s
Ctrl+L

# Chat Command
Type "genie chat" Sleep 500ms Enter
Sleep 5s
Type "How can you help me with coding tasks?" Sleep 500ms Enter
Sleep 25s
Sleep 4s
Type "/history" Sleep 500ms Enter
Sleep 8s
Type "exit" Sleep 500ms Enter
Sleep 3s
Sleep 2s
Ctrl+L

# Scrape Command
Type "genie scrape 'https://crawler-test.com' -e h3" Sleep 500ms Enter
Sleep 12s
Sleep 4s
Ctrl+L

# Engine Command
Type "genie engine" Sleep 500ms Enter
Sleep 6s
Sleep 2s
Ctrl+L

# Switch Command
Type "genie switch" Sleep 500ms Enter
Sleep 6s
Sleep 2s
Ctrl+L

