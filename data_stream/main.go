package main

import (
	"datastream/routes"
)

func main() {

	r := routes.SetupRouter()

	r.Run(":8080")
}
