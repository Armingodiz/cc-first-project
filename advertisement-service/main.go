package main

import (
	"log"
	"cc-first-project/advertisement-service/app"
)

func main() {
	app := app.NewApp()
	log.Fatalln(app.Start())
}
