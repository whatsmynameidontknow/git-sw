Refined Implementation Plan for `git-sw`

I have analyzed the `git-sw` codebase and designed a plan to introduce the mandatory `--no-tui` flag while prioritizing modularity and minimal changes to existing files.

### 1. Strategy Overview

We will introduce a `UserInterface` interface to abstract the interaction layer.
- **Current Behavior (TUI)**: Wrapped in a `TUI` struct that calls existing `render.go` functions.
- **New Behavior (NoTUI)**: Implemented in a `NoTUI` struct that reads from flags/arguments and strictly validates inputs without prompting.

The `main.go` will decide which implementation to use based on the `--no-tui` flag.

### 2. New Files

We will create **3 new files** to encapsulate the changes:

#### A. `types_ext.go`
Defines the interface for user interaction to decouple logic from the UI.
```go
package main

// UserInterface defines the methods required for user interaction.
type UserInterface interface {
    // CreateProfile gathers data to create a new profile.
    CreateProfile() (Profile, error)
    // SelectProfile allows the user to select a profile from the list.
    SelectProfile(profiles []Profile) (Profile, error)
    // ListProfiles displays the list of profiles.
    ListProfiles(profiles []Profile) error
    // ConfirmDelete asks for confirmation before deletion.
    ConfirmDelete() bool
    // EditProfile handles the editing of a profile (e.g., opening editor).
    EditProfile(path string) error
}
```

#### B. `interactive.go`
Wraps the existing `render.go` functions and `util.go` helpers into the `TUI` struct.
```go
package main

type TUI struct{}

func (t *TUI) CreateProfile() (Profile, error) {
    return displayCreateForm()
}

func (t *TUI) SelectProfile(profiles []Profile) (Profile, error) {
    return displayProfileSelector(profiles)
}

func (t *TUI) ListProfiles(profiles []Profile) error {
    return displayProfileList(profiles)
}

func (t *TUI) ConfirmDelete() bool {
    return displayDeleteConfirmation()
}

func (t *TUI) EditProfile(path string) error {
    return openTextEditor(path)
}
```

#### C. `non_interactive.go`
Implements the automated logic. It strictly checks flags/arguments and returns errors if required data is missing.
```go
package main

import (
    "errors"
    "flag"
    "github.com/thansetan/git-sw/pkg/gitconfig"
)

type NoTUI struct{}

func (n *NoTUI) CreateProfile() (Profile, error) {
    // Validate flags (profileName, gitName, gitEmail, etc.)
    // If missing, return error.
    // Construct and return Profile.
}

func (n *NoTUI) SelectProfile(profiles []Profile) (Profile, error) {
    // Check `profileName` flag or `flag.Arg(1)`.
    // Find in `profiles`.
    // Return Profile or error if not found.
}

func (n *NoTUI) ListProfiles(profiles []Profile) error {
    // Print simple text list to stdout.
}

func (n *NoTUI) ConfirmDelete() bool {
    // Always true (automated mode implies intent), 
    // or arguably require a --confirm flag (but prompt said "no-tui" means automated).
    // Safest: return true, as the user invoked the command explicitly.
    return true
}

func (n *NoTUI) EditProfile(path string) error {
    return errors.New("interactive edit is not supported in --no-tui mode")
}
```

### 3. Modifications to Existing Files

#### A. `flag.go`
Add global flags for the non-interactive inputs.
```go
var (
    noTui       bool
    profileName string
    gitName     string
    gitEmail    string
    signingKey  string
    gpgFormatStr string
    gpgProgram  string
)

func parseFlag() {
    // ... existing ...
    flag.BoolVar(&noTui, "no-tui", false, "Disable TUI prompts (automated mode)")
    flag.StringVar(&profileName, "profile", "", "Profile name (for create/use/delete)")
    flag.StringVar(&gitName, "name", "", "Git user name (for create)")
    flag.StringVar(&gitEmail, "email", "", "Git user email (for create)")
    // ... other flags ...
}
```

#### B. `command.go`
Update the logic to use the `ui` interface instead of direct calls.
- Define `var ui UserInterface` at package level (or inject it).
- Replace `displayCreateForm()` with `ui.CreateProfile()`.
- Replace `displayProfileSelector(profiles)` with `ui.SelectProfile(profiles)`.
- Replace `openTextEditor` calls with `ui.EditProfile`.

#### C. `main.go`
Initialize the `ui` variable based on the flag.
```go
func main() {
    // ... parseFlag ...
    
    // Initialize UI
    if noTui {
        ui = &NoTUI{}
    } else {
        ui = &TUI{}
    }

    // ... rest of main ...
}
```

### 4. Hook Points Summary

| File | Hook Point | Action |
| :--- | :--- | :--- |
| `flag.go` | `parseFlag()` | Register `--no-tui` and configuration flags (`--name`, `--email`, etc.). |
| `main.go` | After `parseFlag()` | Instantiate global `ui` variable as `&NoTUI{}` or `&TUI{}`. |
| `command.go` | `CREATE` command | Call `ui.CreateProfile()` instead of `displayCreateForm()`. |
| `command.go` | `USE`, `DELETE` | Call `ui.SelectProfile()` instead of `displayProfileSelector()`. |
| `command.go` | `EDIT` | Call `ui.EditProfile()` instead of `openTextEditor()`. |

### 5. Behavior Specification

- **`--no-tui` Present**:
    - **Interactive prompts**: Disabled.
    - **Missing Info**: Exit with error (e.g., "missing required flag: --email").
    - **`create`**: Requires flags `--profile`, `--name`, `--email` (and optional signing key flags).
    - **`use`/`delete`**: Requires `--profile` flag or profile name as 2nd argument.
    - **`edit`**: Returns error (cannot open vim).
    - **`list`**: prints plain text list.
- **`--no-tui` Absent**:
    - Uses existing TUI behavior (`promptui`).
    - Flags (if provided) are currently ignored by TUI (keeping changes minimal), but TUI continues to work as before.

### 6. Conventional Commits
The final implementation will be committed with:
`feat: add --no-tui flag for automated mode and modularize UI logic`
