// history_test.go
package main

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *HistoryModel {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open the database: %v", err)
	}

	model := &HistoryModel{db: db}
	if err := model.Migrate(); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return model
}

func TestHistoryModelRecord(t *testing.T) {
	model := setupTestDB(t)

	err := model.Record("client1", "some data")
	if err != nil {
		t.Errorf("inserting record produced an error: %v", err)
	}

	var count int
	err = model.db.QueryRow("SELECT COUNT(*) FROM history WHERE client_id = ?", "client1").Scan(&count)
	if err != nil {
		t.Errorf("counting records produced an error: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 record for client1, got %d", count)
	}
}

func TestHistoryModelAll(t *testing.T) {
	model := setupTestDB(t)

	for i := 0; i < 5; i++ {
		err := model.Record("client1", "data")
		if err != nil {
			t.Errorf("inserting record produced an error: %v", err)
		}
	}

	histories, err := model.All("client1", "0")
	if err != nil {
		t.Errorf("retrieving records produced an error: %v", err)
	}
	if len(histories) != 5 {
		t.Errorf("expected 5 records for client1, got %d", len(histories))
	}
}

func TestHistoryModelMigrate(t *testing.T) {
	model := setupTestDB(t)

	err := model.Migrate()
	if err != nil {
		t.Errorf("migration produced an error: %v", err)
	}
}
