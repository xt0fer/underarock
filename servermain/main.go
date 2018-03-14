package main

//
import (
	"github.com/kristofer/underarock"
)

func main() {

	App := &underarock.App{}

	App.Initialize("")

	App.Run(":8085")
}
