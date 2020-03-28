package commands

import (
	"fmt"
	"io"

	"github.com/jorgemsrs/zet/pkg/index"
)

// List is the CLI handler function for the list operation
func List(w io.Writer, index index.Index) {
	for _, zettel := range index.GetAll() {
		fmt.Fprintln(w, zettel.Title)
	}
}
