Output docs/demos/consumer.gif
Output demos/vhs/outputs/consumer.mp4
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
Type "go run ./examples/consumer"
Enter
Wait+Screen /Components/
Show
Sleep 700ms
Down
Down
Sleep 1200ms
