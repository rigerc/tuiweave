# Greenfield mode

1. Inspect repository conventions, Go/tuiweave versions, and tests.
2. Route the app to the closest tuiweave example and matching recipes.
3. Produce a component map, layout rectangles, size modes, focus order,
   state-ownership map, inspection/action map, and command/result flow.
4. Implement the smallest working shell before product behavior.
5. Add snapshots for visual states, paired role-aware goldens where the app
   assigns theme roles, and named scenarios with explicit command delivery.
6. Run build, vet, tests, intentional golden review, and the strict audit.

Use explicit `WindowSizeMsg` layout and `SetSize`; compose one root `tea.View`.
Keep all application side effects and policy outside components.
Build semantic trees and route qualified actions in the app, validating the
current visible tree immediately before invocation.
