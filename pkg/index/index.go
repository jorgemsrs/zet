// Major inspiration taken from https://github.com/sirupsen/dotfiles/blob/master/home/.bin/fts-search.rb
// (https://superorganizers.substack.com/p/how-to-build-a-learning-machine)
package index

import (
	"database/sql"
	"fmt"
	"github.com/jorgemsrs/zet/pkg/api"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const ddl = `
  CREATE VIRTUAL TABLE IF NOT EXISTS zettelkasten 
	USING fts5(title, body, tags, mtime UNINDEXED, prefix = 3, tokenize = "porter unicode61");
`

func processZettel(path string) (error, api.Zettel) {
	var tagRegexp = regexp.MustCompile(`#[\w-]+`)
	var createdDate = regexp.MustCompile(`^(\d+)`)

	zettel := api.Zettel{}
	basename := filepath.Base(path)
	name := strings.TrimSuffix(basename, filepath.Ext(basename))

	zettel.Path = path
	zettel.Title = name

	b, err := ioutil.ReadFile(path)
	if err != nil {
		//FIXME: this method should return and error and a Zettel tuple
		fmt.Print(err)
	}
	zettel.Content = string(b)

	zettel.Tags = tagRegexp.FindAllString(zettel.Content, -1)
	for _, value := range zettel.Tags {
		if value == "#O7" { //FIXME: extract to constant
			zettel.IsOutline = true
		}
	}

	if match := createdDate.FindStringIndex(basename); match != nil {
		var createdDate = basename[match[0]:match[1]]
		t, _:= time.Parse("200601021504", createdDate) // TODO: handle error
		zettel.Created = t.Unix()
	}

	// Last modified date, as reported by the file system
	stat, _ := os.Stat(path) // TODO: handle error
	zettel.LastModified = stat.ModTime().Unix()

	return nil, zettel
}

type Index struct {
	db *sql.DB
	Zettel []api.Zettel
}


func (c Index) Close() {
	c.db.Close()
}

func (c Index) GetAll() []api.Zettel {
	return c.Zettel
}

func (c Index) Search(term string) {

	results, err := c.db.Query("SELECT rank, highlight(zettelkasten, 0, '\x1b[0;41m', '\x1b[0m'), tags FROM zettelkasten WHERE zettelkasten MATCH ? ORDER BY rank;", term)
	checkErr(err)

	for results.Next() {
		var _score string
		var content string
		var _tags string
		err := results.Scan(&_score, &content, &_tags)
		checkErr(err)
		fmt.Println(content)
	}
}

func Refresh(cfg api.Config) (error, Index) {
	index := Index{}

	db, err := sql.Open("sqlite3", cfg.IndexPath)
	checkErr(err)
	//defer db.Close()

	index.db = db

	_, err = db.Exec(ddl)
	checkErr(err)

	// Weigh tags higher, and title a bit higher.
	_, err = db.Exec("INSERT INTO zettelkasten (zettelkasten, rank) VALUES('rank', 'bm25(2.0, 1.0, 5.0, 0.0)');")
	checkErr(err)

	// Fetch entries already indexed
	existing := make(map[string]int64)
	raw_existing, err := db.Query("SELECT title, mtime FROM zettelkasten")
	checkErr(err)
	defer raw_existing.Close()
	for raw_existing.Next() {
		var title string
		var mtime int64
		err := raw_existing.Scan(&title, &mtime)
		checkErr(err)
		existing[title] = mtime
	}

	// Collect and build in-memory graph of entries in filesystem
	matches, _ := filepath.Glob(filepath.Join(cfg.ZettelPath, "*.md"))
	for _, path := range matches {
		_, zettel := processZettel(filepath.Clean(path))
		index.Zettel = append(index.Zettel, zettel)

		// Any file that's been modified since its entry in the full-text search index
		// will get updated (or if it doesn't exist, of course).
		if _, ok := existing[path]; !ok {
			fmt.Println("add")
			db.Exec("INSERT INTO zettelkasten (title, body, tags, mtime) VALUES (?, ?, ?, ?);", zettel.Path, zettel.Content, strings.Join(zettel.Tags, " "), zettel.LastModified)
		} else if zettel.LastModified > existing[path] {
			fmt.Println("update")
			db.Exec("UPDATE zettelkasten SET body = ?, tags = ?, mtime = ? WHERE title = ?", zettel.Content, strings.Join(zettel.Tags, " "), zettel.LastModified, zettel.Path)
		}

		// mark as visited
		existing[zettel.Path] = -1
	}

	// Delete any entries in the full text index that doesn't have a matching file
	for path, mtime := range existing {
		if mtime != -1 {
			fmt.Println("delete")
			_, err := db.Exec("DELETE FROM zettelkasten WHERE title = ?;", path)
			checkErr(err)
		}
	}

	//fmt.Println(fmt.Sprintf("existing: %+v", existing))
	//fmt.Println(fmt.Sprintf("%+v", index))
	return nil, index
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
