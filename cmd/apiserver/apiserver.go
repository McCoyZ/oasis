package main

import (
	"log"

	"zmc.io/oasis/cmd/apiserver/app"
)

func main() {

	cmd := app.NewAPIServerCommand()

	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
