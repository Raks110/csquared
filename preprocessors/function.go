package preprocessors

import (
	"fmt"
	"strings"
)

type Argument struct {
	name     string
	dataType string
}

type Function struct {
	function   string
	name       string
	argList    []Argument
	returnType string
}

func ProcessBlock(block Block) Function {

	if len(block.contentLines) == 0 {
		panic("Block is empty! Function is illegal.")
	}

	line1 := block.contentLines[0]
	metaData := strings.Split(line1, "=")

	fmt.Println(line1)
	if len(metaData) != 2 {
		panic("Function definition does not look right!")
	}

	function := Function{}
	function.name, function.returnType = getFunctionHeader(metaData[0])
	function.argList = getArgList(metaData[1])

	return function
}

func getFunctionHeader(data string) (string, string) {
	splitItem := strings.Split(strings.Trim(data, " "), ":")
	if len(splitItem) != 2 {
		panic("Function name defined illegally!")
	}
	return strings.Trim(splitItem[0], " "), strings.Trim(splitItem[1], " ")
}

func getArgList(data string) []Argument {
	var arr []Argument
	splitItem := strings.Split(trimArgList(data), ",")
	if len(splitItem) > 1 {
		for _, item := range splitItem {
			vals := strings.Split(item, ":")
			if len(vals) != 2 {
				panic("Argument list not processed successfully.")
			}
			arr = append(arr, Argument{vals[0], vals[1]})
		}
	}
	return arr
}

func trimArgList(line string) string {
	i := strings.Index(line, "(")
	if i >= 0 {
		j := strings.Index(line, ")")
		if j >= 0 {
			return strings.ReplaceAll(line[i+1:j], " ", "")
		}
	}
	panic("Illegal argument list definition!")
}

func (function Function) ToString() string {
	args := ""
	for index, item := range function.argList {
		args += item.dataType + " " + item.name
		if index < len(function.argList)-1 {
			args += ", "
		}
	}
	return function.returnType + " " + function.name + "(" + args + ");"
}

// TODO: implement this
func nextArgName(name string) string {
	return "name"
}
