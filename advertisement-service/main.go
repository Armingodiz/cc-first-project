package main

import (
	"log"

	"cc-first-project/user-service/app"
)

func main() {
	app := app.NewApp()
	log.Fatalln(app.Start(":" + "3000"))
}
