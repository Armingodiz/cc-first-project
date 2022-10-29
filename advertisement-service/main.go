package main

import (
	"log"
	"cc-first-project/advertisement-service/app"

	"github.com/joho/godotenv"
)

func main() {
	app := app.NewApp()
	log.Fatalln(app.Start())
}
