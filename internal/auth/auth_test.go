package auth

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestValidateAuthFile_Valid(t *testing.T) {
	af := &AuthFile{
		Target:    "example.com",
		AuthToken: "testtoken123",
		IssuedBy:  "Tester",
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(72 * time.Hour),
		ProgramURL: "https://bugbounty.example.com",
		Scope:     []string{"example.com"},
		Version:   AuthVersion,
	}

	secret := "test-secret"
	payload := af.Target + ":" + af.AuthToken + ":" + af.IssuedBy + ":" + af.IssuedAt.Format(time.RFC3339) + ":" + af.ExpiresAt.Format(time.RFC3339)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	af.Signature = hex.EncodeToString(mac.Sum(nil))

	f, err := os.CreateTemp("", "authfile_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	if err := json.NewEncoder(f).Encode(af); err != nil {
		t.Fatal(err)
	}
	f.Close()

	// Set environment variable for secret key
	os.Setenv(SecretKeyEnv, secret)
	defer os.Unsetenv(SecretKeyEnv)

	err = ValidateAuthFile(f.Name(), af.Target)
	if err != nil {
		t.Errorf("Expected valid authfile, got error: %v", err)
	}
}

func TestValidateAuthFile_Expired(t *testing.T) {
	af := &AuthFile{
		Target:    "example.com",
		AuthToken: "expiredtoken",
		IssuedBy:  "Tester",
		IssuedAt:  time.Now().Add(-48 * time.Hour),
		ExpiresAt: time.Now().Add(-24 * time.Hour),
		ProgramURL: "https://bugbounty.example.com",
		Scope:     []string{"example.com"},
		Version:   AuthVersion,
		Signature: "invalid",
	}

	f, err := os.CreateTemp("", "authfile_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	if err := json.NewEncoder(f).Encode(af); err != nil {
		t.Fatal(err)
	}
	f.Close()

	os.Setenv(SecretKeyEnv, "some-secret")
	defer os.Unsetenv(SecretKeyEnv)

	err = ValidateAuthFile(f.Name(), af.Target)
	if err == nil {
		t.Error("Expected error for expired auth file but got none")
	}
}
