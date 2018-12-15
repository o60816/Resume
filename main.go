// main.go

package main

import (
	"Resume/models"
)

func main() {
	models.InitializeDB()
	initializeRoutes()
}
