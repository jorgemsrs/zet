package main

import (
    "flag"
    "fmt"
    "os"

    _ "github.com/jorgemsrs/zet/pkg/search"
)

func main() {

    listCmd := flag.NewFlagSet("ls", flag.ExitOnError)
/**
    barCmd := flag.NewFlagSet("bar", flag.ExitOnError)
    barLevel := barCmd.Int("level", 0, "level")
*/
    if len(os.Args) < 2 {
        fmt.Println("expected 'list' or 'bar' subcommands")
        os.Exit(1)
    }

    switch os.Args[1] {

    case "list":
        listCmd.Parse(os.Args[2:])
        fmt.Println("subcommand 'foo'")
        fmt.Println("  tail:", listCmd.Args())
/*    case "bar":
        barCmd.Parse(os.Args[2:])
        fmt.Println("subcommand 'bar'")
        fmt.Println("  level:", *barLevel)
        fmt.Println("  tail:", barCmd.Args())
	*/
    default:
        fmt.Println("expected 'foo' or 'bar' subcommands")
        os.Exit(1)
    }
}
