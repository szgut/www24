package score

import "log"
import "database/sql"
import _ "github.com/mattn/go-sqlite3"

type DB struct {
	*sql.DB
}

func ConnectDB(path string) DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	return DB{db}
}

func (db *DB) Exec(query string, args ...interface{}) sql.Result {
	result, err := db.DB.Exec(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (db *DB) Query(query string, args ...interface{}) *sql.Rows {
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	return rows
}

func (db *DB) ScanQuery(dest interface{}, query string, args ...interface{}) {
	err := db.DB.QueryRow(query, args...).Scan(dest)
	if err != nil {
		log.Fatal(err)
	}
}

func (db *DB) Begin() Tx {
	tx, err := db.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	return Tx{tx}
}

type Tx struct {
	*sql.Tx
}

func (tx *Tx) Commit() {
	if err := tx.Tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func (tx *Tx) Exec(query string, args ...interface{}) sql.Result {
	result, err := tx.Tx.Exec(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
