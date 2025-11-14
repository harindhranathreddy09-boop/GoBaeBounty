package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// AuthFile represents the authorization file structure
type AuthFile struct {
	Target        string    `json:"target"`
	AuthToken     string    `json:"auth_token"`
	Signature     string    `json:"signature"`
	IssuedBy      string    `json:"issued_by"`
	IssuedAt      time.Time `json:"issued_at"`
	ExpiresAt     time.Time `json:"expires_at"`
	ProgramURL    string    `json:"program_url"`
	Scope         []string  `json:"scope"`
	Version       int       `json:"version,omitempty"`
}

const (
	// SecretKeyEnv is the environment variable name for secret key
	SecretKeyEnv = "GOBEEB_AUTH_SECRET"
	// AuthVersion is the expected version of the auth file format
	AuthVersion = 1
)

// ValidateAuthFile opens and verifies the authorization file given a target
func ValidateAuthFile(path, target string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read auth file: %w", err)
	}

	var af AuthFile
	if err := json.Unmarshal(data, &af); err != nil {
		return fmt.Errorf("invalid auth file format: %w", err)
	}

	if af.Version != 0 && af.Version != AuthVersion {
		return fmt.Errorf("unsupported auth file version %d; expected %d", af.Version, AuthVersion)
	}

	if af.Target != target {
		return fmt.Errorf("auth file target '%s' does not match scan target '%s'", af.Target, target)
	}

	// Check if target in scope
	inScope := false
	for _, scope := range af.Scope {
		if scope == target || scope == "*" {
			inScope = true
			break
		}
	}
	if !inScope {
		return fmt.Errorf("target '%s' is not in authorized scope %v", target, af.Scope)
	}

	if time.Now().After(af.ExpiresAt) {
		return fmt.Errorf("authorization expired at %v", af.ExpiresAt)
	}

	secret := os.Getenv(SecretKeyEnv)
	if secret == "" {
		// In production, this must be set securely!
		secret = "demo-secret-for-testing"
	}

	if err := verifySignature(&af, secret); err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}

// verifySignature checks the HMAC-SHA256 signature matches the computed signature
func verifySignature(af *AuthFile, secret string) error {
	payload := fmt.Sprintf("%s:%s:%s:%s:%s",
		af.Target,
		af.AuthToken,
		af.IssuedBy,
		af.IssuedAt.Format(time.RFC3339),
		af.ExpiresAt.Format(time.RFC3339),
	)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(af.Signature)) {
		return fmt.Errorf("invalid signature")
	}
	return nil
}
