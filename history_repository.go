// adapted from https://gosamples.dev/sqlite-intro/
package main

import (
  "log"
  "database/sql"
  //"errors"

  _ "github.com/mattn/go-sqlite3"
)

// var (
// ErrDuplicate    = errors.New("record already exists")
// ErrNotExists    = errors.New("row not exists")
// ErrUpdateFailed = errors.New("update failed")
// ErrDeleteFailed = errors.New("delete failed")
// )
type History struct {
  ID int
  Timestamp string
  ClientId string
  Data string
}

type HistoryModel struct {
  db *sql.DB
}

func NewHistoryModel(path string) *HistoryModel {
  db, err := sql.Open("sqlite3", path)
  if err != nil {
    log.Fatal(err)
  }

  newInstance :=  &HistoryModel{
    db: db,
  }

  newInstance.Migrate()
  return newInstance
}

func (r *HistoryModel) Migrate() error {
  query := `
    CREATE TABLE IF NOT EXISTS history(
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      client_id TEXT NOT NULL,
      data TEXT NOT NULL
    );
  `

  _, err := r.db.Exec(query)
  return err
}

func (r *HistoryModel) Record(clientId, data string) error {
  _, err := r.db.Exec("INSERT INTO history(client_id, data) values(?, ?)",clientId, data)
  if err != nil {
    return err
  }

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
    if err := rows.Scan(&history.ID, &history.Timestamp , history.ClientId, &history.Data); err != nil {
      return nil, err
    }
    all = append(all, history)
  }

  return all, nil
}

