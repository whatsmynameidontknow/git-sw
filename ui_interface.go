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

// AppState holds shared application state and dependencies.
// This replaces global mutable state with an explicit dependency injection pattern.
type AppState struct {
	UI UserInterface
}

// NewAppState creates a new AppState with the appropriate UI implementation
// based on the noTUI flag.
func NewAppState(noTUI bool) *AppState {
	var ui UserInterface
	if noTUI {
		ui = &NoTUI{}
	} else {
		ui = &TUI{}
	}
	return &AppState{UI: ui}
}
