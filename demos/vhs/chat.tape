Output docs/demos/chat.gif
Output demos/vhs/outputs/chat.mp4
Require go
Require git
Set Shell "bash"
Set Width 1100
Set Height 700
Set FontSize 17
Set Theme "TokyoNight"
Set TypingSpeed 45ms

Hide
Type "cd $(git rev-parse --show-toplevel)"
Enter
Type "go run ./examples/chat"
Enter
Wait+Screen /mock session/
Show
Sleep 700ms
Type "how do layouts work?"
Enter
Wait+Screen /Good question!/
Sleep 3000ms
