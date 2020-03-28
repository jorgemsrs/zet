package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jorgemsrs/zet/pkg/api"
	"github.com/jorgemsrs/zet/pkg/commands"
	"github.com/jorgemsrs/zet/pkg/index"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	listCmd := flag.NewFlagSet("ls", flag.ExitOnError)
	/**
	  barCmd := flag.NewFlagSet("bar", flag.ExitOnError)
	  barLevel := barCmd.Int("level", 0, "level")
	*/

	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("must specify a command")
		os.Exit(1)
	}

	// TODO: load this from ~/.zetrc?
	cfg := api.Config{
		ZettelPath: "/Users/silvaj/Documents/Google Drive/4 Archive/zettelkasten",
		IndexPath: "/Users/silvaj/.zetdb",
	}

	_, index := index.Refresh(cfg)

	command := os.Args[1]
	switch command {
	case "open":
	case "search":
		searchCmd.Parse(os.Args[2:])
		term:= strings.Join(searchCmd.Args(), " ")
		commands.Search(index, term)
	case "list":
		listCmd.Parse(os.Args[2:])
		commands.List(os.Stdout, index)
		/*    case "bar":
		      barCmd.Parse(os.Args[2:])
		      fmt.Println("subcommand 'bar'")
		      fmt.Println("  level:", *barLevel)
		      fmt.Println("  tail:", barCmd.Args())
		*/
	default:
		fmt.Println(fmt.Sprintf("unknown command '%s'", command))
		os.Exit(1)
	}

	index.Close()
}
