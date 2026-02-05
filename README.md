<h1 align="center">Git Sw(itch)</h1>

<p align="center">
  <img src="https://github.com/user-attachments/assets/e08c8a1d-c00e-4aa8-b31f-338dea710bf7" alt="git-sw demo"/>
</p>

A CLI tool to switch between multiple Git profiles/configs. A package to parse and create a .gitconfig file is also available [here](https://github.com/thansetan/git-sw/tree/main/pkg/gitconfig).

## Installation

Binary releases are available on the [releases page](https://github.com/thansetan/git-sw/releases).

**Go**
```sh
go install github.com/thansetan/git-sw@latest
```

## Usage
```text
usage: git-sw [options] command
```

### Available Commands
| Command | Description |
| :--- | :--- |
| `use` | Select a profile to use. |
| `create` | Create a new profile. |
| `edit` | Edit an existing profile in text editor. |
| `delete` | Delete an existing profile. |
| `list` | List all available profiles. |

### Available Options
| Option | Description |
| :--- | :--- |
| `-g` | Run the command globally (can only be used with 'use', 'edit', and 'delete'). |
| `--no-tui` | Disable interactive TUI prompts (Automated/Agent mode). |
| `--profile <name>` | Specify profile name (for create/use/delete). |
| `--name <name>` | Specify Git user name (for create). |
| `--email <email>` | Specify Git user email (for create). |
| `--signing-key <key>` | Specify GPG key ID or SSH key path. |
| `--key-format <format>` | Key format: `openpgp`, `ssh`, or `x509`. |
| `--gpg-program <prog>` | Path to GPG program (default: `gpg`). |
| `--yes` | Auto-confirm destructive operations (for delete). |

## Agent-Friendly Mode (Non-Interactive)

Use the `--no-tui` flag for automation or when using the tool from an AI agent. In this mode, prompts are disabled, and the tool will error out if required information is missing.

**Example: Create a profile**
```bash
git-sw --no-tui --profile work --name "User Name" --email "user@example.com" create
```

**Example: Switch profile**
```bash
git-sw --no-tui --profile work use
```

**Example: Delete profile**
```bash
git-sw --no-tui --profile work --yes delete
```

## Tips
- Run `git-sw list` to see current profiles and the active one.
- Use `-g` to apply a profile to your global `~/.gitconfig`.
- Use `--no-tui` in scripts or CI/CD pipelines.
