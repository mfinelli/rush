package main

import (
	"embed"
	"log"

	"github.com/mfinelli/rush/cmd"
	"github.com/mfinelli/rush/server"
	"github.com/mfinelli/rush/version"
)

//go:embed dist/*
var staticFS embed.FS

//go:embed src/*.tmpl
var templates embed.FS

//go:embed package.json
var pkgjson []byte

func main() {
	if err := version.ParseVersion(&pkgjson); err != nil {
		log.Fatal(err)
	}

	server.DistFiles = staticFS
	server.Templates = templates
	cmd.Execute()
}
