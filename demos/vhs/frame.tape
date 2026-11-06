Output docs/demos/frame.gif
Output demos/vhs/outputs/frame.mp4
Require go
Require git
Set Shell "bash"
Set Width 1200
Set Height 650
Set FontSize 18
Set Theme "TokyoNight"
Set TypingSpeed 45ms

Hide
Type "cd $(git rev-parse --show-toplevel)"
Enter
Type "go run ./examples/frame"
Enter
Wait+Screen /service overview/
Show
Sleep 700ms
Tab
Down
Tab
Right
Tab
Space
Sleep 1400ms
