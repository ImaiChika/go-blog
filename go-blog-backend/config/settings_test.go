package config

import "testing"

func TestGetConfigFromEnv(t *testing.T) {
	t.Setenv("GO_BLOG_DB_DSN", "custom-dsn")
	t.Setenv("GO_BLOG_REDIS_ADDR", "127.0.0.1:6380")
	t.Setenv("GO_BLOG_REDIS_PASSWORD", "secret")
	t.Setenv("GO_BLOG_REDIS_DB", "2")
	t.Setenv("GO_BLOG_JWT_SECRET", "custom-secret")

	if got := GetDBDSN(); got != "custom-dsn" {
		t.Fatalf("expected custom dsn, got %q", got)
	}
	if got := GetRedisAddr(); got != "127.0.0.1:6380" {
		t.Fatalf("expected custom redis addr, got %q", got)
	}
	if got := GetRedisPassword(); got != "secret" {
		t.Fatalf("expected custom redis password, got %q", got)
	}
	if got := GetRedisDB(); got != 2 {
		t.Fatalf("expected custom redis db, got %d", got)
	}
	if got := GetJWTSecret(); got != "custom-secret" {
		t.Fatalf("expected custom jwt secret, got %q", got)
	}
}

func TestGetConfigUsesDefaultWhenEnvMissingOrInvalid(t *testing.T) {
	t.Setenv("GO_BLOG_DB_DSN", "")
	t.Setenv("GO_BLOG_REDIS_ADDR", "")
	t.Setenv("GO_BLOG_REDIS_PASSWORD", "")
	t.Setenv("GO_BLOG_REDIS_DB", "invalid")
	t.Setenv("GO_BLOG_JWT_SECRET", "")

	if got := GetDBDSN(); got != defaultDBDSN {
		t.Fatalf("expected default dsn, got %q", got)
	}
	if got := GetRedisAddr(); got != defaultRedisAddr {
		t.Fatalf("expected default redis addr, got %q", got)
	}
	if got := GetRedisPassword(); got != defaultRedisPassword {
		t.Fatalf("expected default redis password, got %q", got)
	}
	if got := GetRedisDB(); got != defaultRedisDB {
		t.Fatalf("expected default redis db, got %d", got)
	}
	if got := GetJWTSecret(); got != defaultJWTSecret {
		t.Fatalf("expected default jwt secret, got %q", got)
	}
}
