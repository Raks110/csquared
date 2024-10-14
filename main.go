package main

import (
	"CSquared/preprocessors"
	"fmt"
)

func main() {
	blocks := preprocessors.ProcessToBlocks("./static/test.csq")
	for _, item := range blocks {
		fmt.Println("Block reads: ", item.ToString())
	}
}
