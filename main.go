// main.go

package main

import (
	"Resume/models"

	"google.golang.org/appengine"
)

func main() {
	models.InitializeDB()
	initializeRoutes()

	appengine.Main()
}
