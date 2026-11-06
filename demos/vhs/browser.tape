Output docs/demos/browser.gif
Output demos/vhs/outputs/browser.mp4
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
Type "go run ./examples/browser"
Enter
Wait+Screen /BROWSER/
Show
Sleep 700ms
Type "/"
Sleep 200ms
Type "main"
Wait+Screen /main.go/
Sleep 600ms
Tab
Down
Sleep 1200ms
