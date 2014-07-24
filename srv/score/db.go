package score

import "fmt"
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

func amain() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	db := ConnectDB("../../site/db.sqlite3")
	defer db.Close()
	fmt.Println(db)

	var count int
	db.ScanQuery(&count, "select count(*) from score_team")

	db.Exec("insert into score_team(name, score) values(?,?)", fmt.Sprintf("Team%d", count+1), float64(count)/10)
	//db.Exec("delete from score_team where name like 'b%'")

	rows := db.Query("select name, score from score_team")
	defer rows.Close()
	for rows.Next() {
		var name string
		var score float64
		if err := rows.Scan(&name, &score); err != nil {
			log.Fatal(err)
		}
		fmt.Println(name, score)
	}
}
