// logs-converter-cli is a command-line application that allow to parse files with different log formats and according
// on their basis insert MongoDB documents with a monotonous structure.
package main

import "fmt"

var (
	version   string
	date      string
	commit    string
	goversion string
)

func versionInfo() {
	fmt.Printf("Version info: %s:%s\n", version, date)
	fmt.Printf("commit: %s \n", commit)
	fmt.Printf("go version: %s \n", goversion)
}
