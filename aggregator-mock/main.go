package main

import (
	"aggregator_mock/http"
)

func main() {
	println("TEST")
	http.InitTestController().Router.Run(":8084")
}
