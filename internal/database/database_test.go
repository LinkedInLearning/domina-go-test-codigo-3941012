package database

import (
	_ "embed"
	"testing"
)

func TestNew(t *testing.T) {
	srv := MustNewWithDatabase()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := MustNewWithDatabase()

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	srv := MustNewWithDatabase()

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}
