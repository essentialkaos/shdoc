package script

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import "github.com/essentialkaos/ek/v12/mathutil"

// ////////////////////////////////////////////////////////////////////////////////// //

// VariableType contains variable type
type VariableType uint8

const (
	VAR_TYPE_UNKNOWN  VariableType = 0
	VAR_TYPE_STRING  VariableType = 1
	VAR_TYPE_NUMBER  VariableType = 2
	VAR_TYPE_BOOLEAN VariableType = 3
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	VAR_MOD_DEFAULT         int = 0
	VAR_MOD_UPPERCASE       int = 1
	VAR_MOD_LOWERCASE       int = 2
	VAR_MOD_UPPERCASE_SHORT int = 3
	VAR_MOD_LOWERCASE_SHORT int = 4
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
	IsWildcard bool         `json:"wildcard"` // Wildcard
}

// Variable contains info about variable
type Variable struct {
	Name  string       `json:"name"`  // Name
	Desc  []string     `json:"desc"`  // Description
	Type  VariableType `json:"type"`  // Type
	Value string       `json:"value"` // Value
	Line  int          `json:"line"`  // LOC of definition
}

// Document contains info about all constants, global variables and methods
type Document struct {
	Title     string      `json:"title"`
	About     []string    `json:"about"`
	Constants []*Variable `json:"constants"`
	Variables []*Variable `json:"variables"`
	Methods   []*Method   `json:"methods"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	typeNameString  = []string{"String", "string", "STRING", "S", "s"}
	typeNameNumber  = []string{"Number", "number", "NUMBER", "N", "n"}
	typeNameBoolean = []string{"Boolean", "boolean", "BOOLEAN", "B", "b"}
)

// ////////////////////////////////////////////////////////////////////////////////// //

// IsValid return false if document is nil or doesn't have any content
func (d *Document) IsValid() bool {
	switch {
	case d == nil:
		return false
	case d.Constants != nil && len(d.Constants) != 0,
		d.Variables != nil && len(d.Variables) != 0,
		d.Methods != nil && len(d.Methods) != 0:
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

	return a.Type == VAR_TYPE_UNKNOWN
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

	return v.Type == VAR_TYPE_UNKNOWN
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

// getTypeName returns variable type name
func getTypeName(t VariableType, mod int) string {
	mod = mathutil.Between(mod, 0, 4)

	switch t {
	case VAR_TYPE_STRING:
		return typeNameString[mod]
	case VAR_TYPE_NUMBER:
		return typeNameNumber[mod]
	case VAR_TYPE_BOOLEAN:
		return typeNameBoolean[mod]
	default:
		return ""
	}
}

// mergeDesc merges description lines to one string
func mergeDesc(data []string) string {
	var result string

	dataLen := len(data)

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
