package terminal

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"regexp"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fmtutil"

	"github.com/essentialkaos/shdoc/script"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var varExtractRegex = regexp.MustCompile(`\$\{*[^\}\n\r]+\}*`)

// ////////////////////////////////////////////////////////////////////////////////// //

// Render prints script info into terminal
func Render(doc *script.Document, pattern string) error {
	if pattern != "" {
		renderPart(doc, pattern)
	} else {
		renderAll(doc)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// renderAll renders all document info
func renderAll(doc *script.Document) {
	if doc.HasAbout() {
		fmtutil.Separator(false, "ABOUT")

		for _, l := range doc.About {
			fmtc.Printf("  %s\n", l)
		}
	}

	if doc.HasConstants() {
		fmtutil.Separator(false, "CONSTANTS")

		totalConstants := len(doc.Constants)

		for i, c := range doc.Constants {
			renderConstant(c)

			if i < totalConstants-1 {
				fmtc.NewLine()
			}
		}
	}

	if doc.HasVariables() {
		fmtutil.Separator(false, "GLOBAL VARIABLES")

		totalVariables := len(doc.Variables)

		for i, v := range doc.Variables {
			renderVariable(v)

			if i < totalVariables-1 {
				fmtc.NewLine()
			}
		}
	}

	if doc.HasMethods() {
		fmtutil.Separator(false, "METHODS")

		totalMethods := len(doc.Methods)

		for i, m := range doc.Methods {
			renderMethod(m, false)

			if i < totalMethods-1 {
				fmtc.Println("\n{s-}" + strings.Repeat("-", 88) + "{!}")
				fmtc.NewLine()
			}
		}
	}

	fmtutil.Separator(false)
}

// renderPart renders only part of document (method/variable/constant)
func renderPart(doc *script.Document, pattern string) {
	fmtc.NewLine()

	if doc.Constants != nil {
		for _, c := range doc.Constants {
			if strings.Contains(c.Name, pattern) {
				renderConstant(c)
				fmtc.NewLine()
			}
		}
	}

	if doc.Variables != nil {
		for _, v := range doc.Variables {
			if strings.Contains(v.Name, pattern) {
				renderVariable(v)
				fmtc.NewLine()
			}
		}
	}

	if doc.Methods != nil {
		for _, m := range doc.Methods {
			if strings.Contains(m.Name, pattern) {
				renderMethod(m, true)
				fmtc.NewLine()
			}
		}
	}
}

// renderConstant prints constant info to console
func renderConstant(c *script.Variable) {
	fmtc.Printf("{s}%4d:{!} {m*}%s{!} {s}={!} "+colorizeValue(c.Value)+" "+getVarTypeDesc(c.Type)+"\n", c.Line, c.Name)
	fmtc.Printf("      %s\n", c.UnitedDesc())
}

// renderMethod prints variable info to console
func renderVariable(v *script.Variable) {
	fmtc.Printf("{s}%4d:{!} {c*}%s{!} {s}={!} "+colorizeValue(v.Value)+" "+getVarTypeDesc(v.Type)+"\n", v.Line, v.Name)
	fmtc.Printf("      %s\n", v.UnitedDesc())
}

// renderMethod prints method info to console
func renderMethod(m *script.Method, showExamples bool) {
	fmtc.Printf("{s}%4d:{!} {b*}%s{!} {s}-{!} %s\n", m.Line, m.Name, m.UnitedDesc())

	if len(m.Arguments) != 0 {
		fmtc.NewLine()

		for _, a := range m.Arguments {
			switch {
			case a.IsOptional:
				fmtc.Printf("  {s-}%2s.{!} %s "+getVarTypeDesc(a.Type)+" {s-}[Optional]{!}\n", a.Index, a.Desc)
			case a.IsWildcard:
				fmtc.Printf("  {s-}%2s.{!} %s\n", a.Index, a.Desc)
			default:
				fmtc.Printf("  {s-}%2s.{!} %s "+getVarTypeDesc(a.Type)+"\n", a.Index, a.Desc)
			}
		}
	}

	if m.ResultCode {
		fmtc.NewLine()
		fmtc.Printf("    {*}Code:{!} 0 - ok, 1 - not ok\n")
	}

	if m.ResultEcho != nil {
		fmtc.NewLine()
		fmtc.Printf("  {*}Echo:{!} %s "+getVarTypeDesc(m.ResultEcho.Type)+"\n", strings.Join(m.ResultEcho.Desc, " "))
	}

	if m.Example != nil && showExamples {
		fmtc.NewLine()
		fmtc.Println("  {*}Example:{!}")
		fmtc.NewLine()

		for _, l := range m.Example {
			fmtc.Printf("    {s}%s{!}\n", l)
		}
	}
}

// colorizeValue adds color tags based on variable value
func colorizeValue(value string) string {
	if !varExtractRegex.MatchString(value) {
		return value
	}

	return varExtractRegex.ReplaceAllStringFunc(value, func(v string) string {
		return "{g}" + v + "{!}"
	})
}

// getVarTypeDesc returns type description
func getVarTypeDesc(t script.VariableType) string {
	switch t {
	case script.VAR_TYPE_STRING:
		return "{b}({&}String{!&}){!}"
	case script.VAR_TYPE_NUMBER:
		return "{y}({&}Number{!&}){!}"
	case script.VAR_TYPE_BOOLEAN:
		return "{g}({&}Boolean{!&}){!}"
	default:
		return ""
	}
}
