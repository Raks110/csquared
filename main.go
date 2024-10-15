package main

import (
	"CSquared/preprocessors"
	"fmt"
)

func main() {
	blocks := preprocessors.ProcessToBlocks("./static/test.csq")
	for _, item := range blocks {
		function := preprocessors.ProcessBlock(item)
		fmt.Println("Function reads: ", function.ToString())
	}
}
