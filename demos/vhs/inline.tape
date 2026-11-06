Output docs/demos/inline.gif
Output demos/vhs/outputs/inline.mp4
Require go
Require git
Set Shell "bash"
Set Width 1000
Set Height 780
Set FontSize 18
Set Theme "TokyoNight"
Set TypingSpeed 45ms

Hide
Type "cd $(git rev-parse --show-toplevel)"
Enter
Type "go run ./examples/inline"
Enter
Wait+Screen /Inline fields/
Show
Sleep 700ms
Type "Ada"
Tab
Type "ada@example.com"
Tab
Type "hunter2"
Tab
Down
Tab
Down
Space
Enter
Tab
Space
Down
Space
Tab
Type "8000"
Tab
Space
Tab
Down
Sleep 1600ms
