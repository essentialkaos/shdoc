package parser

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2015 Essential Kaos                         //
//      Essential Kaos Open Source License <http://essentialkaos.com/ekol?en>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/essentialkaos/ek/fsutil"
	"github.com/essentialkaos/ek/mathutil"
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

type VariableType uint8

const (
	VAR_TYPE_UKNOWN  VariableType = 0
	VAR_TYPE_STRING               = 1
	VAR_TYPE_NUMBER               = 2
	VAR_TYPE_BOOLEAN              = 3
)

// ////////////////////////////////////////////////////////////////////////////////// //

type Method struct {
	Name       string      `json:"name"`        // Name
	Desc       []string    `json:"desc"`        // Description
	Arguments  []*Argument `json:"arguments"`   // Arguments
	ResultCode bool        `json:"result_code"` // Method uses exit codes
	ResultEcho *Variable   `json:"result_echo"` // Return argument
	Example    []string    `json:"example"`
	Line       int         `json:"line"` // LOC of definition
}

type Argument struct {
	Index      string       `json:"index"`    // Index
	Desc       string       `json:"desc"`     // Desc
	Type       VariableType `json:"type"`     // Type
	IsOptional bool         `json:"optional"` // Optional
	IsWildcard bool         `json:"wildcard"` // Wilcard
}

type Variable struct {
	Name  string       `json:"name"`  // Name
	Desc  []string     `json:"desc"`  // Description
	Type  VariableType `json:"type"`  // Type
	Value string       `json:"value"` // Value
	Line  int          `json:"line"`  // LOC of definition
}

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

	if d.Constants != nil {
		return true
	}

	if d.Variables != nil {
		return true
	}

	if d.Methods != nil {
		return true
	}

	return false
}

// HasAbout return true if about is present
func (d *Document) HasAbout() bool {
	return d.About != nil
}

// HasConstants return true if doc has constants info
func (d *Document) HasConstants() bool {
	return d.Constants != nil
}

// HasVariables return true if doc has global variables info
func (d *Document) HasVariables() bool {
	return d.Variables != nil
}

// HasMethods return true if doc has methods info
func (d *Document) HasMethods() bool {
	return d.Methods != nil
}

// TypeDesc return type description
func (a *Argument) TypeName(mod int) string {
	return getTypeName(a.Type, mod)
}

// IsString return true if type is string
func (a *Argument) IsString() bool {
	return a.Type == VAR_TYPE_STRING
}

// IsNumber return true if type is number
func (a *Argument) IsNumber() bool {
	return a.Type == VAR_TYPE_NUMBER
}

// IsBoolean return true if type is boolean
func (a *Argument) IsBoolean() bool {
	return a.Type == VAR_TYPE_BOOLEAN
}

// IsUnknown return true if type is unknown
func (a *Argument) IsUnknown() bool {
	return a.Type == VAR_TYPE_UKNOWN
}

// TypeDesc return type description
func (v *Variable) TypeName(mod int) string {
	return getTypeName(v.Type, mod)
}

// IsString return true if type is string
func (v *Variable) IsString() bool {
	return v.Type == VAR_TYPE_STRING
}

// IsNumber return true if type is number
func (v *Variable) IsNumber() bool {
	return v.Type == VAR_TYPE_NUMBER
}

// IsBoolean return true if type is boolean
func (v *Variable) IsBoolean() bool {
	return v.Type == VAR_TYPE_BOOLEAN
}

// IsUnknown return true if type is unknown
func (v *Variable) IsUnknown() bool {
	return v.Type == VAR_TYPE_UKNOWN
}

// UnitedDesc return united description string
func (v *Variable) UnitedDesc() string {
	return strings.TrimRight(strings.Join(v.Desc, " "), " ")
}

// HasArguments return true if method has arguments
func (m *Method) HasArguments() bool {
	return m.Arguments != nil
}

// HasEcho return true if method echoed some data
func (m *Method) HasEcho() bool {
	return m.ResultEcho != nil
}

// UnitedDesc return united description string
func (m *Method) UnitedDesc() string {
	return strings.TrimRight(strings.Join(m.Desc, " "), " ")
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

	reader := bufio.NewReader(fd)
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
				doc.About = buffer
			}

			buffer = nil
			continue
		}

		if strings.Trim(line, "#") == "" {
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
			m.Line = lineNum

			doc.Methods = append(doc.Methods, m)

			buffer = nil

		case ENT_TYPE_VARIABLE, ENT_TYPE_CONSTANT:
			v := parseVariableComment(name, value, buffer)

			v.Line = lineNum

			if t == ENT_TYPE_VARIABLE {
				doc.Variables = append(doc.Variables, v)
			} else {
				doc.Constants = append(doc.Constants, v)
			}

			buffer = nil
		}
	}

	return doc, []error{}
}

// ////////////////////////////////////////////////////////////////////////////////// //

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

			retValue := line[6:]

			if negativeValRegexp.MatchString(retValue) {
				continue
			}

			method.ResultCode = true
		}

		if strings.HasPrefix(line, "Echo:") {
			if method.Desc == nil {
				method.Desc = extractMethodDesc(data, index)
			}

			echoValue := line[6:]

			if negativeValRegexp.MatchString(echoValue) {
				continue
			}

			method.ResultEcho = parseVariableComment("", "", []string{echoValue})
		}

		if strings.HasPrefix(line, "Example:") {
			if method.Desc == nil {
				method.Desc = extractMethodDesc(data, index)
			}

			method.Example = data[index+1:]
			break // Example is last part of comment
		}
	}

	if method.Desc == nil {
		method.Desc = extractMethodDesc(data, len(data))
	}

	return method
}

func extractMethodDesc(data []string, index int) []string {
	if data[index] == "" {
		return data[:index-1]
	}

	return data[:index]
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

	return result, resultType
}

func getTypeName(t VariableType, mod int) string {
	mod = mathutil.Between(mod, 0, 3)

	switch t {
	case VAR_TYPE_STRING:
		return []string{"String", "string", "STRING", "S"}[mod]
	case VAR_TYPE_NUMBER:
		return []string{"Number", "number", "NUMBER", "N"}[mod]
	case VAR_TYPE_BOOLEAN:
		return []string{"Boolean", "boolean", "BOOLEAN", "B"}[mod]
	default:
		return ""
	}
}
