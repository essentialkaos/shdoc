package script

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"testing"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type ScriptSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&ScriptSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *ScriptSuite) TestValidation(c *C) {
	sh1 := &Document{}
	sh2 := &Document{Constants: make([]*Variable, 0)}
	sh3 := &Document{Variables: make([]*Variable, 0)}
	sh4 := &Document{Methods: make([]*Method, 0)}
	sh5 := &Document{Constants: make([]*Variable, 1)}
	sh6 := &Document{Variables: make([]*Variable, 1)}
	sh7 := &Document{Methods: make([]*Method, 1)}

	c.Assert(sh1.IsValid(), Equals, false)
	c.Assert(sh2.IsValid(), Equals, false)
	c.Assert(sh3.IsValid(), Equals, false)
	c.Assert(sh4.IsValid(), Equals, false)
	c.Assert(sh5.IsValid(), Equals, true)
	c.Assert(sh6.IsValid(), Equals, true)
	c.Assert(sh7.IsValid(), Equals, true)
}

func (s *ScriptSuite) TestNil(c *C) {
	var d *Document
	var a *Argument
	var v *Variable
	var m *Method

	c.Assert(d.IsValid(), Equals, false)
	c.Assert(d.HasAbout(), Equals, false)
	c.Assert(d.HasConstants(), Equals, false)
	c.Assert(d.HasVariables(), Equals, false)
	c.Assert(d.HasMethods(), Equals, false)

	c.Assert(a.TypeName(0), Equals, "")
	c.Assert(a.IsString(), Equals, false)
	c.Assert(a.IsNumber(), Equals, false)
	c.Assert(a.IsBoolean(), Equals, false)
	c.Assert(a.IsUnknown(), Equals, false)

	c.Assert(v.TypeName(0), Equals, "")
	c.Assert(v.IsString(), Equals, false)
	c.Assert(v.IsNumber(), Equals, false)
	c.Assert(v.IsBoolean(), Equals, false)
	c.Assert(v.IsUnknown(), Equals, false)
	c.Assert(v.UnitedDesc(), Equals, "")

	c.Assert(m.HasArguments(), Equals, false)
	c.Assert(m.HasEcho(), Equals, false)
	c.Assert(m.HasExample(), Equals, false)
	c.Assert(m.UnitedDesc(), Equals, "")
}

func (s *ScriptSuite) TestHelpers(c *C) {
	d := &Document{
		Title:     "Test title",
		About:     []string{"Test", "description"},
		Constants: []*Variable{&Variable{}},
		Variables: []*Variable{&Variable{}},
		Methods:   []*Method{&Method{}},
	}

	c.Assert(d.HasAbout(), Equals, true)
	c.Assert(d.HasConstants(), Equals, true)
	c.Assert(d.HasVariables(), Equals, true)
	c.Assert(d.HasMethods(), Equals, true)

	a1 := &Argument{"1", "A1", VAR_TYPE_UNKNOWN, false, false}
	a2 := &Argument{"2", "A2", VAR_TYPE_STRING, false, false}
	a3 := &Argument{"3", "A3", VAR_TYPE_NUMBER, false, false}
	a4 := &Argument{"4", "A4", VAR_TYPE_BOOLEAN, false, false}
	a5 := &Argument{"*", "A5", VAR_TYPE_UNKNOWN, true, true}

	c.Assert(a1.TypeName(VAR_MOD_DEFAULT), Equals, "")
	c.Assert(a2.TypeName(VAR_MOD_DEFAULT), Equals, "String")
	c.Assert(a3.TypeName(VAR_MOD_DEFAULT), Equals, "Number")
	c.Assert(a4.TypeName(VAR_MOD_DEFAULT), Equals, "Boolean")
	c.Assert(a5.TypeName(VAR_MOD_DEFAULT), Equals, "")

	c.Assert(a1.IsUnknown(), Equals, true)
	c.Assert(a2.IsString(), Equals, true)
	c.Assert(a3.IsNumber(), Equals, true)
	c.Assert(a4.IsBoolean(), Equals, true)

	v1 := &Variable{"1", []string{"V1", "", "D"}, VAR_TYPE_UNKNOWN, "v1", 1}
	v2 := &Variable{"2", []string{"V2"}, VAR_TYPE_STRING, "v2", 2}
	v3 := &Variable{"3", []string{"V3"}, VAR_TYPE_NUMBER, "v3", 3}
	v4 := &Variable{"4", []string{"V4"}, VAR_TYPE_BOOLEAN, "v4", 4}

	c.Assert(v1.TypeName(VAR_MOD_DEFAULT), Equals, "")
	c.Assert(v2.TypeName(VAR_MOD_DEFAULT), Equals, "String")
	c.Assert(v3.TypeName(VAR_MOD_DEFAULT), Equals, "Number")
	c.Assert(v4.TypeName(VAR_MOD_DEFAULT), Equals, "Boolean")

	c.Assert(v1.IsUnknown(), Equals, true)
	c.Assert(v2.IsString(), Equals, true)
	c.Assert(v3.IsNumber(), Equals, true)
	c.Assert(v4.IsBoolean(), Equals, true)

	c.Assert(v1.UnitedDesc(), Equals, "V1 D")

	m := &Method{
		Name: "m1",
		Desc: []string{"M1", "", "D"},
		Arguments: []*Argument{
			&Argument{"1", "A1", VAR_TYPE_UNKNOWN, false, false},
		},
		ResultCode: true,
		ResultEcho: &Variable{"1", []string{"V1"}, VAR_TYPE_STRING, "v1", 1},
		Example:    []string{"example"},
		Line:       15,
	}

	c.Assert(m.HasArguments(), Equals, true)
	c.Assert(m.HasEcho(), Equals, true)
	c.Assert(m.HasExample(), Equals, true)
	c.Assert(m.UnitedDesc(), Equals, "M1 D")
}
