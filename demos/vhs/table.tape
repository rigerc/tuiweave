Output docs/demos/table.gif
Output demos/vhs/outputs/table.mp4
Require go
Require git
Set Shell "bash"
Set Width 1200
Set Height 700
Set FontSize 18
Set Theme "TokyoNight"
Set TypingSpeed 45ms

Hide
Type "cd $(git rev-parse --show-toplevel)"
Enter
Type "go run ./examples/table"
Enter
Wait+Screen /git-status/
Show
Sleep 700ms
Down
Down
Sleep 700ms
Tab
Down
Sleep 1400ms
