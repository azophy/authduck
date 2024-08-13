// history_test.go
package main

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

const EXAMPLE_PAYLOAD = `{"success":"ok"}`

func setupCodeExchangeTestDB(t *testing.T) *CodeExchangeModel {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open the database: %v", err)
	}

	model := &CodeExchangeModel{db: db}
	if err := model.Migrate(); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return model
}

func TestCodeExchangeModelMigrate(t *testing.T) {
	model := setupCodeExchangeTestDB(t)

	err := model.Migrate()
	if err != nil {
		t.Errorf("migration produced an error: %v", err)
	}
}

func TestCodeExchangeModelAddAndGetOne(t *testing.T) {
	model := setupCodeExchangeTestDB(t)

  err := model.Add("client1", "code1", EXAMPLE_PAYLOAD)
	if err != nil {
		t.Errorf("inserting record produced an error: %v", err)
	}

	var count int
	err = model.db.QueryRow("SELECT COUNT(*) FROM code_exchange WHERE client_id = ? AND code = ?", "client1", "code1").Scan(&count)
	if err != nil {
		t.Errorf("counting records produced an error: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 record for client1, got %d", count)
	}

	item, err := model.GetOne("client1", "code1")
	if err != nil {
		t.Errorf("retrieving records produced an error: %v", err)
	}
  if (item.Payload != EXAMPLE_PAYLOAD) {
    t.Errorf("expected record with payload `%v`, got %v", EXAMPLE_PAYLOAD, item.Payload)
  }
}

