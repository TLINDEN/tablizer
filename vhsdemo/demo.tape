# -*-sh-*-

Output demo.gif
Set FontSize 20
Set Width 1200
Set Height 1000
Set Theme { "name": "Whimsy", "black": "#535178", "red": "#ef6487", "green": "#5eca89", "yellow": "#fdd877", "blue": "#65aef7", "magenta": "#aa7ff0", "cyan": "#43c1be", "white": "#ffffff", "brightBlack": "#535178", "brightRed": "#ef6487", "brightGreen": "#5eca89", "brightYellow": "#fdd877", "brightBlue": "#65aef7", "brightMagenta": "#aa7ff0", "brightCyan": "#43c1be", "brightWhite": "#ffffff", "background": "#29283b", "foreground": "#b3b0d6", "selection": "#3d3c58", "cursor": "#b3b0d6" }
Set WindowBar Colorful
Set BorderRadius 10
Set Shell zsh
Set FontFamily "IBM Plex Mono"
Set CursorBlink false
Set PlaybackSpeed 1
Set TypingSpeed .05

# initialize
Hide
Type `PROMPT=''`
Enter
Type "setopt interactivecomments"
Enter
Type "autoload -U colors && colors"
Enter
Type `PS1="%{$fg[magenta]%}demo> %{$reset_color%}"`
Enter
Type "clear"
Enter
Show

Type "# Our input data"
Enter
Sleep 1s
Type "cat input | head -10"
Enter
Sleep 2s

Enter
Type "# Filter over all rows"
Enter
Sleep 1s
Type "tablizer Central < input"
Enter
Sleep 2s

Enter
Type "# Filter over all rows case insensitive"
Enter
Sleep 1s
Type "tablizer '/penc/i' < input"
Enter
Sleep 2s

Enter
Type "# Filter over specific column"
Enter
Sleep 1s
Type "tablizer -Fcost=4.99 < input"
Enter
Sleep 2s

Enter
Type "# Filter by regex on specific column"
Enter
Sleep 1s
Type "tablizer -Funits=Pen. < input"
Enter
Sleep 2s

Enter
Type "# Output as markdown"
Enter
Sleep 1s
Type "tablizer -Funits=Pen. -M < input"
Enter
Sleep 2s

Enter
Type "# Output as CSV"
Enter
Sleep 1s
Type "tablizer -Funits=Pen. -C < input"
Enter
Sleep 2s

Enter
Type "# Output as shell evaluable"
Enter
Sleep 1s
Type "tablizer -Funits=Pen. -S < input"
Enter
Sleep 2s
Type "bat eval.sh"
Enter
Sleep 2s
Type "tablizer -Funits=Pen. -S < input | ./eval.sh"
Enter
Sleep 2s

Enter
Type "# Reduce columns"
Enter
Sleep 1s
Type "tablizer -Funits=Pen. -c region,customer,units,count < input"
Enter
Sleep 2s

Enter
Type "# Sort by COUNT column numerically "
Enter
Sleep 1s
Type "tablizer -Funits=Pen. -c region,customer,units,count -kcount -i < input"
Enter
Sleep 2s

Enter
Type "# Do further filtering interactively"
Enter
Sleep 1s
Type "tablizer -Funits=Pen. -c region,customer,units,count -I -O < input"
Enter
Sleep 2s
Type "?"
Sleep 2s
Type "/"
Sleep 2s
Type "J"
Sleep 1s
Type "o"
Sleep 1s
Type "n"
Sleep 1s
Type "e"
Sleep 1s
Type "s"
Sleep 1s
Enter
Sleep 2s
Tab
Sleep 1s
Tab
Sleep 1s
Tab
Sleep 1s
Tab
Type "n"
Sleep 2s
Space
Sleep 1s
Down
Sleep 1s
Down
Sleep 1s
Space
Sleep 2s
Type "q"

Sleep 10s
