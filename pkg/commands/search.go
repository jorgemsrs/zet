package commands

import (
	"fmt"
	"github.com/jorgemsrs/zet/pkg/index"
)

// Search is the CLI handler function for the search operation
func Search(index index.Index, term string) {
	fmt.Println(fmt.Sprintf("search: \"%s\"", term) )

	index.Search(term)
	/*
	matches, _ := filepath.Glob(filepath.Join(cfg.ZettelPath, "*.md"))
	for _, path := range matches {
		stat, _ := os.Stat(path)
		stat.ModTime().Unix()
		fmt.Println(filepath.Base(path))

		// Any file that's been modified since its entry in the full-text search index
		// will get updated (or if it doesn't exist, of course).
		/*
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
	*/
}

