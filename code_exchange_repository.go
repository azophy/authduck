// adapted from https://gosamples.dev/sqlite-intro/
package main

import (
	"database/sql"
	"errors"
	"log"

	_ "modernc.org/sqlite"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type CodeExchange struct {
	ID        int
	Timestamp string
	ClientId  string
	Code      string
	Payload   string
}

type CodeExchangeModel struct {
	db *sql.DB
}

func NewCodeExchangeModel(db *sql.DB) (*CodeExchangeModel, error) {
	newInstance := &CodeExchangeModel{
		db: db,
	}

	newInstance.Migrate()
	return newInstance, nil
}

func (r *CodeExchangeModel) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS code_exchange(
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      client_id TEXT NOT NULL,
      code TEXT NOT NULL,
      payload TEXT NOT NULL
    );
  `

	_, err := r.db.Exec(query)
	return err
}

func (r *CodeExchangeModel) Add(clientId, code, payload string) error {
	res, err := r.db.Exec(`
    INSERT INTO code_exchange(client_id, code, payload)
    VALUES(?, ?, ?)
  `, clientId, code, payload)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	log.Printf("inserted code exchange with client id %v & code exchange id %v\n", clientId, id)

	return nil
}

func (r *CodeExchangeModel) GetOne(clientId string, code string) (*CodeExchange, error) {
	var item CodeExchange
	row := r.db.QueryRow("SELECT * FROM code_exchange WHERE client_id = ? AND code = ?", clientId, code)
	if err := row.Scan(
		&item.ID,
		&item.Timestamp,
		&item.ClientId,
		&item.Code,
		&item.Payload,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}

	return &item, nil
}
