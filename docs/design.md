# Product Design Notes

This document captures product design directives and user-facing behavior guidelines that should shape future DevBuddy changes. Keep entries concise, merge duplicates when they appear, and resolve ambiguity during occasional consolidation passes.

## Update Notices

- DevBuddy update checks should run only from `bud up`, not from shell hooks, shell init, shell completion, `bud upgrade`, or unrelated commands.
- The update check should start when `bud up` starts and run in the background. It must not add blocking network latency to the command.
- DevBuddy should print the update notice near command exit only if the background check has completed by then. It should not wait for the check to finish.
- A check attempt should update the cache timestamp before the network request. If the request is slow, times out, or fails, DevBuddy should not immediately retry on the next `bud up`.
- The notice should use the same upgrade plan as `bud upgrade`: Homebrew installs use the Homebrew command, and non-Homebrew installs use the documented install-script command.
