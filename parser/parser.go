package parser

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2018 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"pkg.re/essentialkaos/ek.v9/fsutil"
	"pkg.re/essentialkaos/ek.v9/mathutil"
	"pkg.re/essentialkaos/ek.v9/sliceutil"
	"pkg.re/essentialkaos/ek.v9/strutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type EntityType uint8

const (
	ENT_TYPE_UNKNOWN  EntityType = 0
	ENT_TYPE_METHOD              = 1
	ENT_TYPE_VARIABLE            = 2
	ENT_TYPE_CONSTANT            = 3
)

// ////////////////////////////////////////////////////////////////////////////////// //

// VariableType contains variable type
type VariableType uint8

const (
	VAR_TYPE_UKNOWN  VariableType = 0
	VAR_TYPE_STRING               = 1
	VAR_TYPE_NUMBER               = 2
	VAR_TYPE_BOOLEAN              = 3
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Method contains info about method
type Method struct {
	Name       string      `json:"name"`        // Name
	Desc       []string    `json:"desc"`        // Description
	Arguments  []*Argument `json:"arguments"`   // Arguments
	ResultCode bool        `json:"result_code"` // Method uses exit codes
	ResultEcho *Variable   `json:"result_echo"` // Return argument
	Example    []string    `json:"example"`     // Example
	Line       int         `json:"line"`        // LOC of definition
}

// Argument contains info about method argument
type Argument struct {
	Index      string       `json:"index"`    // Index
	Desc       string       `json:"desc"`     // Desc
	Type       VariableType `json:"type"`     // Type
	IsOptional bool         `json:"optional"` // Optional
	IsWildcard bool         `json:"wildcard"` // Wilcard
}

// Variable contains info about variable
type Variable struct {
	Name  string       `json:"name"`  // Name
	Desc  []string     `json:"desc"`  // Description
	Type  VariableType `json:"type"`  // Type
	Value string       `json:"value"` // Value
	Line  int          `json:"line"`  // LOC of definition
}

// Document contians info about all constants, global variables and methods
type Document struct {
	Title     string      `json:"title"`
	About     []string    `json:"about"`
	Constants []*Variable `json:"constants"`
	Variables []*Variable `json:"variables"`
	Methods   []*Method   `json:"methods"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	methodRegExp      = regexp.MustCompile(`^([a-zA-Z0-9._]{1,})\(\)`)
	variableRegExp    = regexp.MustCompile(`^([a-zA-Z0-9_.\[\]]{1,})=(.*)$`)
	constantRegExp    = regexp.MustCompile(`^[A-Z0-9_]{1,}$`)
	numberRegExp      = regexp.MustCompile(`^[0-9]{1,}$`)
	typeCommentRegExp = regexp.MustCompile(`^(.*) \((Boolean|String|Number)\)`)
	methodArgRegExp   = regexp.MustCompile(`([0-9]{1,}|\*):[ ]{0,}(.*)`)
	negativeValRegexp = regexp.MustCompile(`^((N|n)one|(N|n)o(t|)|(F|f)alse)`)
)

// ////////////////////////////////////////////////////////////////////////////////// //

// IsValid return false if document is nil or doesn't have any content
func (d *Document) IsValid() bool {
	if d == nil {
		return false
	}

	if d.Constants != nil && len(d.Constants) != 0 {
		return true
	}

	if d.Variables != nil && len(d.Variables) != 0 {
		return true
	}

	if d.Methods != nil && len(d.Methods) != 0 {
		return true
	}

	return false
}

// HasAbout return true if about is present
func (d *Document) HasAbout() bool {
	if d == nil {
		return false
	}

	return d.About != nil
}

// HasConstants return true if doc has constants info
func (d *Document) HasConstants() bool {
	if d == nil {
		return false
	}

	return d.Constants != nil
}

// HasVariables return true if doc has global variables info
func (d *Document) HasVariables() bool {
	if d == nil {
		return false
	}

	return d.Variables != nil
}

// HasMethods return true if doc has methods info
func (d *Document) HasMethods() bool {
	if d == nil {
		return false
	}

	return d.Methods != nil
}

// TypeDesc return type description
func (a *Argument) TypeName(mod int) string {
	if a == nil {
		return ""
	}

	return getTypeName(a.Type, mod)
}

// IsString return true if type is string
func (a *Argument) IsString() bool {
	if a == nil {
		return false
	}

	return a.Type == VAR_TYPE_STRING
}

// IsNumber return true if type is number
func (a *Argument) IsNumber() bool {
	if a == nil {
		return false
	}

	return a.Type == VAR_TYPE_NUMBER
}

// IsBoolean return true if type is boolean
func (a *Argument) IsBoolean() bool {
	if a == nil {
		return false
	}

	return a.Type == VAR_TYPE_BOOLEAN
}

// IsUnknown return true if type is unknown
func (a *Argument) IsUnknown() bool {
	if a == nil {
		return false
	}

	return a.Type == VAR_TYPE_UKNOWN
}

// TypeDesc return type description
func (v *Variable) TypeName(mod int) string {
	if v == nil {
		return ""
	}

	return getTypeName(v.Type, mod)
}

// IsString return true if type is string
func (v *Variable) IsString() bool {
	if v == nil {
		return false
	}

	return v.Type == VAR_TYPE_STRING
}

// IsNumber return true if type is number
func (v *Variable) IsNumber() bool {
	if v == nil {
		return false
	}

	return v.Type == VAR_TYPE_NUMBER
}

// IsBoolean return true if type is boolean
func (v *Variable) IsBoolean() bool {
	if v == nil {
		return false
	}

	return v.Type == VAR_TYPE_BOOLEAN
}

// IsUnknown return true if type is unknown
func (v *Variable) IsUnknown() bool {
	if v == nil {
		return false
	}

	return v.Type == VAR_TYPE_UKNOWN
}

// UnitedDesc return united description string
func (v *Variable) UnitedDesc() string {
	if v == nil {
		return ""
	}

	return mergeDesc(v.Desc)
}

// HasArguments return true if method has arguments
func (m *Method) HasArguments() bool {
	if m == nil {
		return false
	}

	return m.Arguments != nil
}

// HasEcho return true if method echoed some data
func (m *Method) HasEcho() bool {
	if m == nil {
		return false
	}

	return m.ResultEcho != nil
}

// HasExample return true if method has code usage example
func (m *Method) HasExample() bool {
	if m == nil {
		return false
	}

	return m.Example != nil
}

// UnitedDesc return united description string
func (m *Method) UnitedDesc() string {
	if m == nil {
		return ""
	}

	return mergeDesc(m.Desc)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Parse method parse given file and return document struct and slice with errors
func Parse(file string) (*Document, []error) {
	if !fsutil.CheckPerms("FR", file) {
		return nil, []error{fmt.Errorf("File %s is not readable or not exist", file)}
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
func readData(file string, reader io.Reader) (*Document, []error) {
	scanner := bufio.NewScanner(reader)

	var buffer []string
	var doc = &Document{Title: filepath.Base(file)}

	var lineNum = 0

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimLeft(line, " ")

		lineNum++

		if lineNum == 1 {
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
func parseVariableComment(name, value string, data []string) *Variable {
	if len(data) == 0 || sliceutil.Contains([]string{"private", "PRIVATE", "-"}, strings.TrimRight(data[0], " ")) {
		return nil
	}

	variable := &Variable{
		Name:  name,
		Value: value,
	}

	data, t := getVariableType(data)

	if t == VAR_TYPE_UKNOWN {
		t = guessVariableType(value)
	}

	variable.Type = t
	variable.Desc = data

	return variable
}

// parseMethodComment method parse method comment data and return
// method struct
func parseMethodComment(name string, data []string) *Method {
	if len(data) == 0 || sliceutil.Contains([]string{"private", "PRIVATE", "-"}, strings.TrimRight(data[0], " ")) {
		return nil
	}

	method := &Method{Name: name}

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
func parseArgumentComment(data string) *Argument {
	argument := &Argument{}

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
			argument.Type = VAR_TYPE_BOOLEAN
		case "(Number)":
			argument.Type = VAR_TYPE_NUMBER
		case "(String)":
			argument.Type = VAR_TYPE_STRING
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
func guessVariableType(data string) VariableType {
	if data == "" {
		return VAR_TYPE_UKNOWN
	}

	if data == "true" {
		return VAR_TYPE_BOOLEAN
	}

	if numberRegExp.MatchString(data) {
		return VAR_TYPE_NUMBER
	}

	return VAR_TYPE_STRING
}

// getVariableType search in comment info about variable type and
// return data without comment and varuable type
func getVariableType(data []string) ([]string, VariableType) {
	var result []string
	var resultType VariableType

	for _, line := range data {
		if resultType == VAR_TYPE_UKNOWN {
			if typeCommentRegExp.MatchString(line) {
				cd := typeCommentRegExp.FindStringSubmatch(line)

				// Append to result first regexp group contains
				// description without type marker
				result = append(result, cd[1])

				switch cd[2] {
				case "Boolean":
					resultType = VAR_TYPE_BOOLEAN

				case "Number":
					resultType = VAR_TYPE_NUMBER

				default:
					resultType = VAR_TYPE_STRING
				}

				continue
			}
		}

		result = append(result, line)
	}

	return getCleanData(result), resultType
}

// getTypeName return variable type name
func getTypeName(t VariableType, mod int) string {
	mod = mathutil.Between(mod, 0, 4)

	switch t {
	case VAR_TYPE_STRING:
		return []string{"String", "string", "STRING", "S", "s"}[mod]
	case VAR_TYPE_NUMBER:
		return []string{"Number", "number", "NUMBER", "N", "n"}[mod]
	case VAR_TYPE_BOOLEAN:
		return []string{"Boolean", "boolean", "BOOLEAN", "B", "b"}[mod]
	default:
		return ""
	}
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

// mergeDesc merge description lines to one string
func mergeDesc(data []string) string {
	var result string
	var dataLen = len(data)

	for index, line := range data {
		if line == "" {
			continue
		}

		result += line

		if index < dataLen-1 {
			result += " "
		}
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
