Output docs/demos/palette.gif
Output demos/vhs/outputs/palette.mp4
Require go
Require git
Set Shell "bash"
Set Width 1000
Set Height 600
Set FontSize 18
Set Theme "TokyoNight"
Set TypingSpeed 45ms

Hide
Type "cd $(git rev-parse --show-toplevel)"
Enter
Type "go run ./examples/palette"
Enter
Wait+Screen /Commands/
Show
Sleep 700ms
Type "format"
Wait+Screen /Format document/
Sleep 700ms
Enter
Wait+Screen /selected/
Sleep 1300ms
