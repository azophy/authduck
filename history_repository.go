// adapted from https://gosamples.dev/sqlite-intro/
package main

import (
	"database/sql"
	"log"
	//"errors"

	_ "modernc.org/sqlite"
)

// var (
// ErrDuplicate    = errors.New("record already exists")
// ErrNotExists    = errors.New("row not exists")
// ErrUpdateFailed = errors.New("update failed")
// ErrDeleteFailed = errors.New("delete failed")
// )
type History struct {
	ID          int
	Timestamp   string
	ClientId    string
	HTTPMethod  string
	Url         string
	Headers     string
	Body        string
	QueryParams string
}

type HistoryModel struct {
	db *sql.DB
}

func NewHistoryModel(db *sql.DB) (*HistoryModel, error) {
	newInstance := &HistoryModel{
		db: db,
	}

	newInstance.Migrate()
	return newInstance, nil
}

func (r *HistoryModel) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS history(
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      client_id TEXT NOT NULL,
      http_method TEXT NOT NULL,
      url TEXT NOT NULL,
      headers TEXT,
      body TEXT,
      query_params TEXT
    );
  `

	_, err := r.db.Exec(query)
	return err
}

func (r *HistoryModel) Record(clientId, httpMethod, url, headers, body, queryParams string) error {
	res, err := r.db.Exec(`
    INSERT INTO history(client_id, http_method, url, headers, body, query_params)
    VALUES(?, ?, ?, ?, ?, ?)
  `, clientId, httpMethod, url, headers, body, queryParams)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	log.Printf("inserted history with client id %v & history id %v\n", clientId, id)

	return nil
}

func (r *HistoryModel) All(clientId string, from string) ([]History, error) {
	rows, err := r.db.Query("SELECT * FROM history WHERE client_id = ? ORDER BY timestamp DESC LIMIT 10 OFFSET ?", clientId, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []History
	for rows.Next() {
		var history History
		if err := rows.Scan(
			&history.ID,
			&history.Timestamp,
			&history.ClientId,
			&history.HTTPMethod,
			&history.Url,
			&history.Headers,
			&history.Body,
			&history.QueryParams,
		); err != nil {
			return nil, err
		}
		all = append(all, history)
	}

	return all, nil
}
