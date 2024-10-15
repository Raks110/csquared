package preprocessors

import (
	"fmt"
	"strings"
	"unicode"
)

type Argument struct {
	name     string
	dataType string
}

type Function struct {
	function   string
	name       string
	argList    map[string]Argument
	variables  map[string]Argument
	returnType string
}

func (function Function) ToString() string {
	args := ""
	index := 0
	for _, item := range function.argList {
		args += item.dataType + " " + item.name
		if index < len(function.argList)-1 {
			args += ", "
		}
		index++
	}
	fmt.Println("[DEBUG] Variable list: ", function.variables)
	return function.returnType + " " + function.name + "(" + args + ");"
}

func ProcessBlock(block Block) Function {

	if len(block.contentLines) == 0 {
		panic("Block is empty! Function is illegal.")
	}

	line1 := block.contentLines[0]
	metaData := strings.Split(line1, "=")
	if len(metaData) != 2 {
		panic("Function definition does not look right!")
	}

	function := Function{}
	function.name, function.returnType = getFunctionHeader(metaData[0])
	function.argList = getArgList(metaData[1])

	function.variables = map[string]Argument{}
	for i := 1; i < len(block.contentLines); i++ {
		processLine(block.contentLines[i], function)
	}

	return function
}

func processLine(data string, function Function) {
	splitItem := strings.Split(data, "=")
	if len(splitItem) == 2 {
		variable := getVariableDetails(splitItem[0])
		function.variables[variable.name] = variable
	} else if len(splitItem) == 1 && strings.Contains(data, ":") {
		variable := getVariableDetails(splitItem[0])
		function.variables[variable.name] = variable
	}
}

func getVariableDetails(data string) Argument {
	data = strings.Trim(data, " ")
	varName := "<TBD>"
	dataType := "<TBD>"
	if strings.Contains(data, ":") {
		varDetails := strings.Split(data, ":")
		varName = strings.Trim(varDetails[0], " ")
		dataType = strings.Trim(varDetails[1], " ")
	} else {
		panic("Variable data type is necessary.")
	}
	if !isValidVariableName(varName) {
		panic("Variable name is invalid!")
	}

	return Argument{varName, dataType}
}

func getFunctionHeader(data string) (string, string) {
	splitItem := strings.Split(strings.Trim(data, " "), ":")
	if len(splitItem) != 2 {
		panic("Function name defined illegally!")
	}
	return strings.Trim(splitItem[0], " "), strings.Trim(splitItem[1], " ")
}

func getArgList(data string) map[string]Argument {
	var arr []Argument
	splitItem := strings.Split(trimArgList(data), ",")
	if len(splitItem) > 1 {
		for _, item := range splitItem {
			vals := strings.Split(item, ":")
			if len(vals) > 2 {
				panic("Argument list not processed successfully.")
			}
			if len(vals) == 1 {
				arr = append(arr, Argument{vals[0], "<TBD>"})
			} else {
				for i := len(arr) - 1; i >= 0; i-- {
					if arr[i].dataType == "<TBD>" {
						arr[i].dataType = vals[1]
					} else {
						break
					}
				}
				arr = append(arr, Argument{vals[0], vals[1]})
			}
		}
	}
	m := make(map[string]Argument)

	for _, argument := range arr {
		m[argument.name] = argument
	}
	return m
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

func isValidVariableName(name string) bool {
	if len(name) == 0 {
		return false
	}
	firstChar := rune(name[0])
	if !unicode.IsLetter(firstChar) && firstChar != '_' {
		return false
	}
	for _, char := range name[1:] {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' {
			return false
		}
	}
	reservedKeywords := []string{"break", "case", "const", "continue", "default", "else", "for", "goto", "if", "class", "interface", "map", "package", "return", "select", "struct", "switch", "type"}
	for _, keyword := range reservedKeywords {
		if name == keyword {
			return false
		}
	}
	return true
}

// TODO: implement this
func nextArgName(name string) string {
	return "name"
}
