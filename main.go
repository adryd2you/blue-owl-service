package main

import (
	"blue-owl-service/route"
	"blue-owl-service/transport"

	_ "github.com/lib/pq"
)

func main() {
	transport.InitDatabase()
	route.InitRoute()
}
