Output docs/demos/ops.gif
Output demos/vhs/outputs/ops.mp4
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
Type "go run ./examples/ops"
Enter
Wait+Screen /OPS/
Show
Sleep 700ms
Ctrl+P
Wait+Screen /Command palette/
Type "restart"
Wait+Screen /Restart service/
Enter
Wait+Screen /Restart service\?/
Sleep 1400ms
