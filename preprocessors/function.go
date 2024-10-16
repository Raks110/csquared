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

/* State Diagram based programs */

func variableState(line string, previous string, sw int, function Function, functionMap map[string]Function, metaData map[string]string) {

	for index, c := range line {
		if sw == 0 {
			if isControl(previous + string(c)) {
				controlState(previous+string(c), line[index+1:])
			} else if isValidVariableName(previous + string(c)) {
				variableState(line[index+1:], previous+string(c), sw, function, functionMap, metaData)
				return
			} else {
				if c == ':' {
					function.variables[previous] = Argument{previous, "<TBD>"}
					metaData["varname"] = previous
					variableState(line[index+1:], "", 1, function, functionMap, metaData)
					return
				}
				if c == '=' {
					function.variables[previous] = Argument{previous, "<TBD>"}
					metaData["varname"] = previous
					variableState(line[index+1:], "", 2, function, functionMap, metaData)
					return
				}
			}
		}
		//for data type
		if sw == 1 {
			//TODO: validate type
			if c == ' ' {
				continue
			}
			if unicode.IsLetter(c) || c == '_' {
				variableState(line[index+1:], previous+string(c), 1, function, functionMap, metaData)
				return
			} else if c == '!' {
				//TODO: check
				function.variables[metaData["varname"]] = Argument{metaData["varname"], previous}
				metaData["vartype"] = previous
				variableState(line[index+1:], "", 0, function, functionMap, metaData)
				return
			}
		}
		//for "="
		if sw == 2 {
			previous = strings.Trim(previous, " ")
			fmt.Println("Var red: ", line, " prev: <", previous, ">", " Is? ", previous == "", " Is Char? <", string(c), ">")
			cc := rune(c)
			if unicode.IsLetter(cc) || c == '_' {
				variableState(line[index+1:], previous+string(c), 2, function, functionMap, metaData)
				return
			} else if c == '(' {
				prevFunct, found := isFunction(previous, functionMap)
				fmt.Println(prevFunct, found)
				if !found {
					panic("Function does not exist!")
				}

				ret, set := metaData[previous]
				if !set {
					function.variables[metaData["varname"]] = Argument{metaData["varname"], prevFunct.returnType}
					metaData["vartype"] = prevFunct.returnType
				} else if ret != prevFunct.returnType {
					panic("Function and variable data type do not match!")
				}
				return
			} else if previous != "" {
				panic("Illegal space!")
			}
		}
	}
}

func controlState(control string, upcoming string) {
	fmt.Println("<TBD>")
}

func isControl(data string) bool {
	return data == "if" || data == "switch" || data == "return"
}

func isNewVariable(name string, function Function) (Argument, bool) {
	itemVar, okv := function.variables[name]
	itemArg, oka := function.argList[name]

	if !okv && !oka {
		return Argument{}, true
	}
	if !okv {
		return itemArg, false
	}
	if !oka {
		return itemVar, false
	}
	panic("Variable invalid!")
}

func isFunction(name string, functionMap map[string]Function) (Function, bool) {

	itemVar, okf := functionMap[name]

	if !okf {
		return Function{}, false
	}
	return itemVar, true
}

/* End of State Diagram based programs */

func ProcessBlock(block Block) Function {

	if len(block.contentLines) == 0 {
		panic("Block is empty! Function is illegal.")
	}

	line1 := block.contentLines[0]
	metaData := strings.Split(line1, "=")
	if len(metaData) != 2 {
		panic("Function definition does not look right!")
	}

	metaDataForVar := make(map[string]string)
	funcMap := make(map[string]Function)

	funcMap["add"] = Function{function: "int add(x, y) {return x + y;}", name: "add", argList: map[string]Argument{"x": Argument{"x", "int"}, "y": Argument{"y", "int"}}, variables: make(map[string]Argument), returnType: "int"}

	function := Function{}
	function.name, function.returnType = getFunctionHeader(metaData[0])
	function.argList = getArgList(metaData[1])

	function.variables = map[string]Argument{}
	for i := 1; i < len(block.contentLines); i++ {
		//processLine(block.contentLines[i], function)
		variableState(block.contentLines[i], "", 0, function, funcMap, metaDataForVar)
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
