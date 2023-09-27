package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

var (
	log         *int
	help        *bool
	showVersion bool
	version     = "dev"
	commit      = "none"
	date        = "unknown"
)

func main() {
	// Parse flags
	// see: https://pkg.go.dev/github.com/spf13/pflag
	log = flag.IntP("verbose", "v", 0, "verbose level")
	flag.Lookup("verbose").NoOptDefVal = "1"
	flag.BoolVar(&showVersion, "version", false, "display version")
	help = flag.BoolP("help", "h", false, "display help")
	flag.Parse()

	if *log > 0 {
		fmt.Fprintln(os.Stderr, "Logging level:", *log)
	}

	if *help {
		flag.Usage()
		return
	}

	if showVersion {
		// GoReleaser automatically sets the version, commit and date
		// see: https://goreleaser.com/cookbooks/using-main.version/
		fmt.Printf("go-aeneas version %s (commit %s, built at %s)\n", version, commit, date)
		return
	}

	fmt.Println("Hello Aeneas")
}
