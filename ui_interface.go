package main

// UserInterface defines the methods required for user interaction.
// This abstraction allows for both interactive (TUI) and non-interactive (flag-based) modes.
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

// ui is the global user interface implementation.
// It is set in main() based on the --no-tui flag.
var ui UserInterface
