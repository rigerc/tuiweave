# Install for other agents

Use the host's Agent Skills discovery directory and copy the complete
`skills/tuiweave-app-builder` directory there. Prefer the cross-host project path
`.agents/skills/tuiweave-app-builder` when supported. The host must preserve the
standard `SKILL.md` frontmatter and make the adjacent `references/` files
available to the agent.

If a host has no Agent Skills support, provide `SKILL.md` as workflow context
and preserve its relative reference paths; do not fork the workflow text.
