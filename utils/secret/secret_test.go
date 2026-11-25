package secret

import (
	"strings"
	"testing"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	if err := Init("unit-test-encryption-key"); err != nil {
		t.Fatalf("Init: %v", err)
	}

	plain := "p@ssw0rd-中文-\n-----BEGIN KEY-----"
	enc, err := Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if !IsEncrypted(enc) {
		t.Fatalf("expected encrypted prefix, got %q", enc)
	}
	if strings.Contains(enc, plain) {
		t.Fatal("ciphertext contains plaintext")
	}

	// 幂等：重复加密不应嵌套
	enc2, err := Encrypt(enc)
	if err != nil {
		t.Fatalf("Encrypt idempotent: %v", err)
	}
	if enc2 != enc {
		t.Fatal("Encrypt should be idempotent")
	}

	dec, err := Decrypt(enc)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if dec != plain {
		t.Fatalf("round trip mismatch: got %q", dec)
	}
}

func TestDecryptLegacyPlaintext(t *testing.T) {
	if err := Init("unit-test-encryption-key"); err != nil {
		t.Fatalf("Init: %v", err)
	}

	dec, err := Decrypt("legacy-plain-password")
	if err != nil {
		t.Fatalf("Decrypt legacy: %v", err)
	}
	if dec != "legacy-plain-password" {
		t.Fatalf("unexpected legacy value: %q", dec)
	}
}

func TestEmptyValue(t *testing.T) {
	if err := Init("unit-test-encryption-key"); err != nil {
		t.Fatalf("Init: %v", err)
	}

	enc, err := Encrypt("")
	if err != nil || enc != "" {
		t.Fatalf("empty encrypt: %q, %v", enc, err)
	}
	dec, err := Decrypt("")
	if err != nil || dec != "" {
		t.Fatalf("empty decrypt: %q, %v", dec, err)
	}
}
