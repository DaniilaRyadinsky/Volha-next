package postgres

import (
	"net/url"
	"strings"
	"testing"
)

func TestURLPreservesPostgresCredentials(t *testing.T) {
	t.Parallel()

	passwords := []string{
		"simple-password",
		"p@ss:word",
		"a/b?c#d",
		"100%secret",
	}
	for _, password := range passwords {
		t.Run(password, func(t *testing.T) {
			t.Parallel()

			dsn := URL("postgres.example", 5432, "volha/db?x", "volha@user", password, "disable")

			parsed, err := url.Parse(dsn)
			if err != nil {
				t.Fatalf("parse dsn: %v", err)
			}
			if got := parsed.User.Username(); got != "volha@user" {
				t.Fatalf("username = %q, want %q", got, "volha@user")
			}
			gotPassword, ok := parsed.User.Password()
			if !ok || gotPassword != password {
				t.Fatalf("password = %q, %t; want %q, true", gotPassword, ok, password)
			}
			if got := strings.TrimPrefix(parsed.Path, "/"); got != "volha/db?x" {
				t.Fatalf("database = %q, want %q", got, "volha/db?x")
			}
			if got := parsed.Query().Get("sslmode"); got != "disable" {
				t.Fatalf("sslmode = %q, want disable", got)
			}
		})
	}
}

func TestURLHandlesIPv6Host(t *testing.T) {
	t.Parallel()

	dsn := URL("::1", 5432, "volha", "user", "simple-password", "disable")
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse dsn: %v", err)
	}
	if got := parsed.Host; got != "[::1]:5432" {
		t.Fatalf("host = %q, want [::1]:5432", got)
	}
}
