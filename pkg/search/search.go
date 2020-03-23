package search

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

const DDL = `
  CREATE VIRTUAL TABLE IF NOT EXISTS zettelkasten 
	USING fts5(title, body, tags, mtime UNINDEXED, prefix = 3, tokenize = "porter unicode61");
`

func main() {
	defer fmt.Println("Hello, world.")

	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)
	defer db.Close()

	_, err = db.Exec(DDL)
	checkErr(err)

	_, err = db.Exec("INSERT INTO zettelkasten (zettelkasten, rank) VALUES('rank', 'bm25(2.0, 1.0, 5.0, 0.0)');")
	checkErr(err)

	raw_existing, _ := db.Query("SELECT title, mtime FROM zettelkasten")
	defer raw_existing.Close()
	var existing map[string]int
	for raw_existing.Next() {
		var title string
		var mtime int
		err2 := raw_existing.Scan(&title, &mtime)
		if err2 != nil { panic(err2) }
		existing[title] = mtime
	}

	matches, err := filepath.Glob("/Users/silvaj/Documents/Google Drive/4 Archive/zettelkasten/*.md")
	for _, path := range matches {
		stat, _ := os.Stat(path)
		fmt.Println(path, stat.ModTime().Unix())

		// Any file that's been modified since its entry in the full-text search index
		// will get updated (or if it doesn't exist, of course).
		if !existing[path]
			contents = File.read(path)
			tags = contents.scan(/#[\w-]+/).join(" ")
			db.execute(<<-SQL, [path, contents, tags, File.stat(path).mtime.to_s])
			INSERT INTO zettelkasten (title, body, tags, mtime) VALUES (?, ?, ?, ?);
			SQL
		elsif mtime.to_i > existing[path] # to_i because the stat may have more precision
			contents = File.read(path)
			tags = contents.scan(/#[\w-]+/).join(" ")
			db.execute(<<-SQL, [contents, tags, mtime.to_s, path])
			UPDATE zettelkasten SET body = ?, tags = ?, mtime = ? WHERE title = ?
			SQL
		end

		existing[path] = 'VISITED'
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
