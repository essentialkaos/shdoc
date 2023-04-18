package parser

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"io/ioutil"
	"testing"

	"github.com/essentialkaos/shdoc/script"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const _SCRIPT = `#!/bin/bash
# shellcheck disable=SC1117

# This is example of shell script.      
# Second line of about info.
# 
# Third line of about info.
#
#

###############################################################################

# Constant #1 without type
# Second line of description
CONST_1=""

# Constant #2 without type
CONST_2=0

# Constant #3 without type
CONST_3=true

# Constant #4 with type (String)
CONST_4=""

# Constant #5 with type (Number)
CONST_5=""

# Constant #6 with type (Boolean)
CONST_6=""

# -
# Private constant
CONST_7=""

###############################################################################

# Variable #1 without type
# Second line of description
var_1=""

# Variable #2 without type
var_2=1

# Variable #3 without type
var_3=true

# Variable #4 with type (String)
var_4=""

# Variable #5 with type (Number)
var_5=""

# Variable #6 with type (Boolean)
var_6=""

# Variable #7 without value
# shellcheck disable=SC1117
var_7=

# Variable #8 with multiline value
var_8="This is 
multiline 
value"

# -
# Private variable
var_9=""

###############################################################################

# This is desc for method #1.
# Second line of description.
# Third line of description.
#
# Fourth line of description.
#
# Code: none
# Echo: none
method1() {
  stub=1
  # Must be ignored
  stub=2
}

# This is desc for method #2.
#
# 1: First argument (String)
# 2: Second argument (Number)
# 3: Third argument (Boolean)
# *: Wildcard argument
#
# Code: Yes
# Echo: Magic value (Boolean)
# 
# Example:
# if [[ -f $file ]] ; then
#   method2 123
# fi
#
method2() {
  stub=1
}

# This is desc for method #3.
#
#
# 1: First argument (String) [Optional]
#
# Code: No
# Echo: No
method3() {
  stub=1
}

# This is desc for method #4.
method4() {
  stub=1
}

# This is desc for method #5.
# Echo: none
method5() {
  stub=1
}

# This is desc for method #6.
# Code: none
method6() {
  stub=1
}

# This is desc for method #7.
#
# 1: First argument
method7() {
  stub=1
}

# This is desc for method #8.
#
# Example:
# method8 123
method8() {
  stub=1
}

#
# This is desc for method #9.
method9() {
  stub=1 
}

# Code: Yes
method10() {
  stub=1 
}

# Echo: Magic value (Number)
method11() {
  stub=1 
}

# -
#
# This is private method
method12() {
  stub=1 
}

# 1: ()
# Code:
# Echo:
# Example:
method13(){
  stub=1
}
`

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type ParseSuite struct {
	TmpDir string
}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&ParseSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *ParseSuite) SetUpSuite(c *C) {
	s.TmpDir = c.MkDir()

	err := ioutil.WriteFile(s.TmpDir+"/script.sh", []byte(_SCRIPT), 0644)

	if err != nil {
		c.Fatal(err.Error())
	}
}

func (s *ParseSuite) TestErrors(c *C) {
	doc, errs := Parse(s.TmpDir + "/script1.sh")

	c.Assert(doc, IsNil)
	c.Assert(errs, Not(HasLen), 0)
	c.Assert(errs[0], ErrorMatches, "File .*/script1.sh doesn't exist or not accessible")
}

func (s *ParseSuite) TestParsing(c *C) {
	doc, errs := Parse(s.TmpDir + "/script.sh")

	c.Assert(doc, NotNil)
	c.Assert(errs, HasLen, 0)

	c.Assert(doc.IsValid(), Equals, true)
	c.Assert(doc.HasAbout(), Equals, true)
	c.Assert(doc.HasConstants(), Equals, true)
	c.Assert(doc.HasVariables(), Equals, true)
	c.Assert(doc.HasMethods(), Equals, true)

	c.Assert(doc.Title, Equals, "script.sh")

	c.Assert(doc.About, HasLen, 4)
	c.Assert(doc.About[0], Equals, "This is example of shell script.")
	c.Assert(doc.About[1], Equals, "Second line of about info.")
	c.Assert(doc.About[2], Equals, "")
	c.Assert(doc.About[3], Equals, "Third line of about info.")

	c.Assert(doc.Constants, HasLen, 6)
	c.Assert(doc.Variables, HasLen, 8)
	c.Assert(doc.Methods, HasLen, 9)

	// //////////////////////////////////////////////////////////////////////////////// //

	c.Assert(doc.Constants[0], NotNil)
	c.Assert(doc.Constants[0].Name, Equals, "CONST_1")
	c.Assert(doc.Constants[0].Desc, DeepEquals, []string{"Constant #1 without type", "Second line of description"})
	c.Assert(doc.Constants[0].Type, Equals, script.VariableType(script.VAR_TYPE_STRING))
	c.Assert(doc.Constants[0].Value, Equals, "\"\"")
	c.Assert(doc.Constants[0].Line, Equals, 15)
	c.Assert(doc.Constants[0].UnitedDesc(), Equals, "Constant #1 without type Second line of description")
	c.Assert(doc.Constants[0].IsString(), Equals, true)
	c.Assert(doc.Constants[0].IsNumber(), Equals, false)
	c.Assert(doc.Constants[0].IsBoolean(), Equals, false)
	c.Assert(doc.Constants[0].IsUnknown(), Equals, false)
	c.Assert(doc.Constants[0].TypeName(0), Equals, "String")
	c.Assert(doc.Constants[0].TypeName(1), Equals, "string")
	c.Assert(doc.Constants[0].TypeName(2), Equals, "STRING")
	c.Assert(doc.Constants[0].TypeName(3), Equals, "S")
	c.Assert(doc.Constants[0].TypeName(4), Equals, "s")

	c.Assert(doc.Constants[1], NotNil)
	c.Assert(doc.Constants[1].Name, Equals, "CONST_2")
	c.Assert(doc.Constants[1].Desc, DeepEquals, []string{"Constant #2 without type"})
	c.Assert(doc.Constants[1].Type, Equals, script.VariableType(script.VAR_TYPE_NUMBER))
	c.Assert(doc.Constants[1].Value, Equals, "0")
	c.Assert(doc.Constants[1].Line, Equals, 18)
	c.Assert(doc.Constants[1].UnitedDesc(), Equals, "Constant #2 without type")
	c.Assert(doc.Constants[1].IsString(), Equals, false)
	c.Assert(doc.Constants[1].IsNumber(), Equals, true)
	c.Assert(doc.Constants[1].IsBoolean(), Equals, false)
	c.Assert(doc.Constants[1].IsUnknown(), Equals, false)
	c.Assert(doc.Constants[1].TypeName(0), Equals, "Number")
	c.Assert(doc.Constants[1].TypeName(1), Equals, "number")
	c.Assert(doc.Constants[1].TypeName(2), Equals, "NUMBER")
	c.Assert(doc.Constants[1].TypeName(3), Equals, "N")
	c.Assert(doc.Constants[1].TypeName(4), Equals, "n")

	c.Assert(doc.Constants[2], NotNil)
	c.Assert(doc.Constants[2].Name, Equals, "CONST_3")
	c.Assert(doc.Constants[2].Desc, DeepEquals, []string{"Constant #3 without type"})
	c.Assert(doc.Constants[2].Type, Equals, script.VariableType(script.VAR_TYPE_BOOLEAN))
	c.Assert(doc.Constants[2].Value, Equals, "true")
	c.Assert(doc.Constants[2].Line, Equals, 21)
	c.Assert(doc.Constants[2].UnitedDesc(), Equals, "Constant #3 without type")
	c.Assert(doc.Constants[2].IsString(), Equals, false)
	c.Assert(doc.Constants[2].IsNumber(), Equals, false)
	c.Assert(doc.Constants[2].IsBoolean(), Equals, true)
	c.Assert(doc.Constants[2].IsUnknown(), Equals, false)
	c.Assert(doc.Constants[2].TypeName(0), Equals, "Boolean")
	c.Assert(doc.Constants[2].TypeName(1), Equals, "boolean")
	c.Assert(doc.Constants[2].TypeName(2), Equals, "BOOLEAN")
	c.Assert(doc.Constants[2].TypeName(3), Equals, "B")
	c.Assert(doc.Constants[2].TypeName(4), Equals, "b")

	c.Assert(doc.Constants[3], NotNil)
	c.Assert(doc.Constants[3].Name, Equals, "CONST_4")
	c.Assert(doc.Constants[3].Desc, DeepEquals, []string{"Constant #4 with type"})
	c.Assert(doc.Constants[3].Type, Equals, script.VariableType(script.VAR_TYPE_STRING))
	c.Assert(doc.Constants[3].Value, Equals, "\"\"")
	c.Assert(doc.Constants[3].Line, Equals, 24)

	c.Assert(doc.Constants[4], NotNil)
	c.Assert(doc.Constants[4].Name, Equals, "CONST_5")
	c.Assert(doc.Constants[4].Desc, DeepEquals, []string{"Constant #5 with type"})
	c.Assert(doc.Constants[4].Type, Equals, script.VariableType(script.VAR_TYPE_NUMBER))
	c.Assert(doc.Constants[4].Value, Equals, "\"\"")
	c.Assert(doc.Constants[4].Line, Equals, 27)

	c.Assert(doc.Constants[5], NotNil)
	c.Assert(doc.Constants[5].Name, Equals, "CONST_6")
	c.Assert(doc.Constants[5].Desc, DeepEquals, []string{"Constant #6 with type"})
	c.Assert(doc.Constants[5].Type, Equals, script.VariableType(script.VAR_TYPE_BOOLEAN))
	c.Assert(doc.Constants[5].Value, Equals, "\"\"")
	c.Assert(doc.Constants[5].Line, Equals, 30)

	// //////////////////////////////////////////////////////////////////////////////// //

	c.Assert(doc.Variables[0], NotNil)
	c.Assert(doc.Variables[0].Name, Equals, "var_1")
	c.Assert(doc.Variables[0].Desc, DeepEquals, []string{"Variable #1 without type", "Second line of description"})
	c.Assert(doc.Variables[0].Type, Equals, script.VariableType(script.VAR_TYPE_STRING))
	c.Assert(doc.Variables[0].Value, Equals, "\"\"")
	c.Assert(doc.Variables[0].Line, Equals, 40)
	c.Assert(doc.Variables[0].UnitedDesc(), Equals, "Variable #1 without type Second line of description")
	c.Assert(doc.Variables[0].IsString(), Equals, true)
	c.Assert(doc.Variables[0].IsNumber(), Equals, false)
	c.Assert(doc.Variables[0].IsBoolean(), Equals, false)
	c.Assert(doc.Variables[0].IsUnknown(), Equals, false)
	c.Assert(doc.Variables[0].TypeName(0), Equals, "String")
	c.Assert(doc.Variables[0].TypeName(1), Equals, "string")
	c.Assert(doc.Variables[0].TypeName(2), Equals, "STRING")
	c.Assert(doc.Variables[0].TypeName(3), Equals, "S")
	c.Assert(doc.Variables[0].TypeName(4), Equals, "s")

	c.Assert(doc.Variables[1], NotNil)
	c.Assert(doc.Variables[1].Name, Equals, "var_2")
	c.Assert(doc.Variables[1].Desc, DeepEquals, []string{"Variable #2 without type"})
	c.Assert(doc.Variables[1].Type, Equals, script.VariableType(script.VAR_TYPE_NUMBER))
	c.Assert(doc.Variables[1].Value, Equals, "1")
	c.Assert(doc.Variables[1].Line, Equals, 43)
	c.Assert(doc.Variables[1].UnitedDesc(), Equals, "Variable #2 without type")
	c.Assert(doc.Variables[1].IsString(), Equals, false)
	c.Assert(doc.Variables[1].IsNumber(), Equals, true)
	c.Assert(doc.Variables[1].IsBoolean(), Equals, false)
	c.Assert(doc.Variables[1].IsUnknown(), Equals, false)
	c.Assert(doc.Variables[1].TypeName(0), Equals, "Number")
	c.Assert(doc.Variables[1].TypeName(1), Equals, "number")
	c.Assert(doc.Variables[1].TypeName(2), Equals, "NUMBER")
	c.Assert(doc.Variables[1].TypeName(3), Equals, "N")
	c.Assert(doc.Variables[1].TypeName(4), Equals, "n")

	c.Assert(doc.Variables[2], NotNil)
	c.Assert(doc.Variables[2].Name, Equals, "var_3")
	c.Assert(doc.Variables[2].Desc, DeepEquals, []string{"Variable #3 without type"})
	c.Assert(doc.Variables[2].Type, Equals, script.VariableType(script.VAR_TYPE_BOOLEAN))
	c.Assert(doc.Variables[2].Value, Equals, "true")
	c.Assert(doc.Variables[2].Line, Equals, 46)
	c.Assert(doc.Variables[2].UnitedDesc(), Equals, "Variable #3 without type")
	c.Assert(doc.Variables[2].IsString(), Equals, false)
	c.Assert(doc.Variables[2].IsNumber(), Equals, false)
	c.Assert(doc.Variables[2].IsBoolean(), Equals, true)
	c.Assert(doc.Variables[2].IsUnknown(), Equals, false)
	c.Assert(doc.Variables[2].TypeName(0), Equals, "Boolean")
	c.Assert(doc.Variables[2].TypeName(1), Equals, "boolean")
	c.Assert(doc.Variables[2].TypeName(2), Equals, "BOOLEAN")
	c.Assert(doc.Variables[2].TypeName(3), Equals, "B")
	c.Assert(doc.Variables[2].TypeName(4), Equals, "b")

	c.Assert(doc.Variables[3], NotNil)
	c.Assert(doc.Variables[3].Name, Equals, "var_4")
	c.Assert(doc.Variables[3].Desc, DeepEquals, []string{"Variable #4 with type"})
	c.Assert(doc.Variables[3].Type, Equals, script.VariableType(script.VAR_TYPE_STRING))
	c.Assert(doc.Variables[3].Value, Equals, "\"\"")
	c.Assert(doc.Variables[3].Line, Equals, 49)

	c.Assert(doc.Variables[4], NotNil)
	c.Assert(doc.Variables[4].Name, Equals, "var_5")
	c.Assert(doc.Variables[4].Desc, DeepEquals, []string{"Variable #5 with type"})
	c.Assert(doc.Variables[4].Type, Equals, script.VariableType(script.VAR_TYPE_NUMBER))
	c.Assert(doc.Variables[4].Value, Equals, "\"\"")
	c.Assert(doc.Variables[4].Line, Equals, 52)

	c.Assert(doc.Variables[5], NotNil)
	c.Assert(doc.Variables[5].Name, Equals, "var_6")
	c.Assert(doc.Variables[5].Desc, DeepEquals, []string{"Variable #6 with type"})
	c.Assert(doc.Variables[5].Type, Equals, script.VariableType(script.VAR_TYPE_BOOLEAN))
	c.Assert(doc.Variables[5].Value, Equals, "\"\"")
	c.Assert(doc.Variables[5].Line, Equals, 55)

	c.Assert(doc.Variables[6], NotNil)
	c.Assert(doc.Variables[6].Name, Equals, "var_7")
	c.Assert(doc.Variables[6].Desc, DeepEquals, []string{"Variable #7 without value"})
	c.Assert(doc.Variables[6].Type, Equals, script.VariableType(script.VAR_TYPE_UNKNOWN))
	c.Assert(doc.Variables[6].Value, Equals, "")
	c.Assert(doc.Variables[6].Line, Equals, 59)
	c.Assert(doc.Variables[6].TypeName(0), Equals, "")
	c.Assert(doc.Variables[6].TypeName(1), Equals, "")
	c.Assert(doc.Variables[6].TypeName(2), Equals, "")
	c.Assert(doc.Variables[6].TypeName(3), Equals, "")
	c.Assert(doc.Variables[6].TypeName(4), Equals, "")

	c.Assert(doc.Variables[7], NotNil)
	c.Assert(doc.Variables[7].Name, Equals, "var_8")
	c.Assert(doc.Variables[7].Desc, DeepEquals, []string{"Variable #8 with multiline value"})
	c.Assert(doc.Variables[7].Type, Equals, script.VariableType(script.VAR_TYPE_STRING))
	c.Assert(doc.Variables[7].Value, Equals, "\"This is multiline value\"")
	c.Assert(doc.Variables[7].Line, Equals, 62)
	c.Assert(doc.Variables[7].IsString(), Equals, true)
	c.Assert(doc.Variables[7].IsNumber(), Equals, false)
	c.Assert(doc.Variables[7].IsBoolean(), Equals, false)
	c.Assert(doc.Variables[7].IsUnknown(), Equals, false)
	c.Assert(doc.Variables[7].TypeName(0), Equals, "String")
	c.Assert(doc.Variables[7].TypeName(1), Equals, "string")
	c.Assert(doc.Variables[7].TypeName(2), Equals, "STRING")
	c.Assert(doc.Variables[7].TypeName(3), Equals, "S")
	c.Assert(doc.Variables[7].TypeName(4), Equals, "s")

	// //////////////////////////////////////////////////////////////////////////////// //

	c.Assert(doc.Methods[0], NotNil)
	c.Assert(doc.Methods[0].Name, Equals, "method1")
	c.Assert(doc.Methods[0].Desc, DeepEquals, []string{
		"This is desc for method #1.",
		"Second line of description.",
		"Third line of description.",
		"",
		"Fourth line of description.",
	})
	c.Assert(doc.Methods[0].Arguments, HasLen, 0)
	c.Assert(doc.Methods[0].ResultCode, Equals, false)
	c.Assert(doc.Methods[0].ResultEcho, IsNil)
	c.Assert(doc.Methods[0].Example, HasLen, 0)
	c.Assert(doc.Methods[0].Line, Equals, 78)
	c.Assert(doc.Methods[0].HasArguments(), Equals, false)
	c.Assert(doc.Methods[0].HasEcho(), Equals, false)
	c.Assert(doc.Methods[0].HasEcho(), Equals, false)
	c.Assert(doc.Methods[0].UnitedDesc(), Equals, "This is desc for method #1. Second line of description. Third line of description. Fourth line of description.")

	c.Assert(doc.Methods[1], NotNil)
	c.Assert(doc.Methods[1].Name, Equals, "method2")
	c.Assert(doc.Methods[1].Desc, DeepEquals, []string{"This is desc for method #2."})
	c.Assert(doc.Methods[1].Arguments, HasLen, 4)
	// //////////////////////////////////////////////////////////////////////////////// //
	c.Assert(doc.Methods[1].Arguments[0], NotNil)
	c.Assert(doc.Methods[1].Arguments[0].Index, Equals, "1")
	c.Assert(doc.Methods[1].Arguments[0].Desc, Equals, "First argument")
	c.Assert(doc.Methods[1].Arguments[0].Type, Equals, script.VariableType(script.VAR_TYPE_STRING))
	c.Assert(doc.Methods[1].Arguments[0].IsOptional, Equals, false)
	c.Assert(doc.Methods[1].Arguments[0].IsWildcard, Equals, false)
	c.Assert(doc.Methods[1].Arguments[0].TypeName(0), Equals, "String")
	c.Assert(doc.Methods[1].Arguments[0].TypeName(1), Equals, "string")
	c.Assert(doc.Methods[1].Arguments[0].TypeName(2), Equals, "STRING")
	c.Assert(doc.Methods[1].Arguments[0].TypeName(3), Equals, "S")
	c.Assert(doc.Methods[1].Arguments[0].TypeName(4), Equals, "s")
	c.Assert(doc.Methods[1].Arguments[0].IsString(), Equals, true)
	c.Assert(doc.Methods[1].Arguments[0].IsNumber(), Equals, false)
	c.Assert(doc.Methods[1].Arguments[0].IsBoolean(), Equals, false)
	c.Assert(doc.Methods[1].Arguments[0].IsUnknown(), Equals, false)
	// //////////////////////////////////////////////////////////////////////////////// //
	c.Assert(doc.Methods[1].Arguments[1], NotNil)
	c.Assert(doc.Methods[1].Arguments[1].Index, Equals, "2")
	c.Assert(doc.Methods[1].Arguments[1].Desc, Equals, "Second argument")
	c.Assert(doc.Methods[1].Arguments[1].Type, Equals, script.VariableType(script.VAR_TYPE_NUMBER))
	c.Assert(doc.Methods[1].Arguments[1].IsOptional, Equals, false)
	c.Assert(doc.Methods[1].Arguments[1].IsWildcard, Equals, false)
	c.Assert(doc.Methods[1].Arguments[1].TypeName(0), Equals, "Number")
	c.Assert(doc.Methods[1].Arguments[1].TypeName(1), Equals, "number")
	c.Assert(doc.Methods[1].Arguments[1].TypeName(2), Equals, "NUMBER")
	c.Assert(doc.Methods[1].Arguments[1].TypeName(3), Equals, "N")
	c.Assert(doc.Methods[1].Arguments[1].TypeName(4), Equals, "n")
	c.Assert(doc.Methods[1].Arguments[1].IsString(), Equals, false)
	c.Assert(doc.Methods[1].Arguments[1].IsNumber(), Equals, true)
	c.Assert(doc.Methods[1].Arguments[1].IsBoolean(), Equals, false)
	c.Assert(doc.Methods[1].Arguments[1].IsUnknown(), Equals, false)
	// //////////////////////////////////////////////////////////////////////////////// //
	c.Assert(doc.Methods[1].Arguments[2], NotNil)
	c.Assert(doc.Methods[1].Arguments[2].Index, Equals, "3")
	c.Assert(doc.Methods[1].Arguments[2].Desc, Equals, "Third argument")
	c.Assert(doc.Methods[1].Arguments[2].Type, Equals, script.VariableType(script.VAR_TYPE_BOOLEAN))
	c.Assert(doc.Methods[1].Arguments[2].IsOptional, Equals, false)
	c.Assert(doc.Methods[1].Arguments[2].IsWildcard, Equals, false)
	c.Assert(doc.Methods[1].Arguments[2].TypeName(0), Equals, "Boolean")
	c.Assert(doc.Methods[1].Arguments[2].TypeName(1), Equals, "boolean")
	c.Assert(doc.Methods[1].Arguments[2].TypeName(2), Equals, "BOOLEAN")
	c.Assert(doc.Methods[1].Arguments[2].TypeName(3), Equals, "B")
	c.Assert(doc.Methods[1].Arguments[2].TypeName(4), Equals, "b")
	c.Assert(doc.Methods[1].Arguments[2].IsString(), Equals, false)
	c.Assert(doc.Methods[1].Arguments[2].IsNumber(), Equals, false)
	c.Assert(doc.Methods[1].Arguments[2].IsBoolean(), Equals, true)
	c.Assert(doc.Methods[1].Arguments[2].IsUnknown(), Equals, false)
	// //////////////////////////////////////////////////////////////////////////////// //
	c.Assert(doc.Methods[1].Arguments[3], NotNil)
	c.Assert(doc.Methods[1].Arguments[3].Index, Equals, "*")
	c.Assert(doc.Methods[1].Arguments[3].Desc, Equals, "Wildcard argument")
	c.Assert(doc.Methods[1].Arguments[3].Type, Equals, script.VariableType(script.VAR_TYPE_UNKNOWN))
	c.Assert(doc.Methods[1].Arguments[3].IsOptional, Equals, false)
	c.Assert(doc.Methods[1].Arguments[3].IsWildcard, Equals, true)
	c.Assert(doc.Methods[1].Arguments[3].TypeName(0), Equals, "")
	c.Assert(doc.Methods[1].Arguments[3].TypeName(1), Equals, "")
	c.Assert(doc.Methods[1].Arguments[3].TypeName(2), Equals, "")
	c.Assert(doc.Methods[1].Arguments[3].TypeName(3), Equals, "")
	c.Assert(doc.Methods[1].Arguments[3].TypeName(4), Equals, "")
	c.Assert(doc.Methods[1].Arguments[3].IsString(), Equals, false)
	c.Assert(doc.Methods[1].Arguments[3].IsNumber(), Equals, false)
	c.Assert(doc.Methods[1].Arguments[3].IsBoolean(), Equals, false)
	c.Assert(doc.Methods[1].Arguments[3].IsUnknown(), Equals, true)
	// //////////////////////////////////////////////////////////////////////////////// /1
	c.Assert(doc.Methods[1].ResultCode, Equals, true)
	c.Assert(doc.Methods[1].ResultEcho, NotNil)
	c.Assert(doc.Methods[1].ResultEcho.Desc, DeepEquals, []string{"Magic value"})
	c.Assert(doc.Methods[1].ResultEcho.Type, Equals, script.VariableType(script.VAR_TYPE_BOOLEAN))
	c.Assert(doc.Methods[1].Example, HasLen, 3)
	c.Assert(doc.Methods[1].Example[0], Equals, "if [[ -f $file ]] ; then")
	c.Assert(doc.Methods[1].Example[1], Equals, "  method2 123")
	c.Assert(doc.Methods[1].Example[2], Equals, "fi")
	c.Assert(doc.Methods[1].Line, Equals, 99)
	c.Assert(doc.Methods[1].HasArguments(), Equals, true)
	c.Assert(doc.Methods[1].HasEcho(), Equals, true)
	c.Assert(doc.Methods[1].HasExample(), Equals, true)
	c.Assert(doc.Methods[1].UnitedDesc(), Equals, "This is desc for method #2.")

	c.Assert(doc.Methods[2], NotNil)
	c.Assert(doc.Methods[2].Name, Equals, "method3")
	c.Assert(doc.Methods[2].Desc, DeepEquals, []string{"This is desc for method #3."})
	c.Assert(doc.Methods[2].Arguments, HasLen, 1)
	c.Assert(doc.Methods[2].Arguments[0], NotNil)
	c.Assert(doc.Methods[2].Arguments[0].Index, Equals, "1")
	c.Assert(doc.Methods[2].Arguments[0].Desc, Equals, "First argument")
	c.Assert(doc.Methods[2].Arguments[0].Type, Equals, script.VariableType(script.VAR_TYPE_STRING))
	c.Assert(doc.Methods[2].Arguments[0].IsOptional, Equals, true)
	c.Assert(doc.Methods[2].Arguments[0].IsWildcard, Equals, false)
	c.Assert(doc.Methods[2].ResultCode, Equals, false)
	c.Assert(doc.Methods[2].ResultEcho, IsNil)
	c.Assert(doc.Methods[2].Example, HasLen, 0)
	c.Assert(doc.Methods[2].Line, Equals, 110)
	c.Assert(doc.Methods[2].HasArguments(), Equals, true)
	c.Assert(doc.Methods[2].HasEcho(), Equals, false)
	c.Assert(doc.Methods[2].HasExample(), Equals, false)
	c.Assert(doc.Methods[2].UnitedDesc(), Equals, "This is desc for method #3.")

	c.Assert(doc.Methods[3], NotNil)
	c.Assert(doc.Methods[3].Name, Equals, "method4")
	c.Assert(doc.Methods[3].Desc, DeepEquals, []string{"This is desc for method #4."})
	c.Assert(doc.Methods[3].Arguments, HasLen, 0)
	c.Assert(doc.Methods[3].ResultCode, Equals, false)
	c.Assert(doc.Methods[3].ResultEcho, IsNil)
	c.Assert(doc.Methods[3].Example, HasLen, 0)
	c.Assert(doc.Methods[3].Line, Equals, 115)
	c.Assert(doc.Methods[3].HasArguments(), Equals, false)
	c.Assert(doc.Methods[3].HasEcho(), Equals, false)
	c.Assert(doc.Methods[3].HasExample(), Equals, false)
	c.Assert(doc.Methods[3].UnitedDesc(), Equals, "This is desc for method #4.")

	c.Assert(doc.Methods[4], NotNil)
	c.Assert(doc.Methods[4].Name, Equals, "method5")
	c.Assert(doc.Methods[4].Desc, DeepEquals, []string{"This is desc for method #5."})
	c.Assert(doc.Methods[4].Arguments, HasLen, 0)
	c.Assert(doc.Methods[4].ResultCode, Equals, false)
	c.Assert(doc.Methods[4].ResultEcho, IsNil)
	c.Assert(doc.Methods[4].Example, HasLen, 0)
	c.Assert(doc.Methods[4].Line, Equals, 121)
	c.Assert(doc.Methods[4].HasArguments(), Equals, false)
	c.Assert(doc.Methods[4].HasEcho(), Equals, false)
	c.Assert(doc.Methods[4].HasExample(), Equals, false)
	c.Assert(doc.Methods[4].UnitedDesc(), Equals, "This is desc for method #5.")

	c.Assert(doc.Methods[5], NotNil)
	c.Assert(doc.Methods[5].Name, Equals, "method6")
	c.Assert(doc.Methods[5].Desc, DeepEquals, []string{"This is desc for method #6."})
	c.Assert(doc.Methods[5].Arguments, HasLen, 0)
	c.Assert(doc.Methods[5].ResultCode, Equals, false)
	c.Assert(doc.Methods[5].ResultEcho, IsNil)
	c.Assert(doc.Methods[5].Example, HasLen, 0)
	c.Assert(doc.Methods[5].Line, Equals, 127)
	c.Assert(doc.Methods[5].HasArguments(), Equals, false)
	c.Assert(doc.Methods[5].HasEcho(), Equals, false)
	c.Assert(doc.Methods[5].HasExample(), Equals, false)
	c.Assert(doc.Methods[5].UnitedDesc(), Equals, "This is desc for method #6.")

	c.Assert(doc.Methods[6], NotNil)
	c.Assert(doc.Methods[6].Name, Equals, "method7")
	c.Assert(doc.Methods[6].Desc, DeepEquals, []string{"This is desc for method #7."})
	c.Assert(doc.Methods[6].Arguments, HasLen, 1)
	c.Assert(doc.Methods[6].Arguments[0], NotNil)
	c.Assert(doc.Methods[6].Arguments[0].Index, Equals, "1")
	c.Assert(doc.Methods[6].Arguments[0].Desc, Equals, "First argument")
	c.Assert(doc.Methods[6].Arguments[0].Type, Equals, script.VariableType(script.VAR_TYPE_UNKNOWN))
	c.Assert(doc.Methods[6].Arguments[0].IsOptional, Equals, false)
	c.Assert(doc.Methods[6].Arguments[0].IsWildcard, Equals, false)
	c.Assert(doc.Methods[6].ResultCode, Equals, false)
	c.Assert(doc.Methods[6].ResultEcho, IsNil)
	c.Assert(doc.Methods[6].Example, HasLen, 0)
	c.Assert(doc.Methods[6].Line, Equals, 134)
	c.Assert(doc.Methods[6].HasArguments(), Equals, true)
	c.Assert(doc.Methods[6].HasEcho(), Equals, false)
	c.Assert(doc.Methods[6].HasExample(), Equals, false)
	c.Assert(doc.Methods[6].UnitedDesc(), Equals, "This is desc for method #7.")

	c.Assert(doc.Methods[7], NotNil)
	c.Assert(doc.Methods[7].Name, Equals, "method8")
	c.Assert(doc.Methods[7].Desc, DeepEquals, []string{"This is desc for method #8."})
	c.Assert(doc.Methods[7].Arguments, HasLen, 0)
	c.Assert(doc.Methods[7].ResultCode, Equals, false)
	c.Assert(doc.Methods[7].ResultEcho, IsNil)
	c.Assert(doc.Methods[7].Example, HasLen, 1)
	c.Assert(doc.Methods[7].Example[0], Equals, "method8 123")
	c.Assert(doc.Methods[7].Line, Equals, 142)
	c.Assert(doc.Methods[7].HasArguments(), Equals, false)
	c.Assert(doc.Methods[7].HasEcho(), Equals, false)
	c.Assert(doc.Methods[7].HasExample(), Equals, true)
	c.Assert(doc.Methods[7].UnitedDesc(), Equals, "This is desc for method #8.")

	c.Assert(doc.Methods[8], NotNil)
	c.Assert(doc.Methods[8].Name, Equals, "method9")
	c.Assert(doc.Methods[8].Desc, DeepEquals, []string{"This is desc for method #9."})
	c.Assert(doc.Methods[8].Arguments, HasLen, 0)
	c.Assert(doc.Methods[8].ResultCode, Equals, false)
	c.Assert(doc.Methods[8].ResultEcho, IsNil)
	c.Assert(doc.Methods[8].Example, HasLen, 0)
	c.Assert(doc.Methods[8].Line, Equals, 148)
	c.Assert(doc.Methods[8].HasArguments(), Equals, false)
	c.Assert(doc.Methods[8].HasEcho(), Equals, false)
	c.Assert(doc.Methods[8].HasExample(), Equals, false)
	c.Assert(doc.Methods[8].UnitedDesc(), Equals, "This is desc for method #9.")
}
