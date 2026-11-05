# Contributing

Issues and pull requests are welcome. For a new playbook, skill, template, or
audit rule, open an issue first so its audience, scope, and alternatives can be
agreed before implementation. Small documentation and fixture fixes may go
directly to a pull request.

Keep application behavior in the consuming repository; this kit owns portable
workflows, templates, installation guidance, and read-only audit tooling. Pin
the exact tuiweave version or commit used by an example and keep local
`replace` directives out of committed application instructions.

Changes must pass `./tests/tuiweave-audit_test.sh` and
`./tests/sync-tuiweave-reference_test.sh`, preserve fixture immutability, and
explain any intentional contract or workflow change. Pull requests should be
focused and include verification details. No CLA or DCO sign-off is required.
Participation is governed by [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
