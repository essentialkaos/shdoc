package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2015 Essential Kaos                         //
//      Essential Kaos Open Source License <http://essentialkaos.com/ekol?en>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/essentialkaos/ek/arg"
	"github.com/essentialkaos/ek/fmtc"
	"github.com/essentialkaos/ek/fmtutil"
	"github.com/essentialkaos/ek/fsutil"
	"github.com/essentialkaos/ek/usage"

	. "github.com/essentialkaos/shdoc/parser"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "SHDoc"
	VER  = "0.1.2"
	DESC = "Tool for viewing and exporting docs for shell scripts"
)

const (
	ARG_OUTPUT   = "o:output"
	ARG_TEMPLATE = "t:template"
	ARG_NAME     = "n:name"
	ARG_NO_COLOR = "nc:no-color"
	ARG_HELP     = "h:help"
	ARG_VER      = "v:version"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var argMap = arg.Map{
	ARG_OUTPUT:   &arg.V{},
	ARG_TEMPLATE: &arg.V{Value: "/usr/local/share/shdoc/templates/markdown.tpl"},
	ARG_NAME:     &arg.V{},
	ARG_NO_COLOR: &arg.V{Type: arg.BOOL},
	ARG_HELP:     &arg.V{Type: arg.BOOL, Alias: "u:usage"},
	ARG_VER:      &arg.V{Type: arg.BOOL, Alias: "ver"},
}

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	args, errs := arg.Parse(argMap)

	if len(errs) != 0 {
		fmtc.Println("Arguments parsing errors:")

		for _, err := range errs {
			fmtc.Printf("  %s\n", err.Error())
		}

		os.Exit(1)
	}

	if arg.GetB(ARG_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if arg.GetB(ARG_VER) {
		showAbout()
		return
	}

	if arg.GetB(ARG_HELP) || len(args) == 0 {
		showUsage()
		return
	}

	switch len(args) {
	case 1:
		process(args[0], "")
	case 2:
		process(args[0], args[1])
	default:
		showUsage()
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

func process(file string, pattern string) {
	if !fsutil.IsExist(file) {
		fmtc.Printf("{r}File %s is not exist{!}\n", file)
		os.Exit(1)
	}

	if !fsutil.IsReadable(file) {
		fmtc.Printf("{r}File %s is not readable{!}\n", file)
		os.Exit(1)
	}

	if !fsutil.IsNonEmpty(file) {
		fmtc.Printf("{r}File %s is empty{!}\n", file)
		os.Exit(1)
	}

	doc, errs := Parse(file)

	if len(errs) != 0 {
		fmtc.Println("{r}Shell script docs parsing errors:{!}")

		for _, e := range errs {
			fmtc.Printf("  {r}%s{!}\n", e.Error())
		}

		os.Exit(1)
	}

	if !doc.IsValid() {
		fmtc.Printf("{y}File %s doesn't contains documentation.{!}\n", file)
		os.Exit(2)
	}

	if arg.GetS(ARG_NAME) != "" {
		doc.Title = arg.GetS(ARG_NAME)
	}

	if arg.GetS(ARG_OUTPUT) == "" {
		if pattern == "" {
			simpleRender(doc)
		} else {
			findInfo(doc, pattern)
		}
	} else {
		renderTemplate(doc)
	}
}

// findInfo search geven pattern in entities names
func findInfo(doc *Document, pattern string) {
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

// simpleRender print all document info to console
func simpleRender(doc *Document) {
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
				fmtc.NewLine()
				fmtc.NewLine()
			}
		}
	}

	fmtutil.Separator(false)
}

// renderTemplate read template and render to file
func renderTemplate(doc *Document) {
	if !fsutil.CheckPerms("FRS", arg.GetS(ARG_TEMPLATE)) {
		fmtc.Printf("{r}Can't read template %s - file is not exist or empty.{!}\n", arg.GetS(ARG_TEMPLATE))
		os.Exit(1)
	}

	if fsutil.IsExist(arg.GetS(ARG_OUTPUT)) {
		os.Remove(arg.GetS(ARG_OUTPUT))
	}

	fd, err := os.OpenFile(arg.GetS(ARG_OUTPUT), os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		printErrorAndExit(err)
	}

	defer fd.Close()

	tpl, err := ioutil.ReadFile(arg.GetS(ARG_TEMPLATE))

	if err != nil {
		printErrorAndExit(err)
	}

	t := template.New("Template")
	t, err = t.Parse(string(tpl[:]))

	err = t.Execute(fd, doc)

	if err != nil {
		printErrorAndExit(err)
	}
}

// renderConstant print constant info to console
func renderConstant(c *Variable) {
	fmtc.Printf("{s}%4d:{!} {m*}%s{!} {s}={!} %s\n", c.Line, c.Name, c.Value)
	fmtc.Printf("      %s "+getVarTypeDesc(c.Type)+"\n", strings.Join(c.Desc, " "))
}

// renderMethod print variable info to console
func renderVariable(v *Variable) {
	fmtc.Printf("{s}%4d:{!} {c*}%s{!} {s}={!} %s\n", v.Line, v.Name, v.Value)
	fmtc.Printf("      %s "+getVarTypeDesc(v.Type)+"\n", v.UnitedDesc())
}

// renderMethod print method info to console
func renderMethod(m *Method, showExamples bool) {
	fmtc.Printf("{s}%4d:{!} {b*}%s{!} {s}-{!} %s\n", m.Line, m.Name, m.UnitedDesc())

	if len(m.Arguments) != 0 {
		fmtc.NewLine()

		for _, a := range m.Arguments {
			switch {
			case a.IsOptional:
				fmtc.Printf("  {s}%s.{!} %s "+getVarTypeDesc(a.Type)+" {s}[Optional]{!}\n", a.Index, a.Desc)
			case a.IsWildcard:
				fmtc.Printf("  {s}%s.{!} %s\n", a.Index, a.Desc)
			default:
				fmtc.Printf("  {s}%s.{!} %s "+getVarTypeDesc(a.Type)+"\n", a.Index, a.Desc)
			}
		}
	}

	if m.ResultCode {
		fmtc.NewLine()
		fmtc.Printf("  {*}Code:{!} 0 - ok, 1 - not ok\n")
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
			fmtc.Printf("    %s\n", l)
		}
	}
}

// getVarTypeDesc return type description
func getVarTypeDesc(t VariableType) string {
	switch t {
	case VAR_TYPE_STRING:
		return "{b}(String){!}"
	case VAR_TYPE_NUMBER:
		return "{y}(Number){!}"
	case VAR_TYPE_BOOLEAN:
		return "{g}(Boolean){!}"
	default:
		return ""
	}
}

// printErrorAndExit print error mesage and exit with exit code 1
func printErrorAndExit(err error) {
	fmtc.Printf("{r}%s{!}\n", err.Error())
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func showUsage() {
	info := usage.NewInfo("", "file")

	info.AddOption(ARG_OUTPUT, "Path to output file", "file")
	info.AddOption(ARG_TEMPLATE, "Path to template file", "file")
	info.AddOption(ARG_NAME, "Overwrite default name", "name")
	info.AddOption(ARG_NO_COLOR, "Disable colors in output")
	info.AddOption(ARG_HELP, "Show this help message")
	info.AddOption(ARG_VER, "Show version")

	info.AddExample(
		"script.sh",
		"Parse shell script and show docs in console",
	)

	info.AddExample(
		"script.sh -t path/to/template.tpl -o my_script.md",
		"Parse shell script and save docs using given export template",
	)

	info.AddExample(
		"script.sh someEntity",
		"Parse shell script and show docs for some constant, variable or method",
	)

	info.Render()
}

func showAbout() {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2009,
		Owner:   "Essential Kaos",
		License: "Essential Kaos Open Source License <https://essentialkaos.com/ekol?en>",
	}

	about.Render()
}
