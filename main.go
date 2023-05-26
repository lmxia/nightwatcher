package main

import (
	"github.com/lmxia/nightwatcher/routers"
)

func main() {
	routersInit := routers.InitRouter()
	routersInit.Run(":8282")
}
