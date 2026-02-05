package main

import (
	"errors"
	"fmt"
	"net/mail"
	"os"
	"path/filepath"
	"strings"

	"github.com/thansetan/git-sw/pkg/gitconfig"
	"golang.org/x/crypto/ssh"
)

// NoTUI implements UserInterface for non-interactive (automated/scripted) usage.
type NoTUI struct{}

var (
	ErrMissingProfile    = errors.New("missing required flag: --profile")
	ErrMissingName       = errors.New("missing required flag: --name")
	ErrMissingEmail      = errors.New("missing required flag: --email")
	ErrProfileNotFound   = errors.New("profile not found")
	ErrEditNoTUI         = errors.New("interactive edit is not supported in --no-tui mode")
	ErrDeleteNoConfirm   = errors.New("delete in --no-tui mode requires --yes flag for safety")
	ErrInvalidKeyFormat  = errors.New("invalid key format: must be 'openpgp', 'ssh', or 'x509'")
	ErrMissingSigningKey = errors.New("--signing-key is required when --key-format is specified")
)

func (n *NoTUI) CreateProfile() (Profile, error) {
	var profile Profile
	profile.Config = gitconfig.New()

	// Validate required flags
	if profileFlag == "" {
		return Profile{}, ErrMissingProfile
	}
	if nameFlag == "" {
		return Profile{}, ErrMissingName
	}
	if emailFlag == "" {
		return Profile{}, ErrMissingEmail
	}

	// Check for duplicate profile
	profileMap := make(map[string]struct{})
	for i := range profiles {
		profileMap[strings.ToLower(profiles[i].Name)] = struct{}{}
	}
	if _, ok := profileMap[strings.ToLower(profileFlag)]; ok {
		return Profile{}, ErrDuplicateProfile
	}

	// Validate email format
	_, err := mail.ParseAddress(emailFlag)
	if err != nil {
		return Profile{}, ErrInvalidEmail
	}

	// Validate gitconfig values
	if err := gitconfig.ValidateValue(nameFlag); err != nil {
		return Profile{}, fmt.Errorf("invalid name: %w", err)
	}
	if err := gitconfig.ValidateValue(emailFlag); err != nil {
		return Profile{}, fmt.Errorf("invalid email: %w", err)
	}

	// Set profile name
	profile.Name = profileFlag

	// Set user.name
	if err := profile.Config.Set("user.name", nameFlag); err != nil {
		return Profile{}, err
	}

	// Set user.email
	if err := profile.Config.Set("user.email", emailFlag); err != nil {
		return Profile{}, err
	}

	// Handle signing key configuration
	if keyFormatFlag != "" || signingKeyFlag != "" {
		if signingKeyFlag == "" {
			return Profile{}, ErrMissingSigningKey
		}

		// Determine key format: use provided format or default to openpgp
		keyFormat := GPGFormat(strings.ToLower(keyFormatFlag))
		if keyFormat == "" {
			keyFormat = OPENPGP
		}

		// Validate key format
		validFormat := false
		for _, f := range gpgFormat {
			if f == keyFormat {
				validFormat = true
				break
			}
		}
		if !validFormat {
			return Profile{}, ErrInvalidKeyFormat
		}

		// Validate signing key based on format
		if err := gitconfig.ValidateValue(signingKeyFlag); err != nil {
			return Profile{}, fmt.Errorf("invalid signing key: %w", err)
		}

		if keyFormat == SSH {
			// Validate SSH public key
			if filepath.Ext(signingKeyFlag) != ".pub" {
				return Profile{}, ErrInvalidPublicKeyExt
			}
			content, err := os.ReadFile(signingKeyFlag)
			if err != nil {
				return Profile{}, fmt.Errorf("cannot read SSH key: %w", err)
			}
			_, _, _, _, err = ssh.ParseAuthorizedKey(content)
			if err != nil {
				return Profile{}, fmt.Errorf("invalid SSH key: %w", err)
			}
		}

		// Set GPG format
		if err := profile.Config.Set("gpg.format", string(keyFormat)); err != nil {
			return Profile{}, err
		}

		// Set signing key
		if err := profile.Config.Set("user.signingKey", signingKeyFlag); err != nil {
			return Profile{}, err
		}

		// Enable commit signing
		if err := profile.Config.Set("commit.gpgsign", "true"); err != nil {
			return Profile{}, err
		}

		// Set GPG program for openpgp format
		if keyFormat == OPENPGP {
			gpgProg := gpgProgramFlag
			if gpgProg == "" {
				gpgProg = "gpg" // default
			}
			if err := profile.Config.Set("gpg.program", gpgProg); err != nil {
				return Profile{}, err
			}
		}
	}

	return profile, nil
}

func (n *NoTUI) SelectProfile(profiles []Profile) (Profile, error) {
	profileName := profileFlag
	if profileName == "" {
		return Profile{}, ErrMissingProfile
	}

	// Find the profile
	for _, p := range profiles {
		if strings.EqualFold(p.Name, profileName) {
			return p, nil
		}
	}

	return Profile{}, fmt.Errorf("%w: %s", ErrProfileNotFound, profileName)
}

func (n *NoTUI) ListProfiles(profiles []Profile) error {
	for _, p := range profiles {
		status := ""
		if p.IsActive {
			status = " (active)"
		}
		fmt.Printf("%s%s\n", p.Name, status)
	}
	return nil
}

func (n *NoTUI) ConfirmDelete() bool {
	// In non-interactive mode, require --yes flag for safety
	if !yesFlag {
		fmt.Fprintln(os.Stderr, ErrDeleteNoConfirm)
		return false
	}
	return true
}

func (n *NoTUI) EditProfile(path string) error {
	return ErrEditNoTUI
}
