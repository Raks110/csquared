package preprocessors

import (
	"bufio"
	"os"
	"path"
	"strings"
)

type Block struct {
	content string
}

func (block Block) ProcessToC() string {
	return "TBD"
}

func (block Block) ToString() string {
	return block.content
}

func processString(line string, foundOpen bool) string {
	var find = "{"

	if foundOpen {
		find = "}"
	}

	if strings.Contains(line, find) {
		return find
	}
	return "-"
}

func ProcessToBlocks(filePath string) []Block {

	readFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileExt := path.Ext(filePath)
	if fileExt != ".csq" {
		panic("File name should end with .csq")
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var blocks []Block
	found := false
	block := Block{}

	for fileScanner.Scan() {
		foundChar := processString(fileScanner.Text(), found)
		found = found || foundChar == "{"

		block.content += fileScanner.Text()
		if foundChar == "}" {
			blocks = append(blocks, block)
			block = Block{}
			found = false
		}
	}
	err = readFile.Close()
	if err != nil {
		panic(err)
	}

	return blocks
}
