package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

var (
	log     = 0
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func processArguments() {
	var (
		showHelp          *bool
		showVersion       bool
		showVersionNumber bool
	)

	// Parse flags
	// see: https://pkg.go.dev/github.com/spf13/pflag
	flag.IntVar(&log, "verbose", 0, "verbose level")
	flag.Lookup("verbose").NoOptDefVal = "1"
	flag.Lookup("verbose").Shorthand = "v"
	flag.BoolVar(&showVersion, "version", false, "display full version information")
	flag.BoolVar(&showVersionNumber, "version-number", false, "display version number")
	// Note: if we use BoolVar for help, we still see "pflag: help requested"
	showHelp = flag.BoolP("help", "h", false, "display help")
	flag.Parse()

	if log > 0 {
		fmt.Fprintln(os.Stderr, "Logging level:", log)
	}

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showVersionNumber {
		fmt.Println(version)
		os.Exit(0)
	}

	if showVersion {
		// GoReleaser automatically sets the version, commit and date
		// see: https://goreleaser.com/cookbooks/using-main.version/
		fmt.Printf("go-aeneas version %s (commit %s, built at %s)\n", version, commit, date)
		os.Exit(0)
	}
}

func main() {
	processArguments()

	fmt.Println("Hello Aeneas")
}
