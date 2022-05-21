package main

import (
	"embed"
	"log"

	"github.com/mfinelli/rush/cmd"
	"github.com/mfinelli/rush/server"
	"github.com/mfinelli/rush/version"
)

//go:embed package.json
var pkgjson []byte

func main() {
	if err := version.ParseVersion(&pkgjson); err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}
