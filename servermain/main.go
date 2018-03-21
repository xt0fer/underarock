package main

//
import (
	"github.com/vvmk/underarock"
)

func main() {

	App := &underarock.App{}

	App.Initialize("")

	App.Run(":8085")
}
