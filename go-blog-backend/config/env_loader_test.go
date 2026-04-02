package config

import (
	"os"
	"path/filepath"
	"testing"
)

func unsetEnvForTest(t *testing.T, key string) {
	t.Helper()

	original, exists := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("failed to unset env %s: %v", key, err)
	}

	t.Cleanup(func() {
		var err error
		if exists {
			err = os.Setenv(key, original)
		} else {
			err = os.Unsetenv(key)
		}
		if err != nil {
			t.Fatalf("failed to restore env %s: %v", key, err)
		}
	})
}

func TestLoadEnvFileLoadsValuesFromDotEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "GO_BLOG_DB_DSN=test-dsn\nGO_BLOG_REDIS_DB=3\nGO_BLOG_JWT_SECRET=\"quoted-secret\"\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}

	unsetEnvForTest(t, "GO_BLOG_DB_DSN")
	unsetEnvForTest(t, "GO_BLOG_REDIS_DB")
	unsetEnvForTest(t, "GO_BLOG_JWT_SECRET")

	if err := LoadEnvFile(path); err != nil {
		t.Fatalf("expected load env success, got %v", err)
	}

	if got := os.Getenv("GO_BLOG_DB_DSN"); got != "test-dsn" {
		t.Fatalf("expected test-dsn, got %q", got)
	}
	if got := os.Getenv("GO_BLOG_REDIS_DB"); got != "3" {
		t.Fatalf("expected redis db 3, got %q", got)
	}
	if got := os.Getenv("GO_BLOG_JWT_SECRET"); got != "quoted-secret" {
		t.Fatalf("expected trimmed quoted secret, got %q", got)
	}
}

func TestLoadEnvFileDoesNotOverrideExistingEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "GO_BLOG_JWT_SECRET=from-dotenv\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}

	t.Setenv("GO_BLOG_JWT_SECRET", "from-system-env")

	if err := LoadEnvFile(path); err != nil {
		t.Fatalf("expected load env success, got %v", err)
	}

	if got := os.Getenv("GO_BLOG_JWT_SECRET"); got != "from-system-env" {
		t.Fatalf("expected existing env to win, got %q", got)
	}
}

func TestLoadEnvFileRejectsInvalidLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "INVALID_LINE\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}

	err := LoadEnvFile(path)
	if err == nil {
		t.Fatal("expected invalid .env line to return error")
	}
}
