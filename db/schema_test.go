package db

import (
	"testing"
)

func TestCreateSchema(t *testing.T) {
	_, err := openSqlite(":memory:")
	if err != nil {
		t.Fatalf("couldn't create schema: %s", err)
	}
}
