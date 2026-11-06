Output docs/demos/statusbar.gif
Output demos/vhs/outputs/statusbar.mp4
Require go
Require git
Set Shell "bash"
Set Width 1000
Set Height 400
Set FontSize 18
Set Theme "TokyoNight"
Set TypingSpeed 45ms

Hide
Type "cd $(git rev-parse --show-toplevel)"
Enter
Type "go run ./examples/statusbar"
Enter
Wait+Screen /tuiweave/
Show
Sleep 700ms
Type "t"
Wait+Screen /Light/
Sleep 900ms
Type "t"
Wait+Screen /Dracula/
Sleep 1200ms
