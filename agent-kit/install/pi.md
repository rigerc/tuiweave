# Install for Pi Coding Agent

Pi supports the standard `SKILL.md` format and discovers skills from locations
including `.agents/skills` and `.pi/skills`. Copy or symlink the unchanged
package into either project location:

```sh
mkdir -p .pi/skills
cp -R /path/to/tuiweave-agent-kit/skills/tuiweave-app-builder .pi/skills/
```

Alternatively use `.agents/skills/tuiweave-app-builder` when the repository is
shared with other Agent Skills-compatible hosts. Do not maintain a Pi-specific
fork. See https://pi.dev/docs/latest/skills for current discovery details.
