# Task contract

Every generated app task or migration slice uses this shape:

```text
Goal:
Reference tuiweave example:
Components:
Application-owned state:
Files allowed:
Behavior to preserve:
Acceptance tests:
Verification commands:
Out of scope:
```

Fill every field before implementation. `Files allowed` must identify shared
model/render/update files explicitly. Acceptance tests name observable behavior
and required snapshot or scenario checkpoints. A reusable library need is
recorded as a candidate tuiweave API gap, not implemented as an app-local fork.
