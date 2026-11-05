# Install for Codex

Install the same standard skill package globally or per project:

```sh
mkdir -p "$CODEX_HOME/skills"
cp -R skills/tuiweave-app-builder "$CODEX_HOME/skills/tuiweave-app-builder"
```

For a project-local installation, copy or symlink the package to
`.agents/skills/tuiweave-app-builder`. Keep `SKILL.md` and its `references/`
together. Restart or refresh skill discovery, then request greenfield,
migration, or review mode explicitly when useful.
