package api

import "fmt"

// Config is a data object containing all the external configuration required to run the program
type Config struct {
	ZettelPath string
	IndexPath string
}

// ConfigResolver is the interface used to load external configuration
type ConfigResolver interface {
	load() Config
}

// Zettel is a TODO:...
type Zettel struct {
	Path         string
	Title        string
	Tags         []string
	Content      string
	Created      int64
	LastModified int64
	IsOutline bool
	//References []Zettel
	//Backlinks []Zettel
}

func (z Zettel) String() string {
	return fmt.Sprintf("Zettel{Title: '%s', Created: '%d', LastModified: '%d', Tags: %v, IsOutline: %t\n", z.Title, z.Created, z.LastModified, z.Tags, z.IsOutline)
}

// Index is TODO:...
type Index interface {
	GetAll() []Zettel
	Search(term string) []Zettel
}

