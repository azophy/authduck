// history_test.go
package main

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupHistoryTestDB(t *testing.T) *HistoryModel {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open the database: %v", err)
	}

	model := &HistoryModel{db: db}
	if err := model.Migrate(); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return model
}

func TestHistoryModelMigrate(t *testing.T) {
	model := setupHistoryTestDB(t)

	err := model.Migrate()
	if err != nil {
		t.Errorf("migration produced an error: %v", err)
	}
}

func TestHistoryModelRecord(t *testing.T) {
	model := setupHistoryTestDB(t)

	err := model.Record("client1", "example.com", "GET", "X-example: wow", `{"success":"ok"}`, "")
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
	model := setupHistoryTestDB(t)

	for i := 0; i < 5; i++ {
		err := model.Record("client1", "example.com", "GET", "X-example: wow", `{"success":"ok"}`, "")
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
