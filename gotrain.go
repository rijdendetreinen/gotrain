package main

import (
	"github.com/rijdendetreinen/gotrain/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Version = cmd.VersionInformation{
		version,
		commit,
		date,
	}
	cmd.Execute()
}
