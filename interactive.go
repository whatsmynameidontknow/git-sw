package main

// TUI implements UserInterface using the existing promptui-based interactive interface.
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
