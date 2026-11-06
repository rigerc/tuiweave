Output docs/demos/autocomplete.gif
Output demos/vhs/outputs/autocomplete.mp4
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
Type "go run ./examples/autocomplete"
Enter
Wait+Screen /Command completion/
Show
Sleep 700ms
Type "git"
Wait+Screen /git checkout/
Sleep 700ms
Down
Sleep 400ms
Enter
Wait+Screen /accepted/
Sleep 1400ms
