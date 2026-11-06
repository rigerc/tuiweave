Output docs/demos/demo.gif
Output demos/vhs/outputs/demo.mp4
Require go
Require git
Set Shell "bash"
Set Width 1100
Set Height 650
Set FontSize 18
Set Theme "TokyoNight"
Set TypingSpeed 45ms

Hide
Type "cd $(git rev-parse --show-toplevel)"
Enter
Type "go run ./examples/demo"
Enter
Wait+Screen /welcome/
Show
Sleep 700ms
Tab
Tab
Type "new item"
Enter
Wait+Screen /new item/
Sleep 800ms
Ctrl+D
Wait+Screen /Quit demo\?/
Sleep 1200ms
Escape
Sleep 500ms
