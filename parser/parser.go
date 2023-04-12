package parser

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/sliceutil"
	"github.com/essentialkaos/ek/v12/strutil"

	"github.com/essentialkaos/shdoc/script"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type EntityType uint8

const (
	ENT_TYPE_UNKNOWN  EntityType = 0
	ENT_TYPE_METHOD   EntityType = 1
	ENT_TYPE_VARIABLE EntityType = 2
	ENT_TYPE_CONSTANT EntityType = 3
)

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	methodRegExp      = regexp.MustCompile(`^([a-zA-Z0-9._]{1,})\(\)`)
	variableRegExp    = regexp.MustCompile(`^([a-zA-Z0-9_.\[\]]{1,})=(.*)$`)
	constantRegExp    = regexp.MustCompile(`^[A-Z0-9_]{1,}$`)
	numberRegExp      = regexp.MustCompile(`^[0-9]{1,}$`)
	typeCommentRegExp = regexp.MustCompile(`^(.*) \((Boolean|String|Number)\)`)
	methodArgRegExp   = regexp.MustCompile(`([0-9]{1,}|\*):[ ]{0,}(.*)`)
	negativeValRegexp = regexp.MustCompile(`^((N|n)one|(N|n)o(t|)|(F|f)alse)`)

	shellcheckRegexp = regexp.MustCompile(`\# +shellcheck +disable\=`)
)

var ignoreTags = []string{"private", "PRIVATE", "-"}

// ////////////////////////////////////////////////////////////////////////////////// //

// Parse method parse given file and return document struct and slice with errors
func Parse(file string) (*script.Document, []error) {
	err := fsutil.ValidatePerms("FRS", file)

	if err != nil {
		return nil, []error{err}
	}

	fd, err := os.Open(file)

	if err != nil {
		return nil, []error{err}
	}

	defer fd.Close()

	return readData(file, bufio.NewReader(fd))
}

// ////////////////////////////////////////////////////////////////////////////////// //

// readData read data
func readData(file string, reader io.Reader) (*script.Document, []error) {
	scanner := bufio.NewScanner(reader)

	var buffer []string
	var methodsSection bool
	var lineNum int

	doc := &script.Document{Title: filepath.Base(file)}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimLeft(line, " ")

		lineNum++

		if lineNum == 1 || shellcheckRegexp.MatchString(line) {
			continue
		}

		if line == "" {
			if buffer != nil && !doc.IsValid() {
				doc.About = getCleanData(buffer)
			}

			buffer = nil
			continue
		}

		if strings.Trim(line, "#") == "" {
			if buffer != nil {
				buffer = append(buffer, "")
			}

			continue
		}

		if line[0] == '#' {
			buffer = append(buffer, line[2:])
			continue
		}

		t, name, value := parseEntity(line)

		if t == ENT_TYPE_UNKNOWN || len(buffer) == 0 {
			buffer = nil
			continue
		}

		// Ignore all var definitions after first method
		if t != ENT_TYPE_METHOD && methodsSection {
			buffer = nil
			continue
		}

		switch t {
		case ENT_TYPE_METHOD:
			m := parseMethodComment(name, buffer)

			if m == nil {
				buffer = nil
				continue
			}

			m.Line = lineNum

			// Methods MUST have description
			if len(m.Desc) != 0 {
				doc.Methods = append(doc.Methods, m)
			}

			if !methodsSection {
				methodsSection = true
			}

			buffer = nil

		case ENT_TYPE_VARIABLE, ENT_TYPE_CONSTANT:
			v := parseVariableComment(name, value, buffer)

			if v == nil {
				buffer = nil
				continue
			}

			// Append multiline parts to value
			if isMultilineValue(value) {
			MULTIPART:
				for scanner.Scan() {
					valuePart := scanner.Text()

					v.Value += valuePart

					if strutil.Tail(valuePart, 1) == "\"" {
						break MULTIPART
					}
				}
			}

			v.Line = lineNum

			// Variables MUST have description
			if len(v.Desc) != 0 {
				if t == ENT_TYPE_VARIABLE {
					doc.Variables = append(doc.Variables, v)
				} else {
					doc.Constants = append(doc.Constants, v)
				}
			}

			buffer = nil
		}
	}

	return doc, []error{}
}

// parseEntity method parse entity and return type, name and value of
// entity
func parseEntity(data string) (EntityType, string, string) {
	if methodRegExp.MatchString(data) {
		md := methodRegExp.FindStringSubmatch(data)
		return ENT_TYPE_METHOD, md[1], ""
	}

	if variableRegExp.MatchString(data) {
		vd := variableRegExp.FindStringSubmatch(data)

		if constantRegExp.MatchString(vd[1]) {
			return ENT_TYPE_CONSTANT, vd[1], vd[2]
		}

		return ENT_TYPE_VARIABLE, vd[1], vd[2]
	}

	return ENT_TYPE_UNKNOWN, "", ""
}

// parseVariableComment method parse variable comment data and return
// variable struct
func parseVariableComment(name, value string, data []string) *script.Variable {
	if len(data) == 0 || sliceutil.Contains(ignoreTags, strings.TrimRight(data[0], " ")) {
		return nil
	}

	variable := &script.Variable{
		Name:  name,
		Value: value,
	}

	data, t := getVariableType(data)

	if t == script.VAR_TYPE_UKNOWN {
		t = guessVariableType(value)
	}

	variable.Type = t
	variable.Desc = data

	return variable
}

// parseMethodComment method parse method comment data and return
// method struct
func parseMethodComment(name string, data []string) *script.Method {
	if len(data) == 0 || sliceutil.Contains(ignoreTags, strings.TrimRight(data[0], " ")) {
		return nil
	}

	method := &script.Method{Name: name}

	for index, line := range data {
		if methodArgRegExp.MatchString(line) {
			if method.Desc == nil {
				method.Desc = extractMethodDesc(data, index)
			}

			method.Arguments = append(method.Arguments, parseArgumentComment(line))

			continue
		}

		if strings.HasPrefix(line, "Code:") {
			if method.Desc == nil {
				method.Desc = extractMethodDesc(data, index)
			}

			retValue := strutil.Substr(line, 6, 99999)

			if negativeValRegexp.MatchString(retValue) {
				continue
			}

			method.ResultCode = true
		}

		if strings.HasPrefix(line, "Echo:") {
			if method.Desc == nil {
				method.Desc = extractMethodDesc(data, index)
			}

			echoValue := strutil.Substr(line, 6, 99999)

			if negativeValRegexp.MatchString(echoValue) {
				continue
			}

			method.ResultEcho = parseVariableComment("", "", []string{echoValue})
		}

		if strings.HasPrefix(line, "Example:") {
			if method.Desc == nil {
				method.Desc = extractMethodDesc(data, index)
			}

			method.Example = getCleanData(data[index+1:])
			break // Example is last part of comment
		}
	}

	if method.Desc == nil {
		method.Desc = extractMethodDesc(data, len(data))
	}

	return method
}

// extractMethodDesc return description from all comment data
func extractMethodDesc(data []string, index int) []string {
	if len(data) <= index {
		return getCleanData(data[:])
	}

	return getCleanData(data[:index])
}

// parseArgumentComment method parse given comment data and return
// argument struct
func parseArgumentComment(data string) *script.Argument {
	argument := &script.Argument{}

	ar := methodArgRegExp.FindStringSubmatch(data)

	argument.Index = ar[1]

	if argument.Index == "*" {
		argument.IsWildcard = true
	}

	ds := strings.Split(ar[2], " ")

	var desc []string

	for _, word := range ds {
		switch word {
		case "(Boolean)":
			argument.Type = script.VAR_TYPE_BOOLEAN
		case "(Number)":
			argument.Type = script.VAR_TYPE_NUMBER
		case "(String)":
			argument.Type = script.VAR_TYPE_STRING
		case "[Optional]":
			argument.IsOptional = true
		default:
			desc = append(desc, word)
		}
	}

	argument.Desc = strings.Join(desc, " ")

	return argument
}

// guessVariableType try to guess variable type by value
func guessVariableType(data string) script.VariableType {
	if data == "" {
		return script.VAR_TYPE_UKNOWN
	}

	if data == "true" {
		return script.VAR_TYPE_BOOLEAN
	}

	if numberRegExp.MatchString(data) {
		return script.VAR_TYPE_NUMBER
	}

	return script.VAR_TYPE_STRING
}

// getVariableType search in comment info about variable type and
// return data without comment and varuable type
func getVariableType(data []string) ([]string, script.VariableType) {
	var result []string
	var resultType script.VariableType

	for _, line := range data {
		if resultType == script.VAR_TYPE_UKNOWN {
			if typeCommentRegExp.MatchString(line) {
				cd := typeCommentRegExp.FindStringSubmatch(line)

				// Append to result first regexp group contains
				// description without type marker
				result = append(result, cd[1])

				switch cd[2] {
				case "Boolean":
					resultType = script.VAR_TYPE_BOOLEAN

				case "Number":
					resultType = script.VAR_TYPE_NUMBER

				default:
					resultType = script.VAR_TYPE_STRING
				}

				continue
			}
		}

		result = append(result, line)
	}

	return getCleanData(result), resultType
}

// getCleanData return removes empty lines and whitespaces at the
// end of the line
func getCleanData(data []string) []string {
	if len(data) == 0 {
		return data
	}

	var result []string

	lastNonEmptyIndex := 0

	// Search index of last non empty line
	for index, line := range data {
		if line != "" {
			lastNonEmptyIndex = index
		}
	}

	// Make result slice with non empty lines without whitespaces
	// at end of the line
	for i := 0; i <= lastNonEmptyIndex; i++ {
		result = append(result, strings.TrimRight(data[i], " "))
	}

	return result
}

// isMultilineValue return true if value is multiline string
func isMultilineValue(value string) bool {
	if strutil.Head(value, 1) == "\"" && strutil.Tail(value, 1) != "\"" {
		return true
	}

	return false
}
