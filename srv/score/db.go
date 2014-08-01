package score

import "log"
import "runtime/debug"
import "database/sql"
import _ "github.com/mattn/go-sqlite3"

type DB struct {
	*sql.DB
}

func check(err error) {
	if err != nil {
		debug.PrintStack()
		log.Fatal(err)
	}
}

func ConnectDB(path string) DB {
	db, err := sql.Open("sqlite3", path)
	check(err)
	return DB{db}
}

func (db *DB) Exec(query string, args ...interface{}) sql.Result {
	result, err := db.DB.Exec(query, args...)
	check(err)
	return result
}

func (db *DB) Query(query string, args ...interface{}) *sql.Rows {
	rows, err := db.DB.Query(query, args...)
	check(err)
	return rows
}

func (db *DB) ScanQuery(dest interface{}, query string, args ...interface{}) {
	err := db.DB.QueryRow(query, args...).Scan(dest)
	check(err)
}

func (db *DB) Begin() Tx {
	tx, err := db.DB.Begin()
	check(err)
	return Tx{tx}
}

type Tx struct {
	*sql.Tx
}

func (tx *Tx) Commit() {
	err := tx.Tx.Commit()
	check(err)
}

func (tx *Tx) Exec(query string, args ...interface{}) sql.Result {
	result, err := tx.Tx.Exec(query, args...)
	check(err)
	return result
}
