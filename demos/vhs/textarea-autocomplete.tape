Output docs/demos/textarea-autocomplete.gif
Output demos/vhs/outputs/textarea-autocomplete.mp4
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
Type "go run ./examples/textarea-autocomplete"
Enter
Wait+Screen /Textarea completion/
Show
Sleep 700ms
Type "run git"
Wait+Screen /git status/
Sleep 700ms
Down
Enter
Wait+Screen /run git log --oneline/
Sleep 1400ms
