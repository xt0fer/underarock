package main

//
import (
	"../../underarock"
)

func main() {

	App := &underarock.App{}

	App.Initialize("")

	App.Run(":8085")
}
