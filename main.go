package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

var (
	log         *int
	showVersion bool
	version     = "dev"
	commit      = "none"
	date        = "unknown"
)

func main() {
	log = flag.IntP("verbose", "v", 0, "verbose level")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.Parse()

	if *log > 0 {
		fmt.Fprintln(os.Stderr, "Logging level:", *log)
	}

	if showVersion {
		// GoReleaser automatically sets the version, commit and date
		// see https://goreleaser.com/cookbooks/using-main.version/
		fmt.Printf("go-aeneas version %s (commit %s, built at %s)\n", version, commit, date)
		return
	}

	fmt.Println("Hello Aeneas")
}
