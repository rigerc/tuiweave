Output docs/demos/form.gif
Output demos/vhs/outputs/form.mp4
Require go
Require git
Set Shell "bash"
Set Width 1000
Set Height 650
Set FontSize 18
Set Theme "TokyoNight"
Set TypingSpeed 45ms

Hide
Type "cd $(git rev-parse --show-toplevel)"
Enter
Type "go run ./examples/form"
Enter
Wait+Screen /Preferences/
Show
Sleep 700ms
Type " Lovelace"
Tab
Type "api"
Tab
Down
Tab
Space
Tab
Enter
Wait+Screen /saved/
Sleep 1400ms
